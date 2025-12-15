.PHONY: frontend-install frontend-dev frontend-build backend-dev

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

backend-dev:
	cd backend && go run .
