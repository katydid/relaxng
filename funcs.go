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
	"github.com/katydid/katydid/relapse/funcs"

	"fmt"
	"strings"
)

//Token is a function used in relapse to validate values as described here
//http://books.xmlschemata.org/relaxng/relax-CHP-7-SECT-4.html
func Token(s funcs.String) funcs.String {
	return &token{s}
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
	funcs.Register("stripTextPrefix", new(stripTextPrefix))
}

//TypeAndTrue is a function used in relapse to return true if the input bool is true and no error was returned.
func TypeAndTrue(b funcs.Bool) funcs.Bool {
	return &typeAndTrue{b}
}

type typeAndTrue struct {
	B funcs.Bool
}

func (this *typeAndTrue) Eval() (bool, error) {
	b, err := this.B.Eval()
	if err != nil {
		return false, nil
	}
	return b, nil
}

func init() {
	funcs.Register("typeAndTrue", new(typeAndTrue))
}
