package middleware

import (
	"sync"
)

type myStopSign struct {
	signed       bool
	dealCountMap map[string]uint32
	rwmutex      sync.RWMutex
}

func NewStopSign() StopSign {
	ss := &myStopSign{
		dealCountMap: make(map[string]uint32),
	}
	return ss
}

func (m *myStopSign) Sign() bool {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	if m.signed {
		return false
	}
	m.signed = true
	return true
}

func (m *myStopSign) Signed() bool {
	return m.signed
}

func (m *myStopSign) Deal(code string) {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	if !m.signed {
		return
	}
	if _, ok := m.dealCountMap[code]; !ok {
		m.dealCountMap[code] = 1
	} else {
		m.dealCountMap[code] += 1
	}
}

//reset sign and counter
func (m *myStopSign) Reset() {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	m.signed = false
	m.dealCountMap = make(map[string]uint32)
}

func (m *myStopSign) DealCount(code string) uint32 {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if v, ok := m.dealCountMap[code]; !ok {
		return 0
	} else {
		return v
	}
}

func (m *myStopSign) DealTotal() uint32 {
	m.rwmutex.RLock()
	m.rwmutex.RUnlock()
	result := uint32(0)
	for _, v := range m.dealCountMap {
		result += v
	}
	return result
}

func (m *myStopSign) Summary() string {
	return ""
}
