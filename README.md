# Бот для телеграм канала Мыльная мама

[![Go](https://github.com/soapmama/telegram-bot/actions/workflows/go.yml/badge.svg)](https://github.com/soapmama/telegram-bot/actions/workflows/go.yml)

## Требуется

- [Go 1.24](https://go.dev/dl/)
- [Telegram API token](https://core.telegram.org/bots/api#authorizing-your-bot)

## Подготовка

- Создать файл `.env` и добавить в него `TOKEN`

## Запуск локально

```bash
go run .
```

## Проверка обновлений пакетов

```bash
# Install the tool
go install github.com/psampaz/go-mod-outdated@latest

# Run it with go list
go list -u -m -json all | go-mod-outdated -update -direct
```

## Деплой

TODO