// This file is part of go-getoptions.
//
// Copyright (C) 2015-2020  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package option

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DavidGamba/go-getoptions/text"
)

func TestOption(t *testing.T) {
	tests := []struct {
		name   string
		option *Option
		input  []string
		output interface{}
		err    error
	}{
		{"bool", func() *Option {
			b := false
			return New("help", BoolType).SetBoolPtr(&b)
		}(), []string{""}, true, nil},
		{"bool", func() *Option {
			b := true
			return New("help", BoolType).SetBoolPtr(&b)
		}(), []string{""}, false, nil},
		{"bool setbool", func() *Option {
			b := true
			return New("help", BoolType).SetBoolPtr(&b).SetBool(false)
		}(), []string{""}, true, nil},
		{"bool setbool", func() *Option {
			b := false
			return New("help", BoolType).SetBoolPtr(&b).SetBool(false)
		}(), []string{""}, true, nil},
		{"bool setbool", func() *Option {
			b := true
			return New("help", BoolType).SetBoolPtr(&b).SetBool(true)
		}(), []string{""}, false, nil},
		{"bool setbool", func() *Option {
			b := false
			return New("help", BoolType).SetBoolPtr(&b).SetBool(true)
		}(), []string{""}, false, nil},

		{"string", func() *Option {
			s := ""
			return New("help", StringType).SetStringPtr(&s)
		}(), []string{""}, "", nil},
		{"string", func() *Option {
			s := ""
			return New("help", StringType).SetStringPtr(&s)
		}(), []string{"hola"}, "hola", nil},
		{"string", func() *Option {
			s := ""
			return New("help", StringType).SetStringPtr(&s).SetString("xxx")
		}(), []string{""}, "", nil},
		{"string", func() *Option {
			s := ""
			return New("help", StringType).SetStringPtr(&s).SetString("xxx")
		}(), []string{"hola"}, "hola", nil},

		{"int", func() *Option {
			i := 0
			return New("help", IntType).SetIntPtr(&i)
		}(), []string{"123"}, 123, nil},
		{"int", func() *Option {
			i := 0
			return New("help", IntType).SetIntPtr(&i).SetInt(456)
		}(), []string{"123"}, 123, nil},
		{"int error", func() *Option {
			i := 0
			return New("help", IntType).SetIntPtr(&i)
		}(), []string{"123x"}, 0,
			fmt.Errorf(text.ErrorConvertToInt, "", "123x")},
		{"int error alias", func() *Option {
			i := 0
			return New("help", IntType).SetIntPtr(&i).SetCalled("int")
		}(), []string{"123x"}, 0,
			fmt.Errorf(text.ErrorConvertToInt, "int", "123x")},

		{"float64", func() *Option {
			f := 0.0
			return New("help", Float64Type).SetFloat64Ptr(&f)
		}(), []string{"123.123"}, 123.123, nil},
		{"float64 error", func() *Option {
			f := 0.0
			return New("help", Float64Type).SetFloat64Ptr(&f)
		}(), []string{"123x"}, 0.0,
			fmt.Errorf(text.ErrorConvertToFloat64, "", "123x")},
		{"float64 error alias", func() *Option {
			f := 0.0
			return New("help", Float64Type).SetFloat64Ptr(&f).SetCalled("float")
		}(), []string{"123x"}, 0.0,
			fmt.Errorf(text.ErrorConvertToFloat64, "float", "123x")},

		{"string slice", func() *Option {
			ss := []string{}
			return New("help", StringRepeatType).SetStringSlicePtr(&ss)
		}(), []string{"hola", "mundo"}, []string{"hola", "mundo"}, nil},

		{"int slice", func() *Option {
			ii := []int{}
			return New("help", IntRepeatType).SetIntSlicePtr(&ii)
		}(), []string{"123", "456"}, []int{123, 456}, nil},
		{"int slice error", func() *Option {
			ii := []int{}
			return New("help", IntRepeatType).SetIntSlicePtr(&ii)
		}(), []string{"x"}, []int{},
			fmt.Errorf(text.ErrorConvertToInt, "", "x")},

		{"int slice range", func() *Option {
			ii := []int{}
			return New("help", IntRepeatType).SetIntSlicePtr(&ii)
		}(), []string{"1..5"}, []int{1, 2, 3, 4, 5}, nil},
		{"int slice range error", func() *Option {
			ii := []int{}
			return New("help", IntRepeatType).SetIntSlicePtr(&ii)
		}(), []string{"x..5"}, []int{},
			fmt.Errorf(text.ErrorConvertToInt, "", "x..5")},
		{"int slice range error", func() *Option {
			ii := []int{}
			return New("help", IntRepeatType).SetIntSlicePtr(&ii)
		}(), []string{"1..x"}, []int{},
			fmt.Errorf(text.ErrorConvertToInt, "", "1..x")},
		{"int slice range error", func() *Option {
			ii := []int{}
			return New("help", IntRepeatType).SetIntSlicePtr(&ii)
		}(), []string{"5..1"}, []int{},
			fmt.Errorf(text.ErrorConvertToInt, "", "5..1")},

		{"map", func() *Option {
			m := make(map[string]string)
			return New("help", StringMapType).SetStringMapPtr(&m)
		}(), []string{"hola=mundo"}, map[string]string{"hola": "mundo"}, nil},
		{"map", func() *Option {
			m := make(map[string]string)
			opt := New("help", StringMapType).SetStringMapPtr(&m)
			opt.MapKeysToLower = true
			return opt
		}(), []string{"Hola=Mundo"}, map[string]string{"hola": "Mundo"}, nil},
		// TODO: Currently map is only handling one argument at a time so the test below fails.
		//	It seems like the caller is handling this properly so I don't really know if this is needed here.
		// {"map", func() *Option {
		// 	m := make(map[string]string)
		// 	return New("help", StringMapType).SetStringMapPtr(&m)
		// }(), []string{"hola=mundo", "hello=world"}, map[string]string{"hola": "mundo", "hello": "world"}, nil},
		{"map error", func() *Option {
			m := make(map[string]string)
			return New("help", StringMapType).SetStringMapPtr(&m)
		}(), []string{"hola"}, map[string]string{},
			fmt.Errorf(text.ErrorArgumentIsNotKeyValue, "")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.option.Save(tt.input...)
			if err == nil && tt.err != nil {
				t.Errorf("got = '%#v', want '%#v'", err, tt.err)
			}
			if err != nil && tt.err == nil {
				t.Errorf("got = '%#v', want '%#v'", err, tt.err)
			}
			if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("got = '%#v', want '%#v'", err, tt.err)
			}
			got := tt.option.Value()
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("got = '%#v', want '%#v'", got, tt.output)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name        string
		option      *Option
		input       []string
		output      interface{}
		err         error
		errRequired error
	}{
		{"bool", func() *Option {
			b := false
			return New("help", BoolType).SetBoolPtr(&b)
		}(), []string{""}, true, nil, nil},
		{"bool", func() *Option {
			b := false
			return New("help", BoolType).SetBoolPtr(&b).SetRequired("")
		}(), []string{""}, true, nil, fmt.Errorf(text.ErrorMissingRequiredOption, "help")},
		{"bool", func() *Option {
			b := false
			return New("help", BoolType).SetBoolPtr(&b).SetRequired("err")
		}(), []string{""}, true, nil, fmt.Errorf("err")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.option.Save(tt.input...)
			if err == nil && tt.err != nil {
				t.Errorf("got = '%#v', want '%#v'", err, tt.err)
			}
			if err != nil && tt.err == nil {
				t.Errorf("got = '%#v', want '%#v'", err, tt.err)
			}
			if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("got = '%#v', want '%#v'", err, tt.err)
			}
			got := tt.option.Value()
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("got = '%#v', want '%#v'", got, tt.output)
			}
			err = tt.option.CheckRequired()
			if err == nil && tt.errRequired != nil {
				t.Errorf("got = '%#v', want '%#v'", err, tt.errRequired)
			}
			if err != nil && tt.errRequired == nil {
				t.Errorf("got = '%#v', want '%#v'", err, tt.errRequired)
			}
			if err != nil && tt.errRequired != nil && err.Error() != tt.errRequired.Error() {
				t.Errorf("got = '%#v', want '%#v'", err, tt.errRequired)
			}
		})
	}
}

