-- ============================================================
-- Grove reservation system — reference schema
-- This file reflects the cumulative state of all migrations.
-- ============================================================

-- Extensions
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Users
CREATE TABLE public.users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  name  TEXT NOT NULL,
  role  TEXT NOT NULL DEFAULT 'MEMBER',
  villa_number INT NOT NULL,
  phone_number TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),

  CONSTRAINT users_id_fk_auth FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Resources
CREATE TABLE public.resources (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  location TEXT,
  slot_minutes INT NOT NULL DEFAULT 60,
  buffer_minutes INT NOT NULL DEFAULT 0,
  max_advance_days INT NOT NULL DEFAULT 14,
  open_hours JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- Reservations
CREATE TABLE public.reservations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resource_id UUID NOT NULL REFERENCES public.resources(id) ON DELETE CASCADE,
  user_id     UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  starts_at   TIMESTAMPTZ NOT NULL,
  ends_at     TIMESTAMPTZ NOT NULL,
  status      TEXT NOT NULL CHECK (status IN ('HELD','CONFIRMED','CANCELLED','NOSHOW','COMPLETED')),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ DEFAULT now(),
  time_range  tstzrange GENERATED ALWAYS AS (tstzrange(starts_at, ends_at, '[)')) STORED,
  CONSTRAINT starts_before_ends CHECK (starts_at < ends_at)
);

CREATE INDEX reservations_resource_time_gist ON public.reservations USING GIST (resource_id, time_range);

ALTER TABLE public.reservations ADD CONSTRAINT reservations_no_overlap
  EXCLUDE USING GIST (resource_id WITH =, time_range WITH &&)
  WHERE (status IN ('HELD','CONFIRMED'));

-- Audit logs
CREATE TABLE public.audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  action TEXT NOT NULL,
  resource_id UUID REFERENCES public.resources(id) ON DELETE SET NULL,
  resource_type TEXT,
  details TEXT,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_user_id ON public.audit_logs(user_id);
CREATE INDEX idx_audit_logs_timestamp ON public.audit_logs(timestamp DESC);

-- ============================================================
-- Triggers: auto-create public.users on Supabase Auth signup
-- ============================================================
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO public.users (id, email, name, role, villa_number, phone_number, active)
  VALUES (
    NEW.id,
    COALESCE(NEW.email, ''),
    COALESCE(
      NEW.raw_user_meta_data->>'name',
      NEW.raw_user_meta_data->>'full_name',
      split_part(COALESCE(NEW.email, ''), '@', 1)
    ),
    'MEMBER',
    0,
    COALESCE(NEW.phone, ''),
    true
  )
  ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email,
    name = CASE
      WHEN public.users.name = '' OR public.users.name IS NULL
      THEN EXCLUDED.name
      ELSE public.users.name
    END;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();

-- Sync email changes
CREATE OR REPLACE FUNCTION public.handle_user_updated()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE public.users SET email = NEW.email WHERE id = NEW.id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_auth_user_updated
  AFTER UPDATE OF email ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_user_updated();

-- Cascade delete
CREATE OR REPLACE FUNCTION public.handle_user_deleted()
RETURNS TRIGGER AS $$
BEGIN
  DELETE FROM public.users WHERE id = OLD.id;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_auth_user_deleted
  AFTER DELETE ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_user_deleted();

-- updated_at auto-update
CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_users
  BEFORE UPDATE ON public.users
  FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER set_updated_at_resources
  BEFORE UPDATE ON public.resources
  FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER set_updated_at_reservations
  BEFORE UPDATE ON public.reservations
  FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- ============================================================
-- Row Level Security
-- ============================================================

ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_select_own ON public.users FOR SELECT TO authenticated USING (auth.uid() = id);
CREATE POLICY users_update_own ON public.users FOR UPDATE TO authenticated USING (auth.uid() = id) WITH CHECK (auth.uid() = id);

ALTER TABLE public.resources ENABLE ROW LEVEL SECURITY;
CREATE POLICY resources_select_all ON public.resources FOR SELECT TO authenticated USING (true);

ALTER TABLE public.reservations ENABLE ROW LEVEL SECURITY;
CREATE POLICY reservations_select_own ON public.reservations FOR SELECT TO authenticated USING (auth.uid() = user_id);
CREATE POLICY reservations_insert_own ON public.reservations FOR INSERT TO authenticated WITH CHECK (auth.uid() = user_id);
CREATE POLICY reservations_update_own ON public.reservations FOR UPDATE TO authenticated USING (auth.uid() = user_id) WITH CHECK (auth.uid() = user_id);
CREATE POLICY reservations_delete_own ON public.reservations FOR DELETE TO authenticated USING (auth.uid() = user_id);

ALTER TABLE public.audit_logs ENABLE ROW LEVEL SECURITY;
