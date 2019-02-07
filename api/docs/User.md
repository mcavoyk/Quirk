*[Quirk](../../README.md) > [Quirk API](../README.md) > Authentication*

---

## Users

The authentication API needs no input and returns a session token. The
token then should be used in the Authorization header as a bearer token.

Actions: 

* [Create user account](#create-user-account)

* [Login as user](#login-as-user)

* [Delete user account](#delete-user-account)

### Create user account

`POST /user`

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `username` | string | no |  The username of the account (globally unique) |
| `display_name` | string | no |  The display name of the account |
| `password` | string | no | The password for the user account |
| `email` | string | no | An email account to be associated with the user account |

Example Response:

```http
HTTP/1.1 200
Content-Type: application/json

{
    "id": "1GmcAx0HGGJcSed8cvlMmoK26fQ",
    "created_at": "2019-02-06T00:47:54Z",
    "updated_at": "2019-02-06T00:47:54Z",
    "username": "Falcon",
    "display_name": "Falcon",
    "email": ""
}
```

### Login as user

`POST /user/login`

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `username` | string | yes |  The username of the account |
| `password` | string | yes |  The password for the user account |
| `latitude` | number | yes |  The latitude of the user in degrees|
| `longitude` | number | yes |  The longitude of the user in degrees|

Example Response: 

```http
HTTP/1.1 200
Content-Type: application/json

{
    "expires": "2019-02-19T03:19:10Z",
    "token": "1Gk5Ru8uaTrLy4JQAfPRkHIOXrW"
}
```

### Delete user account

`DELETE /user/:id`