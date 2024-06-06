package example_test

import (
	"fmt"
	"sdqorm/example"
	"sdqorm/pkg/parser"
	"testing"
)

func TestQueryTypeA(t *testing.T) {
	spanDividedQueryParser := parser.NewSDQParser(" ")
	v := example.QueryTypeA{
		SDQParser: spanDividedQueryParser,
		ID:        -1,
		Name:      "",
	}
	type testSet struct {
		Query  string
		Answer example.QueryTypeA
	}

	testSuccessCases := make([]testSet, 0)
	testFailCases := make([]string, 0)

	testSuccessCases = append(testSuccessCases, testSet{
		Query: "1 test1",
		Answer: example.QueryTypeA{
			SDQParser: v.SDQParser,
			ID:        1,
			Name:      "test1",
		},
	})

	testFailCases = append(testFailCases,
		"-1 test2",
		"1.23 test2",
		"1 test-case",
	)

	for _, successCase := range testSuccessCases {
		if err := v.String2Struct(successCase.Query, &v); err != nil {
			t.Error(err)
			continue
		}
		if v != successCase.Answer {
			t.Error(fmt.Errorf("value mismatch: %+v vs %+v", v, successCase.Answer))
			continue
		}
	}

	for _, failCase := range testFailCases {
		if err := v.String2Struct(failCase, &v); err == nil {
			t.Error(fmt.Errorf("testcase query %s must fail", failCase))
			continue
		}
	}
}

func TestQueryTypeB(t *testing.T) {
	//defaultParser := parser.NewSDQParser("")
	spanDividedQueryParser := parser.NewSDQParser(" ")
	colonDividedQueryParser := parser.NewSDQParser(":")
	v := example.QueryTypeB{
		SDQParser: spanDividedQueryParser,
		ID:        -1,
		QueryA1: example.QueryTypeA{
			SDQParser: colonDividedQueryParser,
		},
		QueryA2: example.QueryTypeA{
			SDQParser: colonDividedQueryParser,
		},
	}
	type testSet struct {
		Query  string
		Answer example.QueryTypeB
	}

	testSuccessCases := make([]testSet, 0)

	testSuccessCases = append(testSuccessCases, testSet{
		Query: "1 2:hoge 3:fuga",
		Answer: example.QueryTypeB{
			SDQParser: v.SDQParser,
			ID:        1,
			QueryA1: example.QueryTypeA{
				SDQParser: v.QueryA1.SDQParser,
				ID:        2,
				Name:      "hoge",
			},
			QueryA2: example.QueryTypeA{
				SDQParser: v.QueryA2.SDQParser,
				ID:        3,
				Name:      "fuga",
			},
		},
	})

	for _, successCase := range testSuccessCases {
		if err := v.String2Struct(successCase.Query, &v); err != nil {
			t.Error(err)
			continue
		}
		if v != successCase.Answer {
			t.Error(fmt.Errorf("value mismatch: %+v vs %+v", v, successCase.Answer))
			continue
		}
	}
}
