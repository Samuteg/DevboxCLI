package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

var jsonOutput bool

type JSONResult struct {
	Success bool        `json:"success"`
	Command string      `json:"command"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func printJSON(result JSONResult) {
	out, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"success":false,"command":"%s","message":"json marshal error"}\n`, result.Command)
		return
	}
	fmt.Println(string(out))
}

func initJSONFlag() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "saída em formato JSON (máquina)")
}
