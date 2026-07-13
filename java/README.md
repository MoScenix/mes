# MES Java Workspace

This is the Java Maven aggregation workspace for the MES backend rewrite.

## Modules

```text
common/                Shared protobuf, generated gRPC stubs, and Java gRPC runtime
java/
  pom.xml              Maven aggregation root
app/
  user/                Spring Boot gRPC user microservice
  workorder/           Spring Boot gRPC workorder microservice
  inventory/           Spring Boot gRPC inventory microservice
  frontend/            Existing frontend project
```

## Build

```bash
mvn -Dmaven.repo.local=.m2/repository -f java/pom.xml test
```

## Run Services

User service:

```bash
mvn -Dmaven.repo.local=.m2/repository -f java/pom.xml -pl :mes-user-service -am spring-boot:run
```

Default local gRPC port from `start.sh`: `19091`, override with `USER_GRPC_PORT`.
Docker keeps the container port at `9091`.

Workorder service:

```bash
mvn -Dmaven.repo.local=.m2/repository -f java/pom.xml -pl :mes-workorder-service -am spring-boot:run
```

Default local gRPC port from `start.sh`: `19092`, override with `WORKORDER_GRPC_PORT`.
Docker keeps the container port at `9092`.

Inventory service:

```bash
mvn -Dmaven.repo.local=.m2/repository -f java/pom.xml -pl :mes-inventory-service -am spring-boot:run
```

Default local gRPC port from `start.sh`: `19093`, override with `INVENTORY_GRPC_PORT`.
Docker keeps the container port at `9093`.

Services use the existing MySQL database by default. Override with `MYSQL_DSN`, `MYSQL_USER`, and `MYSQL_PASSWORD`. Set `CONSUL_ADDRESS` to register Java gRPC services into Consul.

## Observability

Java services follow the same discovery shape as the Go services:

- Prometheus metrics are exposed on `/metrics` and registered into Consul as service `prometheus`.
- The Consul tag `service:<name>` identifies `user`, `workorder`, or `inventory`.
- gRPC server calls are timed with Micrometer under `rpc.server.duration`.
- Jaeger traces are exported through OTLP HTTP configured by the active Spring profile.

Default metrics ports:

```text
user        9995
workorder   9998
inventory   9999
```

For `dev` and `test`, OTLP uses `http://127.0.0.1:4318/v1/traces`. For `online`, OTLP uses `http://jaeger:4318/v1/traces`.

## Profiles

Each Java service uses Spring profiles selected by `SPRING_PROFILES_ACTIVE`.

```text
application.yml          Shared service settings
application-dev.yml      Local development defaults
application-test.yml     Local integrated test defaults
application-online.yml   Docker/online defaults
```

`start.sh` defaults Java services to `test`; `docker-compose.yml` sets them to `online`. If no environment value is provided, Spring falls back to `test`.
