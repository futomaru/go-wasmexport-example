//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type record struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	successHandler js.Func
	errorHandler   js.Func
)

func main() {
	doc := js.Global().Get("document")
	status := doc.Call("getElementById", "status")
	list := doc.Call("getElementById", "records")

	successHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		resp := args[0]
		return resp.Call("json").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			renderRecords(doc, status, list, args[0])
			return nil
		}))
	})

	errorHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		status.Set("textContent", fmt.Sprintf("fetch error: %v", args[0]))
		return nil
	})

	status.Set("textContent", "loading records...")
	js.Global().Call("fetch", "/api/records").Call("then", successHandler).Call("catch", errorHandler)

	select {}
}

func renderRecords(doc, status, list js.Value, payload js.Value) {
	raw := js.Global().Get("JSON").Call("stringify", payload).String()
	var records []record
	if err := json.Unmarshal([]byte(raw), &records); err != nil {
		status.Set("textContent", fmt.Sprintf("decode error: %v", err))
		return
	}

	list.Set("innerHTML", "")
	for _, record := range records {
		item := doc.Call("createElement", "li")
		item.Set("textContent", fmt.Sprintf("#%d %s — %s", record.ID, record.Title, record.Content))
		list.Call("appendChild", item)
	}
	status.Set("textContent", fmt.Sprintf("%d records", len(records)))
}
