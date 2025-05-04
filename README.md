# Labs-HR-Go

## Installation

### Prerequisites

- Go 1.24 or higher
- [GVM](https://github.com/moovweb/gvm) (optional, for managing Go versions)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/WangWilly/labs-gin.git
   cd labs-gin
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the development server:
   ```bash
   ./scripts/dev.sh
   ```

### Docker Setup

1. Build the Docker image using the provided build script:
   ```bash
   ./scripts/build.sh
   ```
   This will create a Docker image named `labs-hr-go:latest`.

2. Run the service using Docker Compose:
   ```bash
   cd deployments
   docker compose up -d
   ```
   This will start the service in detached mode, listening on port 8080.

3. To stop the service:
   ```bash
   docker compose down
   ```

4. Monitor logs:
   ```bash
   docker compose logs -f
   ```

### Docker Environment Variables

When running with Docker Compose, you can configure the following environment variables in the `deployments/docker-compose.yml` file:

| Name | Description | Default |
|------|-------------|---------|
| PORT | The port on which the service listens | `8080` |

## Migrations

### Prerequisites

- https://github.com/pressly/goose?tab=readme-ov-file#install

Use the docker compose file to run setup the database:
```bash
./scripts/local_setup.sh
```

### Development

1. Create a new migration ([example](./database/migrations/00001_init.go))

2. Edit the migration file to define the `Up` and `Down` functions. The `Up` function should contain the SQL statements to apply the migration, while the `Down` function should contain the SQL statements to revert it.

3. Run the migration using the following command:
   ```bash
   go run database/cmd/main.go
   ```

### Reference

- https://stackoverflow.com/questions/64510093/gorm-migration-using-golang-migrate-migrate

## API Documentation (TBD)

### Employee

```bash
curl --location 'http://localhost:8080/employee' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Will",
    "age": 39,
    "address": "united states",
    "phone": "654321232",
    "email": "test@goooo.co",
    "position": "tester",
    "department": "tech",
    "salary": 4000,
    "start_date": 1746365072
}'
```

```json
{
   "employee_id": 1,
   "position_id": 1
}
```

```bash
curl --location 'http://localhost:8080/employee/1'
```

```json
{
   "employee_id": 1,
   "name": "Will",
   "age": 39,
   "phone": "654321232",
   "email": "test@goooo.co",
   "address": "united states",
   "created_at": "2025-05-04 13:26:51",
   "updated_at": "2025-05-04 13:26:51",
   "position_id": 1,
   "position": "tester",
   "department": "tech",
   "salary": 4000,
   "start_date": "2025-05-04 00:00:00"
}
```

```bash
curl --location --request PUT 'http://localhost:8080/employee/1' \
--header 'Content-Type: application/json' \
--data '{
    "address": "taiwan"
}'
```

```json
{
   "id": 1,
   "name": "Will",
   "age": 39,
   "address": "taiwan",
   "phone": "654321232",
   "email": "test@goooo.co"
}

```bash
curl --location 'http://localhost:8080/promote/1' \
--header 'Content-Type: application/json' \
--data '{
    "position": "tester2",
    "department": "tech",
    "salary": 5000,
    "start_date": 1747365072
}'
```

```json
{
   "position_id": 5,
   "start_date": "2025-05-16 11:11:12"
}
```

### Attendance

```bash
curl --location 'http://localhost:8080/attendance' \
--header 'Content-Type: application/json' \
--data '{
    "employee_id": 1
}'
```

```json
{
   "attendance_id": 1,
   "position_id": 3,
   "clock_in_time": "2025-05-04 13:41:15",
   "clock_out_time": ""
}
```

```bash
curl --location 'http://localhost:8080/attendance/1'
```

```json
{
    "attendance_id": 1,
    "position_id": 3,
    "clock_in_time": "2025-05-04 13:41:15",
    "clock_out_time": "2025-05-04 13:41:15"
}
```


## All Environment Variables

| Name | Description | Default |
|------|-------------|---------|
| PORT | The port on which the service listens | `8080` |
| HOST | The host address for the service | `0.0.0.0` |

## Development Resources

- [Go Modules Documentation](https://go.dev/wiki/Modules#quick-start)
- https://github.com/smartystreets/goconvey
- https://github.com/uber-go/mock
