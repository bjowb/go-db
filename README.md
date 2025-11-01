# go-db: Redis-Inspired In-Memory Key-Value Store

**`go-db`** is a **lightweight, Redis-inspired in-memory key-value store** written entirely in Go.

It supports basic string storage with `SET`/`GET` and hash maps with `HSET`/`HGET`/`HGETALL`. All operations are **concurrency-safe** using Go’s `sync.RWMutex`.

> **Project Goal:** This project was built to gain a deeper understanding of how systems like Redis store data, handle client requests, and utilize mutexes for thread-safe in-memory storage.

---

## 🚀 Quick Start

### What is go-db?

`go-db` is a **lightweight, Redis-inspired, in-memory key-value store** written in Go. It supports basic key-value (`SET`/`GET`) and hash map (`HSET`/`HGET`) commands, all designed to be concurrency-safe.

### ⚙️ Setup and Running

1.  **Clone the Repository**

    Fetch the project code to your local machine:

    ```bash
    git clone [https://github.com/bjowb/go-db.git](https://github.com/bjowb/go-db.git)
    cd go-db
    ```

2.  **Run the Server**

    Run the server directly using the Go runtime:

    ```bash
    go run .
    ```

3.  **Connect to the Database**

    Connect to the server (e.g., using `redis-client` if the server is listening on port 6379):

    ```bash
    PING
    ```

    You can now send commands like `GET key` or `SET key value`.
