format:
	@gofmt -w . && goimports -w .
dev:
	@go run cmd/main.go
test:
	@godotenv -f ./.env go test -v ./... | grep -v '\[no test files\]'
gen:
	@flatc --go --gen-object-api -o ./gen ./proto/fb/schema.fbs