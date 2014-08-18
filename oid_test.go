package oid

import (
	"log"
	"testing"
	"time"
)

func TestDistribution(t *testing.T) {
	Nodes = &NodesCounter{Total: 2}

	for i := 0; i < 10; i++ {
		oid := NewObjectId()
		if i%2 == 0 {
			if oid.NodeId() != 0 {
				t.Fatalf("[%d] Expecting nodeId = 0, %v got.", i, oid.NodeId())
			}
		}
	}
}

func TestFlags(t *testing.T) {
	Nodes = &NodesCounter{Total: 5}

	cases := map[int]int{0: 0, 1: 1, 2: 2, 255: 255, 256: 0}
	for key, val := range cases {
		oid := NewObjectId(key)
		if oid.Flag() != val {
			t.Fatalf("Expected flag = %v, %v got.", val, oid.Flag())
		}
	}
}

func TestUniq(t *testing.T) {
	Nodes = &NodesCounter{Total: 5}
	result := make([]string, 0)

	for i := 0; i < 100; i++ {
		go func() { result = append(result, NewObjectId().String()) }()
	}

	time.Sleep(time.Duration(time.Second)) // just a hack

	d_count := map[string]int{}

	for _, val := range result {
		if _, found := d_count[val]; found {
			d_count[val] += 1
		} else {
			d_count[val] = 1
		}
	}

	for key, val := range d_count {
		if val > 1 {
			t.Errorf("%v generated %d times\n", key, val)
		}
	}
}

func TestGivenId(t *testing.T) {
	Nodes = &NodesCounter{Total: 2}

	oid := NewObjectId(0, 13)

	if oid.NodeId() != 13 {
		t.Fatalf("Expected nodeId = 13, %v got.", oid.NodeId())
	}
}

func TestHex(t *testing.T) { //
	Nodes = &NodesCounter{Total: 5}

	oid := NewObjectId(12)

	if !IsObjectIdHex(oid.String()) {
		t.Fatalf("Expecting valid hex representaion of objectid, %v got.\n", oid.String())
	}

	if ObjectIdHex(oid.String()) != oid {
		t.Fatalf("Expecting %v, %v got\n", oid.String(), ObjectIdHex(oid.String()).String())
	}
}

func BenchmarkOids(b *testing.B) {
	Nodes = &NodesCounter{Total: 5}

	for i := 0; i < b.N; i++ {
		NewObjectId()
	}
}

// 51e406db582b944e51000027,51e406db582b944e5100002b,51e406db582b944e5100002f,51e406db582b944e51000033,51e406db582b944e51000037,51e406db582b944e5100003b
func TestSome(t *testing.T) {
	ids := []string{"51e406db582b944e51000027", "51e406db582b944e5100002b", "51e406db582b944e5100002f", "51e406db582b944e51000033", "51e406db582b944e51000037", "51e406db582b944e5100003b"}

	for _, id := range ids {
		log.Println(id, ObjectIdHex(id).NodeId())
	}
}
