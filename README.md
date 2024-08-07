# Prerequisite
1. Docker
2. make (brew install make)
3. MongoDB instance

# Configuration
1. Set the value of `MONGODB_URI` in `docker-compose.yml` to the uri shared over email. It's a temporary user's credentials which will expire in a week.

# Setup
1. Start Docker
2. `make network` - to create the docker network for the containers to interact with each other.
3. `make redis` - to start the redis container.
4. `sudo make run` - It gets the required images, makes the required images and finallys runs the containers. If you get permission denied, try with sudo because it might need permission to create the docker image
5. `make restart` (optional) - The rest or worker container might exit because the connection with kafka-broker might timeout. This happens if they try to make a connection before kafka-broker is fully functional. You can restart the container manually once zookeeper and broker is up.

**Note** - There are other intermediate commands as well which makes up the above commands. You can check them out inside the `Makefile`. 

# Endpoints
1. Log User Login Events
```
curl --location 'http://127.0.0.1:8080/api/v1/activity/logins' \
--header 'Content-Type: application/json' \
--data '{
    "userId": "user-4",
    "action": "logout"
}'

{
    "status": "activity logged"
}
```

2. Log Post Interactions
```
curl --location 'http://127.0.0.1:8080/api/v1/activity/posts' \
--header 'Content-Type: application/json' \
--data '{
    "userId": "user-2",
    "postId": "post-1",
    "action": "like"
}'

{
    "status": "activity logged"
}
```

3. Active Users
```
curl --location 'http://127.0.0.1:8080/api/v1/analysis/active-users'

{
    "last_day": 4,
    "last_hour": 4,
    "last_minute": 4
}
```

4. Popular Posts
```
curl --location 'http://127.0.0.1:8080/api/v1/analysis/popular-posts?limit=2'

{
    "last_day": [
        {
            "postId": "post-2",
            "count": 30
        },
        {
            "postId": "post-1",
            "count": 12
        }
    ],
    "last_hour": [
        {
            "postId": "post-2",
            "count": 13
        },
        {
            "postId": "post-1",
            "count": 5
        }
    ],
    "last_minute": []
}
```

# Performance
## Metrics
1. POST "/api/v1/activity/posts"
- average: 14.959152 ms
- p90: 8.107458ms

2. POST "/api/v1/activity/logins"
- average: 4.764479 ms
- p90: 5.110875ms

3. GET "/api/v1/analysis/active-users"
- average: 1.276875 ms
- p90: 3.318625 ms

4. GET "/api/v1/analysis/popular-posts"
- average: 1.345014 ms
- p90: 3.471292 ms

## Optimizations
1. *Data Ingestion Stream*: All posts requests push the activity data to thier respective kafka topics which can be consumed async.
2. DB indexing on timestamps, userId and postId.
3. *Real-Time Queries are optimised using Redis In-Memory Data Structures*:
- Hash Sets for Active Users:
Maintaining a hash set in Redis that gets updated in real-time with user activity.
- Sorted Sets for Popular Posts:
Using sorted sets to keep track of the most popular posts based on likes, comments, and shares.
4. On warmup, the app backfills the cache from the database for persistent data in cache.

## Logs
```
2024-07-29 03:20:24 [GIN] 2024/07/28 - 21:50:24 | 200 |   55.834916ms |    192.168.65.1 | POST     "/api/v1/activity/posts"
2024-07-29 03:20:29 [GIN] 2024/07/28 - 21:50:29 | 200 |    2.627625ms |    192.168.65.1 | POST     "/api/v1/activity/posts"
2024-07-29 03:20:30 [GIN] 2024/07/28 - 21:50:30 | 200 |    3.539792ms |    192.168.65.1 | POST     "/api/v1/activity/posts"
2024-07-29 03:20:33 [GIN] 2024/07/28 - 21:50:33 | 200 |    3.685959ms |    192.168.65.1 | POST     "/api/v1/activity/posts"
2024-07-29 03:20:48 [GIN] 2024/07/28 - 21:50:48 | 200 |    8.107458ms |    192.168.65.1 | POST     "/api/v1/activity/logins"
2024-07-29 03:20:50 [GIN] 2024/07/28 - 21:50:50 | 200 |    4.106208ms |    192.168.65.1 | POST     "/api/v1/activity/logins"
2024-07-29 03:20:51 [GIN] 2024/07/28 - 21:50:51 | 200 |    5.110875ms |    192.168.65.1 | POST     "/api/v1/activity/logins"
2024-07-29 03:20:58 [GIN] 2024/07/28 - 21:50:58 | 200 |    7.522208ms |    192.168.65.1 | POST     "/api/v1/activity/logins"
2024-07-29 03:21:13 [GIN] 2024/07/28 - 21:51:13 | 200 |    3.318625ms |    192.168.65.1 | GET      "/api/v1/analysis/active-users"
2024-07-29 03:21:17 [GIN] 2024/07/28 - 21:51:17 | 200 |     706.958µs |    192.168.65.1 | GET      "/api/v1/analysis/popular-posts?limit=1"
2024-07-29 03:31:04 [GIN] 2024/07/28 - 22:01:04 | 200 |    3.471292ms |    192.168.65.1 | GET      "/api/v1/analysis/popular-posts?limit=2"
2024-07-29 03:40:27 [GIN] 2024/07/28 - 22:10:27 | 200 |      3.0865ms |    192.168.65.1 | GET      "/api/v1/analysis/popular-posts?limit=2"
2024-07-29 03:40:28 [GIN] 2024/07/28 - 22:10:28 | 200 |         369µs |    192.168.65.1 | GET      "/api/v1/analysis/popular-posts?limit=2"
2024-07-29 03:40:30 [GIN] 2024/07/28 - 22:10:30 | 200 |         449µs |    192.168.65.1 | GET      "/api/v1/analysis/popular-posts?key=last_hour&limit=2"
2024-07-29 03:40:31 [GIN] 2024/07/28 - 22:10:31 | 200 |     487.333µs |    192.168.65.1 | GET      "/api/v1/analysis/popular-posts?key=last_hour&limit=2"
2024-07-29 03:40:33 [GIN] 2024/07/28 - 22:10:33 | 200 |     642.084µs |    192.168.65.1 | GET      "/api/v1/analysis/active-users"
2024-07-29 03:40:34 [GIN] 2024/07/28 - 22:10:34 | 200 |     576.333µs |    192.168.65.1 | GET      "/api/v1/analysis/active-users"
2024-07-29 03:40:34 [GIN] 2024/07/28 - 22:10:34 | 200 |     570.458µs |    192.168.65.1 | GET      "/api/v1/analysis/active-users"
```

## Coverage
92.9 %

## Run Tests
1. `make tests` - running tests
2. `make coverage` - command line coverage report
3. `make coverage/report` - GUI (HTML) based report
