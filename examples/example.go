package examples

import "github.com/snakesneaks/sdqorm/pkg/parser"

// QueryTypeA is parsed from query `ID Name`
// simplest struct
type QueryTypeA struct {
	//parser.Parser
	ID   int    `sdqorm:"index:0"`
	Name string `sdqorm:"index:1"`
}

// QueryTypeB is parsed from query `ID|ID Name|ID Name`
// struct which contains struct
type QueryTypeB struct {
	parser.Parser            //this can be used for implement custom query parser
	ID            int        `sdqorm:"index:0,delimiter:|"`
	QueryA1       QueryTypeA `sdqorm:"index:1,delimiter:|"`
	QueryA2       QueryTypeA `sdqorm:"index:2,delimiter:|"`
}

// QueryTypeC is parsed from query `ID QueryA_ID QueryA_Name QueryB_queries...`
// complicated schema
type QueryTypeC struct {
	ID           int        `sdqorm:"index:0,delimiter: ,ignore:false"`
	QueryA1      QueryTypeA `sdqorm:"index:1-3,delimiter: ,custom:false"`
	QueryB1      QueryTypeB `sdqorm:"index:3-,delimiter: ,custom:true"`
	IgnoredValue string     `sdqorm:"index:0,delimiter: ,ignore:true"`
}
