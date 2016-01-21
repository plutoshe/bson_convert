package bson_convert

import (
	"reflect"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestConvertToBSON(t *testing.T) {
	type substructure struct {
		E string `bson:"e,omitempty"`
		F []int  `bson:"f,omitempty"`
	}
	type testStructure struct {
		A int          `bson:"a,omitempty"`
		B string       `bson:"b,omitempty"`
		C []int        `bson:"c,omitempty"`
		D substructure `bson:"d"`
	}

	testcases := []struct {
		source   testStructure
		omitted  bool
		required []string
		ignored  []string
		res      bson.M
	}{
		{
			testStructure{
				3,
				"BBB",
				[]int{32, 123, 32},
				substructure{"3", []int{3, 2}},
			},
			false,
			[]string{},
			[]string{},
			bson.M{"a": 3, "b": "BBB", "c": []int{32, 123, 32}, "d": bson.M{"e": "3", "f": []int{3, 2}}},
		},
		{
			testStructure{
				3,
				"BBB",
				[]int{32, 123, 32},
				substructure{"3", []int{3, 2}},
			},
			true,
			[]string{},
			[]string{},
			bson.M{"a": 3, "b": "BBB", "c": []int{32, 123, 32}, "d": bson.M{"e": "3", "f": []int{3, 2}}},
		},
		{
			testStructure{
				0,
				"",
				[]int{32, 123, 32},
				substructure{"", []int{}},
			},
			true,
			[]string{},
			[]string{},
			bson.M{"c": []int{32, 123, 32}},
		},
		// require more specific attr after involving its fathers attr into required field
		{
			testStructure{
				0,
				"",
				[]int{32, 123, 32},
				substructure{"", []int{}},
			},
			true,
			[]string{"a", "b", "d", "e"},
			[]string{},
			bson.M{"a": 0, "b": "", "c": []int{32, 123, 32}, "d": bson.M{"e": ""}},
		},
		{
			testStructure{
				3,
				"BBB",
				[]int{32, 123, 32},
				substructure{"3", []int{3, 2}},
			},
			false,
			[]string{},
			[]string{"c", "d"},
			bson.M{"a": 3, "b": "BBB"},
		},
		{
			testStructure{
				0,
				"BBB",
				[]int{32, 123, 32},
				substructure{"3", []int{3, 2}},
			},
			false,
			[]string{"a"},
			[]string{"c", "e"},
			bson.M{"a": 0, "b": "BBB", "d": bson.M{"f": []int{3, 2}}},
		},
	}
	for i, v := range testcases {
		res := Convert2BSON(v.source, v.omitted, v.required, v.ignored)
		if !reflect.DeepEqual(v.res, res) {
			t.Errorf("#%d: Convert result is wrong! Want = %v, Get = %v.\n", i, v.res, res)
		}
	}
}
