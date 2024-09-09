# Test task BackDev

Тестовое задание на позицию Junior Backend Developer

**Используемые технологии:**

- GO
- JWT
- PostgresSQL

## /signup - регистрация пользователей с валидации

```json
{
  "email": "your email",
  "password": "your password"
}
```

## /signin

авторизация пользователей. Выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса

1. Создает пара токенов с user_id и ip пользователя.
2. Сохраняет refresh токен в БД виде bcrypt хеша
3. Сохраняет refresh токен в cookie в формате base64.

## /refresh

1. Принимает Bearer "access token".
2. Получает refresh tokens из cookie
3. Декодирует из base64 формата в refresh token.
4. Берет user_id из claim токена
5. сравнивает с токенам пользователя на базе данных
6. Проверка ip, если не совпадает отправляет email на пользователя.
7. выдает новый access token.
