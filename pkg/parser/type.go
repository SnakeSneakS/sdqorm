package parser

import "errors"

const SdqormTagName = "sdqorm"

type TagPartKey string
type IntTagPartKey TagPartKey
type StringTagPartKey TagPartKey

const (
	IndexKey TagPartKey = "index" //int

	IntMinKey IntTagPartKey = "min" //int
	IntMaxKey IntTagPartKey = "max" //int

	StringRegexpKey StringTagPartKey = "regexp" //string
)

var requiredTagPartKeys = []TagPartKey{
	IndexKey,
}

var errInvalidKey = errors.New("invalid key")
