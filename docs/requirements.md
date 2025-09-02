## Functional Requirements (v1 focus)
User Management & Access
1. Residents register/login (email + password, or OTP).
  - Roles: Member, Admin/Manager (to manage resources/blackouts).
  - Each user tied to one community (single-tenant).
2. Resource Catalog
  - Admin can create and manage resources (e.g., “Tennis Court 1”).
  - Define rules: open hours, slot length, min/max booking duration, max advance days, buffer time, quotas per user/day/week.
3. Reservations
  - Members can browse availability calendar per resource.
  - Book a single slot (create reservation).
  - Cancel or modify their own reservations (within rules).
  - Conflict-safe booking (no overlapping reservations for same resource).
  - Statuses: HELD, CONFIRMED, CANCELLED, NOSHOW, COMPLETED.
4. Notifications
  - Email (and optional SMS) confirmations + reminders.
  - Admins notified of cancellations or no-shows if needed.
5. Realtime Updates
  - Calendar view updates instantly if another user books/cancels.
  - User sees slot locked immediately on their screen.
6. Audit & Reporting
  - Track who created/cancelled a booking, with timestamps.
  - Basic utilization reports (e.g., % occupancy per resource per week).

## Non-Functional Requirements
1. Performance
  - Booking transaction should complete in < 1 second.
  - Availability view should load < 2 seconds for a week of slots.
2. Concurrency Safety
  - No double-bookings even under high contention (enforced at DB level).
3. Usability
  - Mobile-first design (residents will mostly book via phone).
  - Simple flow: 2–3 taps to book or cancel.
4. Security
  - HTTPS everywhere.
  - Role-based access enforced (Admins vs Members).
  - Audit logs immutable.
5. Reliability
  - System uptime 99% (acceptable for community use).
  - Daily automated DB backups with point-in-time recovery.
6. Scalability (local scope)
  - Initially designed for one community, ~200–500 users, ~10–20 resources.
  - Architecture extendable later to multi-community if needed.

## Deferred / Extension Features
1. Blackouts & Maintenance
  - Admins can set blackout periods for resources (maintenance, events).
  - Blackouts override user bookings.
2. Waitlist
3. Recurring bookings (RRULE expansion).
4. Payments/penalties (integrate later if community decides).
5. Check-ins (QR codes / no-show handling).
6. Advanced analytics dashboards.
7. Multi-tenant SaaS model.