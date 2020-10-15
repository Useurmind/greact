build_test_app:
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" test_app/assets/
	go build -o bin/test_app_server test_app/server/main.go
	GOOS=js GOARCH=wasm go build -o test_app/assets/index.wasm ./test_app/wasm/main.go

run_test_app:
	bin/test_app_server 127.0.0.1:8080 test_app/assets