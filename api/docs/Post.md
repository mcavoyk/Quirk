## Post

Posts on Quirk represent pieces of content viewable to users.
A post on Quirk is therefore used both for top level threads as well
as for comments to any other posts.

### Create Vote
`POST /api/post` Creates a top level post

`POST /api/post/:parentID/post` Creates post in reply to another post

Example Request:

```http
POST /api/post/1FEWViwSeKkQ8hqaVkM2crOezbj/post HTTP/1.1
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
	"ParentID": "1FEWViwSeKkQ8hqaVkM2crOezbj"
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