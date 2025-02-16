# Инструкция по запуску

1. Для миграции сделать:

```
migrate -path ./migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:5432/$POSTGRES_DB?sslmode=disable" up
```
