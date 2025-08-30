
# KeyVal

**KeyVal** is a lightweight, distributed key-value store implemented in Go. It leverages gRPC for network communication, LevelDB as local storage, and supports sharding and a simple routing mechanism to distribute data across nodes.

---

## Table of Contents

- [Features](#features)  
- [Architecture Overview](#architecture-overview)  
- [Prerequisites](#prerequisites)  
- [Installation & Setup](#installation--setup)  
- [Usage](#usage)  
  - [Running the Server](#running-the-server)  
  - [Client Operations](#client-operations)  
- [Testing](#testing)  
- [Project Structure](#project-structure)  
- [Why I’m Proud of This Project](#why-im-proud-of-this-project)  
- [Contributing](#contributing)  
- [License](#license)

---

## Features

- Distributed key-value storage with optional **sharding** for horizontal scaling.  
- **gRPC**-based API ensuring fast and type-safe communication between clients and servers.  
- **LevelDB** integration for reliable, high-performance local persistence.  
- Supports basic **PUT** and **GET** operations, making it ideal for learning distributed storage foundations.

---

## Architecture Overview

The project comprises separate **client**, **router**, and **server** components:

- **Client**: Sends `PUT`/`GET` requests to the router via gRPC.  
- **Router**: Determines the appropriate shard/server for each key, forwarding requests accordingly.  
- **Server**: Receives operations via gRPC, interacts with LevelDB (under `dbs/`), and responds to clients.

This separation allows clear modularity—enhancing maintainability and extensibility.

---

## Prerequisites

Ensure you have the following installed:

- Go 1.19 or higher  
- Protocol Buffers compiler (`protoc`)  
- gRPC Go plugin (`protoc-gen-go` and `protoc-gen-go-grpc`)  
- LevelDB 

---

## Installation & Setup

```bash
# Clone the repository
git clone https://github.com/1byinf8/KeyVal.git
cd KeyVal

# Generate protobuf code
protoc --go_out=. --go-grpc_out=. proto/badies.proto

# Optionally, build binaries
go build ./router
go build ./server
go build ./client
````

---

## Usage

### Running the Server

Start a node on a given port:

```bash
go run server/main.go --port=5001
```

### Running the Router

Start the router with information about available servers (e.g., ports):

```bash
go run router/main.go --servers=localhost:5001,localhost:5002
```

### Client Operations

Use client commands to interact with KeyVal:

* **PUT operation**:

  ```bash
  go run client/main.go put <key> <value>
  ```

* **GET operation**:

  ```bash
  go run client/main.go get <key>
  ```

---

## Testing

Automated tests are available to verify core functionality:

```bash
go test ./test_hash.go
go test ./test_put.go
go test ./test_get.go
```

These ensure consistent sharding logic, key insertion, and retrieval functionality.

---

## Project Structure

```
.
├── client/               # Client-side code for issuing requests
├── router/               # Routing logic for request forwarding
├── server/               # Server-side logic with LevelDB persistence
├── proto/
│   └── badies.proto      # Protocol Buffers definitions for gRPC interfaces
├── dbs/                  # Directory for LevelDB storage files
├── test_hash.go          # Tests for hashing/rebalancing logic
├── test_put.go           # Tests for PUT operations
├── test_get.go           # Tests for GET operations
├── instance_creator.go   # Utility for initializing instances (if applicable)
├── go.mod                # Module definition
└── go.sum                # Dependency locks
```

---

## Why I’m Proud of This Project

KeyVal is a hands-on exploration of distributed systems—integrating sharding, gRPC, persistence, and modular architecture. Implementing a router-based design required mastering data partitioning logic, client-server communication, and fault-tolerant data access patterns. It’s a solid foundation for building scalable, real-world systems.

---

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests. Ideas for future improvements:

* Add **replication** for redundancy and high availability
* Implement a **consensus algorithm** (e.g., Raft) for strong consistency
* Build a **monitoring dashboard** for metrics and performance insights

---

## License

This project is released under the **[MIT License](LICENSE)**. Feel free to use, modify, and distribute it freely.

---

