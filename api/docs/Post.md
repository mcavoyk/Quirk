*[Quirk](../../README.md) > [Quirk API](../README.md) > Posts*

---

## Posts

Posts on Quirk represent pieces of content viewable to users.
A post on Quirk is therefore used both for top level threads as well
as for comments to any other posts.

Actions: 

* [Create post](#create-post)

* [Create post as reply](#create-post-as-reply)

* [Update post](#update-post)

* [Delete post](#delete-post)

* [Vote on a post](#vote-on-a-post)

* [Search posts by location](#search-posts-by-location)

* [Get post replies](#get-post-replies)


---

### Create post

This action creates a top level post which will be viewable to other users
based on their distance to where the post was created.

`POST /post` 

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `latitude` | number | yes |  The latitude in degrees of the user making the post |
| `longitude` | number | yes |  The longitude in degrees of the user making the post|
| `access_type` | string | yes | Visibility of post, `public` or `private` |
| `content` | string | yes | The content of the post |


### Create post as reply

This action creates a post which will be a reply to another post.

`POST /post/:id/post` 

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `id` | string | yes | The ID of the post to reply |
| `latitude` | number | yes |  The latitude in degrees of the user making the post |
| `longitude` | number | yes |  The longitude in degrees of the user making the post|
| `access_type` | string | yes | Visibility of post, `public` or `private` |
| `content` | string | yes | The content of the post |

### Update post

This action updates the post to contain the new content.

`PATCH /post/:id` 

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `id` | string | yes | The ID of the post to update |
| `latitude` | number | yes |  The latitude in degrees of the user making the post |
| `longitude` | number | yes |  The longitude in degrees of the user making the post|
| `access_type` | string | no | Visibility of post, `public` or `private` |
| `content` | string | no | The content of the post |


### Delete post

This action requires being the same user who created the post.

`DELETE /post/:id`

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `id` | string | yes | The ID of the post |


### Vote on a post

Posts on Quirk support being upvoted, abstained, or downvoted, which 
corresponds to a query string parameters `state` of 1, 0, or -1 respectively.
The score of a post will reflect the aggregate votes of all users who
have voted on the post.

`POST /posts/:id/vote`

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `state` | int | yes | An integer (1, 0, or -1) representing the vote |


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

`GET /post/:id/posts`

| **Attributes** | **Type** | **Required** | **Description** |
| ---------- | ---- | -------- | ----------- |
| `id` | string | yes | The ID of the post |
| `page` | integer | no | The page number for posts |
| `per_page` | integer | no | The number of posts to return |

