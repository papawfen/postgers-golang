### Для запсука сервера введите команду
```bash
go run main/main.go
```

### Чтобы протестировать можно запустить тесты
```bash
go test ./...
```
### Или выполнить этот запрос
```bash
curl -X POST http://localhost:8080/commands \
     -H "Content-Type: application/json" \
     -d '{"command": "ls -a"}'

```

В результате придет json в виде:

```json
{
    "id": "some-uuid",
    "command": "ls -a",
    "output": "",
    "status": "running",
    "created_at": "2024-05-14T12:00:00Z"
}

```

Для проверки статуса команды

```bash
curl http://localhost:8080/commands/uuid
```

Вывод должен быть 

```json
{
    "id": "some-uuid",
    "command": "ls -a",
    "output": ".\n..\n.git\nREADME.md\n",
    "status": "success",
    "created_at": "2024-05-14T12:00:00Z"
}
```