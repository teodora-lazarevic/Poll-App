# Poll App

## Prerequisites

- Docker installed
- Go installed



## Start the Database

Navigate to the project directory and start the database containers:

```bash
cd ~/Desktop/Poll-App
sudo docker compose up -d
```



## Start the Go Application

Run the API server:

```bash
go run cmd/api/main.go
```



## Test the Application

Verify that the application is running:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```



## Stop the Database

When finished, stop and remove the database containers:

```bash
cd ~/Desktop/Poll-App
sudo docker compose down
```

### Everytime you change schema, you need to regenerate the Go code
```bash
go generate ./ent
```