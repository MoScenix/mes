# MES Monolith

This repository now runs as a Spring Boot monolith plus a Vue frontend.

- Backend: `app/backend/`
- Frontend: `app/frontend/`
- Runtime secrets: root `.env`
- Local document data: `static/document/`

The former Go microservices, BFF, IDL, RPC generation, Consul, Jaeger, Prometheus, and Grafana
runtime paths are retired. Public HTTP routes are served directly by the Spring Boot application
under package `com.team10.mes`.

## Local Development

Start required middleware and the monolith with Docker:

```bash
docker-compose -p mes up --build
```

Run the backend directly:

```bash
./start.sh
```

Run the frontend dev server:

```bash
cd app/frontend
pnpm dev
```

## Checks

```bash
make test
cd app/frontend && pnpm type-check
```
