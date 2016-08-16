package middleware

type IdGenerator interface {
	GetUint32() uint32
}

type Entity interface {
	Id() uint32
}

type Pool interface {
	Take() (*Entity, error)
	Return(entity *Entity) error
	Total() uint32
	Used() uint32
}

type StopSign interface {
	// set the stop sign, if the stop sign has been signed, return false
	Sign() bool
	//get the info whether the stop sign has been sent
	Signed() bool
	//reset the stop sign
	Reset()
	//code represents the code of stop sign processor
	Deal(code string) uint32
	//get the stop sign processed count from one stop sign processor
	DealCount(code string) uint32
	//get the total count of stop sign processed by all stop sign processors
	DealTotal() uint32
	//get the summary info
	Summary() string
}
