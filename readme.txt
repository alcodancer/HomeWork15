# сгенерировать sqlc-код
sqlc generate

# запустить сервер (миграции применятся автоматически)
go run ./cmd/server

# пример транзакционного запроса
curl -X POST localhost:8080/users/with-order \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"a@x.com","amount":100,"note":"first"}'

# тесты и покрытие
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
