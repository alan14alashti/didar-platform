# Didar Platform

**Didar** ("دیدار" — meeting/encounter) is an Airbnb-Experiences-style marketplace for
Iran: guests discover and book local experiences (cooking classes, handicraft
workshops, nature tours, etc.) hosted by verified local hosts. Built to be
white-label ready so the same codebase can power multiple branded storefronts.

## Stack

- **Frontend**: Vue 3, TypeScript, Vite, Tailwind CSS, Pinia, Vue Router, Vue I18n (fa/en, RTL/LTR), Leaflet/OpenStreetMap
- **Backend**: Go microservices — **auth service** (OTP, users, JWT) + **main API** (experiences, bookings, hosts)
- **Auth**: Email + OTP (SMTP; mock mode logs codes locally). Phone/SMS (Kavenegar) kept ready to restore via `OTP_CHANNEL=sms`.
- **Infra**: Docker Compose for local dev (Postgres/PostGIS + auth + API + frontend)

## Repository layout

```
backend/           Main marketplace API (experiences, bookings, host profiles)
services/auth/     Auth microservice (OTP login, users, JWT, sessions)
pkg/authkit/       Shared JWT validation library used by both services
pkg/platform/      Shared HTTP helpers
frontend/          Vue 3 SPA (guest + host views)
docker-compose.yml
.env.example
```

## Local development

```bash
cp .env.example .env   # adjust secrets/ports as needed
docker compose up --build
```

- Frontend: http://localhost:5173
- Main API: http://localhost:8080/api/v1
- Auth service: http://localhost:8081/api/v1
- Postgres/PostGIS: localhost:5432 (`didar` + `didar_auth` databases)

On first boot each service runs its migrations. The main API seeds sample brand,
categories, hosts, and experiences automatically. Demo host users are seeded in
the auth service with matching UUIDs.

### Running without Docker

```bash
# Create auth database (once)
psql postgres://didar:didar@localhost:5432/didar -c "CREATE DATABASE didar_auth;"

# auth service
cd services/auth && go run ./cmd/server

# main API
cd backend && go run ./cmd/server

# frontend
cd frontend && npm install && npm run dev
```

### Kavenegar OTP setup

1. Create an account at [kavenegar.com](https://kavenegar.com)
2. Create an OTP template (e.g. `didar-otp`) in the Kavenegar panel
3. Set in `.env`:
   - `KAVENEGAR_API_KEY=your-api-key`
   - `KAVENEGAR_OTP_TEMPLATE=didar-otp`
   - `OTP_MOCK_MODE=false`

OTP codes are sent only by the auth service. With an empty API key, mock mode
logs codes to the auth service stdout instead.

## White-labeling

Branding (name, colors, logo, locales, domain) lives in the `brands` table and
is resolved at runtime by the `/api/v1/brand` endpoint; the frontend applies it
as CSS variables + document metadata on boot (`src/stores/brand.ts`). Adding a
new brand/storefront is a data change, not a code change.

## Deployment

- **Frontend** is deployed automatically to **GitHub Pages** on every push to
  `main` that touches `frontend/**` (see `.github/workflows/deploy-pages.yml`).
- **Backend services** need separate deployment (or an API gateway routing
  `/auth/*` to the auth service and everything else to the main API). Set
  `VITE_API_BASE_URL` and `VITE_AUTH_BASE_URL` to their public URLs.
