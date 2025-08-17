# go-demoservice
PostgresSQL-Kafka-Cache microservice project | Wildberries Tech School

# Используемые технологии
- Go 1.24.5
- PostgreSQL (pgx + database/sql)
- Docker Compose
- Echo (backend)
- Kafka (kafka-go)
- golangci-lint (линтер)

# Установка
  1. `git clone https://github.com/KenyaProgressive/go-demoservice`
  2. `cd go-demoservice`

# Запуск проекта
**Перед запуском необходимо переименовать .env.example в .env и установить переменные окружения (DB_PORT ДОЛЖЕН БЫТЬ 5433)**

` DB_USERNAME="user"
  DB_PASSWORD="user"
  DB_HOST="localhost"
  DB_PORT=5433
  DATABASE_NAME="testdb"
`

Запуск проекта рекомендовано производить с помощью make:
  1. `make run` -- запуск проекта без генератора сообщений
  2. `make run-gen` -- запуск проекта с генерацией сообщений (генерирует 450 сообщений

# Завершение работы проекта
  1. CTRL+C
  2. `docker compose down`

# Пример запроса к API
  `curl http://localhost:8080/order/262d7788-92fe-48ed-beb8-e79120f4dce3`
