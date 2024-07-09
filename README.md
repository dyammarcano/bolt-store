# Bolt Store

A Go package for persisting data in a BoltDB database.

## Installation

To install the package, use `go get`:

```sh
go get github.com/dyammarcano/bolt-store
```

## Usage

To use the package, import it in your code:

```go
import (
    store "github.com/dyammarcano/bolt-store"
)

func main() {
    if err := store.NewBoltStore("example.db"); err != nil {
        panic(err)
    }
    
    if err := store.RegisterBucket("example"); err != nil {
        panic(err)
    }

    store.SetString("example", fmt.Sprintf("%x", sha256.Sum256([]byte("testing"))))
    store.SetString("example", fmt.Sprintf("%x", sha256.Sum256([]byte("testing"))))
    store.SetString("example", fmt.Sprintf("%x", sha256.Sum256([]byte("testing"))))
    
    fmt.Println(store.Total("example"))
    
    store.SetString("example", fmt.Sprintf("%x", sha256.Sum256([]byte("testing"))))
    
    fmt.Println(store.Total("example"))
}
```

# Contributing
Feel free to open issues or submit pull requests with improvements.

# License
This project is licensed under the MIT License. See the LICENSE file for details.

# Acknowledgements
This project uses the following libraries:

- [github.com/segmentio/ksuid](github.com/segmentio/ksuid)
- [go.etcd.io/bbolt](go.etcd.io/bbolt)
