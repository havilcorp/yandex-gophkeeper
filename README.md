# Менеджер паролей GophKeeper

![demonstration](https://github.com/havilcorp/yandex-gophkeeper/assets/58453931/6a8c4dd8-c506-4bb1-83ed-d03a6ecb52ed)

## О проекте

Решил в этой дипломной работе попробовать модульную реализацию чистой архитектуры где auth и storage имеют свои controller, usecase и repository.

- Авторизация и регистрация работает по протоколу HTTP.
- Сохранение и получение данных работает по протоколу GRPC, так же предусмотрел HTTP.
- Синхронизация работает вызовом команды sync. Хотел сделать по websocket, но не хватило времени.
- Протокол HTTP и GRPC обернуты протоколом TLS
- Данные шифруются паролем на стороне клиента

P.S. Первый раз работаю с TLS, поэтому мог реализовать что-то не так.

## Сервер

Роуты

- POST auth/login - Авторизация
- POST auth/registration - Регистрация
- GRPC Save - Добавить данные
- GRPC GetAll - Получить все данные

HTTP и GRPC работают по протоколу TLS

В качестве аутентификации используется стандарт JWT

## Клиент

Команды

- login: Авторизация
- registration: Регистрация
- add: Добавить данные
- list: Получить все данные
- sync: Синхронизация данных между сервером и клиентом

## Покрытие

- Сервер: 82.9%
- Клиент: 23.6%

P.S. Клиент не успел покрыть из-за его сложной структуры, я пытался ее упростить и у меня это получалось, но сильно не успевал.

```shell
go test ./... -coverprofile cover.tmp.out
cat cover.tmp.out | grep -v "/mocks/" > cover.out
go tool cover -func=cover.out | grep total:
go tool cover -html=cover.out
```

## Build

#### Client

```shell
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o build/macos cmd/main.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -o build/win.exe cmd/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc  CXX=x86_64-linux-musl-g++ go build -o build/linux cmd/main.go
```

#### Server

```shell
GOOS=darwin GOARCH=amd64 go build -o build/macos cmd/main.go
GOOS=windows GOARCH=amd64 go build -o build/win.exe cmd/main.go
GOOS=linux GOARCH=amd64 go build -o build/linux cmd/main.go
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
- при больших файлах спрашивать пользователя об синхронизации
- статус связи с сверером и показ последней синхронизации
