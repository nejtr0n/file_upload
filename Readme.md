# Readme
Для запуска приложения необходимо выполнить команду 
```
docker-compose up
```
Запускается 5 контейнеров:
1) consul - для service discovery
2) api - прокси для http запросов
3) minio - amazon s3 подобное хранилище (https://www.minio.io/)
4) storage_web - web интерфейс для приема http запросов на загрузку файлов
5) storage_srv - обработчик rpc запросов на хранение файлов

Хранилище minio имеет web интерфейс по адресу
http://localhost:9000/minio

Логин/пароль - это связка MINIO_ACCESS_KEY / MINIO_SECRET_KEY из файла docker-compose.yml

Для тестирование сервиса воспользоваться тестовойстраницей загрузки по адресу
http://localhost:8080/storage/test

## Загрузка multipart
Выбираем несколько файлов через ctrl в поле файил формы и жмем Upload Multiple Image

Все загруженные картинки будут доступны в веб интерфейсе http://localhost:9000/minio/images/

Публичный доступ к картинке осуществляется по адресу
http://localhost:9000/images/имя_картики, например
http://localhost:9000/images/1.png

Превью картинки расположен по адресу 
http://localhost:9000/thumbs/1.png

## Загрузка через url
Для теста вводим https://staticaltmetric.s3.amazonaws.com/uploads/2015/12/Altmetric_rgb.png
и жмем Upload by link

После загрузки фаил будет доступен по адресу 
http://localhost:9000/images/Altmetric_rgb.png

Превью картинки расположен по адресу 
http://localhost:9000/thumbs/Altmetric_rgb.png

## Загрузка json
Выбираем несколько файлов через ctrl в поле фаил формы и жмем Upload json

Все загруженные картинки будут доступны в веб интерфейсе http://localhost:9000/minio/images/

Превью картинки расположен по адресу http://localhost:9000/minio/thumbs/

# P.S
Собранные бинарные файлы располагаются в папке bin
Собрать можно через make
```
cd storage && make build
```
Собрать protobuf
```
cd storage && make buf
```