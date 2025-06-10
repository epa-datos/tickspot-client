# Tickspot Client

A simple Golang client for interacting with the [Tickspot API](https://github.com/tick/tick-api).

## Features

- Easy authentication with Tickspot API
- Fetch and manage entries

## Installation

```bash
go get github.com/epa-datos/tickspot-client
```

## Usage

```go
import "github.com/epa-datos/tickspot-client"

client := NewTickspotClient("MyProject", os.Getenv("MyToken"), "MyAppAgent")

tasks, err := client.GetTasks("userID")

```

## API

- `GetUsers`
- `GetTasks`
- `UploadTask`
- `DeleteTask`


## License

MIT