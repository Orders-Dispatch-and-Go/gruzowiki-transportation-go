# Локальное поднятие

1. Запустите контейнеры:

```
$ make local
```

2. Подключитель к базе данных по параметрам из .docker/local/config.json.

3. Выполните в базе данных миграции из ./internal/db/pg/migrations.

4. Сервис доступен по http://localhost:8000/

# Запросы

/register

email
name
password
bio
role ("writer" или "reader")
avatar_url

получаете 201 статус если все ок, если данные фиговые 400, если сервер упал 500

/login

email
password

на ок 200 статус + token поле а ответе с токеном, если данные кал 400 статус, если сервер одурел 500

получить профиль надо токен вставить в заголовок Authorization с Bearer схемой.

Текущий профиль /current-user

вернётся

id
email
name
role
bio
avatar_url
created_at
