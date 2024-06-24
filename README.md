[![progress-banner](https://backend.codecrafters.io/progress/redis/6bd37c05-2efe-46a1-a075-e0edbd0f571d)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is my solution for the
[Build Your Own Redis Challenge](https://codecrafters.io/challenges/redis) by [Codecrafters](https://codecrafters.io/)!

In this challenge, I build a toy Redis clone that's capable of handling
basic commands like `PING`, `SET` and `GET`. We then subsequently expand the capacilities of the same via extensions like:
- Replication
- Persistence and RDB formats
- Streaming

All of the implementations are based on the actual Redis protocols and thus are 100% compliant with the same. 

To find a detailed descirption of the capacilities of this Redis implementation, head over to the [Codecrafters' course](https://app.codecrafters.io/courses/redis/overview)!


## Running Locally

You can spawn the Redis server instances by using the [`spawn_redis_server.sh`](./spawn_redis_server.sh) shell script with the appropiate arguments. To test the same, you can use the client, and start the same by using `go run client/main.go --port PORT` where `PORT` is the port of the local Redis server. It defaults to `6379`.