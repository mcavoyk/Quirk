## Vote

The Vote API allows a user to send an upvote, abstain, or downvote (1, 0, -1)
for a particular post.

### Create Vote
`POST /api/vote`

Example Request:

```http
POST /api/auth/keys HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization:Bearer 1FEToRKxL7aNTgGtYR91WszCvA

{
	"state": 1,
	"postID": "1FEWViwSeKkQ8hqaVkM2crOezbj"
}
```

Example Response:

```http
HTTP/1.1 200
Content-Type: application/json

{
	"OK"
}
```