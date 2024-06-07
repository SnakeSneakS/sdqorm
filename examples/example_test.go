package examples

import (
	"fmt"
	"testing"

	"github.com/snakesneaks/sdqorm/pkg/parser"
)

func TestQueryTypeA(t *testing.T) {
	callOnlyParser := parser.NewSDQParser()
	var v QueryTypeA
	type testSet struct {
		Query  string
		Answer QueryTypeA
	}

	testSuccessCases := make([]testSet, 0)
	testFailCases := make([]string, 0)

	testSuccessCases = append(testSuccessCases, testSet{
		Query: "1 test1",
		Answer: QueryTypeA{
			ID:   1,
			Name: "test1",
		},
	})

	testFailCases = append(testFailCases,
		//	"-1 test2", //validationは今はやめた
		"1.23 test2",
		//"1 test-case", //validationは今はやめた
	)

	for _, successCase := range testSuccessCases {
		if err := callOnlyParser.Parse(successCase.Query, &v); err != nil {
			t.Error(err)
			continue
		}
		if v != successCase.Answer {
			t.Error(fmt.Errorf("value mismatch: %+v vs %+v", v, successCase.Answer))
			continue
		}
	}

	for _, failCase := range testFailCases {
		if err := callOnlyParser.Parse(failCase, &v); err == nil {
			t.Error(fmt.Errorf("testcase query %s must fail", failCase))
			continue
		}
	}
}

func TestQueryTypeB(t *testing.T) {
	callOnlyParser := parser.NewCallOnlyParser()
	baseParser := parser.NewSDQParser()
	v := QueryTypeB{
		Parser: baseParser,
	}
	type testSet struct {
		Query  string
		Answer QueryTypeB
	}

	testSuccessCases := make([]testSet, 0)

	testSuccessCases = append(testSuccessCases, testSet{
		Query: "1|2 hoge|3 fuga",
		Answer: QueryTypeB{
			Parser: v.Parser,
			ID:     1,
			QueryA1: QueryTypeA{
				ID:   2,
				Name: "hoge",
			},
			QueryA2: QueryTypeA{
				ID:   3,
				Name: "fuga",
			},
		},
	})

	for _, successCase := range testSuccessCases {
		if err := callOnlyParser.Parse(successCase.Query, &v); err != nil {
			t.Error(err)
			continue
		}
		if v != successCase.Answer {
			t.Error(fmt.Errorf("value mismatch: %+v vs %+v", v, successCase.Answer))
			continue
		}
	}
}

func TestQueryTypeC(t *testing.T) {
	baseParser := parser.NewSDQParser()
	customizedParser := &customParser{
		baseParser: parser.NewSDQParser(),
	}
	v := QueryTypeC{
		ID:      -1,
		QueryA1: QueryTypeA{},
		QueryB1: QueryTypeB{
			Parser: customizedParser,
		},
	}
	type testSet struct {
		Query  string
		Answer QueryTypeC
		Check  func(result QueryTypeC) error
	}

	testSuccessCases := make([]testSet, 0)

	testSuccessCases = append(testSuccessCases, testSet{
		Query: "1 2 hoge 3|4 fuga|5 piyo",
		Answer: QueryTypeC{
			ID: 1,
			QueryA1: QueryTypeA{
				ID:   2,
				Name: "hoge",
			},
			QueryB1: QueryTypeB{
				Parser: customizedParser,
				ID:     3,
				QueryA1: QueryTypeA{
					ID:   4,
					Name: "fuga",
				},
				QueryA2: QueryTypeA{
					ID:   5,
					Name: "piyo-customized",
				},
			},
			IgnoredValue: "",
		},
	})

	for _, successCase := range testSuccessCases {
		if err := baseParser.Parse(successCase.Query, &v); err != nil {
			t.Error(err)
			continue
		}
		if v != successCase.Answer {
			t.Error(fmt.Errorf("value mismatch: %+v vs %+v", v, successCase.Answer))
			continue
		}
	}
}

type customParser struct {
	baseParser parser.Parser
}

var _ parser.Parser = &customParser{}

func (cp *customParser) Parse(query string, target interface{}) error {
	return cp.baseParser.Parse(query+"-customized", target)
}
