test-backend:
	cd backend && go test -v ./tests/

run-backend:
	cd backend && go run cmd/server/main.go
