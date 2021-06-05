# Twidder
Imitating twitter backend API using Gin, GORM, and MySQL, just for learning purpose. 

## Features
- Create, deactivate, and reactivate account.
- Create tweets. Search a tweet with a keyword.
- Reply to a tweet. See replies to a tweet.
- Like and unlike a tweet. See users who like a tweet.
- Follow and unfollow other users. See other users followers and followings.
- See tweets from people you follow (timeline).
- See notifications when other users followed you, replied to your tweet, and liked your tweet.

## Models
![Models](./docs/db_models.png)

## Getting Started
### Using Local GO Installation
1. Create a `.env` file, set your env variables there, see the example in `.env.example`.
2. Install the dependencies.
```
go mod download
```
3. Start the service.
```
go run main.go
```
### Using Docker Compose
```
docker-compose build
docker-compose up
```

## Documentation
After running the server, you can checkout the docs at `/swagger/index.html`

## Test
```
go test -v ./tests/
```
