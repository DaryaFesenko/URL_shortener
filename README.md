# URL_shortener


Добавить ссылку
POST запрос /link с параметрами  longLink и токеном авторизации

Получить длинную ссылку по коротой
curl -X GET "localhost:9000?shortLink=:9000/2zFBZ8olHz"

Удалить ссылку
curl -X DELETE "localhost:9000/link?linkID=7b5b872e-e587-47fe-9e4e-c5e8edfde8e0" -H"Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzA1OTk0MTIuMDU3NzY3OSwiaWF0IjoxNjI5OTk0NjEyLjA1Nzc2OSwidXNlcklkIjoiOTVhMDhiOTAtYmY3My00ZWZiLTkxYWMtODkxNGZjMTUyYzhjIn0.sHquETD2USwsTjoIEOfDiDqgOw9r-TMEO91CMQdd8Co"

Регистрация юзера
curl -X POST "localhost:9000/auth/signUp" -d"{\"login\":\"test1\", \"password\":\"12345\"}"

Статистика о ссылке
curl -X GET "localhost:9000/link?id=cb6f7891-b319-43f1-baa4-849d8d787e0b" -H"Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzE4MzI1ODcuODA1NzA2LCJpYXQiOjE2MzEyMjc3ODcuODA1NzA4LCJ1c2VySWQiOiJiNDdhNGNlZS0yMTY5LTRlYTItODEwNi1kOTY5NTMwY2IyODkifQ.Kp8izWdP_hc8rxp-aqPKoMflCj7WT7HsJwpYlNT4N4s"

Вход


- сваггер
- добавить конфиг


- ну и вылить на хероку

