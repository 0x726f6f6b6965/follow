# Follow Feature

## About it

This is a repository for following each other on social media features.

It uses the singleflight package to address a large number of requests in a short period and leverages the counting bloom filter to handle invalid users.

## Entity Relationship Diagram

Because the followers are also users, this will be a many-to-many model.

<img src="assets/img/erd.png" width="300"> 

## Features

1. Follow user
   - Users should be allowed to follow another user.
   - URL: `POST /v1/relationship/follow`
2. Unfollow user
   - Users should be allowed to unfollow other users.
   - URL: `POST /v1/relationship/unfollow`
3. Get followers
   - Users should get their followers.
   - The default page size is 50.
   - URL: `GET /v1/relationship/list/followers/${usrname}`
   - Query parameters: `page_token`, `size`
4. Get following
   - Users should get their following.
   - The default page size is 50.
   - URL: `GET /v1/relationship/list/following/${usrname}`
   - Query parameters: `page_token`, `size`
5. Get friends
   - The definition of friends is users who follow each other.
   - Users should get their friends.
   - The default page size is 50.
   - URL: `GET /v1/relationship/list/friends/${usrname}`
   - Query parameters: `page_token`, `size`

## How to Run it
1. Set `SERVICE_NAME` in the `.env` file, and then run `make service-build` to build an image
2. Set db and Redis relevant information in the `.env` file and run `make service-up` to run the service.

## Tech
- Go: 1.22.1
- PostgreSQL: 16.1
- Redis: 7.2.3
