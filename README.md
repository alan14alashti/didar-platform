# Didar Platform

**Didar** ("دیدار" — meeting/encounter) is an Airbnb-Experiences-style marketplace for
Iran: guests discover and book local experiences (cooking classes, handicraft
workshops, nature tours, etc.) hosted by verified local hosts. Built to be
white-label ready so the same codebase can power multiple branded storefronts.

## Stack

- **Frontend**: Vue 3, TypeScript, Vite, Tailwind CSS, Pinia, Vue Router, Vue I18n (fa/en, RTL/LTR), Leaflet/OpenStreetMap
- **Backend**: Go, chi router, pgx, PostgreSQL + PostGIS, JWT auth
- **Auth**: Phone number + OTP (via [Kavenegar](https://kavenegar.com), with a mock sender for local dev)
- **Infra**: Docker Compose for local dev (Postgres/PostGIS + backend + frontend)

## Repository layout

```
backend/    Go API server, DB migrations (embedded SQL), business logic
frontend/   Vue 3 SPA (guest + host views)
docker-compose.yml
.env.example
```

## Local development

```bash
cp .env.example .env   # adjust secrets/ports as needed
docker compose up --build
```

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api/v1
- Postgres/PostGIS: localhost:5432

On first boot the backend runs its migrations and seeds sample brand,
categories, hosts, and experiences (with photos) automatically.

### Running without Docker

```bash
# backend
cd backend && go run ./cmd/server

# frontend
cd frontend && npm install && npm run dev
```

## White-labeling

Branding (name, colors, logo, locales, domain) lives in the `brands` table and
is resolved at runtime by the `/api/v1/brand` endpoint; the frontend applies it
as CSS variables + document metadata on boot (`src/stores/brand.ts`). Adding a
new brand/storefront is a data change, not a code change.

## Deployment

- **Frontend** is deployed automatically to **GitHub Pages** on every push to
  `main` that touches `frontend/**` (see `.github/workflows/deploy-pages.yml`).
- **Backend + Postgres/PostGIS** are not hosted by GitHub Pages (it only
  serves static files). To get a fully working live demo, deploy the `backend`
  (with a Postgres/PostGIS database) to any Go-friendly host (Fly.io, Railway,
  Render, a VPS, etc.), then set the `VITE_API_BASE_URL` repository variable
  (Settings → Secrets and variables → Actions → Variables) to that backend's
  public URL and re-run the workflow. Until then, the deployed static site
  will render but API calls (listings, login, bookings) will fail.
