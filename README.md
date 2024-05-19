![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zuzi90/tz-enricher)
![Go](https://img.shields.io/badge/-Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![SQL](https://img.shields.io/badge/-SQL-4479A1?style=flat-square&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/-Redis-DC382D?style=flat-square&logo=redis&logoColor=white)
![Kafka](https://img.shields.io/badge/-Kafka-231F20?style=flat-square&logo=apache-kafka&logoColor=white)
![Grafana](https://img.shields.io/badge/-Grafana-F46800?style=flat-square&logo=grafana&logoColor=white)
![Prometheus](https://img.shields.io/badge/-Prometheus-E6522C?style=flat-square&logo=prometheus&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-85EA2D?style=flat-square&logo=swagger&logoColor=white)

# *Сервис для валидации и обогащения данных*

## Задача
Реализовать сервис, который будет получать поток ФИО, из открытых API обогащать ответ наиболее вероятными возрастом, полом, 
национальностью и сохранять данные в БД. По запросу выдавать инфу о найденных людях. 

Необходимо реализовать следующее:

    1. Сервис слушает очередь кафки FN, в котором приходит информация с ФИО в формате 

    {
		"name": "Andrey",
		"surname": "Sahorov",
		"patronymic": "Dmitrievich" // необязательно
	}
    
    2. В случае некорректного сообщения, обогатить его причиной ошибки (нет обязательного поля, некорректный формат...) и 
       отправить в очередь кафки WRONG_FN

    3. Корректное сообщение обогатить
        1. Возрастом - https://api.agify.io/?name=Dmitriy
        2. Полом - https://api.genderize.io/?name=Dmitriy		
        3. Национальностью - https://api.nationalize.io/?name=Dmitriy 

    4. Обогащенное сообщение положить в БД Postgres (структура БД должна быть создана путем миграций)

    5. Выставить REST методы
        1. Для получения данных с различными фильтрами и пагинацией 
        2. Для добавления новых людей	
        3. Для удаления по идентификатору	
        4. Для изменения сущности
    
    6. Предусмотреть кэширование данных в Redis

    7. Покрыть код логами

    8. Покрыть бизнес-логику unit-тестами

### Описание сервиса
Реализован CRUD API пользователя, Kafka consumer для чтения из топика FN, producer для
записи не валидных запросов в топик WRONG_FN.

### Установка и запуск
Клонировать репозиторий. В корне проекта выполнить:

+ make up
+ make run

Сервер запустится по адресу http://localhost:5005

Swagger доступен по адресу http://localhost:5005/swagger/index.html после запуска сервера

Список команд описан в Makefile в корне проекта

<br>

#### Пример POST запроса на url http://localhost:5005/api/v1/users:
```json
    {
      "name":"Michael",
      "surname":"Jackson",
      "patronymic":"Joseph",
      "age":55,
      "gender":"male",
      "nationality":"american"
    }
```

#### Пример ответа:
```json
    {
      "id": 181,
      "name": "Michael",
      "surname": "Jackson",
      "patronymic": "Joseph",
      "age": 55,
      "gender": "male",
      "nationality": "american",
      "isDeleted": false,
      "createdAt": "2024-05-19T21:06:10.900143+03:00",
      "updatedAt": "2024-05-19T21:06:10.900143+03:00"
    }
```















