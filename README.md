# Template

## Используемые библиотеки

- golang 1.23
- golangci-lint (https://github.com/golangci/golangci-lint)
- oapi-codegen (https://github.com/oapi-codegen/oapi-codegen)

## Работа с приватным репозиторием

После проделанного вы можете пользоватся приватными проектами как пакетами `go get <package>`

## Go Packages

Для сборки локально образа надо запустить команду `make CI_JOB_LOGIN="Your x5 user name used by gitlab" CI_JOB_TOKEN="Your x5 user access_token(generate) used by gitlab" docker.build` .
Go умеет раобтать с приватными репозиториям для пакетов, но для этого надо сделать несколько операции:

- Указываем пременную окружения `GOPRIVATE=scm.x5.ru/x5m/go-backend/packages/*`
- Дальше в файле `~/.netrc` указываем параметры без ковычек `machine scm.x5.ru login <username> password <Access Token>`

После проделанного вы можете пользоватся приватными проектами как пакетами `go get <package>`

## Docker & Docker Compose
Для успешного стягивания зависимостей - укажите ваши креды в рамках `devops/local/docker-compose.yml` файла. `CI_JOB_LOGIN`- ваш GitLab login(без @x5.ru), `CI_JOB_TOKEN` - AccessToken сгенерированный вами в настройках вышего аккаунта GitLab(галочка напротив api).
После чего нужно войти в docker registry: [Ссылка](https://wiki.x5.ru/pages/viewpage.action?pageId=1209378771)

Далее вы должны получить возможность запустить `make docker.start` & `make docker.build` без проблем

## Кодогенерация

Для того чтобы обновить файлы с заголовками DO NOT EDIT, нужно сперва удостовериться что доступ к бинарникам устанавливаемых через `go install ...` пакетов присутствует у вашей OS, добавьте в PATH путь к ним. Example for *nix: `export PATH=$PATH:$SOME_PATH_TO_DIRECTORY/go/bin`. Затем запустите команду `make tools` - она установит все необходимые для запуска утилиты и обновит пакеты. После чего запустите команду `go generate path/to/file.go`. Обязательно просмотрите изменения.
