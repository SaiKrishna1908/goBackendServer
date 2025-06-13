+-------------------+
|    main.go        |
+-------------------+
         |
         v
+-------------------+
| app.NewApplication()  <--- uses ---+
+-------------------+                |
         |                           |
         v                           |
+-------------------+                |
| routes.SetUpRoutes(*app) <---------+
+-------------------+
         |
         v
+-------------------+
|   chi Router      |
+-------------------+
         |
         v
+-------------------------------+
|   Handler Functions           |
|  (HealthCheck, WorkoutHandler)|
+-------------------------------+
         |
         v
+-------------------+
|   Data Store      |  (planned: workout_store.go)
+-------------------+

### Run tests

```
go test ./...
```

### Run the application

```
docker-compose up -d
```

```
go run main.go
```

### Run goose migrations

```
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```