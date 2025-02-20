# notification-service

Сервис для рассылки и уведомлений пользователей по почте, в телеграме и в SMS. 

# Краткое описание

С помощью этого сервиса можно хранить группы людей и их контакты, чтобы уведомлять их о чем то или делать рассылку по почте/боте в телеграме. В будущем планируется добавить рассылку SMS. 

Сервис предоставляет эндпоинты для добавления контактов пользователей, объединения их в группы и рассылки определенных сообщений конкретным группам.

Подробная документация по API описана в `notification_service/README.md`

# Стек

Сервис написан на Go. Веб-фреймворк - `gorilla/mux`.
Для хранения информации о личностях и их групп используется
БД `PostgreSQL`. Для посылки задач для воркеров используется `RabbitMQ`. Для хранения метаданных используется `etcd`.

# Инструкция по локальному запуску

### 1. Запустить необходимые инфраструктурные компоненты(PostgreSQL, etcd, RabbitMQ) с помощью docker-compose:
* Создать `.env` файл, где будут заданы переменные `POSTGRES_USER`, `POSTGRES_DB`, `POSTGRES_PASSWORD`, `POSTGRES_HOST=localhost`
* Запустить контейнеры командой `docker-compose up -d`

### 2. Сделать миграции БД с помощью инструмента [golang-migrate](https://github.com/golang-migrate/migrate):

Установка `golang-migrate`:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Миграции:
```bash
migrate -path ./migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:5432/$POSTGRES_DB?sslmode=disable" up
```

### 3. Создать конфиги:
Нужно создать `notification_service/config.yml` с конфигурацией базы данных и очереди сообщений. Конфигурацию базы данных нужно
указать такой же, как и в `.env` файле.
Пример конфигурации:

```yaml
db:
  host: localhost
  port: 5432
  user: <postgres user>
  password: <postgres password>
  db_name: <postgres db>
rabbit:
  host: localhost
  port: 5672
  user: <rabbit user>
  password: <rabbit password>
  email_queue: email
  tg_queue: tg
```

Также нужно создать `notification_service/worker-config.yml` с конфигурацией воркеров. Пример конфигурации:

```yaml
email:
  smtp_host: <your smtp server address>
  smtp_port: 587
  sender_address: <your address>
  sender_password: <secret key for smtp server auth>
tg:
  token: <telegram bot token>
  etcd:
    host: localhost
    port: 2379
```
> Если у вас есть gmail почта, можно использовать их smtp сервер. Тогда адрес smtp сервера будет `smtp.gmail.com`, адрес отправителя - ваша почта, а секретный ключ - сгенерированный App password. (Подробнее - [тут](https://support.google.com/mail/answer/185833?hl=en))

### 4. Запустить API сервис и воркеров:

Надо запустить API сервис и воркеров (Запускать нужно из директории notification_service)

```
go run cmd/api/main.go
```

```
go run cmd/worker/main.go
```

Теперь всё должно заработать!
