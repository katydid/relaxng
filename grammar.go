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
	NotAllowed *Empty         `xml:"notAllowed"`
	Empty      *Empty         `xml:"empty"`
	Text       *Empty         `xml:"text"`
	Data       *Data          `xml:"data"`
	Value      *Value         `xml:"value"`
	List       *NameOrPattern `xml:"list"`
	Attribute  *NameOrPattern `xml:"attribute"`
	Ref        *Ref           `xml:"ref"`
	OneOrMore  *NameOrPattern `xml:"oneOrMore"`
	Choice     *Pair          `xml:"choice"`
	Group      *Pair          `xml:"group"`
	Interleave *Pair          `xml:"interleave"`

	AnyName *AnyNameClass  `xml:"anyName"`
	NsName  *NsNameClass   `xml:"nsName"`
	Name    *NameNameClass `xml:"name"`
}

// func (this *NameOrPattern) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	name := start.Name.Local
// 	_, isNameOrPattern := nameOrPattern[name]
// 	if isNameOrPattern {
// 		return this.unmarshalXML(d, start)
// 	}
// 	var s xml.StartElement
// 	var ok bool
// 	for !ok {
// 		t, err := d.Token()
// 		if err != nil {
// 			return err
// 		}
// 		s, ok = t.(xml.StartElement)
// 	}
// 	if err := this.unmarshalXML(d, s); err != nil {
// 		return err
// 	}
// 	return d.Skip()
// }

var nameOrPattern = map[string]struct{}{
	"notAllowed": struct{}{},
	"empty":      struct{}{},
	"text":       struct{}{},
	"data":       struct{}{},
	"value":      struct{}{},
	"list":       struct{}{},
	"attribute":  struct{}{},
	"ref":        struct{}{},
	"oneOrMore":  struct{}{},
	"choice":     struct{}{},
	"group":      struct{}{},
	"interleave": struct{}{},
	"anyName":    struct{}{},
	"nsName":     struct{}{},
	"name":       struct{}{},
}

func (this *NameOrPattern) unmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "notAllowed":
		this.NotAllowed = &Empty{}
		return d.DecodeElement(this.NotAllowed, &start)
	case "empty":
		this.Empty = &Empty{}
		return d.DecodeElement(this.Empty, &start)
	case "text":
		this.Text = &Empty{}
		return d.DecodeElement(this.Text, &start)
	case "data":
		this.Data = &Data{}
		return d.DecodeElement(this.Data, &start)
	case "value":
		this.Value = &Value{}
		return d.DecodeElement(this.Value, &start)
	case "list":
		this.List = &NameOrPattern{}
		return d.DecodeElement(this.List, &start)
	case "attribute":
		this.Attribute = &NameOrPattern{}
		return d.DecodeElement(this.Attribute, &start)
	case "ref":
		this.Ref = &Ref{}
		return d.DecodeElement(this.Ref, &start)
	case "oneOrMore":
		this.OneOrMore = &NameOrPattern{}
		return d.DecodeElement(this.OneOrMore, &start)
	case "choice":
		this.Choice = &Pair{}
		return d.DecodeElement(this.Choice, &start)
	case "group":
		this.Group = &Pair{}
		return d.DecodeElement(this.Group, &start)
	case "interleave":
		this.Interleave = &Pair{}
		return d.DecodeElement(this.Interleave, &start)
	case "anyName":
		this.AnyName = &AnyNameClass{}
		return d.DecodeElement(this.AnyName, &start)
	case "nsName":
		this.NsName = &NsNameClass{}
		return d.DecodeElement(this.NsName, &start)
	case "name":
		this.Name = &NameNameClass{}
		return d.DecodeElement(this.Name, &start)
	}
	return fmt.Errorf("unknown pattern " + start.Name.Local)
}

type Empty struct{}

type Data struct {
	Type            string         `xml:"type,attr"`
	DatatypeLibrary string         `xml:"datatypeLibrary,attr"`
	Param           []Param        `xml:"param"`
	Except          *NameOrPattern `xml:"except"`
}

// func (this *Data) unmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	switch start.Name.Local {
// 	case "param":
// 		p := &Param{}
// 		if err := d.DecodeElement(p, &start); err != nil {
// 			return err
// 		}
// 		this.Param = append(this.Param, *p)
// 	case "except":
// 		this.Except = &NameOrPattern{}
// 		return d.DecodeElement(this.Except, &start)
// 	}
// 	panic("unknown data " + start.Name.Local)
// }

// func (this *Data) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	name := start.Name.Local
// 	if name != "data" {
// 		panic("wtf")
// 	}
// 	for _, a := range start.Attr {
// 		if a.Name.Local == "type" {
// 			this.Type = a.Value
// 		} else if a.Name.Local == "datatypeLibrary" {
// 			this.DatatypeLibrary = a.Value
// 		}
// 	}
// 	for {
// 		t, err := d.Token()
// 		if err != nil {
// 			if err == io.EOF {
// 				return nil
// 			}
// 			return err
// 		}
// 		s, ok := t.(xml.StartElement)
// 		if !ok {
// 			continue
// 		}
// 		if err := this.unmarshalXML(d, s); err != nil {
// 			return err
// 		}
// 	}
// 	panic("unreachable")
// }

type Value struct {
	DatatypeLibrary string `xml:"datatypeLibrary,attr"`
	Type            string `xml:"type,attr"`
	Ns              string `xml:"ns,attr"`
	Text            string `xml:",chardata"`
}

type Ref struct {
	Name string `xml:"name,attr"`
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

type Param struct {
	Name string `xml:",attr"`
	Text string `xml:",chardata"`
}

type AnyNameClass struct {
	Except *ExceptNameClass `xml:"except"`
}

type NsNameClass struct {
	Ns     string           `xml:"ns,attr"`
	Except *ExceptNameClass `xml:"except"`
}

type NameNameClass struct {
	Ns   string `xml:"ns,attr"`
	Text string `xml:",chardata"`
}

type ExceptNameClass struct {
	Text string `xml:",chardata"`
}
