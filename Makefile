swag-init:
	swag init -g ./cmd/main.go -o ./docs
go-mod:
	go mod download
	go mod tidy
