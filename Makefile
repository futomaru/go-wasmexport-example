WASM=host/add.wasm
GO_WASM_ENV=GOOS=wasip1 GOARCH=wasm

.PHONY: wasm run clean

wasm: $(WASM)

$(WASM): $(shell find guest -name '*.go') go.mod
	$(GO_WASM_ENV) go build -buildmode=c-shared -o $(WASM) ./guest

run: $(WASM)
	go run ./host

clean:
	rm -f $(WASM)
