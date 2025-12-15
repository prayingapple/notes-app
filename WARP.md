# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Repository overview
This repo is a minimal “notes” app split into two independently-run parts:
- `backend/`: Go `net/http` JSON API (in-memory notes store).
- `frontend/`: React + TypeScript app built with Vite.

In dev, the frontend talks to the backend via Vite’s proxy (`/api` and `/healthz` are forwarded to `http://localhost:8080`).

## Common commands
### Backend (Go)
Run the API server (listens on `:8080`):
```sh
cd backend
go run .
```

Run tests (none currently exist, but this is the standard entrypoint):
```sh
cd backend
go test ./...
```

### Frontend (Vite + React)
Install dependencies:
```sh
cd frontend
npm install
```

Run dev server (typically on `http://localhost:5173`):
```sh
cd frontend
npm run dev
```

Lint:
```sh
cd frontend
npm run lint
```

Build (TypeScript build + Vite build; output in `frontend/dist/`):
```sh
cd frontend
npm run build
```

Preview the production build:
```sh
cd frontend
npm run preview
```

### Makefile shortcuts (repo root)
The repo root `Makefile` provides a few wrappers:
```sh
make frontend-install
make frontend-dev
make frontend-build
make backend-dev
```

## High-level architecture
### Backend request routing
All backend logic currently lives in `backend/main.go`.
- Router: `http.NewServeMux()` with handlers for:
  - `GET /healthz`
  - `GET /api/notes` (list)
  - `POST /api/notes` (create)
  - `GET /api/notes/:id` (fetch)
  - `PUT /api/notes/:id` (update)
  - `DELETE /api/notes/:id` (delete)
- Storage: `noteStore` is an in-memory `map[string]Note` guarded by an `RWMutex`.
  - Restarting the backend clears all notes.
- JSON helpers: `writeJSON` sets `Content-Type: application/json` and encodes the response.
- CORS: `corsMiddleware` allows requests from the Vite dev origins (`http://localhost:5173` / `http://127.0.0.1:5173`).

Update semantics:
- `PUT /api/notes/:id` uses `updateNoteRequest` pointer fields to detect omitted fields.
- `noteStore.update(id, title, content)` treats empty strings as “don’t update”, so the API cannot currently set `title`/`content` to the empty string via `PUT`.

### Frontend data flow
The frontend is intentionally simple and mostly lives in `frontend/src/App.tsx`.
- API client: `api<T>()` wraps `fetch`, forces `Content-Type: application/json`, and throws on non-2xx.
- Data model: `type Note` mirrors backend JSON (`id`, `title`, `content`, `createdAt`, `updatedAt`). If you change backend JSON fields, update this type (and any UI assumptions) accordingly.
- State management:
  - `refresh()` loads `GET /api/notes` and updates `notes`.
  - `createNote()` posts to `POST /api/notes` and prepends the created note into local state.
- Rendering:
  - Notes are sorted by `updatedAt` descending.
  - There is currently no UI for update/delete; only create + list.

### Frontend ↔ backend integration point
The integration contract is:
- Routes under `/api/*` from the browser.
- Vite proxies these requests to the Go server via `frontend/vite.config.ts`.

If you need to change ports or the backend base URL in dev, update `frontend/vite.config.ts` (proxy config) and the backend `addr` in `backend/main.go` together.