package main

// //go:wasmexport のすぐ下に「エクスポートしたい関数」を書く
// エクスポート名は "Add"（ホスト側からこの名前で呼び出す）

//go:wasmexport Add
func Add(a int32, b int32) int32 {
	return a + b
}

// Reactorモジュールにするので main は何もしなくてOK。
// （-buildmode=c-shared の場合、mainは自動で実行されない）
func main() {}
