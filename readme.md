CrudBooks Backend App :books:
=============================

Crud books - приложение для хранения книг, и любых других файлов в формате PDF.


![image](./image%20for%20readme.jpg)

### Содержание :earth_africa:
1. [Установка](#установка-factory)
2. [Использование](#использование-paw_prints)
3. [Глоссарий](#глоссарий-blue_book)
4. [Используемые технологии](#используемые-технологии-computer)
5. [Покрытие тестами](#покрытие-тестами-mag_right)

## Установка :factory:

> *Для запуска приложения нужно **обязательно** получить *gofile API_KEY*, данный сервис выступает в роли хранилища, отвечает за сохранение и выдачу файлов на скачивание: https://gofile.io/*


### Docker-compose :cake:

> Перед установкой убедитесь, что система имеет установленный docker engine

1. Выполните в консоли следующие команды: 

```
$ go mod tidy
$ docker pull mongo
```

2. Замените следующие переменные окружения в файле "Docker-compose.yml", тег "environment":
```yml
- GOFILE_SERVICE_API_KEY=YOUR_GOFILE_SERVICE_API_KEY
- JWT_SECRET=YOUR_JWT_SECRET
``` 

3. Запустите приложение при помощи команды 
```
$ docker-compose up
```


## Использование :paw_prints:

> Общение с приложением, происходит посредством отправки HTTP запросов следующего содержания: 

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


### Upload File

POST /files

Обязательно наличие [Authorization заголовка](#authorization-заголовок) в запросе.   

Content-Type: multipart/form-data

Пример запроса:  

```
$ curl -X POST localhost:4001
-H "Content-Type: multipart/form-data" 
-d "file=your_file" 
```

Response 200,
```json
{
    "fileToken": "5dfb4383-a312-438c-8dbf-1f3ce6fb1060"
}
```

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
Заголовок должен быть во всех запросах, оптравляемых к путям, требующим авторизации. Заголовок несет в себе JWT токен авторизации для идентификации пользователя.  

Пример заголовка
```yml
Authorization: Bearer YOUR_JWT_TOKEN
```

### Система фильтров и сортировок
Запросы книг GetBooks (для авторизованных и не авторизованных пользователей) позволяют использовать дополнительные параметры для фильтрации, сортировки, лимитирования и сдвига выдачи.  
Пример запроса с поиском по названию книги, или части названия.
> GET /books?search=НАЗВАНИЕ_КНИГИ  

Добавление лимита выдачи  

> GET /books?limit=15  

Параметр offset, для исключения из выдачи части книг (полезен при создании пагинации)

> GET /books?offset=15

Пример запроса с поиском по названию, лимитом выдачей и offset параметром

> GET /books?search=Сойер&limit=10&offset=10

## Используемые технологии :computer:

- [Golang](https://go.dev/)  

- [Echo](https://echo.labstack.com/)  

- [MongoDB](https://www.mongodb.com/)  

- [JWT](https://jwt.io/)  

## Покрытие тестами :mag_right:

Total Coverage: 72.5%