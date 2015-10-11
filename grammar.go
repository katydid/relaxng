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
	"fmt"
	"reflect"
)

func ParseGrammar(buf []byte) (*Grammar, error) {
	g := &Grammar{}
	err := xml.Unmarshal(buf, g)
	return g, err
}

func (g *Grammar) String() string {
	data, err := xml.MarshalIndent(g, "", "\t")
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
	Start  *NameOrPattern `xml:"start"`
	Define []Define       `xml:"define"`
}

type Define struct {
	Name    string         `xml:"name,attr"`
	Element *NameOrPattern `xml:"element"`
}

type NameOrPattern struct {
	NotAllowed *NotAllowed `xml:"notAllowed"`
	Empty      *Empty      `xml:"empty"`
	Text       *Text       `xml:"text"`
	Data       *Data       `xml:"data"`
	Value      *Value      `xml:"value"`
	List       *List       `xml:"list"`
	Attribute  *Attribute  `xml:"attribute"`
	Ref        *Ref        `xml:"ref"`
	OneOrMore  *OneOrMore  `xml:"oneOrMore"`
	Choice     *Pair       `xml:"choice"`
	Group      *Pair       `xml:"group"`
	Interleave *Pair       `xml:"interleave"`

	AnyName *AnyNameClass  `xml:"anyName"`
	NsName  *NsNameClass   `xml:"nsName"`
	Name    *NameNameClass `xml:"name"`
}

func (this *NameOrPattern) unmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t := reflect.TypeOf(this).Elem()
	numFields := t.NumField()
	v := reflect.ValueOf(this).Elem()
	for i := 0; i < numFields; i++ {
		f := t.Field(i)
		xmlTag := f.Tag.Get("xml")
		if xmlTag == start.Name.Local {
			n := reflect.New(f.Type)
			err := d.DecodeElement(n.Interface(), &start)
			if err != nil {
				return err
			}
			v.Field(i).Set(n.Elem())
			return nil
		}
	}
	return fmt.Errorf("unknown pattern " + start.Name.Local)
}

func (this *NameOrPattern) marshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := reflect.ValueOf(this).Elem()
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		if !v.Field(i).IsNil() {
			return e.Encode(v.Field(i).Interface())
		}
	}
	return fmt.Errorf("unset pattern")
}

type NotAllowed struct {
	XMLName xml.Name `xml:"notAllowed"`
}

type Empty struct {
	XMLName xml.Name `xml:"empty"`
}

type Text struct {
	XMLName xml.Name `xml:"text"`
}

type Data struct {
	XMLName         xml.Name       `xml:"data"`
	Type            string         `xml:"type,attr"`
	DatatypeLibrary string         `xml:"datatypeLibrary,attr"`
	Param           []Param        `xml:"param"`
	Except          *NameOrPattern `xml:"except"`
}

type Value struct {
	XMLName         xml.Name `xml:"value"`
	DatatypeLibrary string   `xml:"datatypeLibrary,attr"`
	Type            string   `xml:"type,attr"`
	Ns              string   `xml:"ns,attr"`
	Text            string   `xml:",chardata"`
}

type List struct {
	XMLName xml.Name `xml:"list"`
	*NameOrPattern
}

type Attribute struct {
	XMLName xml.Name `xml:"attribute"`
	*NameOrPattern
}

type OneOrMore struct {
	XMLName xml.Name `xml:"oneOrMore"`
	*NameOrPattern
}

type Ref struct {
	XMLName xml.Name `xml:"ref"`
	Name    string   `xml:"name,attr"`
}

type Pair struct {
	Left  *NameOrPattern
	Right *NameOrPattern
}

func skipToStart(d *xml.Decoder) (*xml.StartElement, error) {
	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		s, ok := t.(xml.StartElement)
		if ok {
			return &s, nil
		}
	}
}

func (this *Pair) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fmt.Printf("unmarshaling pair %s\n", start.Name.Local)
	s, err := skipToStart(d)
	if err != nil {
		return err
	}
	fmt.Printf("\t %s\n", s.Name.Local)
	this.Left = &NameOrPattern{}
	if err := this.Left.unmarshalXML(d, *s); err != nil {
		return err
	}
	s, err = skipToStart(d)
	if err != nil {
		return err
	}
	fmt.Printf("\t %s\n", s.Name.Local)
	this.Right = &NameOrPattern{}
	if err := this.Right.unmarshalXML(d, *s); err != nil {
		return err
	}
	return d.Skip()
}

func (this *Pair) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	fmt.Printf("marshaling pair\n")
	if err := this.Left.marshalXML(e, start); err != nil {
		return err
	}
	if err := this.Right.marshalXML(e, start); err != nil {
		return err
	}
	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return nil
}

type Param struct {
	Name string `xml:",attr"`
	Text string `xml:",chardata"`
}

type AnyNameClass struct {
	XMLName xml.Name         `xml:"anyName"`
	Except  *ExceptNameClass `xml:"except"`
}

type NsNameClass struct {
	XMLName xml.Name         `xml:"nsName"`
	Ns      string           `xml:"ns,attr"`
	Except  *ExceptNameClass `xml:"except"`
}

type NameNameClass struct {
	XMLName xml.Name `xml:"name"`
	Ns      string   `xml:"ns,attr"`
	Text    string   `xml:",chardata"`
}

type ExceptNameClass struct {
	Text string `xml:",chardata"`
}
