package main

import (
	"encoding/json"
	"os"
	"sort"
)

type Record struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	var records []Record
	json.NewDecoder(os.Stdin).Decode(&records)
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})
	json.NewEncoder(os.Stdout).Encode(records)
}
