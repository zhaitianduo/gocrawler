package crawler

type parseResponse func(res *Response) ([]Data, []error)

type Analyzer interface {
	Id() uint32
	Analyze(parser []parseResponse, res *Response) ([]Data, []error)
}

type AnalyzerPool interface {
	Take() (Analyzer, error)
	Return(*Analyzer) error
	Total() uint32
	Used() uint32
}
