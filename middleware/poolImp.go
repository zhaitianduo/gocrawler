package middleware

import (
	"errors"
	"reflect"
	"sync"
)

type myPool struct {
	total       uint32
	etype       reflect.Type
	genEntity   func() Entity
	container   chan Entity
	idContainer map[uint32]bool
	mutex       sync.Mutex
}

func NewPool(
	total uint32,
	entityType reflect.Type,
	genEntity func() Entity,
) (Pool, error) {
	if total == 0 {
		errMsg := "Pool size can not be 0!"
		return nil, errors.New(errMsg)
	}
	size := int(total)
	container := make(chan Entity, size)
	idContainer := make(map[uint32]bool)
	for i := 0; i < size; i++ {
		newEntity := genEntity()
		if entityType != reflect.TypeOf(newEntity) {
			errMsg := "Entity doesn't match"
			return nil, errors.New(errMsg)
		}
		container <- newEntity
		idContainer[newEntity.Id()] = true
	}

	pool := &myPool{
		total:       total,
		etype:       entityType,
		genEntity:   genEntity,
		container:   container,
		idContainer: idContainer,
	}
	return pool, nil
}

//what if the container is empty?
func (m *myPool) Take() (Entity, error) {
	entity, ok := <-m.container
	if !ok {
		errMsg := "The container has been closed!"
		return nil, errors.New(errMsg)
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.idContainer[entity.Id()] = false
	return entity, nil
}

//return -1: id doesn't exist
//return 0: fail to process
//return 1: success
//narrow down the scope of critical section
func (m *myPool) compareAndSetForIdContainer(entityId uint32, oldValue bool, newValue bool) int8 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	id, ok := m.idContainer[entityId]
	if !ok {
		//errMsg := "Return entity doesn't belong to pool!"
		return -1
	}

	if id {
		//errMsg := "Return entity is already in the pool!"
		return 0
	}
	m.idContainer[entityId] = newValue
	return 1

}

func (m *myPool) Return(entity Entity) error {
	if entity == nil {
		errMsg := "Return entity should not be nil!"
		return errors.New(errMsg)
	}
	if reflect.TypeOf(entity) != m.etype {
		errMsg := "Return entity type doesn't match with pool entity type"
		return errors.New(errMsg)
	}

	s := m.compareAndSetForIdContainer(entity.Id(), false, true)
	if s == 1 {
		m.container <- entity
		return nil
	} else if s == -1 {
		errMsg := "Return entity doesn't belong to pool!"
		return errors.New(errMsg)
	} else {
		errMsg := "Return entity is already in the pool!"
		return errors.New(errMsg)
	}

}

func (m *myPool) Total() uint32 {
	return m.total
}

func (m *myPool) Used() uint32 {
	return m.total - uint32(len(m.container))
}
