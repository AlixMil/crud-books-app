CrudBooks Backend App 📚
=============================

Crud books - приложение для хранения книг, и любых других файлов в формате PDF.


![image](./image%20for%20readme.jpg)


## Установка :factory:

[Docker-compose](#docker-compose-cake)  

[Vanilla Golang](#vanilla-golang-icecream)


### Docker-compose :cake:

> Перед установкой убедитесь, что система имеет установленный docker engine

1. Выполните в консоли следующие команды: 

- ***go mod tidy***

- ***docker pull mongo***

2. Создайте в корневой директории проекта файл .env со следующим содержанием: 

> SERVER_HOST=localhost  
> SERVER_PORT=4001  
> GOFILE_SERVICE_API_KEY=GOFILE_SERVICE_API_KEY  
> GOFILE_FOLDER_TOKEN=GOFILE_FOLDER_TOKEN  
> DB_NAME=crudbooks  
> DB_LOGIN=admin  
> DB_PWD=0000  
> JWT_SECRET=JWT_SECRET_WORD  
> JWT_TOKEN_TTL=10  


3. Запустите приложение при помощи команды ***docker-compose up***

### Vanilla Golang :icecream:

> Данный способ запуска приложения подраздумевает, наличия у вас в системе запущенного mongoDB экземпляра. 

Указать порт и хост вашего экземпляра MongoDB следует в файле docker-compose в теге environments, ключи: DB_HOST, DB_PORT

1. Выполните в консоли следующие команды: 

- ***go mod tidy***

2. Создайте в корневой директории проекта файл .env со следующим содержанием: 

> SERVER_HOST=localhost  
> SERVER_PORT=4001  
> GOFILE_SERVICE_API_KEY=GOFILE_SERVICE_API_KEY  
> GOFILE_FOLDER_TOKEN=GOFILE_FOLDER_TOKEN  
> DB_NAME=crudbooks  
> DB_LOGIN=admin  
> DB_PWD=0000  
> JWT_SECRET=JWT_SECRET_WORD  
> JWT_TOKEN_TTL=10  

3. Запустите приложение при помощи вызова команды ***go run ./cmd/main.go*** 

Используемые технологии 👨🏻‍💻
-----------

- [Golang](https://go.dev/)  

- [Echo](https://echo.labstack.com/)  

- [MongoDB](https://www.mongodb.com/)  

- [JWT](https://jwt.io/)  