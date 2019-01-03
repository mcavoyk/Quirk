## Post

The Post API allows a user to create, update, read, or delete a post.

### Create Vote
`POST /api/post`

Example Request:

```http
POST /api/auth/keys HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization:Bearer 1FEToRKxL7aNTgGtYR91WszCvA

{
	"accessType": "public",
	"latitude": 0.0,
	"longitude": 0.0,
	"content": 
    	{ 
    		"title": "New post who dis",
    		"body": "ayy"
    	}
}
```

Example Response:

```http
HTTP/1.1 200
Content-Type: application/json

{
	"ID": "1FEj3nI39qajqfiVmxrpz9eexMQ"
}
```

### Delete Post
`DELETE /api/post/1FEj3nI39qajqfiVmxrpz9eexMQ`

Example Request:

```http
DELETE /api/auth/keys HTTP/1.1
Accept: application/json
Authorization:Bearer 1FEToRKxL7aNTgGtYR91WszCvA
```

Example Response:

```http
HTTP/1.1 200
Content-Type: application/json

{
	"OK"
}
```