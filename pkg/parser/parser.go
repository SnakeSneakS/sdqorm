package parser

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Parser Something-Devided Query Parser
type Parser interface {
	Parse(query string, target interface{}) error
	//Stringify(target interface{}) (string,error)
}

// CallOnlyParser doesn't have ability to parse itself, but just find and run other custom parser parsing function
type CallOnlyParser struct{}

func NewCallOnlyParser() Parser {
	return &CallOnlyParser{}
}

func (p *CallOnlyParser) Parse(query string, target interface{}) error {
	if parser, ok := target.(Parser); ok {
		return parser.Parse(query, target)
	} else {
		return errors.New("custom parser must be set for each struct")
	}
}

const delimiterDefault = " "

// BaseSDQParser Something-Devided Query Parser
type BaseSDQParser struct{}

func NewSDQParser() Parser {
	return &BaseSDQParser{}
}

func (p *BaseSDQParser) Parse(query string, target interface{}) error {
	if err := handleFuncOnSpecificTag(SdqormTagName, target, func(field reflect.Value, typeField reflect.StructField, tag string) error {
		return p.parseSdq2Struct(query, field, typeField, tag)
	}); err != nil {
		return err
	}
	return nil
}

// parseSdq2Struct parse something-devided query into struct
func (p *BaseSDQParser) parseSdq2Struct(query string, field reflect.Value, typeField reflect.StructField, tag string) error {
	// Split the input string by whitespace
	tag, _ = strings.CutSuffix(tag, ",")
	tagParts := strings.Split(tag, ",")
	if len(tagParts) == 0 {
		return nil
	}

	//mapping
	tagPartsKV := make(map[string]string)
	for _, part := range tagParts {
		split := strings.Split(part, ":")
		if len(split) != 2 {
			return fmt.Errorf("invalid tag: %s", part)
		}
		tagPartsKV[split[0]] = split[1]
	}

	//ignore checking
	ignoreBool := false
	ignoreBoolStr, ok := tagPartsKV[string(IgnoreKey)]
	if ok {
		var err error
		ignoreBool, err = strconv.ParseBool(ignoreBoolStr)
		if err != nil {
			return fmt.Errorf("failed to parse value for field %s, invalid value \"%s\" for key \"%s\": %v", typeField.Name, ignoreBoolStr, IgnoreKey, err)
		}
	}
	if ignoreBool {
		return nil
	}

	//Find Delimiter
	delimiter, ok := tagPartsKV[string(DelimiterKey)]
	if !ok {
		delimiter = delimiterDefault
	}

	// Find the index part
	indexExp, ok := tagPartsKV[string(IndexKey)]
	if !ok {
		return fmt.Errorf("setting \"index\" is required when using %s for field %s", SdqormTagName, typeField.Name)
	}
	targetQuery, err := p.extractPartFromQuery(query, delimiter, indexExp)
	if err != nil {
		return err
	}

	// Set the field value based on the type
	switch field.Kind() {
	case reflect.Int:
		intVal, err := strconv.Atoi(targetQuery)
		if err != nil {
			return fmt.Errorf("failed to parse int for field %s: %v", typeField.Name, err)
		}
		field.SetInt(int64(intVal))
	case reflect.String:
		field.SetString(targetQuery)
	case reflect.Interface:
	case reflect.Struct:
		useCustomBool := false
		useCustomStr, ok := tagPartsKV[string(UseCustomKey)]
		if ok {
			useCustomBool, err = strconv.ParseBool(useCustomStr)
			if err != nil {
				return fmt.Errorf("failed to parse value for field %s, invalid value \"%s\" for key \"%s\": %v", typeField.Name, useCustomStr, UseCustomKey, err)
			}
		}

		if !useCustomBool {
			//If you don't use custom parser, use this parser
			if err := p.Parse(targetQuery, field.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to parse value for field %s, invalid value %s: %v", typeField.Name, targetQuery, err)
			}
			field.Set(reflect.ValueOf(field.Interface()))
			break
		} else {
			//if you use custom parser
			if !field.CanInterface() {
				return fmt.Errorf("failed to parse value for field %s: you must use struct or interface for this field", typeField.Name)
			}

			//try to use custom parser
			if parser, ok := field.Interface().(Parser); ok {
				if err := parser.Parse(targetQuery, field.Addr().Interface()); err != nil {
					return fmt.Errorf("failed to parse value for field %s: invalid value %s: %v", typeField.Name, targetQuery, err)
				}
				field.Set(reflect.ValueOf(field.Interface()))
				break
			}

			return fmt.Errorf("failed to parse value for field %s: trying to use custom parser, but not found", typeField.Name)
		}
	default:
		return fmt.Errorf("unsupported field type %s", field.Kind().String())
	}

	return nil
}

// extractQueryPart extract specific part from query based on delimiter & indexFrom & indexTo
func (p *BaseSDQParser) extractPartFromQuery(query, delimiter string, indexExp string) (string, error) {
	queryParts := strings.Split(query, delimiter)

	indexStart, indexEnd, err := p.parseIndexExp(indexExp)
	if err != nil {
		return "", err
	}
	if indexEnd == -1 {
		indexEnd = len(queryParts)
	}

	if indexStart < 0 || indexEnd > len(queryParts) {
		return "", fmt.Errorf("index out of range: start=%d, end=%d, query=%s, delimiter=%s", indexStart, indexEnd, query, delimiter)
	}

	return strings.Join(queryParts[indexStart:indexEnd], delimiter), nil
}

// parseIndexExp parses indexExp as indexStart & indexEnd
//
//	indexExp: int or int- or -int or int-int
//	e.g.)
//		if "0-2", indexStart is 0 and indexEnd is 2
//		if "-10", indexStart is 0 and indexEnd is 10
//		if "2-", indexStart is 2 and indexEnd is -1
func (p *BaseSDQParser) parseIndexExp(indexExp string) (indexStart int, indexEnd int, err error) {
	delimiter := "-"
	parts := strings.Split(indexExp, delimiter)
	if len(parts) == 2 {
		//set indexFrom
		if parts[0] == "" {
			indexStart = 0
		} else {
			indexStart, err = strconv.Atoi(parts[0])
			if err != nil {
				err = fmt.Errorf("failed to parse indexExp %s: %v", indexExp, err)
				return
			}
		}
		//set indexTo
		if parts[1] == "" {
			indexEnd = -1
		} else {
			indexEnd, err = strconv.Atoi(parts[1])
			if err != nil {
				err = fmt.Errorf("failed to parse indexExp %s: %v", indexExp, err)
				return
			}
		}
	} else if len(parts) == 1 {
		indexStart, err = strconv.Atoi(parts[0])
		if err != nil {
			err = fmt.Errorf("failed to parse indexExp %s: %v", indexExp, err)
			return
		}
		indexEnd = indexStart + 1
	} else {
		err = fmt.Errorf("invalid indexExp format %s", indexExp)
		return
	}

	if indexEnd != -1 && indexStart >= indexEnd {
		err = fmt.Errorf("invalid indexExp format %s. index out of range: start=%d, end=%d", indexExp, indexStart, indexEnd)
		return
	}
	return
}
