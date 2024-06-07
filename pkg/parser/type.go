package parser

const SdqormTagName = "sdqorm"

type TagPartKey string
type IntTagPartKey TagPartKey
type StringTagPartKey TagPartKey

const (
	IgnoreKey    TagPartKey = "ignore"    //[boolean] this field is ignored
	IndexKey     TagPartKey = "index"     //[string] indexExpr: "index" or "indexStart-" or "-indexEnd" or "indexStart-indexEnd"
	DelimiterKey TagPartKey = "delimiter" //[string] default: " "
	UseCustomKey TagPartKey = "custom"    //[boolean] if using custom parser or not
)
