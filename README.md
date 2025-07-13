# MediaSoft Internship - Practice 2025

Проект, написанный для практики в MediaSoft.

> [!NOTE]
> Требования к выполнению описаны в файле [требований](docs/requirements.md).
> Функции системы описаны в файле [функционала](docs/features.md).


## Запуск приложения
Приложение запускается в контейнере Docker.

Перед запуском необходимо заполнить файл с переменными `docker.env` по образу `docker.env.template`. Все комментарии, начинающиеся с `//` следует удалить. Итоговый файл должен `docker.env` должен выглядеть следующим образом:
```env
POSTGRES_PASSWORD=mediasoft
POSTGRES_USER=mediasoft
POSTGRES_DB=mediasoft
DBNAME=mediasoft
DBUSER=mediasoft
...
```

Для запуска перейдите в корневую директорию проекта и выполните команду
```bash
make docker-up
```
