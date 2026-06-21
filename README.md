Event Booking System - Backend (Go)

Design decisions
- Language: Go (per assignment).
- In-memory store: `store.go` implements a thread-safe in-memory store (maps + mutex). This keeps the demo self-contained and simple to run.
- Auth: token-based simple auth. `POST /signup` returns a token (UUID). Clients supply `X-User-Token` header for authenticated endpoints. Tokens are stored in-memory on the `User` record.
- Roles: two roles (`organizer`, `customer`) captured on signup. Role checks enforced via `RequireRole` middleware.
- Routing: `chi` is used for concise routing and path parameter extraction.
- Background jobs: `jobs.go` implements a simple job queue over a buffered channel processed by a worker goroutine. Jobs simulate external actions by printing to the console.

Background tasks implemented
- Booking Confirmation: triggered when a customer successfully books tickets. Enqueued as `booking_confirmation` job and prints a message simulating an email.
- Event Update Notification: triggered when an organizer updates an event. Enqueued as `event_update_notification` job and prints a message simulating notifying booked customers.

API Endpoints
- `POST /signup` {name, role} -> returns created user with `token`.
- `GET /events/` -> list events (public).
- `POST /events/` -> create event (organizer only).
- `GET /events/{id}/` -> get event details.
- `PUT /events/{id}/` -> update event (organizer only) triggers event update job.
- `POST /events/{id}/book` -> book tickets (customer only) triggers booking confirmation job.
- `GET /bookings` -> list bookings for the authenticated user (customers see own bookings; organizers see bookings for their events).

How to run
1. Ensure Go 1.20+ is installed.
2. From project root:

```bash
go mod tidy
go run .
```

Example usage (curl)
- Signup organizer:

```bash
curl -s -XPOST -d '{"name":"Alice","role":"organizer"}' -H 'Content-Type: application/json' localhost:8080/signup
```

- Signup customer:

```bash
curl -s -XPOST -d '{"name":"Bob","role":"customer"}' -H 'Content-Type: application/json' localhost:8080/signup
```

Use returned `token` value in subsequent requests as header `X-User-Token`.

Notes and limitations
- This implementation uses an in-memory store and is not persistent. It's intended for demonstration and testing only.
- The job queue is simple: single-process, in-memory, and prints simulated notifications.
