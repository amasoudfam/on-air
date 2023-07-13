# Go Redis Package

The Go Redis package is a client library for Redis, a popular in-memory data structure store. This package provides a convenient and efficient way to interact with Redis servers using the Go programming language.

## Installation

To install the Go Redis package, use the following command:

```bash
go get -u github.com/go-redis/redis/v9
```

### Usage

To use the Go Redis package in your Go project, import it as follows:

go
Copy code

```go
import "github.com/go-redis/redis/v9"
```

### Features

The Go Redis package offers the following features:

* Connection management: Establish connections with Redis servers and manage connection pools.
* Data operations: Perform various operations on Redis data structures, such as strings, lists, sets, hashes, sorted sets, and more.
* Pipelining: Optimize performance by sending multiple commands to Redis in a single roundtrip.
* Transactions: Execute atomic operations on Redis using transactions.
* Pub/Sub: Publish and subscribe to channels for real-time messaging.
* Lua scripting: Execute Lua scripts on Redis servers.
* Cluster support: Connect to Redis clusters and perform operations across multiple nodes.
* Connection options: Configure connection parameters, such as timeouts and authentication.
* Monitoring and debugging: Observe Redis commands and responses for debugging and performance analysis.

## Why we use the Go Redis package

We use the Go Redis package in this project for efficient communication and interaction with Redis servers. Redis is widely used for caching, queuing, and data storage purposes, and the Go Redis package provides a reliable and feature-rich client library to interface with Redis in our Go applications.

By using the Go Redis package, we can easily connect to Redis servers, perform various data operations, and leverage advanced functionalities like pub/sub and Lua scripting. The package also offers performance optimizations such as pipelining and transaction support, which can significantly improve the efficiency of our Redis interactions.

Overall, the Go Redis package simplifies the integration of Redis into our Go project, allowing us to harness the power and flexibility of Redis in a seamless manner.

## Examples

* Connecting to a Redis Server

```go
    import (
        "context"
        "fmt"
        "github.com/go-redis/redis/v8"
    )

func main() {
 // Create a new Redis client
 client := redis.NewClient(&redis.Options{
  Addr:     "localhost:6379", // Redis server address
  Password: "",               // Redis server password (if any)
  DB:       0,                // Redis database index
 })

 // Ping the Redis server to check the connection
 pong, err := client.Ping(context.Background()).Result()
 if err != nil {
  fmt.Println("Failed to connect to Redis:", err)
  return
 }

 fmt.Println("Connected to Redis:", pong)
}
```

* Setting and Retrieving Values

```go
import (
 "context"
 "fmt"
 "github.com/go-redis/redis/v8"
)

func main() {
 // Create a new Redis client
 client := redis.NewClient(&redis.Options{
  Addr:     "localhost:6379",
  Password: "",
  DB:       0,
 })

 // Set a key-value pair
 err := client.Set(context.Background(), "mykey", "myvalue", 0).Err()
 if err != nil {
  fmt.Println("Failed to set value:", err)
  return
 }

 // Retrieve the value for a key
 value, err := client.Get(context.Background(), "mykey").Result()
 if err != nil {
  fmt.Println("Failed to get value:", err)
  return
 }

 fmt.Println("Retrieved value:", value)
}
```

* Working with Lists

```go
import (
 "context"
 "fmt"
 "github.com/go-redis/redis/v8"
)

func main() {
 // Create a new Redis client
 client := redis.NewClient(&redis.Options{
  Addr:     "localhost:6379",
  Password: "",
  DB:       0,
 })

 // Push values to a list
 err := client.RPush(context.Background(), "mylist", "value1", "value2", "value3").Err()
 if err != nil {
  fmt.Println("Failed to push values to list:", err)
  return
 }

 // Retrieve all values from the list
 values, err := client.LRange(context.Background(), "mylist", 0, -1).Result()
 if err != nil {
  fmt.Println("Failed to retrieve list values:", err)
  return
 }

 fmt.Println("List values:", values)
}
```

## What we do with this package

We utilize this package to cache flights when we receive the first flight response from the FlightsApiMock. The package allows us to store this information in our local database using the go-redis library. This library is an open-source project that facilitates caching and data storage in memory or on disk via the REDIS protocol.

After caching the flights, we retrieve them from the Redis cache until the cache timeout period expires. If the flights are not found in the Redis cache, we initiate a new request to retrieve flights from the FlightsApiMock and store them in the cache again.
