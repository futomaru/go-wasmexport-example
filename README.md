# wasm フルスタック最小サンプル

Go でフロントエンド・バックエンドを両方とも WebAssembly 化した最小構成アプリ。

## ゴール
WebAssembly（ブラウザ／サーバ）のフルスタック構成を最小構成で体験する
DB（PostgreSQL）のレコードを取得し、ブラウザで表示できることを確認する
「ホスト（ネイティブ）」と「WASMモジュール」の責務分離を理解する
## 非ゴール
認証／認可、トランザクション管理など本番レベルの機能
複雑なフロントエンド UI
高トラフィックや耐障害性を意識したアーキテクチャ設計

## リクエストフロー
1. ユーザがブラウザで index.html にアクセス
2. ブラウザが wasm_exec.js と app.wasm を読み込む
3. フロントエンド WASM が起動し、/api/records に HTTP GET
4. ホストの HTTP サーバがリクエストを受け取る
5. ホストがバックエンド WASM モジュールを呼び出し、「レコード一覧を取得する」関数を実行
6. バックエンド WASM モジュール内のロジックがホスト経由で DB からデータを取得し、JSON 文字列を返す
7. ホストがレスポンスとして JSON を返却
8. フロントエンド WASM が JSON をパースし、ブラウザ DOM に表示

## 技術スタック
言語: Go 1.22 以降を想定
WASM ランタイム: wazero
WASI: wasip1
DB: PostgreSQL
HTTP サーバ: 標準 net/http  
ブラウザランタイム: wasm_exec.js（Go 同梱）

## 構成
- **バックエンド wasm**: `backend` ディレクトリのコードを WASI (wasip1) 向けにビルド。apiが叩かれると、PostgreSQL へ接続し、JSON を返す。
- **ホスト / API サーバ**: `host/main.go` が wazero で wasm を実行する。簡易的なサーバーを立てる。
- **フロントエンド wasm**: `frontend` ディレクトリのコードを `GOOS=js GOARCH=wasm` でビルド。`/api/records` を叩いて DOM に情報を表示する。

## ディレクトリ
- `backend` … WASI で動かす DB モジュール。apiが叩かれたときのレスポンスを返す。
- `build` …　backendを wasm化したwasmファイルを置く。
- `host` … HTTP サーバ
- `frontend` … ブラウザ用
- `web` … 公開するディレクトリ。`app.wasm`, `index.html`, `wasm_exec.js` の3つを置く。

## ビルド / 実行

Makefile から簡単にビルドできます:

```
make backend      # backend wasm を build/backend.wasm に
make frontend     # frontend wasm を web/app.wasm に
make wasm_exec    # wasm_exec.js を web/wasm_exec.js にコピー
make all         # backend + frontend をまとめてビルド
make run         # build + go run ./host
```

`make run` は backend wasm を `build/backend.wasm` に配置し、host サーバを `:8080` で起動します。`DATABASE_URL` を環境変数で渡すと、WASM モジュールも同じ値を利用して PostgreSQL に接続します。未設定時は内部でサンプルデータを返します。

## PostgreSQL の準備例

```sql
CREATE TABLE records (
  id serial PRIMARY KEY,
  title text NOT NULL,
  content text NOT NULL
);

INSERT INTO records (title, content) VALUES
  ('entry one', 'first record'),
  ('entry two', 'second record');
```

ホスト環境の `DATABASE_URL` に `postgres://user:pass@host:port/dbname?sslmode=disable` などを設定してから `make run` してください。
