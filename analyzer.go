package crawler

import (
	"errors"
	"fmt"
	"gocrawler/base"
	"gocrawler/middleware"
	"reflect"
)

type parseResponse func(res base.Response) ([]base.Data, []error)

type Analyzer interface {
	Id() uint32
	Analyze(parser []parseResponse, res base.Response) ([]base.Data, []error)
}

type AnalyzerPool interface {
	Take() (Analyzer, error)
	Return(Analyzer) error
	Total() uint32
	Used() uint32
}

type myAnalyzer struct {
	id uint32
}

var analyzerIdGenerator middleware.IdGenerator = middleware.NewIdGenerator()

func genAnalyzerId() uint32 {
	return analyzerIdGenerator.GetUint32()
}

func NewAnalyzer() (Analyzer, error) {
	id := genAnalyzerId()
	return &myAnalyzer{id: id}, nil
}

func (m *myAnalyzer) Id() uint32 {
	return m.id
}

func appendDataList(dataList []base.Data, data base.Data, depth uint32) []base.Data {
	if data == nil {
		return dataList
	}
	req, ok := data.(*base.Request)
	if !ok {
		return append(dataList, data)
	}
	newDepth := depth + 1
	if req.Depth() != newDepth {
		req = base.NewRequest(req.Get(), newDepth)
	}
	return append(dataList, req)
}

func appendErrorList(errorList []error, err error) []error {
	if err == nil {
		return errorList
	}
	return append(errorList, err)
}

func (m *myAnalyzer) Analyze(parser []parseResponse, res base.Response) ([]base.Data, []error) {
	if parser == nil {
		errMsg := "The response parser is nil!"
		return nil, []error{errors.New(errMsg)}
	}
	if !res.Valid() {
		errMsg := "The response is valid!"
		return nil, []error{errors.New(errMsg)}
	}
	result := make([]base.Data, 0)
	errResult := make([]error, 0)
	for i, p := range parser {
		if p == nil {
			err := errors.New(fmt.Sprintf("The response parser is nil! Index: %d", i))
			errResult = append(errResult, err)
			continue
		}
		datas, errs := p(res)
		for _, data := range datas {
			result = appendDataList(result, data, res.Depth())
		}
		//errResult = append(errResult, errs)
		for _, err := range errs {
			errResult = appendErrorList(errResult, err)
		}
	}
	return result, errResult
}

type myAnalyzerPool struct {
	pool  middleware.Pool
	etype reflect.Type
}

type GenAnalyzer func() Analyzer

func NewAnalyzerPool(
	total uint32,
	gen GenAnalyzer,
) (AnalyzerPool, error) {
	genEntity := func() middleware.Entity { return gen() }
	etype := reflect.TypeOf(gen())
	pool, err := middleware.NewPool(total, etype, genEntity)
	if err != nil {
		return nil, err
	}
	myPool := &myAnalyzerPool{
		pool:  pool,
		etype: etype,
	}
	return myPool, nil
}

func (m *myAnalyzerPool) Take() (Analyzer, error) {
	ana, err := m.pool.Take()
	if err != nil {
		return nil, err
	}
	an, ok := ana.(Analyzer)
	if !ok {
		errMsg := "Entity type doesn't match!"
		panic(errors.New(errMsg))
	}
	return an, nil
}

func (m *myAnalyzerPool) Return(a Analyzer) error {
	return m.pool.Return(a)
}

func (m *myAnalyzerPool) Total() uint32 {
	return m.pool.Total()
}

func (m *myAnalyzerPool) Used() uint32 {
	return m.pool.Used()
}
