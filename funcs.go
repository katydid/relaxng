// Copyright 2015 Walter Schulze
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
	"fmt"
	"regexp"
	"strings"

	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/combinator"
	"github.com/katydid/katydid/relapse/funcs"
)

func stripTextPrefix(s string) (string, error) {
	if !strings.HasPrefix(s, "text_") {
		return "", fmt.Errorf("%q is not of type text", s)
	}
	return strings.Replace(s, "text_", "", 1), nil
}

func newTokenValue(t string) *ast.Pattern {
	return combinator.Value(&token{funcs.StringVar(), funcs.StringConst(t), ""})
}

// token is a function used in relapse to validate values as described here
// http://books.xmlschemata.org/relaxng/relax-CHP-7-SECT-4.html
type token struct {
	S funcs.String
	C funcs.ConstString
	c string
}

func (this *token) Init() error {
	c, err := this.C.Eval()
	if err != nil {
		return err
	}
	this.c = c
	return nil
}

func (this *token) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, nil
	}
	s, err = stripTextPrefix(s)
	if err != nil {
		return false, nil
	}
	ss := tokenize(s)
	s = strings.Join(ss, " ")
	return s == this.c, nil
}

func tokenize(s string) []string {
	ss := []string{}
	s1 := strings.Split(s, " ")
	for _, ss1 := range s1 {
		s2 := strings.Split(ss1, "\n")
		for _, ss2 := range s2 {
			s3 := strings.Split(ss2, "\r")
			for _, ss3 := range s3 {
				s4 := strings.Split(ss3, "\t")
				for _, ss4 := range s4 {
					s5 := strings.TrimSpace(ss4)
					if len(s5) > 0 {
						ss = append(ss, s5)
					}
				}
			}
		}
	}
	return ss
}

func init() {
	funcs.Register("token", new(token))
}

type whitespace struct {
	S funcs.String
}

func (this *whitespace) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, nil
	}
	s, err = stripTextPrefix(s)
	if err != nil {
		return false, nil
	}
	return len(strings.TrimSpace(s)) == 0, nil
}

func init() {
	funcs.Register("whitespace", new(whitespace))
}

type anytext struct {
	S funcs.String
}

func (this *anytext) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, nil
	}
	_, err = stripTextPrefix(s)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func init() {
	funcs.Register("anytext", new(anytext))
}

func newWhitespace() *ast.Pattern {
	return combinator.Value(&whitespace{funcs.StringVar()})
}

func newEmptyValue() *ast.Pattern {
	return combinator.Value(&whitespace{funcs.StringVar()})
}

func newAnyValue() *ast.Pattern {
	return combinator.Value(&anytext{funcs.StringVar()})
}

func newTextValue(value string) *ast.Pattern {
	return combinator.Value(&text{funcs.StringVar(), funcs.StringConst(value), ""})
}

type text struct {
	S funcs.String
	C funcs.ConstString
	c string
}

func (this *text) Init() error {
	c, err := this.C.Eval()
	if err != nil {
		return err
	}
	this.c = c
	return nil
}

func (this *text) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, nil
	}
	s, err = stripTextPrefix(s)
	if err != nil {
		return false, nil
	}
	return s == this.c, nil
}

func init() {
	funcs.Register("text", new(text))
}

type list struct {
	r    *regexp.Regexp
	S    funcs.String
	Expr funcs.ConstString
}

func (this *list) Init() error {
	e, err := this.Expr.Eval()
	if err != nil {
		return err
	}
	r, err := regexp.Compile(e)
	if err != nil {
		return err
	}
	this.r = r
	return nil
}

func (this *list) Eval() (bool, error) {
	s, err := this.S.Eval()
	if err != nil {
		return false, nil
	}
	s, err = stripTextPrefix(s)
	if err != nil {
		return false, nil
	}
	ss := tokenize(s)
	s = strings.Join(ss, " ")
	return this.r.MatchString(s), nil
}

func init() {
	funcs.Register("list", new(list))
}

func newList(nameOrPattern *NameOrPattern) *ast.Pattern {
	regexStr, nullable, err := listToRegex(nameOrPattern)
	if err != nil {
		return ast.NewNot(ast.NewZAny())
	}
	val := combinator.Value(&list{nil, funcs.StringVar(), funcs.StringConst("^" + regexStr + "$")})
	if !nullable {
		return val
	}
	return ast.NewOr(val, ast.NewEmpty())
}

func listToRegex(p *NameOrPattern) (string, bool, error) {
	if p.Empty != nil {
		return "", true, nil
	}
	if p.Data != nil {
		if p.Data.Except == nil {
			return `(\S)*`, false, nil
		}
	}
	if p.Value != nil {
		if len(p.Value.Ns) > 0 {
			panic("list value ns not supported")
		}
		if strings.Contains(p.Value.Text, " ") {
			return "", false, fmt.Errorf("unable to match a list")
		}
		return p.Value.Text, len(p.Value.Text) == 0, nil
	}
	if p.OneOrMore != nil {
		s, nullable, err := listToRegex(p.OneOrMore.NameOrPattern)
		return `(\s)?` + s + "(\\s" + s + ")*", nullable, err
	}
	if p.Choice != nil {
		l, nl, errl := listToRegex(p.Choice.Left)
		r, nr, errr := listToRegex(p.Choice.Right)
		if errl != nil && errr != nil {
			return "", false, errl
		}
		if errl != nil {
			return r, nr, nil
		}
		if errr != nil {
			return l, nl, nil
		}
		return "(" + l + "|" + r + ")", nl || nr, nil
	}
	if p.Group != nil {
		l, nl, errl := listToRegex(p.Group.Left)
		r, nr, errr := listToRegex(p.Group.Right)
		var err error = nil
		if errl != nil {
			err = errl
		}
		if errr != nil {
			err = errr
		}
		return l + `\s` + r, nl && nr, err
	}
	panic(fmt.Sprintf("unsupported list %v", p))
}
