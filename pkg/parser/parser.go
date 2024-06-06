package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// SDQParser Something-Devided Query Parser
type SDQParser interface {
	String2Struct(query string, target interface{}) error
}

// BaseSDQParser Something-Devided Query Parser
type BaseSDQParser struct {
	delimiter string
}

func NewSDQParser(
	delimiter string,
) SDQParser {
	return &BaseSDQParser{
		delimiter: delimiter,
	}
}

//func (p *BaseSDQParser) compareAddress(i1, i2 interface{}) bool {
//	//return fmt.Sprintf("%p", i1) == fmt.Sprintf("%p", i2)
//	return reflect.ValueOf(i1).Pointer() == reflect.ValueOf(i2).Pointer()
//}

func (p *BaseSDQParser) String2Struct(query string, target interface{}) error {
	//targetがデフォルトでパーサを持っている(かつ自身ではない)場合targetのパーサを利用する
	//その際「関数のアドレス」を比較している
	//parser, ok := target.(SDQParser)
	//
	//if ok && !isFuncAddrSame {
	//	fmt.Println("HOGE")
	//	return parser.String2Struct(query, target)
	//} else {
	//持っていない場合は特定のタグのフィールドに対して独自のパーサを利用する
	// Iterate over the fields of the struct
	if err := handleFuncOnSpecificTag(SdqormTagName, target, func(field reflect.Value, typeField reflect.StructField, tag string) error {
		return p.parseSdq2Struct(query, field, typeField, tag)
	}); err != nil {
		return err
	}
	//}
	//isFuncAddrSame := p.compareAddress(target.(SDQParser).String2Struct, p.String2Struct)
	//fmt.Println("hoge: ", isFuncAddrSame)
	return nil
}

// parseSdq2Struct parse something-devided query into struct
func (p *BaseSDQParser) parseSdq2Struct(query string, field reflect.Value, typeField reflect.StructField, tag string) error {
	// Split the input string by whitespace
	queryParts := strings.Split(query, p.delimiter) //Fields(query)
	if len(queryParts) == 0 {
		return fmt.Errorf("input string is empty or invalid")
	}

	tagParts := strings.Split(tag, ",")
	if len(tagParts) == 0 {
		return nil
	}

	// check
	if err := checkRequiredKeys(tagParts, requiredTagPartKeys); err != nil {
		return err
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

	// Find the index part
	var index int
	indexString, ok := tagPartsKV[string(IndexKey)]
	if !ok {
		return fmt.Errorf("setting \"index\" is required when using %s for field %s", SdqormTagName, typeField.Name)
	}
	index, err := strconv.Atoi(indexString)
	if err != nil {
		return fmt.Errorf("failed to parse int for index of field %s: %v", typeField.Name, err)
	}
	if index < 0 || index >= len(queryParts) {
		return fmt.Errorf("index out of range for field %s", typeField.Name)
	}
	targetQuery := queryParts[index]

	// Set the field value based on the type
	switch field.Kind() {
	case reflect.Int:
		intVal, err := strconv.Atoi(targetQuery)
		if err != nil {
			return fmt.Errorf("failed to parse int for field %s: %v", typeField.Name, err)
		}
		for k, v := range tagPartsKV {
			switch k {
			case string(IndexKey):
				continue
			case string(IntMinKey):
				min, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("failed to parse int for %s on field %s: %v", IntMinKey, typeField.Name, err)
				}
				if intVal < min {
					return fmt.Errorf("invalid value for field %s: value %d cannot be less than %d", typeField.Name, intVal, min)
				}
			case string(IntMaxKey):
				max, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("failed to parse int for %s on field %s: %v", IntMinKey, typeField.Name, err)
				}
				if intVal > max {
					return fmt.Errorf("invalid value for field %s: value %d cannot be larger than %d", typeField.Name, intVal, max)
				}
			default:
				return fmt.Errorf("%w: %s, available keys are below: %v", errInvalidKey, k, []string{string(IndexKey), string(IntMinKey), string(IntMaxKey)})
			}
		}
		field.SetInt(int64(intVal))
	case reflect.String:
		for k, v := range tagPartsKV {
			switch k {
			case string(IndexKey):
				continue
			case string(StringRegexpKey):
				re, err := regexp.Compile(v)
				if err != nil {
					return fmt.Errorf("failed to parse regexp for %s on field %s: %v", StringRegexpKey, typeField.Name, err)
				}
				if !re.MatchString(targetQuery) {
					return fmt.Errorf("invalid value for field %s: value \"%s\" must satisfy regexp \"%s\"", typeField.Name, targetQuery, v)
				}
			default:
				return fmt.Errorf("%w: %s, available keys are below: %v", errInvalidKey, k, []string{string(IndexKey), string(StringRegexpKey)})
			}
		}
		field.SetString(targetQuery)
	case reflect.Interface:
	case reflect.Struct:
		if !field.CanInterface() {
			break
		}

		//TODO: Custom Parserが存在する場合はこれを利用する。これをString2Structの方でジッ素王する
		if parser, ok := field.Interface().(SDQParser); ok {
			if err := parser.String2Struct(targetQuery, field.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to parse value for field %s, invalid value %s: %v", typeField.Name, targetQuery, err)
			}
			field.Set(reflect.ValueOf(field.Interface()))
		} else {
			//存在しない場合にはデフォルトのパーサを利用する
			//return fmt.Errorf("struct/interface which cannot be converted into %v is not able to be handled", p)
			if err := p.String2Struct(targetQuery, field.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to parse value for field %s, invalid value %s: %v", typeField.Name, targetQuery, err)
			}
			field.Set(reflect.ValueOf(field.Interface()))
		}

	default:
		return fmt.Errorf("unsupported field type %s", field.Kind().String())
	}

	return nil
}
