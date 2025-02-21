# Документация по API

## Эндпоинты для людей

### GET /api/persons/all

Выводит список всех людей в БД.

Формат ответа:
```
[
    {
        "person_id": 1,
        "email": "a@mail.com",
        "telegram_id": null,
        "phone_number": "+71234567890"
    },
    ...
]
```

### GET /api/persons/{personId}

Выводит информацию о человеке с конкретным id.

Формат ответа:
```
{
    "person_id": int,
    "email": string | null,
    "telegram_id": string | null,
    "phone_number": string | null
}
```

### POST /api/persons/create

Создает нового человека в БД.

Формат запроса:
```

{
    "email": string | null,
    "telegram_id": string | null,
    "phone_number": string | null
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

### DELETE /api/persons/delete

Удаляет людей с заданными айдишниками.

Формат запроса:
```

{
    "person_ids": []int
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

## Эндпоинты для групп

### GET /api/groups/all

Выводит список всех групп в БД (только названия).

Формат ответа:
```
[]string
```

### GET /api/persons/{groupName}

Выводит информацию о группе с конкретным именем.

Формат ответа:
```
{
    "group_id": int,
    "group_name": string,
    "participant_ids": []int
}
```

### POST /api/groups/create

Создает новую группу.

Формат запроса:
```

{
    "group_name": string,
    "participant_ids": []int
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

### POST /api/groups/add_participants

Добавляет людей в группу.

Формат запроса:
```

{
    "group_name": string,
    "participant_ids": []int
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

### DELETE /api/persons/delete/{groupName}

Удаляет группы с заданным именем.

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

## Эндпоинты для отправки сообщений

### POST /api/send/email/addresses

Отправялет email на заданные адреса

Формат запроса:
```
{
    "email": {
        "subject": string,
        "body": string
    },
    "recipients": []string
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

### POST /api/send/email/groups

Отправялет email заданным группам

Формат запроса:
```
{
    "email": {
        "subject": string,
        "body": string
    },
    "group_names": []string
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

### POST /api/send/tg/addresses

Отправялет сообщение в телеграме заданным пользователям

Формат запроса:
```
{
    "content": string,
    "recipients": []string
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```

### POST /api/send/tg/groups

Отправялет сообщение в телеграме заданным группам

Формат запроса:
```
{
    "content": string,
    "group_names": []string
}
```

Формат ответа:
```
{
    "status": "ok"/"error",
    "message": string
}
```
