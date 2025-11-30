# Upsales Golang Backend

## Getting Started

### 1. Install Go (Golang)

- Download Go from [https://go.dev/dl/](https://go.dev/dl/)
- Verify installation:
  ```bash
  go version
  ```

---

### 2. Set up PostgreSQL with Docker (Version 16)

Run the following command to start PostgreSQL with the required credentials:

```bash
docker run --name crackers-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Avaali1$ -e POSTGRES_DB=vinayakacrackers -p 5432:5432 -d postgres:16
```

- **DB URL to use in the application**:
  ```
  DB_URL="host=localhost user=postgres password=Avaali1$ dbname=vinayakacrackers port=5432 sslmode=disable"
  ```

---

### 3. If port 5432 is already in use

If your local machine already has PostgreSQL running on port **5432**, you can run the container on a different port (e.g., 5433):

```bash
docker run --name crackers -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Avaali1$ -e POSTGRES_DB=vinayakacrackers -p 5433:5432 -d postgres:16
```

In this case, update your `DB_URL` variable in the `.env` file accordingly:

```
DB_URL="host=localhost user=postgres password=Avaali1$ dbname=vinayakacrackers port=5433 sslmode=disable"
```

---

### 4. Running the Application

- Install dependencies:
  ```bash
  go mod tidy
  ```
- Run the application:
  ```bash
  go run main.go
  ```
- Run using CompileDaemon:
  ```
  CompileDaemon -command="./crackers"
  ```

---

✅ This way, you’ll have a running Go backend with Postgres 16 in Docker.
