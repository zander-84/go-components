package CHelperType

import (
	"math/rand"
	"time"
)

type Slice struct{}

func NewSlice() interface{} { return new(Slice) }

func (*Slice) Shuffle(slice interface{}) {
	rand.Seed(time.Now().UnixNano())

	switch data := slice.(type) {
	case []string:
		for i := 0; i < len(data); i++ {
			a := rand.Intn(len(data))
			b := rand.Intn(len(data))
			data[a], data[b] = data[b], data[a]
		}
	case []int:
		for i := 0; i < len(data); i++ {
			a := rand.Intn(len(data))
			b := rand.Intn(len(data))
			data[a], data[b] = data[b], data[a]
		}
	}
}
