# Итоговое задание go_final_project

## HTTP API базы данных задач пользователя

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Директория `store` содержит функции для работы с БД и модель Task.

Директория `handlers` содержит функции обработчики веб сервера и прочие вспомогательные функции.

### Задания с повышенной сложностью * выполнены почти все, кроме FullNextDate.

Для тестов используются параметры (FullNextDate = false обязательно):
```
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.ia-D1URCxsnnYdbyP5FvIqr0SdXmQ5kp52_tUihWNQc`
```

Такой токен работает с паролем **qwerty**

С этими параметрами и если не менять пароль в переменных окружения, то все тесты выполнялись успешно.
```
go test ./tests
```

В докер образе поменял порт на 8080:
```
ENV TODO_PORT=8080
ENV TODO_DBFILE=scheduler.db
```
Здесь все параметры можно менять как угодно.

Пароль передается через переменную TODO_PASSWORD при запуске контейнера.

### Запустить локально:
```
docker build go_final_project .
docker run --env TODO_PASSWORD=qwerty -d -p 8080:8080 go_final_project
```

### Или запустить Docker-образ из докерхаба:
```
docker pull alexeybs90/go_final_project:v1
docker run --env TODO_PASSWORD=qwerty -d -p 8080:8080 alexeybs90/go_final_project:v1
```

Приложение можно открыть в браузере по адресу http://localhost:7540/ или http://localhost:8080/ в зависимости от переменных окружения в Вашей системе или в докере
