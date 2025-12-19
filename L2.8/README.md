# Получение точного времени (NTP)

Запуск программы (NTP сервер по умолчанию - 0.beevik-ntp.pool.ntp.org):
``` bash
go run main.go
```
С указанием NTP сервера:
``` bash
go run main.go -server "time.google.com"
```
