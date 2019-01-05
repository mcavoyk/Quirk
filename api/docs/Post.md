*[Quirk](../../README.md) > [Quirk API](../README.md) > Posts*

---

## Posts

Posts on Quirk represent pieces of content viewable to users.
A post on Quirk is therefore used both for top level threads as well
as for comments to any other posts.

Actions: 

* [Create post](#create-post)

* [Delete post](#delete-post)

* [Search posts by location](#search-posts-by-location)

* [Get post replies](#get-post-replies)

---

### Create post
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

### Delete post
`DELETE /api/post/:id`

Example Request:

```http
DELETE /api/post/1FEj3nI39qajqfiVmxrpz9eexMQ HTTP/1.1
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

### Search posts by location

Search for posts by location.

`GET /posts`

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `lat` | string | yes | The latitude in degrees for the given search request |
| `lon` | string | yes | The longitude in degrees for the given search request |
| `page` | integer | no | The page number for posts |
| `per_page` | integer | no | The number of posts to return |


### Get post replies

Get a list of posts in which are descendants of a given post.

`GET /post/:post_id/posts`

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `post_id` | string | yes | The ID of the post |

