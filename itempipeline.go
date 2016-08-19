package crawler

import (
	"errors"
	"fmt"
	"gocrawler/base"
	"sync/atomic"
)

type ItemPipeline interface {
	//send item
	Send(item base.Item) []error
	//FailFast means when failed to process one item, whether or not ignore the coming items
	FailFast() bool
	//set failfast
	SetFailFast(failFast bool)
	//get sent, received and processed items count
	Count() (uint64, uint64, uint64)
	//get the count of processing items
	ProcessingNumber() uint64
	//get summary info
	Summary() string
}

type ProcessItem func(item base.Item) (base.Item, error)

type myItemPipeline struct {
	itemProcessors []ProcessItem
	//identify whether the process should fail immediately
	failFast bool
	//how many items been sent
	sent uint64
	//how many items accepeted
	accepted uint64
	//how many items processed
	processed uint64
	//how many items beeing processed
	processing uint64
}

func NewItemPipeline(itemProcessors []ProcessItem) ItemPipeline {
	if itemProcessors == nil {
		panic(errors.New("ProcessItem list should not be nil!"))
	}
	innerItemProcessors := make([]ProcessItem, 0, len(itemProcessors))
	for i, itemProcessor := range itemProcessors {
		if itemProcessor == nil {
			errMsg := fmt.Sprintf("itemProcessor should not be null! Index: %d", i)
			panic(errors.New(errMsg))
		}
		innerItemProcessors = append(innerItemProcessors, itemProcessor)
	}

	return &myItemPipeline{
		itemProcessors: innerItemProcessors,
	}
}

//Whether there is a need to set failFast? Normally, the ItemProcess func will return nil,err if err is not nil
func (m *myItemPipeline) Send(item base.Item) []error {
	atomic.AddUint64(&m.sent, 1)
	//defer atomic.AddUint64(&addr, ^uint64(0))
	errs := make([]error, 0)
	if !item.Valid() {
		//TODO, not aware which item is valid
		err := errors.New("item is not valid!")
		return append(errs, err)
	}
	atomic.AddUint64(&m.accepted, 1)
	currentItem := item
	//var processedItem base.Item
	for _, processor := range m.itemProcessors {
		atomic.AddUint64(&m.processing, 1)
		//defer func is executed after this function end, aka after return in this function. reasonable?
		defer atomic.AddUint64(&m.processing, ^uint64(0))
		processedItem, err := processor(currentItem)
		defer atomic.AddUint64(&m.processed, 1)
		if err != nil {
			errs = append(errs, err)
			if m.failFast {
				return errs
			} else {
				continue
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		} else {
			errMsg := "The processed item is nil, could not continue!"
			return append(errs, errors.New(errMsg))
		}
	}
	return errs
}

func (m *myItemPipeline) FailFast() bool {
	return m.failFast
}

func (m *myItemPipeline) SetFailFast(failFast bool) {
	m.failFast = failFast
}

//get sent, received and processed items count
func (m *myItemPipeline) Count() (uint64, uint64, uint64) {
	sent := atomic.LoadUint64(&m.sent)
	accepted := atomic.LoadUint64(&m.accepted)
	processed := atomic.LoadUint64(&m.processed)
	return sent, accepted, processed
}

func (m *myItemPipeline) ProcessingNumber() uint64 {
	processing := atomic.LoadUint64(&m.processing)
	return processing
}
func (m *myItemPipeline) Summary() string {
	summaryTemplate := "failFast: %v," +
		"processorNumber: %d, sent: %d, accepted: %d, processed: %d, processingNumber: %d"
	return fmt.Sprintf(summaryTemplate, m.failFast, len(m.itemProcessors), m.sent, m.accepted, m.processed, m.ProcessingNumber())
}
