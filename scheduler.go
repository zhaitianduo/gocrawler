package crawler

import (
	"gocrawler/base"
	"net/http"
)

type GenHttpClient func() *http.Client

type SchedSummary interface {
	String() string
	Detail() string
	Same(other *SchedSummary) bool
}

type Scheduler interface {
	//start scheduler
	//create and init scheduler and each component, after that, scheduler will activate crawling process
	//channelLen will be used to initialize the length of data transmission channel
	//poolSize will be used to initialize the size of downloader pool and analyzer pool
	//crawDepth, the page which depth is larger than this number will be ignored
	//httpClientGenerator represents the func to generate http client
	//resParsers is a slice of func to parse the http response
	//itemProcessors is a slice of func to process the items parsed from the http response
	//firstHttpRequest means the entrance of the crawling process
	Start(channelLen uint32,
		poolSize uint32,
		crawlDepth uint32,
		httpClientGenerator GenHttpClient,
		resParsers []parseResponse,
		itemProcessors []base.ProcessItem,
		firstHttpRequest base.Request,
	) error
	// stop the crawling process and return if the stop process succeed
	Stop() bool
	//whether the scheduler is running
	Running() bool
	//get the error chan which stores the error in each component
	ErrorChan() <-chan error
	//justify whether each component is idle
	Idle() bool
	//get summary info
	Summary(prefix string) SchedSummary
}
