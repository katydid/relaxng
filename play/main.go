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

package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	. "github.com/katydid/relaxng"
)

func main() {
	js.Global.Set("gofunctions", map[string]interface{}{
		"ValidateRelaxNG":  ValidateRelaxNG,
		"TranslateRelaxNG": TranslateRelaxNG,
	})
}

func ValidateRelaxNG(relaxngStr, xmlStr string) string {
	v, err := validate(relaxngStr, xmlStr)
	if err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("%v", v)
}

func TranslateRelaxNG(relaxngStr string) string {
	v, err := translate(relaxngStr)
	if err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("%s", v)
}

func validate(relaxngStr, xmlStr string) (bool, error) {
	v := &validator{nil}
	b, err := v.validate(relaxngStr, xmlStr)
	if err != nil {
		return false, err
	}
	if v.err != nil {
		return false, err
	}
	return b, nil
}

func translate(relaxngStr string) (string, error) {
	v := &validator{nil}
	s, err := v.translate(relaxngStr)
	if err != nil {
		return "", err
	}
	if v.err != nil {
		return "", err
	}
	return s, nil
}

type validator struct {
	err error
}

func (this *validator) translate(relaxngStr string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			this.err = fmt.Errorf("%v", r)
		}
	}()
	g, err := ParseGrammar([]byte(relaxngStr))
	if err != nil {
		return "", err
	}
	katy, err := Translate(g)
	if err != nil {
		return "", err
	}
	return katy.String(), nil
}

func (this *validator) validate(relaxngStr, xmlStr string) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			this.err = fmt.Errorf("%v", r)
		}
	}()
	g, err := ParseGrammar([]byte(relaxngStr))
	if err != nil {
		return false, err
	}
	var m interface{}
	if err := xml.Unmarshal([]byte(xmlStr), &m); err != nil {
		return false, err
	}
	katy, err := Translate(g)
	if err != nil {
		return false, err
	}
	err = Validate(katy, []byte(xmlStr))
	if err != nil {
		return false, nil
	}
	return true, nil
}
