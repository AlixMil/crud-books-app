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

```
SERVER_HOST=localhost  
SERVER_PORT=4001  
GOFILE_SERVICE_API_KEY=GOFILE_SERVICE_API_KEY  
GOFILE_FOLDER_TOKEN=GOFILE_FOLDER_TOKEN  
DB_NAME=crudbooks  
DB_LOGIN=admin  
DB_PWD=0000  
JWT_SECRET=JWT_SECRET_WORD  
JWT_TOKEN_TTL=10  
```      

> *Для запуска приложения нужно **обязательно** указать *gofile API_KEY*, найти его можно, авторизовавшись в одноименном сервисе, выступающим в роли хранилища: https://gofile.io/*

3. Запустите приложение при помощи команды ***docker-compose up***

### Vanilla Golang :icecream:

> Данный способ запуска приложения подраздумевает, наличия у вас в системе запущенного mongoDB экземпляра. 

Указать порт и хост вашего экземпляра MongoDB следует в файле docker-compose в теге environments, ключи: DB_HOST, DB_PORT

1. Выполните в консоли следующие команды: 

- ***go mod tidy***

2. Создайте в корневой директории проекта файл .env со следующим содержанием: 

```
SERVER_HOST=localhost  
SERVER_PORT=4001  
GOFILE_SERVICE_API_KEY=GOFILE_SERVICE_API_KEY  
GOFILE_FOLDER_TOKEN=GOFILE_FOLDER_TOKEN  
DB_NAME=crudbooks  
DB_LOGIN=admin  
DB_PWD=0000  
JWT_SECRET=JWT_SECRET_WORD  
JWT_TOKEN_TTL=10  
```  

> *Для запуска приложения нужно **обязательно** указать *gofile API_KEY*, найти его можно, авторизовавшись в одноименном сервисе, выступающим в роли хранилища: https://gofile.io/*

3. Запустите приложение при помощи вызова команды ***go run ./cmd/main.go*** 

## Использование :paw_prints:

> Общение с приложением, происходит посредством выполнения HTTP запросов следующего содержания: 

### Login

POST /login  

```json
{
	"email": "test@gmail.com",
	"password": "123"
}
```

Response 200,  
```json
{
	"token": "YOUR JWT TOKEN"
}
```

### Register

POST /register  

```json
{
	"email": "test@gmail.com",
	"password": "123"
}
```

Response 200,  
```json
{
	"token": "YOUR JWT TOKEN"
}
```

### GetBooks (для не авторизованных пользователей)

GET /books  

Response 200,  
```json
[
    {
        "Id": "643ff7ec8ddc071105dc4842",
        "Title": "Book",
        "Description": "Book description",
        "FileToken": "FILETOKEN",
        "Url": "https://gofile.io/",
        "OwnerEmail": "test@gmail.com"
    }
]
```

Поддерживается система фильтров и сортировки, [подробнее](#система-фильтров-и-сортировок)  

### GetBooks (для авторизованных пользователей)

GET /books  

Обязательно наличие [Authorization заголовка](#authorization-заголовок) в запросе.   

Response 200,  
```json
[
    {
        "Id": "643ff7ec8ddc071105dc4842",
        "Title": "Book",
        "Description": "Book description",
        "FileToken": "FILETOKEN",
        "Url": "https://gofile.io/",
        "OwnerEmail": "test@gmail.com"
    }
]
```

Поддерживается система фильтров и сортировки, [подробнее](#система-фильтров-и-сортировок)  

### Upload File

POST /files

Обязательно наличие [Authorization заголовка](#authorization-заголовок) в запросе.   

Content-Type: multipart/form-data

Тело запроса:  

```
FormData  
key - "file"
value - your_file
```

Response 200,
```json
{
    "fileToken": "5dfb4383-a312-438c-8dbf-1f3ce6fb1060"
}
```

### Create Book

POST /books

Обязательно наличие [Authorization заголовка](#authorization-заголовок) в запросе.   

Тело запроса в формате JSON
```json
{
	"fileToken": "5a70f95c-d7ff-4cd9-ae05-dcce2d68860e",
	"title": "Book",
	"description": ""
}
```

Response 200, "FILE_TOKEN"

### Get Book

GET /books/{:bookID}

Обязательно наличие [Authorization заголовка](#authorization-заголовок) в запросе.   

Response 200, ""
```json
{
    "fileURL": "https://gofile.io/",
    "title": "Book",
    "description": "Book description"
}
```

### Update Book

POST /books/{:bookID}

Обязательно наличие [Authorization заголовка](#authorization-заголовок) в запросе.   

Тело запроса в формате JSON
```json
{
	"fileToken": "5a70f95c-d7ff-4cd9-ae05-dcce2d68860e",
	"title": "Book",
	"description": ""
}
```

Response 200, "Book data successfully updated!"

### Delete Book

DELETE /books/{:bookID}

Обязательно наличие [Authorization заголовка](#authorization-заголовок) в запросе.     

Response 200, ""  



## Глоссарий :blue_book:

### Authorization заголовок  
Заголовок должен быть во всех запросах, оптравляемых к путям, требующих авторизации. Заголовок несет в себе JWT токен авторизации для идентификации пользователя.  

Пример заголовка
```
Authorization: Bearer YOUR_JWT_TOKEN
```

### Система фильтров и сортировок
Запросы книг GetBooks (для авторизованных и не авторизованных пользователей) позволяют использовать дополнительные параметры для фильтрации и сортировки результатов.  
Пример запроса с поиском по названию книги, или части названия.
> GET /books?search=НАЗВАНИЕ_КНИГИ  

Добавление лимита выдачи  

> GET /books?limit=15  

Параметр offset, для исключения из выдачи части книг (полезен при создании пагинации)

> GET /books?offset=15

Пример запроса с поиском по названию, лимитом выдачей и offset параметром

> GET /books?search=Сойер&limit=10&offset=10

## Используемые технологии 👨🏻‍💻

- [Golang](https://go.dev/)  

- [Echo](https://echo.labstack.com/)  

- [MongoDB](https://www.mongodb.com/)  

- [JWT](https://jwt.io/)  