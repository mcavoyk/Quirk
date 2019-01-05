*[Quirk](../../README.md) > [Quirk API](../README.md) > Voting*

---

## Voting

A vote on Quirk represents an up vote, abstain vote, or downvote
(1, 0, -1) on some post. 

### Create Vote
`POST /api/post/:postID/vote?state=1`

Example Request:

```http
POST /api/post/:postID/vote?state=1 HTTP/1.1
Accept: application/json
Content-Type: application/json
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