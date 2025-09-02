# Community Reservation System – API Specification (v1)

## Auth & Users

* **POST `/auth/register`** – Create user (email/password or OTP). *(Optional if admin-managed onboarding)*
* **POST `/auth/login`** – Login.
* **POST `/auth/logout`** – Logout/invalidate session.
* **GET `/users/me`** – Current user profile.
* **GET `/users`** *(admin)* – List all users.
* **PATCH `/users/:id`** *(admin)* – Update role/status.
* **DELETE `/users/:id`** *(admin)* – Deactivate user.

---

## Resources

* **GET `/resources`** – List all resources.
* **POST `/resources`** *(admin)* – Create resource `{ name, type, location, rules }`.
* **GET `/resources/:id`** – Get resource details and rules.
* **PATCH `/resources/:id`** *(admin)* – Update resource/rules.
* **DELETE `/resources/:id`** *(admin)* – Remove resource.

---

## Availability

* **GET `/availability?resource_id=&from=&to=`**
  Returns array of time slots with status (`open`, `booked`, `blackout`).

---

## Reservations

* **GET `/reservations/my`** – List my reservations.
* **GET `/reservations?resource_id=&from=&to=`** *(admin)* – List reservations for a resource.
* **POST `/reservations`** – Create reservation `{ resource_id, starts_at, ends_at }`.

  * Response: reservation object OR `409 Conflict`.
* **PATCH `/reservations/:id`** – Modify reservation (time change, cancel).
* **DELETE `/reservations/:id`** – Cancel reservation.

---

## Audit & Reporting

* **GET `/audit`** *(admin)* – List audit logs (filter by user, action, resource).
* **GET `/reports/utilization?resource_id=&from=&to=`** – Usage stats (% occupancy). (Optional)

---

## Error Codes

* **400** – Invalid request (e.g., outside open hours).
* **401** – Unauthorized (not logged in).
* **403** – Forbidden (wrong role).
* **404** – Not found.
* **409** – Conflict (slot already booked).

---