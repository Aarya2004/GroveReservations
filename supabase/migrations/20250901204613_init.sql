-- users (minimal)
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  name  TEXT NOT NULL,
  role  TEXT NOT NULL DEFAULT 'MEMBER',
  villa_number INT NOT NULL,
  phone_number TEXT NOT NULL
);

-- resources
CREATE TABLE resources (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  location TEXT,
  slot_minutes INT NOT NULL DEFAULT 60,
  buffer_minutes INT NOT NULL DEFAULT 0,
  max_advance_days INT NOT NULL DEFAULT 14,
  open_hours JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- reservations
CREATE TABLE reservations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  starts_at   TIMESTAMPTZ NOT NULL,
  ends_at     TIMESTAMPTZ NOT NULL,
  status      TEXT NOT NULL CHECK (status IN ('HELD','CONFIRMED','CANCELLED','NOSHOW','COMPLETED')),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT starts_before_ends CHECK (starts_at < ends_at)
);
CREATE EXTENSION IF NOT EXISTS btree_gist;
ALTER TABLE reservations
  ADD COLUMN time_range tstzrange GENERATED ALWAYS AS (tstzrange(starts_at, ends_at, '[)')) STORED;
CREATE INDEX reservations_resource_time_gist ON reservations USING GIST (resource_id, time_range);
ALTER TABLE reservations ADD CONSTRAINT reservations_no_overlap
  EXCLUDE USING GIST (resource_id WITH =, time_range WITH &&)
  WHERE (status IN ('HELD','CONFIRMED'));
ALTER TABLE users     ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT now(), ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();
ALTER TABLE resources ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT now(), ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();
