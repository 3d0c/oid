package oid

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ObjectId string

type NodesCounter struct {
	Total int32
	next  int32
	sync.Mutex
}

var objectIdCounter uint32 = 0
var seq int32

// Don't forget to intialize it. Default is 0.
var Nodes *NodesCounter

func init() {
	if Nodes == nil {
		Nodes = &NodesCounter{}
	}
}

func NewObjectId(args ...int) ObjectId {
	var b [12]byte
	var flag int = 0

	// NodeId, sequential distribution.
	var nextId int32 = Nodes.Next()

	if len(args) > 0 {
		if args[0] <= 255 {
			flag = args[0]
		}

		if len(args) > 1 {
			nextId = int32(args[1])
		}
	}

	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))

	b[4] = byte(nextId >> 8)
	b[5] = byte(nextId)

	// Flags 0-255
	b[6] = byte(int16(flag))

	// 2 bytes Reserved
	b[7] = byte(0 >> 8)
	b[8] = byte(0)

	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	if i >= 16777215 {
		objectIdCounter = 0
		i = 0
	}

	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)

	return ObjectId(b[:])
}

func (id ObjectId) String() string {
	return fmt.Sprintf("%x", string(id))
}

func (id ObjectId) NodeId() int32 {
	i := (id[4] << 8) | id[5]
	return int32(i)
}

func (id ObjectId) Flag() int {
	return int(id[6])
}

func (id ObjectId) Hex() string {
	return hex.EncodeToString([]byte(id))
}

func ObjectIdHex(s string) ObjectId {
	d, err := hex.DecodeString(s)
	if err != nil || len(d) != 12 {
		panic(fmt.Sprintf("Invalid input to ObjectIdHex: %q", s))
	}
	return ObjectId(d)
}

func IsObjectIdHex(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func (this *NodesCounter) Next() int32 {
	if this.Total == 0 {
		return 0
	}

	retval := this.next

	this.Lock()

	if this.next == this.Total-1 {
		this.next = 0
	} else {
		this.next += 1
	}

	this.Unlock()

	return retval
}
