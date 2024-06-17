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
2. Unfollow user
   - Users should be allowed to unfollow other users.
3. Get followers
   - Users should get their followers.
   - The default page size is 50.
4. Get following
   - Users should get their following.
   - The default page size is 50.
5. Get friends
   - The definition of friends is users who follow each other.
   - Users should get their friends.
   - The default page size is 50.

