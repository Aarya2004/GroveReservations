-- Allow all authenticated users to SELECT reservations (needed for availability).
-- The existing reservations_select_own policy only lets users see their own,
-- but availability computation needs to see all HELD/CONFIRMED reservations
-- for a given resource. Drop the restrictive policy and replace with a broader one.

DROP POLICY IF EXISTS reservations_select_own ON public.reservations;

CREATE POLICY reservations_select_all ON public.reservations
  FOR SELECT TO authenticated
  USING (true);
