[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

Микросервис для мониторинга времени.

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- Gin (веб фреймворк)
- golang-migrate/migrate (для миграций БД)
- pgx (драйвер для работы с PostgreSQL)

Сервис был написан с Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его.
Также был реализован Graceful Shutdown для корректного завершения работы сервиса.

# Usage

Запустить сервис можно с помощью команды `make compose-up`

Документацию после запуска сервиса можно посмотреть по адресу `http://localhost:8080/swagger/index.html`
с портом 8080 по умолчанию.
