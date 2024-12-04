
test: test-backend

dev: run-backend dev-frontend


test-backend:
	cd backend && go test -v ./tests/

run-backend:
	doppler -c dev -- cd backend && air

dev-frontend:
	doppler -c dev -- cd frontend && bun run dev

