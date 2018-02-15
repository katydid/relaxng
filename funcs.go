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
	"strconv"
	"strings"

	"github.com/katydid/katydid/relapse/ast"
	c "github.com/katydid/katydid/relapse/combinator"
	"github.com/katydid/katydid/relapse/funcs"
)

func stripTextPrefix(s string) (string, error) {
	if !strings.HasPrefix(s, "text_") {
		return "", fmt.Errorf("%q is not of type text", s)
	}
	return strings.Replace(s, "text_", "", 1), nil
}

func newTokenValue(t string) *ast.Pattern {
	return c.Value(ast.NewFunction("token", c.StringVar(), c.StringConst(t)))
}

// token is a function used in relapse to validate values as described here
// http://books.xmlschemata.org/relaxng/relax-CHP-7-SECT-4.html
type token struct {
	S           funcs.String
	c           string
	hash        uint64
	hasVariable bool
}

func Token(S funcs.String, C funcs.ConstString) (funcs.Bool, error) {
	c, err := C.Eval()
	if err != nil {
		return nil, err
	}
	return funcs.TrimBool(&token{
		S:           S,
		c:           c,
		hash:        funcs.Hash("token", C, S),
		hasVariable: S.HasVariable(),
	}), nil
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

func (this *token) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*token); ok {
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		if c := strings.Compare(this.c, other.c); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func (this *token) HasVariable() bool {
	return this.hasVariable
}

func (this *token) String() string {
	return "token(" + this.S.String() + "," + strconv.Quote(this.c) + ")"
}

func (this *token) Hash() uint64 {
	return this.hash
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
	funcs.Register("token", Token)
}

type whitespace struct {
	S           funcs.String
	hash        uint64
	hasVariable bool
}

func Whitespace(S funcs.String) funcs.Bool {
	return funcs.TrimBool(&whitespace{
		S:           S,
		hash:        funcs.Hash("whitespace", S),
		hasVariable: S.HasVariable(),
	})
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

func (this *whitespace) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*whitespace); ok {
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func (this *whitespace) HasVariable() bool {
	return this.hasVariable
}

func (this *whitespace) String() string {
	return "whitespace(" + this.S.String() + ")"
}

func (this *whitespace) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("whitespace", Whitespace)
}

type anytext struct {
	S           funcs.String
	hash        uint64
	hasVariable bool
}

func AnyText(S funcs.String) funcs.Bool {
	return funcs.TrimBool(&anytext{
		S:           S,
		hash:        funcs.Hash("anytext", S),
		hasVariable: S.HasVariable(),
	})
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

func (this *anytext) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*anytext); ok {
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func (this *anytext) HasVariable() bool {
	return this.hasVariable
}

func (this *anytext) String() string {
	return "anytext(" + this.S.String() + ")"
}

func (this *anytext) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("anytext", AnyText)
}

func newWhitespace() *ast.Pattern {
	return c.Value(ast.NewFunction("whitespace", c.StringVar()))
}

func newEmptyValue() *ast.Pattern {
	return c.Value(ast.NewFunction("whitespace", c.StringVar()))
}

func newAnyValue() *ast.Pattern {
	return c.Value(ast.NewFunction("anytext", c.StringVar()))
}

func newTextValue(value string) *ast.Pattern {
	return c.Value(ast.NewFunction("text", c.StringVar(), c.StringConst(value)))
}

type text struct {
	S           funcs.String
	c           string
	hash        uint64
	hasVariable bool
}

func TextFunc(S funcs.String, C funcs.ConstString) (funcs.Bool, error) {
	c, err := C.Eval()
	if err != nil {
		return nil, err
	}
	return funcs.TrimBool(&text{
		S:           S,
		c:           c,
		hash:        funcs.Hash("token", C, S),
		hasVariable: S.HasVariable(),
	}), nil
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

func (this *text) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*text); ok {
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		if c := strings.Compare(this.c, other.c); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func (this *text) HasVariable() bool {
	return this.hasVariable
}

func (this *text) String() string {
	return "text(" + this.S.String() + "," + strconv.Quote(this.c) + ")"
}

func (this *text) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("text", TextFunc)
}

type list struct {
	r           *regexp.Regexp
	S           funcs.String
	Expr        funcs.ConstString
	hash        uint64
	hasVariable bool
}

func ListFunc(S funcs.String, Expr funcs.ConstString) (funcs.Bool, error) {
	e, err := Expr.Eval()
	if err != nil {
		return nil, err
	}
	r, err := regexp.Compile(e)
	if err != nil {
		return nil, err
	}
	return funcs.TrimBool(&list{
		S:           S,
		r:           r,
		Expr:        Expr,
		hash:        funcs.Hash("list", S, Expr),
		hasVariable: S.HasVariable(),
	}), nil
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

func (this *list) Compare(that funcs.Comparable) int {
	if this.Hash() != that.Hash() {
		if this.Hash() < that.Hash() {
			return -1
		}
		return 1
	}
	if other, ok := that.(*list); ok {
		if c := this.S.Compare(other.S); c != 0 {
			return c
		}
		if c := this.Expr.Compare(other.Expr); c != 0 {
			return c
		}
		return 0
	}
	return strings.Compare(this.String(), that.String())
}

func (this *list) HasVariable() bool {
	return this.hasVariable
}

func (this *list) String() string {
	return "list(" + this.S.String() + "," + this.Expr.String() + ")"
}

func (this *list) Hash() uint64 {
	return this.hash
}

func init() {
	funcs.Register("list", ListFunc)
}

func newList(nameOrPattern *NameOrPattern) *ast.Pattern {
	regexStr, nullable, err := listToRegex(nameOrPattern)
	if err != nil {
		return ast.NewNot(ast.NewZAny())
	}
	val := c.Value(ast.NewFunction("list", c.StringVar(), c.StringConst("^"+regexStr+"$")))
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
