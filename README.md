# Toy Redis

This is my implementation of a toy Redis clone that's capable of handling basic commands like `PING`, `SET` and `GET`. Later I expanded the capabilities of the same via extensions like:
- Replication
- Persistence and RDB formats
- Streaming

All implementations are based on the actual Redis protocols and, thus, are 100% compliant. The code is written in Go and has no external dependencies. The code, including the binary serialization and deserialization, is written from scratch! I have also created a simple client to test the same in the [client](./client/) directory to ensure the server works as expected.

To find a detailed description of the capabilities of this Redis implementation, head over to the [Codecrafters' challenge descriptions](https://app.codecrafters.io/courses/redis/overview), which greatly helped me in defining meaningful milestones for the project. Huge shoutout to them for creating such an amazing platform!


## Running Locally

You can spawn the Redis server instances using the [`spawn_redis_server.sh`](./spawn_redis_server.sh) shell script with the appropriate arguments. To test the same, you can use the client and start the same by using `go run client/main.go --port PORT` where `PORT` is the port of the local Redis server. It defaults to `6379`.
