#### ObjectId Package.

There is 12 bytes:  

- 4-byte - Unix timestamp.
- 2-byte - Nodeid. (safely generated sequence)
- 1-byte flag  ```0 - clean, 1 - encyprt, 2-255 - whatever you want```
- 2-byte - Reserved
- 3-byte - Counter

#### Usage.
Init node counter.

```go
func init() {
	oid.Nodes = &oid.NodesCounter{Total: 2}
}
```

Create new object id:

```go
id := oid.NewObjectId()
```

Create from hex:

```go
id := oid.ObjectIdHex("52e0eba500000c0000000074")
```

Get Node Id and Flag:

```go
id.NodeId()
id.Flag()
```

