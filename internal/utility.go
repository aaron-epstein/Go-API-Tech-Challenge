package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Write(w http.ResponseWriter, v any) error {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		Out("ERROR", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	return err
}

func Out(a ...any) {
	fmt.Println(a...)
}

func Outf(f string, a ...any) {
	fmt.Printf(f+"\n", a...)
}
