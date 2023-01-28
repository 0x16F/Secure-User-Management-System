<h2>Secure User Management System</h2>

<p>Этот проект представляет собой сервер, написанный на Golang, который использует базу данных Postgresql и токены JWT для авторизации. Это позволяет администраторам выполнять операции CRUD с таблицей пользователей, в то время как все авторизованные пользователи могут считывать данные. Кроме того, он включает в себя запрещенную роль, которая не имеет никаких прав доступа вообще.</p>

<p>После включения сервера будет доступ к Swagger по этой ссылке <a href="http://localhost:8080/swagger/index.html#/users/update%20user">*клик*</a></p>

<h2>Стэк технологий</h2>
<ul>
    <li>Golang</li>
    <li>Postgres</li>
    <li>Docker</li>
    <li>Docker-compose</li>
    <li>Swagger</li>
    <li>Makefile</li>
</ul>

<h2>Дополнительные библиотеки</h2>
<ul>
    <li><a href="https://echo.labstack.com/">Echo</a> - web framework</li>
    <li><a href="https://github.com/allegro/bigcache">bigcache</a> - чтобы в случае бана пользователя, он банился сразу, а не по истечению времени жизни токена</li>
    <li><a href="https://github.com/golang-migrate/migrate">Migrate</a> - миграции бд</li>
    <li><a href="https://github.com/spf13/viper">viper</a> - чтение данных из конфига</li>
    <li><a href="https://github.com/go-pg/pg">pg</a> - для работы с БД</li>
</ul>

<h2>Данные от админки</h2>

<p>Логин: admin</p>
<p>Пароль: 12345678</p>

<h2>Запуск</h2>

```
docker compose up
```