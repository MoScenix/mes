# MES Monolith

The backend is a single Spring Boot application. The former `user`, `app`, `ai`, `document`,
`workorder`, and `inventory` services now run in one JVM under `com.team10.mes`. The old BFF is
not a Java module; its public HTTP routes and session/permission rules live in each domain's
controller and service.

The old application-management semantics have been removed; assistant sessions and chat/message
records now live in `history`. Each domain
primarily contains `controller`, `service`, and `dal`, with focused `tools`, `llm`, `state`, or
`utils` packages where the original behavior needs them. Java mapper interfaces contain no SQL;
relational queries are stored under `src/main/resources/mapper/<domain>/*.xml`.

Run locally:

```bash
mvn spring-boot:run
```

The default MySQL connection targets `127.0.0.1:3306/industrial_fault_tree_ai`. Override it with
`MYSQL_DSN`, `MYSQL_USER`, and `MYSQL_PASSWORD`. The application listens on port `8080`, and its
lightweight health endpoint is `GET /health`.

From the repository root, run the application and its required middleware:

```bash
docker-compose -p mes up --build
```

This monolith intentionally has no gRPC, Consul, Actuator, Prometheus, Jaeger, or OpenTelemetry
dependencies.

Spring AI uses the DashScope OpenAI-compatible endpoint while AI state, events, controls, and
checkpoints remain in Redis. Set `DASHSCOPE_API_KEY` and optionally `AI_MODEL`. Document files are
stored under `DOCUMENT_ROOT`; full-text and vector indexes remain in Elasticsearch and Milvus.

Runtime secrets are loaded from the repository root `.env`. The committed
`app/backend/.env.example` documents all supported variables without replacing the local secret file.
