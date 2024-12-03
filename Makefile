
test: test-backend

dev: run-backend dev-frontend


test-backend:
	cd backend && go test -v ./tests/

run-backend:
	cd backend && air

dev-frontend:
	cd frontend && bun run dev
