# Менеджер паролей GophKeeper

## О проекте

Решил в этой дипломной работе попробовать другую реализацию чистой архитектуры,
где функционал будет разделен на модули.

GRPC реализовал в storage, авторизацию решил оставить на HTTP

## Покрытие

```shell
go test -v -race ./... -coverprofile cover.out
go tool cover -html=cover.out
```

## migration

```shell
migrate create -ext sql -dir db/migrations -seq create_users_table

export POSTGRESQL_URL='postgres://postgres:password@localhost:5433/postgres?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

## todo

- 429 обработка
- данные хранить в зашифрованном виде
- шифрование данных в бд, либо клиент сам помнит ключ, либо ключ храним на сервере
- хранение бинарных данных в s3 minio
- при больших файлах спрашивать пользователя об синхронизации
- статус связи с сверером и показ последней инхронизации
- локальные хранилище
