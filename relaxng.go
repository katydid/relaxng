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

//Translates RelaxNG to Katydid Relapse
package relaxng

import (
	"errors"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/interp"
	"github.com/katydid/katydid/serialize/xml"
	"reflect"
)

//Translates a parsed RelaxNG Grammar into a Katydid Relapse Grammar.
func Translate(g *Grammar) (*relapse.Grammar, error) {
	return translate(g)
}

//The function removes the ns attributes with value TODO.
//These ns="TODO" attributes can become present
//after converting from RelaxNG to Simplified RelaxNG
//using rng2srng.jar
func RemoveTODOs(g *Grammar) {
	removeTODOs(reflect.ValueOf(g).Elem())
}

//Validates input xml against a Katydid Relapse Grammar.
func Validate(katydid *relapse.Grammar, xmlContent []byte) error {
	p := xml.NewXMLParser()
	if err := p.Init(xmlContent); err != nil {
		return err
	}
	if !interp.Interpret(katydid, p) {
		return errors.New("not valid")
	}
	return nil
}