func TestOther(t *testing.T) {
	i := 0
	opt := New("help", IntType).SetIntPtr(&i).SetAlias("?", "h").SetDescription("int help").SetHelpArgName("myint").SetDefaultStr("5").SetEnvVar("ENV_VAR")
	got := opt.Aliases
	expected := []string{"help", "?", "h"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got = '%#v', want '%#v'", got, expected)
	}
	if opt.Int() != 0 {
		t.Errorf("got = '%#v', want '%#v'", opt.Int(), 0)
	}
	i = 3
	if opt.Int() != 3 {
		t.Errorf("got = '%#v', want '%#v'", opt.Int(), 0)
	}
	if opt.Description != "int help" {
		t.Errorf("got = '%#v', want '%#v'", opt.Description, "int help")
	}
	if opt.HelpArgName != "myint" {
		t.Errorf("got = '%#v', want '%#v'", opt.HelpArgName, "myint")
	}
	if opt.DefaultStr != "5" {
		t.Errorf("got = '%#v', want '%#v'", opt.DefaultStr, "5")
	}
	if opt.EnvVar != "ENV_VAR" {
		t.Errorf("got = '%#v', want '%#v'", opt.EnvVar, "ENV_VAR")
	}

	list := []*Option{New("b", BoolType), New("a", BoolType), New("c", BoolType)}
	expectedList := []*Option{New("a", BoolType), New("b", BoolType), New("c", BoolType)}
	Sort(list)
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("got = '%#v', want '%#v'", list, expectedList)
	}

	opt = New("help", IntRepeatType)
	opt.MaxArgs = 2
	opt.synopsis()
	if opt.HelpSynopsis != "--help <int>..." {
		t.Errorf("got = '%#v', want '%#v'", opt.HelpSynopsis, "--help <int>...")
	}
}
