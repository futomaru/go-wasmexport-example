# go:wasmexport 最小サンプル

`//go:wasmexport` でエクスポートした Go 関数を、wazero を使ってホスト側から呼び出す最小構成です。guest でビルドした `add.wasm` を host で埋め込み、`Add(1, 41)` の結果を確認できます。

## プロジェクト構成

- `guest/add.go` … WebAssembly モジュール（`Add` 関数をエクスポート）
- `host/main.go` … wazero ランタイム／モジュールを `addClient` にまとめ、`Add` 呼び出しをカプセル化
- `add.wasm` … guest のビルド成果物（`//go:embed` で埋め込み）

## 前提

- Go 1.25.x 以降（`go env GOVERSION` を合わせてください）
- `GOOS=wasip1` `GOARCH=wasm` 向けビルドが可能な環境

## ビルド手順

1. guest を Reactor モジュール（`-buildmode=c-shared`）としてビルドし、`host` ディレクトリに `add.wasm` を生成します。

   ```bash
   GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o host/add.wasm ./guest
   ```

2. host を実行して、埋め込んだモジュール経由で `Add` を呼び出します。

   ```bash
   go run ./host
   ```

   期待される出力：

   ```
   Add(1, 41) = 42
   ```

## 補足

- `add.wasm` は `host/add.wasm` に配置してください（`//go:embed add.wasm` の前提）。Reactor モジュールの初期化関数（`_initialize`）はホスト側で自動実行します。
- host 側は内部で WASI を初期化し、`addClient` が runtime／module／function を一括管理します。ホストから他の関数を呼びたい場合は、`addClient` を拡張するだけで済みます。
- 複数バージョンの Go を使い分けたい場合は、`golang.org/dl` などを利用して目的のバージョンをインストールしてください。
