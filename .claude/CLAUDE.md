# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Community reservation system ("Grove") for tennis courts and shared amenities. Supabase-only architecture — no custom backend server. The frontend talks directly to Supabase (PostgreSQL + Auth + RLS).

## Repository Structure

- `software/clients/web/` — Next.js 15 web client (React 19, Tailwind v4, shadcn/ui)
- `supabase/` — Supabase config and migrations
- `docs/` — Requirements and API specification
- `software/groveapi/` — **DEPRECATED** Go API server (kept for reference, not used at runtime)

## Build & Run

### Frontend (web client)

```bash
cd software/clients/web
bun install
bun run dev          # dev server with Turbopack
bun run build        # production build
bun run lint         # Prettier + ESLint
```

### Database

Managed via Supabase CLI. Migration files live in `supabase/migrations/`.

```bash
supabase start       # local dev (API on 54321, DB on 54322)
supabase db reset    # apply all migrations fresh
```

## Architecture

### Data Flow

Frontend → Supabase JS Client → PostgreSQL (with RLS)

- `src/lib/api.ts` — CRUD functions using `supabase.from()` queries
- `src/lib/availability.ts` — Client-side slot computation (ported from Go)
- `src/lib/supabase.ts` — Browser and server Supabase clients
- `src/hooks/useAuth.tsx` — Auth state management
- React Query for caching and invalidation

### Auth

Supabase Auth (GoTrue) with email/password and Google OAuth. The `users` table has a FK to `auth.users(id)` with triggers that auto-sync on signup/update/delete.

### Authorization (RLS)

- **users**: SELECT/UPDATE own profile only
- **resources**: SELECT by all authenticated users; write via service_role only
- **reservations**: SELECT all (for availability); INSERT/UPDATE/DELETE own only
- **audit_logs**: service_role only

### Conflict Prevention

Reservations use a PostgreSQL exclusion constraint (`EXCLUDE USING GIST`) on `(resource_id, time_range)` to prevent overlapping bookings at the DB level. The frontend detects PG error code `23P01` to show conflict errors.

### Key Models

- **User**: UUID PK, linked to Supabase Auth, roles: `MEMBER`/`ADMIN`
- **Resource**: configurable `slot_minutes`, `buffer_minutes`, `max_advance_days`, `open_hours` (JSONB)
- **Reservation**: statuses: `CONFIRMED` → `CANCELLED`/`NOSHOW`/`COMPLETED`
- **AuditLog**: tracks user actions

### Frontend

Next.js 15 App Router, React Query for data fetching, Supabase JS for auth + DB, react-hook-form + zod for forms.
