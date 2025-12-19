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

	WafWasmPath := "./modules/scloud-eg-waf/target/wasm32-wasip1/release/scloud_eg_waf.wasm"
	RateLimitWasmPath := "./modules/scloud-eg-rate-limit/target/wasm32-wasip1/release/scloud_eg_rate_limit.wasm"

	//========================
	// WASM files bytes
	//========================
	WafWasmBytes, err := os.ReadFile(WafWasmPath)
	if err != nil {
		log.Fatalf("Failed to read WASM file: %v", err)
	}
	RateLimitWasmBytes, err := os.ReadFile(RateLimitWasmPath)
	if err != nil {
		log.Fatalf("Failed to read WASM file: %v", err)
	}

	//========================
	// WASM files modules
	//========================
	WafWasmModule, err := runtime.LoadWasmModule(ctx, r, WafWasmBytes)
	if err != nil {
		log.Fatalf("Failed to load WASM module: %v", err)
	}
	RateLimitWasmModule, err := runtime.LoadWasmModule(ctx, r, RateLimitWasmBytes)
	if err != nil {
		log.Fatalf("Failed to load WASM module: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		requestData := map[string]string{
			"path":   req.URL.Path,
			"method": req.Method,
		}
		jsonData, _ := json.Marshal(requestData)

		decision, err := runtime.CallWaf(ctx, WafWasmModule, jsonData)
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

	http.HandleFunc("/rate-limit", func(w http.ResponseWriter, req *http.Request) {

		decision, err := runtime.CallRatelimit(ctx, RateLimitWasmModule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "WAF error: %v", err)
			return
		}

		if decision == 1 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "RateLimit OK")
		} else {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Request blocked by RateLimit")
		}
	})

	fmt.Println("Edge Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
