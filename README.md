# bank-api
Bank-API is an API-based bank application that can be used for:
- Manage accounts based on the currency used.
- Transfer money in concurrency and record the transaction in Entries.

## Installation
1. Copy [env.example](./env.example) to .env
```
cp env.example .env
```

2. Edit the values in .env with your database source
```
USERNAME=postgres
PASSWORD=1234
HOST=172.17.0.2
PORT=5432
DB_NAME=bank_api
```

3. [Install docker](https://docs.docker.com/engine/install/)

4. [Install cli golang migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

5. Migrate up
```
make migrateup
```

6. Build image docker
```
docker build -t bank-api:latest . 
```

7. Run the container
```
docker run --name bank-api -p 8080:8080 -e GIN_MODE=release bank-api:latest 
```

We can add or remove the currency support of the bank by changing the file: [utils/currency](./utils/currency.go)
```
package utils

const (
	IDR = "IDR"
	USD = "USD"
	EUR = "EUR"
    NewCurrency = "NewCurrency"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, IDR, NewCurrency:
		return true
	}
	return false
}

```

## REST-API
### User Registration
POST: /user
```
curl -i -X POST -H "Content-Type: application/json" -d '{ "username": "fulan1234","password": "hardpassword1234","full_name": "Fulan Fulano","email": "fulana@email.com"}' localhost:8080/user
```
Response
```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 10:44:07 GMT
Content-Length: 161

{"username":"fulan1234","full_name":"Fulan Fulano","email":"fulana@email.com","pwd_changed_at":"0001-01-01T00:00:00Z","created_at":"2023-10-26T10:44:07.312905Z"}
```

### Login
POST: /user/login
```
curl -i -X POST -H "Content-Type: application/json" -d '{ "username": "fulan1234","password": "hardpassword1234"}' localhost:8080/user/login
```
Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 10:46:46 GMT
Content-Length: 965

{"sessions_id":"9a03e7a7-a054-4996-8376-5a327674b6d6","access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2ZDRkYTI0LTE0N2UtNDBlZi04ZDE1LWE1YTJhZTZiZTAzYiIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMDo0Njo0NS4wMjk5ODk2MTRaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MDY6NDUuMDI5OTg5NzI3WiJ9.UfVqQl41GL6QKVqImLhEyuB6dGs0-Pellns2xqgMyto","access_token_expires_at":"2023-10-26T11:06:45.029989727Z","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjlhMDNlN2E3LWEwNTQtNDk5Ni04Mzc2LTVhMzI3Njc0YjZkNiIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMDo0Njo0NS4wMzAwNTYwMDhaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjdUMTA6NDY6NDUuMDMwMDU2MDk0WiJ9.MqyavTS9WdPYNuODCXVDR50EdYvnTAQpT113HwvKs0Q","refresh_token_expires_at":"2023-10-27T10:46:45.030056094Z","user":{"username":"fulan1234","full_name":"Fulan Fulano","email":"fulana@email.com","pwd_changed_at":"0001-01-01T00:00:00Z","created_at":"2023-10-26T10:44:07.312905Z"}}
```

### Logout
POST: /user/logout
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2ZDRkYTI0LTE0N2UtNDBlZi04ZDE1LWE1YTJhZTZiZTAzYiIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMDo0Njo0NS4wMjk5ODk2MTRaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MDY6NDUuMDI5OTg5NzI3WiJ9.UfVqQl41GL6QKVqImLhEyuB6dGs0-Pellns2xqgMyto" localhost:8080/user/logout
```
Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 10:52:03 GMT
Content-Length: 33

{"message":"successfully logout"}
```

### Create account
POST: /account
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMyZTQ1OTE3LWMyZTktNDkzNS05MjM4LTI3ZjFiODdhZWEyMyIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMTowMzoyMC45NTU5OTcxNTJaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MjM6MjAuOTU1OTk3Mjc0WiJ9.aRxZqIlgcgzDbLzLl31vlgkHhvDuAeUkACKqMnnXjF4" -H "Content-Type: application/json" -d '{"currency": "IDR"}' localhost:8080/account
```
Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:04:12 GMT
Content-Length: 136

{"id":"cf4177e5-9a09-47a7-89c3-e6143a32a2d7","owner":"fulan1234","balance":0,"currency":"IDR","created_at":"2023-10-26T11:04:12.06307Z"}
```

### Get Account
GET: /account/:id
```
curl -i -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMyZTQ1OTE3LWMyZTktNDkzNS05MjM4LTI3ZjFiODdhZWEyMyIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMTowMzoyMC45NTU5OTcxNTJaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MjM6MjAuOTU1OTk3Mjc0WiJ9.aRxZqIlgcgzDbLzLl31vlgkHhvDuAeUkACKqMnnXjF4" localhost:8080/account/cf4177e5-9a09-47a7-89c3-e6143a32a2d7
```
Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:07:19 GMT
Content-Length: 94

{"id":"cf4177e5-9a09-47a7-89c3-e6143a32a2d7","owner":"fulan1234","balance":0,"currency":"IDR"}
```

### List Account
GET: /account?page=
```
curl -i -L -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImViNjNmZjlmLTQ4NzUtNDJlZS1hZGVmLTE2MjczMmI3NjI3NSIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMToxNjowNC42NzM4MzIyNzFaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MzY6MDQuNjczODMyMzk3WiJ9.IwZlTqNdM10Ps1o_NbllOOY7FhgA6AEHW06wfiD-FOo" 'localhost:8080/account?page=1'
```
Response
```
HTTP/1.1 301 Moved Permanently
Content-Type: text/html; charset=utf-8
Location: /account/?page=1
Date: Thu, 26 Oct 2023 11:17:14 GMT
Content-Length: 51

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:17:14 GMT
Content-Length: 286

[{"id":"135418bc-067d-45ed-8286-e64867bee809","owner":"fulan1234","balance":0,"currency":"EUR"},{"id":"6148f8e0-24c0-4b8d-9c5c-31ad01ef16a8","owner":"fulan1234","balance":0,"currency":"USD"},{"id":"cf4177e5-9a09-47a7-89c3-e6143a32a2d7","owner":"fulan1234","balance":0,"currency":"IDR"}]
```

### Transfer money
POST: /transfer
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjU3MmVjNjQ0LWU3ZjEtNDQxYi1iZjM1LTZjOTZhM2FjNTUzMCIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMToyMjoxOS4xODEzNTIzODNaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6NDI6MTkuMTgxMzUyNloifQ.Vrlo6flMYtzj5nOMJEjDDSAlSjo3L3vUTHbBFkzFAM8" -H "Content-Type: application/json" -d '{"sender_id": "cf4177e5-9a09-47a7-89c3-e6143a32a2d7","receiver_id": "ad20fcd5-66b7-402d-9d66-289ab74b206a","amount": 500,"currency": "IDR"}' localhost:8080/transfer
```
Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:24:56 GMT
Content-Length: 853

{"transfer":{"id":"cc752f7f-2c52-45e2-8a9e-ded36d5f2db5","sender_id":"cf4177e5-9a09-47a7-89c3-e6143a32a2d7","receiver_id":"ad20fcd5-66b7-402d-9d66-289ab74b206a","amount":500,"created_at":"2023-10-26T11:24:56.75861Z"},"sender":{"id":"cf4177e5-9a09-47a7-89c3-e6143a32a2d7","owner":"fulan1234","balance":99500,"currency":"IDR","created_at":"2023-10-26T11:04:12.06307Z"},"receiver":{"id":"ad20fcd5-66b7-402d-9d66-289ab74b206a","owner":"gizka","balance":500,"currency":"IDR","created_at":"2023-10-26T00:54:41.610874Z"},"sender_entry":{"id":"ebf84af6-0fa1-40b9-9c48-5a95cd55a083","account_id":"cf4177e5-9a09-47a7-89c3-e6143a32a2d7","amount":-500,"created_at":"2023-10-26T11:24:56.75861Z"},"receiver_entry":{"id":"ac1ac727-603a-4f4d-8a55-1a38eb09b634","account_id":"ad20fcd5-66b7-402d-9d66-289ab74b206a","amount":500,"created_at":"2023-10-26T11:24:56.75861Z"}}
```

### Authorization check

#### Create account
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMyZTQ1OTE3LWMyZTktNDkzNS05MjM4LTI3ZjFiODdhZWEyMyIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMTowMzoyMC45NTU5OTcxNTJaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MjM6MjAuOTU1OTk3Mjc0WiJ9.aRxZqIlgcgzDbLzLl31vlgkHhvDuAeUkACKqMnnXjF4" -H "Content-Type: application/json" -d '{"currency": "IDR"}' localhost:8080/account

HTTP/1.1 401 Unauthorized
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:37:39 GMT
Content-Length: 25

{"error":"invalid token"}
```

#### Access account
```
curl -i -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjMyZTQ1OTE3LWMyZTktNDkzNS05MjM4LTI3ZjFiODdhZWEyMyIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMTowMzoyMC45NTU5OTcxNTJaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MjM6MjAuOTU1OTk3Mjc0WiJ9.aRxZqIlgcgzDbLzLl31vlgkHhvDuAeUkACKqMnnXjF4" localhost:8080/account/cf4177e5-9a09-47a7-89c3-e6143a32a2d7

HTTP/1.1 401 Unauthorized
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:39:38 GMT
Content-Length: 25

{"error":"invalid token"}
```

#### List account
```
curl -i -L -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImViNjNmZjlmLTQ4NzUtNDJlZS1hZGVmLTE2MjczMmI3NjI3NSIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMToxNjowNC42NzM4MzIyNzFaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MzY6MDQuNjczODMyMzk3WiJ9.IwZlTqNdM10Ps1o_NbllOOY7FhgA6AEHW06wfiD-FOo" 'localhost:8080/account?page=1'
HTTP/1.1 301 Moved Permanently
Content-Type: text/html; charset=utf-8
Location: /account/?page=1
Date: Thu, 26 Oct 2023 11:40:50 GMT
Content-Length: 51

HTTP/1.1 401 Unauthorized
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:40:50 GMT
Content-Length: 25

{"error":"invalid token"}
```

#### Transfer
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjU3MmVjNjQ0LWU3ZjEtNDQxYi1iZjM1LTZjOTZhM2FjNTUzMCIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMToyMjoxOS4xODEzNTIzODNaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6NDI6MTkuMTgxMzUyNloifQ.Vrlo6flMYtzj5nOMJEjDDSAlSjo3L3vUTHbBFkzFAM8" -H "Content-Type: application/json" -d '{"sender_id": "cf4177e5-9a09-47a7-89c3-e6143a32a2d7","receiver_id": "ad20fcd5-66b7-402d-9d66-289ab74b206a","amount": 500,"currency": "IDR"}' localhost:8080/transfer

HTTP/1.1 403 Forbidden
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:41:30 GMT
Content-Length: 33

{"error":"you must log in first"}
```

#### Logout
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2ZDRkYTI0LTE0N2UtNDBlZi04ZDE1LWE1YTJhZTZiZTAzYiIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMDo0Njo0NS4wMjk5ODk2MTRaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTE6MDY6NDUuMDI5OTg5NzI3WiJ9.UfVqQl41GL6QKVqImLhEyuB6dGs0-Pellns2xqgMyto" localhost:8080/user/logout
HTTP/1.1 401 Unauthorized
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:42:58 GMT
Content-Length: 25

{"error":"invalid token"}
```

### Business logic
#### Money transfers must be to accounts that have the same currency
account IDR to account EUR
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImZhMTBmOGQ0LWUzMzEtNDdlYy04YjhmLTc4ZjliMzllYTA4OCIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMTo0NTowMy45NzAwNTY5OTJaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTI6MDU6MDMuOTcwMDU3MTE0WiJ9.kbchWfeqzXOeYZcD2u2SMHCj0U0b96OTsUG8A9muazI" -H "Content-Type: application/json" -d '{"sender_id": "cf4177e5-9a09-47a7-89c3-e6143a32a2d7","receiver_id": "135418bc-067d-45ed-8286-e64867bee809","amount": 500,"currency": "IDR"}' localhost:8080/transfer

HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:47:56 GMT
Content-Length: 28

{"error":"invalid currency"}
```

#### the currency in the json body must match the currency of the sending account
Sender IDR but the currency in the json body is changed to another currency
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImZhMTBmOGQ0LWUzMzEtNDdlYy04YjhmLTc4ZjliMzllYTA4OCIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMTo0NTowMy45NzAwNTY5OTJaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTI6MDU6MDMuOTcwMDU3MTE0WiJ9.kbchWfeqzXOeYZcD2u2SMHCj0U0b96OTsUG8A9muazI" -H "Content-Type: application/json" -d '{"sender_id": "cf4177e5-9a09-47a7-89c3-e6143a32a2d7","receiver_id": "135418bc-067d-45ed-8286-e64867bee809","amount": 500,"currency": "EUR"}' localhost:8080/transfer

HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:50:01 GMT
Content-Length: 28

{"error":"invalid currency"}
```

#### Unable to transfer money over account balance
```
curl -i -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImZhMTBmOGQ0LWUzMzEtNDdlYy04YjhmLTc4ZjliMzllYTA4OCIsInVzZXJuYW1lIjoiZnVsYW4xMjM0IiwiaXNzdWVkX2F0IjoiMjAyMy0xMC0yNlQxMTo0NTowMy45NzAwNTY5OTJaIiwiZXhwaXJlZF9hdCI6IjIwMjMtMTAtMjZUMTI6MDU6MDMuOTcwMDU3MTE0WiJ9.kbchWfeqzXOeYZcD2u2SMHCj0U0b96OTsUG8A9muazI" -H "Content-Type: application/json" -d '{"sender_id": "cf4177e5-9a09-47a7-89c3-e6143a32a2d7","receiver_id": "ad20fcd5-66b7-402d-9d66-289ab74b206a","amount": 100000000000000,"currency": "IDR"}' localhost:8080/transfer

HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=utf-8
Date: Thu, 26 Oct 2023 11:52:22 GMT
Content-Length: 62

{"error":"insufficient balance: 99000 \u003c 100000000000000"}
```

