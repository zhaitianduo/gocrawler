package crawler

type IdGenerator interface {
	GetUint32() uint32
}

type ChannelManagerStatus uint8

const (
	CHANNEL_MANAGER_STATUS_UNINITIALIZED ChannelManagerStatus = 0
	CHANNEL_MANAGER_STATUS_INITIALIZED   ChannelManagerStatus = 1
	CHANNEL_MANAGER_STATUS_CLOSED        ChannelManagerStatus = 2
)

type ChannelManager interface {
	//channelLen represents the initial length of channels
	//reset means whether or not to reset the channel manager
	Init(channelLen uint32, reset bool) bool
	//close the channel manager
	Close() bool
	//get the request channel
	ReqChan() (chan *Request, error)
	//get the response channel
	ResChan() (chan *Response, error)
	//get the item channel
	ItemChan() (chan *Item, error)
	//get the error channel
	ErrorChan(chan error, error)
	//get the length of channel
	ChannelLen() uint32
	//get the status of channel manager
	Status() ChannelManagerStatus
	//get the summary info
	Summary() string
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
