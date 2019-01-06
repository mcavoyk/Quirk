*[Quirk](../../README.md) > [Quirk API](../README.md) > Authentication*

---

## Authentication

The authentication API needs no input and returns a session token. The
token then should be used in the Authorization header as a bearer token.

Actions: 

* [Create auth token](#create-auth-token)

* [Validate auth token](#validate-auth-token)

### Create auth token

`GET /api/auth/token`

Example Response:

```http
HTTP/1.1 200
Content-Type: application/json

{
  "token": "1FEToRKxL7aNTgGtYR91WszCvAb"
}
```

### Validate auth token

`GET /api/auth/validate`

Example Request:

```http
GET /api/auth/keys HTTP/1.1
Accept: application/json
Authorization:Bearer 1FEToRKxL7aNTgGtYR91WszCvAb
```

Example Response: 

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
