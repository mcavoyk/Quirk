## Authentication

The authentication API needs no input and returns a session token.


### Create Auth Token
`GET /api/auth/token`

Example Request:

```http
GET /api/auth/keys HTTP/1.1
Accept: application/json
```

Example Response:

```http
HTTP/1.1 200
Content-Type: application/json

{
	"token": "1FEToRKxL7aNTgGtYR91WszCvAb"
}
```

### Validate Auth Token

`GET /api/auth/validate`

Example Request:

```http
GET /api/auth/keys HTTP/1.1
Accept: application/json
Authorization:Bearer 1FEToRKxL7aNTgGtYR91WszCvAb
```

Example Response: Valid token

```http
HTTP/1.1 200
Content-Type: application/json

{
	"CreatedAt": "2019-01-03T00:19:24Z",
	"ID": "1FEWViwSeKkQ8hqaVkM2crOezbj",
	"IP": "",
	"UsedAt": "2019-01-03T00:19:24Z"
}
```

Example Response: Invalid token
```http
HTTP/1.1 403
Content-Type: application/json

```