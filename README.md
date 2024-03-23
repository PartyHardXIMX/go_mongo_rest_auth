# go_mongo_rest_auth
1. Скачать образ из DockerHub.
docker pull partyhardximx/mongodb:latest
Зайти в Docker Desktop, выбрать Images, запустить скачанный образ, в Optional Settings ввести порт 27017. Запустить контейнер.
2. Запустить сервер.
Зайти в папку server и запустить server.go # go run server.go Сервер запустится на http://localhost:8080/
3. Запустить клиент.
Открыть папку client в IDE. Открыть терминал и ввести go run client.go
