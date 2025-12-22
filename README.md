# Notes App (Go + React)
A minimal starting point for a notes app.

## Structure
- `backend/`: Go HTTP API (in-memory notes store)
- `frontend/`: React (Vite + TypeScript)

## Prereqs
- Node.js + npm
- Go (recommended: Go 1.22+)

## Run (dev)
In separate terminals:

Backend:
```sh
cd backend
go run .
```

Frontend:
```sh
cd frontend
npm install
npm run dev
```

The frontend calls the backend via Vite proxy (`/api` -> `http://localhost:8080`).

## API
- `GET /healthz`
- `GET /api/notes`
- `POST /api/notes` body: `{ "title": string, "content": string }`
- `GET /api/notes/:id`
- `PUT /api/notes/:id` body: `{ "title"?: string, "content"?: string }`
- `DELETE /api/notes/:id`

## Test
```sh
cd backend
go test ./domain/items
```