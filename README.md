
# Distributed Cache with Pub/Sub Functionality

This is an implementation of a distributed cache system with Pub/Sub functionality using Go and Redis. It includes basic cache operations like setting, getting, and deleting cache entries, along with the ability to publish and subscribe to cache updates. The project is designed to provide a simple yet powerful caching solution that can be used in various distributed systems.

## Installation Instructions

To set up the project locally, follow these steps:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Ashwinr-07/Golang_DistributedCache_PubSub_Redis.git
   ```

2. **Navigate to the project directory:**

   ```bash
   cd Golang_DistributedCache_PubSub_Redis
   ```

3. **Install the required dependencies:**

   ```bash
   go get github.com/go-redis/redis/v8
   go get github.com/alicebob/miniredis/v2
   ```

## Usage

### Running the project:

1. **Ensure Redis is running locally on the default port (6379).**

2. **Start the project by running the main Go file:**

   ```bash
   go run main.go
   ```

### Example usage:

The `main.go` file provides an example of how to use the cache system. It demonstrates creating a cache, setting and getting values, and handling Pub/Sub messages.

## Features

- **Basic Cache Operations:** Set, Get, and Delete methods for managing cache entries.
- **Pub/Sub Functionality:** Updates to the cache are published to subscribers, ensuring consistency across distributed systems.
- **Unit Tests:** Basic tests for cache operations, including Set, Get, Delete, and distributed cache testing.

## License

This is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the creators of [go-redis](https://github.com/go-redis/redis) and [miniredis](https://github.com/alicebob/miniredis) for their excellent libraries.

