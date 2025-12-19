package runtime

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func LoadWasmModule(ctx context.Context, r wazero.Runtime, wasmBytes []byte) (api.Module, error) {
	compiledModule, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile WASM module: %w", err)
	}

	module, err := r.InstantiateModule(ctx, compiledModule, wazero.NewModuleConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate WASM module: %w", err)
	}

	return module, nil
}

func CallWaf(ctx context.Context, module api.Module, _request []byte) (int64, error) {
	handle := module.ExportedFunction("handle")
	if handle == nil {
		return -1, fmt.Errorf("handle function not found in WASM module")
	}

	results, err := handle.Call(ctx, 0, 0)
	if err != nil {
		return -1, fmt.Errorf("failed to call WASM function: %w", err)
	}

	return int64(results[0]), nil
}
