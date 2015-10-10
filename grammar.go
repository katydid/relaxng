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
	"encoding/xml"
)

func ParseGrammar(buf []byte) (*Grammar, error) {
	g := &Grammar{}
	err := xml.Unmarshal(buf, g)
	return nil, err
}

func (g *Grammar) String() string {
	data, err := xml.Marshal(g)
	if err != nil {
		panic(err)
	}
	return string(data)
}

/*
The simplified RelaxNG Grammar as specified in
http://relaxng.org/spec-20011203.html
	grammar	  		::=  <grammar> <start> top </start> define* </grammar>
	define	  		::=  <define name="NCName"> <element> nameClass top </element> </define>
	top	  			::=  <notAllowed/>
					| Pattern
	Pattern	  		::=  <empty/>
					| Pattern
	Pattern	  		::=  <text/>
					| <data type="NCName" datatypeLibrary="anyURI"> param* [exceptPattern] </data>
					| <value datatypeLibrary="anyURI" type="NCName" ns="string"> string </value>
					| <list> Pattern </list>
					| <attribute> nameClass Pattern </attribute>
					| <ref name="NCName"/>
					| <oneOrMore> Pattern </oneOrMore>
					| <choice> Pattern Pattern </choice>
					| <group> Pattern Pattern </group>
					| <interleave> Pattern Pattern </interleave>
	param	  		::=  <param name="NCName"> string </param>
	exceptPattern	::=  <except> Pattern </except>
	nameClass	  	::=  <anyName> [exceptNameClass] </anyName>
					| <nsName ns="string"> [exceptNameClass] </nsName>
					| <name ns="string"> NCName </name>
					| <choice> nameClass nameClass </choice>
	exceptNameClass	::= <except> nameClass </except>
*/
type Grammar struct {
	Start  Start    `xml:"start"`
	Define []Define `xml:"define"`
}

type Start struct {
	*Top
}

type Define struct {
	Name    string   `xml:"name,attr"`
	Element *Element `xml:"element"`
}

type Element struct {
	*NameClass
	*Top
}

type Top struct {
	NotAllowed *NotAllowed `xml:"notAllowed"`
	*Pattern
}

type NotAllowed struct{}

type Pattern struct {
	Empty      *Empty      `xml:"empty"`
	Text       *Text       `xml:"text"`
	Data       *Data       `xml:"data"`
	Value      *Value      `xml:"value"`
	Attribute  *Attribute  `xml:"attribute"`
	Ref        *Ref        `xml:"ref"`
	OneOrMore  *OneOrMore  `xml:"oneOrMore"`
	Choice     *Choice     `xml:"choice"`
	Group      *Group      `xml:"group"`
	Interleave *Interleave `xml:"interleave"`
}

type Patterns struct {
	Empty      []*Empty      `xml:"empty"`
	Text       []*Text       `xml:"text"`
	Data       []*Data       `xml:"data"`
	Value      []*Value      `xml:"value"`
	Attribute  []*Attribute  `xml:"attribute"`
	Ref        []*Ref        `xml:"ref"`
	OneOrMore  []*OneOrMore  `xml:"oneOrMore"`
	Choice     []*Choice     `xml:"choice"`
	Group      []*Group      `xml:"group"`
	Interleave []*Interleave `xml:"interleave"`
}

type Empty struct{}

type Text struct{}

type Data struct {
	Type            string         `xml:"type,attr"`
	DatatypeLibrary string         `xml:"datatypeLibrary,attr"`
	Param           []Param        `xml:"param"`
	Except          *ExceptPattern `xml:"except"`
}

type Value struct {
	DatatypeLibrary string `xml:"datatypeLibrary,attr"`
	Type            string `xml:"type,attr"`
	Ns              string `xml:"ns,attr"`
	Text            string `xml:",chardata"`
}

type List struct {
	*Pattern
}

type Attribute struct {
	*NameClass
	*Pattern
}

type Ref struct {
	Name string `xml:"name,attr"`
}

type OneOrMore struct {
	*Pattern
}

type Choice []*Patterns

type Group []*Patterns

type Interleave []*Patterns

type Param struct {
	Name string `xml:",attr"`
	Text string `xml:",chardata"`
}

type ExceptPattern struct {
	*Pattern
}

type NameClass struct {
	AnyName *AnyNameClass    `xml:"anyName"`
	NsName  *NsNameClass     `xml:"nsName"`
	Name    *NameNameClass   `xml:"name"`
	Choice  *ChoiceNameClass `xml:"choice"`
}

type AnyNameClass struct {
	*ExceptNameClass
}

type NsNameClass struct {
	Ns string `xml:"ns,attr"`
	*ExceptNameClass
}

type NameNameClass struct {
	Ns   string `xml:"ns,attr"`
	Text string `xml:",chardata"`
}

type ChoiceNameClass []*NameClass

type ExceptNameClass struct {
	Text string `xml:",chardata"`
}
