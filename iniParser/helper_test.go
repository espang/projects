package main

import (
	"fmt"
	"testing"
)

var emptyAndComment = map[string]bool{
	"":              true,
	" ":             true,
	" [":            false,
	"      k=v    ": false,
	"; asdasdas":    true,
	"  ; ;  ; asd":  true,
}

func TestEmpty(t *testing.T) {

	for k, v := range emptyAndComment {
		result := IsEmptyOrComment([]byte(k))
		if v != result {
			t.Errorf("IsEmptyOrComment(%s): expect %v got %v", k, v, result)
		}
	}

}

type secTest struct {
	name string
	ok   bool
}

func (s secTest) String() string {
	return fmt.Sprintf("(%s|%v)", s.name, s.ok)
}

var sections = map[string]secTest{
	"[section1]":  {"section1", true},
	"[]":          {"", false},
	"[1]":         {"1", true},
	"[[12]":       {"[12", true},
	"[asd":        {"", false},
	"asd]":        {"", false},
	"key=value[]": {"", false},
}

func TestSection(t *testing.T) {
	for k, v := range sections {
		name, ok := GetSectionFromLine([]byte(k))
		if name != v.name || ok != v.ok {
			t.Errorf("GetSectionFromLine(%s): expect %s got %s",
				k,
				v,
				secTest{name, ok},
			)
		}
	}
}

type keyValuePair struct {
	key       string
	value     string
	erroneous bool
}

func (p keyValuePair) String() string {
	return fmt.Sprintf("(%s->%s)", p.key, p.value)
}

var testMapKeyValue = map[string]keyValuePair{
	" a= b":     keyValuePair{"a", "b", false},
	"a=b=c":     keyValuePair{"", "", true},
	"a a:=c3d]": keyValuePair{"a a:", "c3d]", false},
}

func TestKeyValue(t *testing.T) {
	for k, v := range testMapKeyValue {
		key, value, err := ParseKeyValuePair(1, []byte(k))
		if err != nil && v.erroneous == false {
			t.Errorf("ParseKeyValuePair(%s) should be ok but returns error %v",
				k,
				err,
			)
			continue
		}
		if err == nil && v.erroneous == true {
			t.Errorf("ParseKeyValuePair(%s) should be erroneous but returns no error",
				k,
			)
			continue
		}
		if key != v.key || value != v.value {
			t.Errorf("ParseKeyValuePair(%s): expect %s got %s",
				k,
				v,
				keyValuePair{key, value, false},
			)
		}

	}
}
