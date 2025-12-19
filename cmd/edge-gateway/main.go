package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"2scloud-edge-gateway/internal/runtime"

	"github.com/tetratelabs/wazero"
)

func main() {
	ctx := context.Background()

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	wasmPath := "../../modules/scloud-eg-waf/target/wasm32-wasip1/release/scloud_eg_waf.wasm"
	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		log.Fatalf("Failed to read WASM file: %v", err)
	}

	wasmModule, err := runtime.LoadWasmModule(ctx, r, wasmBytes)
	if err != nil {
		log.Fatalf("Failed to load WASM module: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		requestData := map[string]string{
			"path":   req.URL.Path,
			"method": req.Method,
		}
		jsonData, _ := json.Marshal(requestData)

		decision, err := runtime.CallWaf(ctx, wasmModule, jsonData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "WAF error: %v", err)
			return
		}

		if decision == 0 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Request allowed")
		} else {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Request blocked by WAF")
		}
	})

	fmt.Println("Edge Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
