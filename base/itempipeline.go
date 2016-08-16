package base

type ItemPipeline interface {
	//send item
	Send(item *Item) []error
	//FailFast means when failed to process one item, whether or not ignore the coming items
	FailFast() bool
	//set failfast
	SetFailFast(failFast bool)
	//get sent, received and processed items count
	Count() (uint32, uint32, uint32)
	//get the count of processing items
	ProcessingNumber() uint32
	//get summary info
	Summary() string
}

type ProcessItem func(item *Item) (*Item, error)
