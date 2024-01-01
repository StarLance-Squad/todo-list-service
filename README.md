[![codecov](https://codecov.io/gh/StarLance-Squad/todo-list-service/graph/badge.svg?token=J0WE99LHAE)](https://codecov.io/gh/StarLance-Squad/todo-list-service)

# todo-list-service

Todo List Application Service. FullStack application

### Docs: https://echo.labstack.com/docs

> Echo - High performance, extensible, minimalist Go web framework

---

### Up local environment

Make sure if you have a `Golang` installed. If not install it - https://go.dev/dl/

```shell
go version
```

1. Update .env file with your local environment variables
2. Run the following command to start the server

```shell
go mod tidy
```

```shell
go list -m all
```

```shell
go run cmd/main.go
```

---

### API requests

Create a User

```shell
curl -X POST http://localhost:8000/users \
-H 'Content-Type: application/json' \
-d '{"Username": "newuser", "Email": "newuser@example.com", "Password": "mypassword"}'
```

User Login

```shell
curl -X POST http://localhost:8000/login \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     -d 'username=newuser&password=mypassword'
```

Create a Todo

```shell
curl -X POST http://localhost:8000/todos \
     -H 'Content-Type: application/json' \
     -H 'Authorization: Bearer [Your_JWT_Token]' \
     -d '{"Title": "Sample Todo", "Description": "This is a sample todo item", "Completed": false}'
```

Get all Todos

```shell
curl -X POST http://localhost:8000/todos/get-all \
     -H 'Content-Type: application/json' \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJuYW1lIjoiIiwiYWRtaW4iOmZhbHNlLCJzdWIiOiJuZXd1c2VyIn0.P8QU9JJt5IpB2DCJ5fyp-lJj3TyYQgS2tD665ztLReA"
```
