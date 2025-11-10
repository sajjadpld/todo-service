# Simple Todo List Code Challenge

---

✅ Executed per assessment guidelines:

- Used GO version: 1.23
- Used PostgreSQL Database
- Designed Hexagonal Architecture
- Unit Tests Implemented
- Dockerized via `Dockerfile` and `Docker Compose`
- Initialized by `Makefile` and `make` command
- Well Documented by OpenAPI(Swagger)

---

## Overview

- This service creates, shows, and lists the todo items.
- The Unit Tests are implemented just for the main functionality of create todo item in the `repository` and `usecase` layer.
- All related service clients such as Database, Logger, Locale, Registry, etc were mocked by `mockgen` to be used in the Uint Tests.
- The database migrations will be handled automatically by running the service.
- To import APIs in the `POSTMAN`, download the swagger `json` file and import that.(http://localhost:8080/public/swagger/doc.json)

---

## Quick Start

To start service by docker-compose:
```shell
make run
```

To Stop:
```shell
make down
```

---

the documented APIs are accessible at the link below:

```http request
http://localhost:8080/public/swagger/index.html
```
Credentials
```text
username: admin
password: admin
```

---

## Local Setup


To run the services locally, follow these steps:

1. **Disable the Docker app configuration:**
    - Comment out the `- docker/app.yml` line in `docker-compose.yml`.

2. **Start Docker containers:**
    - Run the following command to bring up the Docker environment(run database):
      ```bash
      make run
      ```
      
3. **Navigate to application directory:**
   - run:
     ```bash
     cd ./api
     ```

4. **Prepare environment files:**
    - Generate `.env` files for both services:
      ```bash
      make env
      ```
    - Set the `APP_DEBUG="true"` and `DB_DEBUG=true` in the `.env` file

5. **Install Go dependencies:**
    - In both services, run:
      ```bash
      go mod tidy
      ```

6. **Generate Swagger documentation (API service only):**
    - Install swagger as needed:
      ```bash
      go install github.com/swaggo/swag/cmd/swag@latest
      ```
    - Run the following in the API service:
      ```bash
      make swag
      ```

7. **Run the Unit Tests:**
   - run:
     ```bash
     make tests
     ```

8. **Run the services:**
    - Start each service manually:
      ```bash
      go run ./cmd/main.go
      ```