# Менеджер паролей GophKeeper

## Кратко

### Сервер

Роуты

- POST auth/login # Авторизация
- POST auth/registration # Регистрация
- GRPC Save # Добавить данные
- GRPC GetAll # Получить все данные

HTTP и GRPC работают по протоколу TLS
В качестве аутентификации используется стандарт JWT

### Клиент

Команды

- login: Авторизация
- registration: Регистрация
- add: Добавить данные
- list: Получить все данные
- sync: Синхронизация данных между сервером и клиентом

## О проекте

Решил в этой дипломной работе попробовать модульную реализацию чистой архитектуры.

Авторизация и регистрация работает по протоколу HTTP
Сохранение и получение данных работает по протоколу GRPC, так же предусмотрел HTTP
Синхронизация работает вызовом команды sync. Хотел сделать по websocket, но не хватило времени.

Протокол HTTP и GRPC обернуты протоколом TLS
P.S. Первый раз работаю с TLS, поэтому мог реализовать что-то не так.

## Покрытие

P.S. Сервер покрыл, клиент не успел покрыть тестами, принцип мне понятен как это делается.
Сервер: 82.9%
Клиент: 23.6%

```shell
go test ./... -coverprofile cover.tmp.out
cat cover.tmp.out | grep -v "/mocks/" > cover.out
go tool cover -func=cover.out | grep total:
go tool cover -html=cover.out
```

## migration

```shell
migrate create -ext sql -dir db/migrations -seq create_users_table

export POSTGRESQL_URL='postgres://postgres:password@localhost:5433/postgres?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

## mockery

Предоставляет возможность легко создавать макеты для интерфейсов Golang

```shell
docker run -v "$PWD":/src -w /src vektra/mockery --all --dir internal/auth --output internal/auth/mocks
docker run -v "$PWD":/src -w /src vektra/mockery --all --dir internal/storage --output internal/storage/mocks
```

## todo

- 429 обработка
- данные хранить в зашифрованном виде
- шифрование данных в бд, либо клиент сам помнит ключ, либо ключ храним на сервере
- хранение бинарных данных в s3 minio
- при больших файлах спрашивать пользователя об синхронизации
- статус связи с сверером и показ последней инхронизации
- локальные хранилище
