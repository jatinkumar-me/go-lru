# LRU Cache Implementation in Go

This is a simple implementation of an LRU (Least Recently Used) cache in Go.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your machine
- [Nodejs](https://nodejs.org/en/download) installed on your machine.

## Running project locally

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/jatinkumar-me/go-lru.git
    ```

   Or download the repository as a ZIP file and extract it.

2. Change directory to the project folder:

    ```bash
    cd go-lru/server/
    go run .
    # This will start the server on http://localhost:8888 by default.
    # To run the react app
    cd go-lru/client/
    npm run dev
    # This will start the dev server on http://localhost:5173/
    ```

3. To Build the project:

    ```bash
    go build
    ```

## Testing

To run the unit tests, use the following command:

```bash
go test ./lru/
```
