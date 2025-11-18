package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// add.wasm をバイナリとして埋め込む例
//
//go:embed add.wasm
var wasmBytes []byte

func main() {
	ctx := context.Background()

	client, err := newAddClient(ctx, wasmBytes)
	if err != nil {
		log.Fatalf("failed to prepare add.wasm: %v", err)
	}
	defer client.Close(ctx)

	const a, b = 1, 41

	sum, err := client.Add(ctx, uint32(a), uint32(b))
	if err != nil {
		log.Fatalf("failed to call Add: %v", err)
	}

	fmt.Printf("Add(%d, %d) = %d\n", a, b, sum)
}

type addClient struct {
	runtime wazero.Runtime
	module  api.Module
	addFn   api.Function
}

func newAddClient(ctx context.Context, wasm []byte) (*addClient, error) {
	runtime := wazero.NewRuntime(ctx)
	if err := instantiateWASI(ctx, runtime); err != nil {
		runtime.Close(ctx)
		return nil, err
	}

	module, err := runtime.Instantiate(ctx, wasm)
	if err != nil {
		runtime.Close(ctx)
		return nil, fmt.Errorf("instantiate module: %w", err)
	}
	if err := initializeReactor(ctx, module); err != nil {
		module.Close(ctx)
		runtime.Close(ctx)
		return nil, err
	}

	addFn := module.ExportedFunction("Add")
	if addFn == nil {
		module.Close(ctx)
		runtime.Close(ctx)
		return nil, fmt.Errorf("exported function Add not found")
	}

	return &addClient{
		runtime: runtime,
		module:  module,
		addFn:   addFn,
	}, nil
}

func (c *addClient) Close(ctx context.Context) {
	if c == nil {
		return
	}
	if c.module != nil {
		if err := c.module.Close(ctx); err != nil {
			log.Printf("failed to close module: %v", err)
		}
	}
	if c.runtime != nil {
		if err := c.runtime.Close(ctx); err != nil {
			log.Printf("failed to close runtime: %v", err)
		}
	}
}

func (c *addClient) Add(ctx context.Context, a, b uint32) (uint32, error) {
	results, err := c.addFn.Call(ctx, uint64(a), uint64(b))
	if err != nil {
		return 0, fmt.Errorf("call Add(%d, %d): %w", a, b, err)
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("Add returned no results")
	}
	return uint32(results[0]), nil
}

func instantiateWASI(ctx context.Context, runtime wazero.Runtime) error {
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		return fmt.Errorf("instantiate WASI: %w", err)
	}
	return nil
}

func initializeReactor(ctx context.Context, module api.Module) error {
	initFn := module.ExportedFunction("_initialize")
	if initFn == nil {
		return fmt.Errorf("_initialize function not found")
	}
	if _, err := initFn.Call(ctx); err != nil {
		return fmt.Errorf("initialize reactor: %w", err)
	}
	return nil
}
