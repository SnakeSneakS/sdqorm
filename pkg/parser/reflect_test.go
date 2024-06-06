package parser

import (
	"reflect"
	"testing"
)

type Test struct{}

func TestHandleFuncOnSpecificTag(t *testing.T) {
	var test Test
	i := 0
	if err := handleFuncOnSpecificTag("hoge", test, func(field reflect.Value, typeField reflect.StructField, tag string) error { return nil }); err == nil {
		t.Errorf("function for interface which is not pointer must fail: %v", err)
	}
	if err := handleFuncOnSpecificTag("hoge", &i, func(field reflect.Value, typeField reflect.StructField, tag string) error { return nil }); err == nil {
		t.Errorf("function for interface which is not pointer for struct must fail: %v", err)
	}
	if err := handleFuncOnSpecificTag("hoge", &test, func(field reflect.Value, typeField reflect.StructField, tag string) error { return nil }); err != nil {
		t.Errorf("function for interface which is pointer for struct must success: %v", err)
	}
}
