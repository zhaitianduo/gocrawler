package middleware

import (
	"math"
	"sync"
)

type myIdGenerator struct {
	sn    uint32
	ended bool
	mutex sync.Mutex
}

func NewIdGenerator() IdGenerator {
	return &myIdGenerator{}
}

func (m *myIdGenerator) GetUint32() uint32 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.ended {
		defer func() { m.ended = false }()
		m.sn = 0
		return uint32(0)
	}
	if m.sn < math.MaxUint32 {
		m.sn++
	} else {
		m.ended = true
	}
	return m.sn
}
