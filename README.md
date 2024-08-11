# Go rate limiter application

A basic application where there's a middleware rate limiting requests by using a redis key:value storage.

### To run it you just have to exec a 'docker compose up' command.

There'll be a exposed endpoint on localhost:8080/hello

To change the rate limits and the cooldown wait time you can do it by changing
the Environment variables inside the docker-compose.yaml file:

      - LIMIT_IP=7
      - LIMIT_TOKEN=5
      - OVER_LIMIT_COOLDOWN=7