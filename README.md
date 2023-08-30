# AvitoTask

## Инструкция по запуску
1) Скачать репозиторий
2) docker-compose build
3) docker-compose up
4) В папке /internal/db/migrations лежит .sql файл с SQL запросами для инициализации базы данных. Нужно зайти в админер на localhost:8081 и выполнить их.

## Запросы
### 1) Метод создания сегмента. Принимает slug (название) сегмента.
### Запрос:
```
  curl --location 'localhost:8080/create' \
  --header 'Content-Type: application/json' \
  --data '{
      "slug": "a5"
  }'
```
### Ответ:
```
  200 OK
```
### 2) Метод удаления сегмента. Принимает slug (название) сегмента.
### Запрос:
 ```
   curl --location --request DELETE 'localhost:8080/delete' \
   --header 'Content-Type: application/json' \
   --data '{
     "slug": "a5"
   }'
 ```
### Ответ:
```
  200 OK
```
### 3) Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.
### Запрос:
 ```
   curl --location 'localhost:8080/add_user' \
   --header 'Content-Type: application/json' \
   --data '{
       "user_id": 4,
       "segments_to_add": ["a1", "a2"],
       "segments_to_delete": ["a2"]
   }'
 ```
### Ответ:
```
  200 OK
```
### 4) Метод получения активных сегментов пользователя. Принимает на вход id пользователя.
   
### Запрос:
```
 curl --location --request GET 'localhost:8080/user_active_segments' \
 --header 'Content-Type: application/json' \
 --data '{
     "user_id": 4
 }'
```
### Ответ:
```
{
  "active_segments": [
      "a1"
  ]
}
```

