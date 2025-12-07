package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	backenddb "example.com/go-wasm-fullstack/backend/db"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"gorm.io/gorm"
)

const (
	backendWasm = "build/backend.wasm"
	appWasm     = "web/app.wasm"
	wasmExec    = "web/wasm_exec.js"
	indexHTML   = "web/index.html"
)

type server struct {
	runtime wazero.Runtime
	module  wazero.CompiledModule
	db      *gorm.DB
}

func main() {
	loadEnv(".env")

	ctx := context.Background()
	runtime := wazero.NewRuntime(ctx)
	defer runtime.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	raw, err := os.ReadFile(backendWasm)
	if err != nil {
		panic(err)
	}
	module, err := runtime.CompileModule(ctx, raw)
	if err != nil {
		panic(err)
	}
	defer module.Close(ctx)

	db, closeDB, err := backenddb.Open()
	if err != nil {
		panic(err)
	}
	defer closeDB()

	srv := &server{runtime: runtime, module: module, db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/records", srv.recordsHandler)
	mux.HandleFunc("/app.wasm", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, appWasm)
	})
	mux.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, wasmExec)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexHTML)
	})

	http.ListenAndServe(":8080", mux)
}

func (s *server) recordsHandler(w http.ResponseWriter, r *http.Request) {
	records, err := backenddb.LoadRecords(r.Context(), s.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload, _ := json.Marshal(records)
	var stdout bytes.Buffer
	cfg := wazero.NewModuleConfig().WithStdout(&stdout).WithStdin(bytes.NewReader(payload))

	mod, err := s.runtime.InstantiateModule(r.Context(), s.module, cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer mod.Close(r.Context())

	out := bytes.TrimSpace(stdout.Bytes())
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		os.Setenv(strings.TrimSpace(parts[0]), strings.Trim(strings.TrimSpace(parts[1]), "\""))
	}
}
