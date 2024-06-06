package example

import "sdqorm/pkg/parser"

// QueryTypeA parses query `ID Name `
type QueryTypeA struct {
	parser.SDQParser
	ID   int    `sdqorm:"index:0,min:0,max:10000"`
	Name string `sdqorm:"index:1,regexp:^[0-9a-zA-Z]+$"`
}

type QueryTypeB struct {
	parser.SDQParser
	ID      int        `sdqorm:"index:0,min:0,max:10000"`
	QueryA1 QueryTypeA `sdqorm:"index:1"`
	QueryA2 QueryTypeA `sdqorm:"index:2"`
}

//TODO:
// 「"custom":true」の場合には、パース処理をカスタマイズできるようにしたい。
