# casbin-go-example

## Dependencies

This uses has the following dependencies:

- [do](https://github.com/samber/do) for dependency injection
- [gin](https://github.com/gin-gonic/gin) for HTTP
- [gorm](https://github.com/go-gorm/gorm) for interacting with the database

## Run in dev mode

```sh
docker compose up -d
go run main.go
```

```
A user can        read/write/delete my own items
An admin can      read/write all items
A super admin can read/write/delete all items
```
