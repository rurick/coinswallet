# Описание API

По умолчанию сервер API запускается на всех сетевых интерфейсах (порт 8081) по протоколу HTTP Параметры запроса
передаются в uri запроса и теле запроса

В случае ошибки в теле ответа возвращается значение "error" с пояснением ошибки. Пример ошибки при неправильном
синтаксисе запроса:

```json
{
  "error": "invalid character '}' looking for beginning of object key string"
}
```

---------------------------

## Создание аккаунта

* Метод: POST
* URI: /account/
* Тело запроса:

```json
{
  "name": "accountName"
}
```

Параметры:

* **name** - имя аккаунта. Строка содержащая латинские буквы и цифры длиною 4-32 символа. Имя является уникальным
  значение для каждого аккаунта. Имя аккаунта является регистрозависимым.

Пример:

```http request
POST http://localhost:8081/account/
content-type: application/json

{
  "name": "wallet1"
}
```

### Ответы

Успешное создание аккаунта:

```http request
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 08:48:04 GMT
Content-Length: 25

{
  "account_id": "wallet1"
}
```

Ошибка: аккаунт с таким именем существует

```http request
HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 08:49:12 GMT
Content-Length: 49

{
  "error": "create account error: duplicate name"
}
```

Ошибка в формате имени аккаунта

```http request
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 08:52:13 GMT
Content-Length: 32

{
  "error": "invalid name format"
}
```

-------------------

## Пополнение баланса аккаунта

* Метод: PATCH
* URI: /account/deposit/
* Тело запроса:

```json
{
  "name": "accountName",
  "amount": 31.2314
}
```

Параметры:

* **name** - имя аккаунта.
* **amount** - сумма пополнения (точность 4 знака после запятой).

Пример:

```http request
PATCH http://localhost:8081/account/deposit/
content-type: application/json

{
  "name": "wallet1",
  "amount": 31.2314
}
```

### Ответы

Успешное пополнение аккаунта:

```http request
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:03:33 GMT
Content-Length: 19

{
  "balance": 62.463
}
```

Аккаунт не найден:

```http request
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:04:02 GMT
Content-Length: 30

{
  "error": "account not found"
}
```

Некорректное значение суммы:

```http request
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:06:01 GMT
Content-Length: 34

{
  "error": "error in amount value"
}
```

-------------------

## Перевод между двумя аккаунтами

* Метод: PATCH
* URI: /account/transfer/
* Тело запроса:

```json
{
  "from": "wallet1",
  "to": "wallet2",
  "amount": 0.5
}
```

Параметры:

* **from** - имя аккаунта источника.
* **to** - имя аккаунта получателя.
* **amount** - сумма перевода (точность 4 знака после запятой).

Пример:

```http request
PATCH http://localhost:8081/account/transfer/
content-type: application/json

{
  "from": "wallet1",
  "to": "wallet2",
  "amount": 0.5
}
```

### Ответы

Успешное пополнение аккаунта:

```http request
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:18:11 GMT
Content-Length: 93

{
  "payment": {
    "account": "wallet1",
    "to_account": "wallet2",
    "amount": 0.5,
    "direction": "outgoing"
  }
}
```

Аккаунт не найден:

```http request
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:18:39 GMT
Content-Length: 35

{
  "error": "from account not found"
}
```

Некорректное значение суммы:

```http request
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:19:12 GMT
Content-Length: 34

{
  "error": "error in amount value"
}
```

Перевод между одним и тем же аккаунтом:
```http request
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:19:29 GMT
Content-Length: 45

{
  "error": "disable transfer to self account"
}
```


-------------------

## Получить список аккаунтов
Возвращает список всех существующих аккаунтов отсортированных по порядку создания

* Метод: GET
* URI: accounts/:offset/:limit/

Параметры:

* **offset** - целое значение, указывающее смещение, начиная с какого аккаунта возвращать результат
* **limit** - целое значение, указывающее сколько аккаунтов возвращать в результате. 
  Если limit=-1, то будут возвращены все аккаунты без ограничения.

Пример:

```http request
GET http://localhost:8081/accounts/0/-1/
```

### Ответы

Успешное выполнение запроса:

```http request
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:30:11 GMT
Content-Length: 263

{
  "list": [
    {
      "id": "wallet1",
      "balance": 61.963,
      "currency": "usd"
    },
    {
      "id": "wallet2",
      "balance": 0.5,
      "currency": "usd"
    }
  ]
}
```

Ошибка в параметрах:
```http request
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:32:01 GMT
Content-Length: 42

{
  "error": "error in offset, limit params"
}
```
-------------------

## Получить список платежей
Возвращается список всех платежей отсортированных по порядку обратном созданию

* Метод: GET
* URI: payments/:offset/:limit/

Параметры:

* **offset** - целое значение, указывающее смещение, начиная с какого платежа возвращать результат
* **limit** - целое значение, указывающее сколько платежей возвращать в результате. 
  Если limit=-1, то будут возвращены все платежи без ограничения.

Пример:

```http request
GET http://localhost:8081/payments/0/-1/
```

### Ответы

Успешное выполнение запроса:

```http request
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 10:02:45 GMT
Content-Length: 248

{
  "list": [
    {
      "account": "wallet1",
      "to_account": "wallet2",
      "amount": 0.5,
      "direction": "outgoing"
    },
    {
      "account": "",
      "to_account": "wallet1",
      "amount": 31.2315,
      "direction": "outgoing"
    },
    {
      "account": "",
      "to_account": "wallet1",
      "amount": 31.2315,
      "direction": "outgoing"
    }
  ]
}
```

Ошибка в параметрах:
```http request
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:32:01 GMT
Content-Length: 42

{
  "error": "error in offset, limit params"
}
```
-------------------

## Получить список платежей аккаунта
Возвращается список всех платежей для заданного аккаунта отсортированных по порядку обратном созданию

* Метод: GET
* URI: payments/:name/:offset/:limit/

Параметры:

* **name** - имя аккаунта
* **offset** - целое значение, указывающее смещение, начиная с какого платежа возвращать результат
* **limit** - целое значение, указывающее сколько платежей возвращать в результате. 
  Если limit=-1, то будут возвращены все платежи без ограничения.

Пример:

```http request
GET http://localhost:8081/payments/wallet2/0/-1/
```

### Ответы

Успешное выполнение запроса:
direction указывает направление перевода. 
В примере ниже в списке указан входящий платеж с аккаунта wallet1
```http request
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 10:05:11 GMT
Content-Length: 92

{
  "list": [
    {
      "account": "wallet2",
      "to_account": "wallet1",
      "amount": 0.5,
      "direction": "incoming"
    }
  ]
}

```

Ошибка в параметрах:
```http request
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Fri, 21 May 2021 09:32:01 GMT
Content-Length: 42

{
  "error": "error in offset, limit params"
}
```