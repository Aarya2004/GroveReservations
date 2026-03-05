# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Community reservation system ("Grove") for tennis courts and shared amenities. Monorepo with a Go API backend and a Next.js web frontend.

## Repository Structure

- `software/groveapi/` — Go API server (Fiber + GORM + Supabase)
- `software/clients/web/` — Next.js 15 web client (React 19, Tailwind v4, shadcn/ui pattern)
- `docs/` — Requirements and API specification

## Build & Run

### Backend (groveapi)

```bash
cd software/groveapi
go build -o groveapi ./cmd/api     # build
go run ./cmd/api                    # run (reads .env via godotenv)
```

Required env vars in `software/groveapi/.env`:
- `DATABASE_URL` — Postgres connection string
- `SUPABASE_URL` — Supabase project URL
- `SUPABASE_SERVICE_KEY` — Supabase service role key
- `HTTP_ADDR` — Listen address (default `:8080`)

### Frontend (web client)

```bash
cd software/clients/web
npm install
npm run dev          # dev server with Turbopack
npm run build        # production build
npm run lint         # ESLint
```

### Database Migrations

Managed via Supabase CLI. Migration files live in `supabase/migrations/`. A reference copy is at `software/groveapi/internal/migrations/init.sql`.

## Architecture

### Backend Layering

Requests flow: **Routes → Handlers → Logic → Store (GORM)**

- `internal/http/routes/` — Route registration per feature (auth, resources, reservations). Each file calls `Register*Routes(router, db, sb)`.
- `internal/http/handlers/` — HTTP handler structs with methods per endpoint. Parse request, call logic, return JSON.
- `internal/logic/` — Business logic (validation, transactions). Input/DTO types defined here.
- `internal/store/` — GORM models (`User`, `Resource`, `Reservation`, `AuditLog`) and DB connection.
- `internal/config/` — Env-based configuration.
- `internal/sb/` — Supabase client initialization.

All API routes are versioned under `/api/v1`. Health check at `GET /health`.

### Auth

Authentication uses Supabase Auth (GoTrue). The `users` table has a foreign key to `auth.users(id)`. Supabase client is passed to auth handlers for login/register/logout. Admin middleware is stubbed but not yet wired.

### Conflict Prevention

Reservations use a PostgreSQL exclusion constraint (`EXCLUDE USING GIST`) on `(resource_id, time_range)` to prevent overlapping bookings at the DB level. The Go code detects PG error code `23P01` to return conflict errors.

### Key Models

- **User**: UUID PK, linked to Supabase Auth, roles: `MEMBER`/`ADMIN`
- **Resource**: configurable `slot_minutes`, `buffer_minutes`, `max_advance_days`, `open_hours` (JSONB)
- **Reservation**: statuses: `HELD` → `CONFIRMED` → `CANCELLED`/`NOSHOW`/`COMPLETED`
- **AuditLog**: tracks user actions

### Frontend

Next.js 15 app using App Router, React Query for data fetching, Supabase JS client for auth, react-big-calendar for schedule views, react-hook-form + zod for forms.
