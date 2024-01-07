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

### API requests for Unix Shell with cURL

#### Register a new User

```shell
curl -X POST http://localhost:8000/authentication/register \
     -H 'Content-Type: application/json' \
     -d '{
           "username": "zelenchuk",
           "email": "zelenchuk@gmail.com",
           "password": "root_admin"
         }'
```

---

#### Login User

```shell
curl -X POST http://localhost:8000/authentication/login \
     -H "Content-Type: application/json" \
     -d '{
           "username": "newuser",
           "password": "mypassword"
         }'
```

---

#### Who I am

```shell
curl -X GET -H "Authorization: Bearer [Your_JWT_Token]" http://localhost:8000/authentication/whoiam
```

JSON response:

```json
{
  "admin": false,
  "exp": 1704704665,
  "userId": 52,
  "username": "zelenchuk"
}
```

---

#### Create a Todo

```shell
curl -X POST http://localhost:8000/todos \
     -H 'Content-Type: application/json' \
     -H 'Authorization: Bearer [Your_JWT_Token]' \
     -d '{"Title": "Sample Todo", "Description": "This is a sample todo item", "Completed": false}'
```

---

Get all Todos

```shell
curl -X GET http://localhost:8000/todos \
     -H 'Content-Type: application/json' \
     -H "Authorization: Bearer [Your_JWT_Token]"
```

### Execution Methods for Windows PowerShell

---

#### Creating a User

```shell
$body = @{
    Username = "example"
    Email = "example@example.com"
    Password = "password"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8000/users" -Method Post -ContentType "application/json" -Body $body
```

#### User Login

```shell
$body = @{
    username = "example"
    password = "password"
}

Invoke-WebRequest -Uri "http://localhost:8000/login" -Method Post -ContentType "application/x-www-form-urlencoded" -Body $body
```

#### Creating a Todo

```shell
$headers = @{
    Authorization = "Bearer <token>"
    ContentType   = "application/json"
}

$body = @{
    Title       = "Sample Todo"
    Description = "This is a sample todo item"
    Completed   = $false
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8000/todos" -Method Post -Headers $headers -Body $body
```

#### Getting all Todos

```shell
$headers = @{
    Authorization = "Bearer <token>"
}

Invoke-WebRequest -Uri "http://localhost:8000/todos" -Headers $headers
```

#### Deleting a Todo

```shell
$headers = @{
    Authorization = "Bearer <token>"
}

Invoke-WebRequest -Uri "http://localhost:8000/todos/<todo_id>/<user_id>" -Method Delete -Headers $headers

```

#### Update Todo

```shell
$headers = @{
    Authorization = "Bearer <token>"
    "Content-Type" = "application/json"
}

$body = @{
    Title = "Updated Title"
    Description = "Updated Description"
    Completed = $true
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8000/todos/<todo_id>/<user_id>" -Method Put -Headers $headers -Body $body
```