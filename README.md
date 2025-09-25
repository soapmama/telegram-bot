# Бот для телеграм канала Мыльная мама

[![Go](https://github.com/soapmama/telegram-bot/actions/workflows/go.yml/badge.svg)](https://github.com/soapmama/telegram-bot/actions/workflows/go.yml)

## Требуется

- [Go 1.25](https://go.dev/dl/)
- [Telegram API token](https://core.telegram.org/bots/api#authorizing-your-bot)

## Подготовка

- Создать файл `.env` и добавить в него `TOKEN`

## Запуск локально

```bash
go run .
```

## Обновление пакетов

Рекомендуется использовать [gomod-upgrade](https://github.com/oligot/go-mod-upgrade).

```bash
# Install the tool
go install github.com/oligot/go-mod-upgrade@latest

# Run it with go list
go-mod-upgrade
```

## Тесты

Рекомендуется использовать gotestsum.

```bash
# Install gotestsum
go install gotest.tools/gotestsum@latest

# Run tests
gotestsum
```

## Деплой

TODO
