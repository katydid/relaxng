// Copyright 2018 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package relaxng

import (
	"testing"

	"github.com/katydid/katydid/parser/debug"
	"github.com/katydid/katydid/relapse/ast"
	c "github.com/katydid/katydid/relapse/combinator"
	"github.com/katydid/katydid/relapse/compose"
)

func TestToken(t *testing.T) {
	expr := ast.NewFunction("token", c.StringVar(), c.StringConst("TheStreet"))
	b, err := compose.NewBool(expr)
	if err != nil {
		t.Fatal(err)
	}
	f, err := compose.NewBoolFunc(b)
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.Eval(debug.NewStringValue("text_TheStreet"))
	if err != nil {
		t.Fatal(err)
	}
	if r != true {
		t.Fatalf("expected true")
	}
	r, err = f.Eval(debug.NewStringValue("text_ThatStreet"))
	if err != nil {
		t.Fatal(err)
	}
	if r != false {
		t.Fatalf("expected false")
	}
}

func TestWhitespace(t *testing.T) {
	expr := ast.NewFunction("whitespace", c.StringVar())
	b, err := compose.NewBool(expr)
	if err != nil {
		t.Fatal(err)
	}
	f, err := compose.NewBoolFunc(b)
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.Eval(debug.NewStringValue("text_   	 "))
	if err != nil {
		t.Fatal(err)
	}
	if r != true {
		t.Fatalf("expected true")
	}
	r, err = f.Eval(debug.NewStringValue("text_   a  "))
	if err != nil {
		t.Fatal(err)
	}
	if r != false {
		t.Fatalf("expected false")
	}
}

func TestAnyText(t *testing.T) {
	expr := ast.NewFunction("anytext", c.StringVar())
	b, err := compose.NewBool(expr)
	if err != nil {
		t.Fatal(err)
	}
	f, err := compose.NewBoolFunc(b)
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.Eval(debug.NewStringValue("text_bla"))
	if err != nil {
		t.Fatal(err)
	}
	if r != true {
		t.Fatalf("expected true")
	}
	r, err = f.Eval(debug.NewStringValue("a"))
	if err != nil {
		t.Fatal(err)
	}
	if r != false {
		t.Fatalf("expected false")
	}
}

func TestText(t *testing.T) {
	expr := ast.NewFunction("text", c.StringVar(), c.StringConst("bla"))
	b, err := compose.NewBool(expr)
	if err != nil {
		t.Fatal(err)
	}
	f, err := compose.NewBoolFunc(b)
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.Eval(debug.NewStringValue("text_bla"))
	if err != nil {
		t.Fatal(err)
	}
	if r != true {
		t.Fatalf("expected true")
	}
	r, err = f.Eval(debug.NewStringValue("a"))
	if err != nil {
		t.Fatal(err)
	}
	if r != false {
		t.Fatalf("expected false")
	}
}

func TestList(t *testing.T) {
	expr := ast.NewFunction("list", c.StringVar(), c.StringConst("bla"))
	b, err := compose.NewBool(expr)
	if err != nil {
		t.Fatal(err)
	}
	f, err := compose.NewBoolFunc(b)
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.Eval(debug.NewStringValue("text_bla"))
	if err != nil {
		t.Fatal(err)
	}
	if r != true {
		t.Fatalf("expected true")
	}
	r, err = f.Eval(debug.NewStringValue("a"))
	if err != nil {
		t.Fatal(err)
	}
	if r != false {
		t.Fatalf("expected false")
	}
}
