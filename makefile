run:
	go run ./cmd/cli/

tidy: 
	go mod tidy
	go mod vendor