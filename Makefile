BACKEND_WASM=build/backend.wasm
FRONTEND_WASM=web/app.wasm
WASM_EXEC=web/wasm_exec.js

.PHONY: backend frontend wasm_exec all run clean

backend: $(BACKEND_WASM)

$(BACKEND_WASM): $(shell find backend -name '*.go') go.mod
	mkdir -p build
	GOOS=wasip1 GOARCH=wasm go build -o $(BACKEND_WASM) ./backend

frontend: $(FRONTEND_WASM)

$(FRONTEND_WASM): $(shell find frontend -name '*.go') go.mod
	mkdir -p web
	GOOS=js GOARCH=wasm go build -o $(FRONTEND_WASM) ./frontend

wasm_exec:
	mkdir -p web
	if [ -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_EXEC); \
	else \
		cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WASM_EXEC); \
	fi

all: backend frontend wasm_exec

run: all
	go run ./host

clean:
	rm -f $(BACKEND_WASM) $(FRONTEND_WASM) $(WASM_EXEC)
	rm -rf build
