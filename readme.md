# Quick start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Util787/test-task
   cd test-task
   ```

2. **Create and configure the `.env` file**:
   ```bash
   cp .env.example .env
   ```
   Example of .env:
   ```env
   # Postgres configuration
   POSTGRES_HOST=postgres
   POSTGRES_PORT=5432
   POSTGRES_DB=app_db
   POSTGRES_USER=user
   POSTGRES_PASSWORD=password

   # HTTP server configuration
   HTTP_SERVER_HOST=0.0.0.0
   HTTP_SERVER_PORT=8000
   ```

3. **Build and start the application**:
   Use Docker Compose to build and run the application:
   ```bash
   docker-compose up --build
   ```

4. **Access the API**:
   The HTTP server will be available at `http://localhost:8000`.

## API Endpoints

- **POST** `/api/v1/sort/`: Save a number and get the sorted array.
  - Request body:
    ```json
    {
      "num": 42
    }
    ```
  - Response:
    ```json
    [1, 2, 42]
    ```

## Database Migrations

The database migrations are automatically applied when the `migrate` service runs in Docker Compose.