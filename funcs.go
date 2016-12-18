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
	"strings"

	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/combinator"
	"github.com/katydid/katydid/relapse/funcs"
)

//Token is a function used in relapse to validate values as described here
//http://books.xmlschemata.org/relaxng/relax-CHP-7-SECT-4.html
func Token(s funcs.String) funcs.String {
	return &token{StripTextPrefix(s)}
}

type token struct {
	S funcs.String
}

func (this *token) Eval() (string, error) {
	s, err := this.S.Eval()
	if err != nil {
		return "", err
	}
	ss := tokenize(s)
	return strings.Join(ss, " "), nil
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

//StripTextPrefix is a function used in relapse to remove the text prefix added by the xml parser.
func StripTextPrefix(s funcs.String) funcs.String {
	return &stripTextPrefix{s}
}

type stripTextPrefix struct {
	S funcs.String
}

func (this *stripTextPrefix) Eval() (string, error) {
	s, err := this.S.Eval()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(s, "text_") {
		return strings.Replace(s, "text_", "", 1), nil
	}
	return "", fmt.Errorf("%q is not of type text", s)
}

func init() {
	funcs.Register("text", new(stripTextPrefix))
}

//Assert is a function used in relapse to return true if the input bool is true and no error was returned.
func Assert(b funcs.Bool) funcs.Bool {
	return &assert{b}
}

type assert struct {
	B funcs.Bool
}

func (this *assert) Eval() (bool, error) {
	b, err := this.B.Eval()
	if err != nil {
		return false, nil
	}
	return b, nil
}

func init() {
	funcs.Register("assert", new(assert))
}

func newWhitespace() *ast.Pattern {
	return combinator.Value(Assert(funcs.Regex(funcs.StringConst("^(\\s)+$"), StripTextPrefix(funcs.StringVar()))))
}

func newEmptyValue() *ast.Pattern {
	return combinator.Value(Assert(funcs.StringEq(Token(funcs.StringVar()), funcs.StringConst(""))))
}

func newTextValue() *ast.Pattern {
	return combinator.Value(funcs.TypeString(StripTextPrefix(funcs.StringVar())))
}

func newValue(value string) *ast.Pattern {
	return combinator.Value(Assert(
		funcs.StringEq(
			StripTextPrefix(funcs.StringVar()),
			funcs.StringConst(value),
		),
	))
}

func newToken(token string) *ast.Pattern {
	return combinator.Value(Assert(
		funcs.StringEq(
			Token(funcs.StringVar()),
			funcs.StringConst(token),
		),
	))
}

func newList(nameOrPattern *NameOrPattern) *ast.Pattern {
	regexStr, nullable, err := listToRegex(nameOrPattern)
	if err != nil {
		return ast.NewNot(ast.NewZAny())
	}
	val := combinator.Value(Assert(funcs.Regex(funcs.StringConst("^"+regexStr+"$"), Token(funcs.StringVar()))))
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
