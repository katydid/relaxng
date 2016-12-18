#Translates RelaxNG to Katydid

Translates Simplified [RelaxNG](http://relaxng.org/) Grammars to the [Katydid](https://github.com/katydid/katydid) Relapse Grammar.

There is a playground available here http://katydid.github.io/relaxng/play/

Simplification is defined at these two resources:
  - http://debeissat.nicolas.free.fr/relaxng_simplification.php
  - http://books.xmlschemata.org/relaxng/relax-CHP-15-SECT-1.html

## Usage

```
simplifiedRelaxNG := `
<grammar>
    <start>
        <ref name="element1"></ref>
    </start>
    <define name="element1">
        <element>
            <name>foo</name>
            <empty></empty>
        </element>
    </define>
</grammar>`
relaxing, _ := ParseGrammar([]byte(simplifiedRelaxNG))
relapse, _ := Translate(relaxing)
input := "<foo/>"
if err := Validate(relapse, []byte(input)); err != nil {
    fmt.Println("invalid")
}
fmt.Println("valid")
// Output: valid
```

For more see the [go package documentation](https://godoc.org/github.com/katydid/relaxng)

### RelaxNG Test Suite

[![Build Status](https://travis-ci.org/katydid/relaxng.svg?branch=master)](https://travis-ci.org/katydid/relaxng)

```
passed: 146
failed: 0
namespace tests skipped: 13
datatypeLibrary tests skipped: 1
incorrect grammars skipped: 213
```

Steps:
 - The RelaxNG Test Suite is run through the simplifier which also eliminates all the incorrect grammars.
 - Next these simplified RelaxNG XML Grammars are parsed and translated to Katydid Relapse.
 - Finally the translated Relapse is used to validate the XML.

Viewing all tests can be done in the playground by going to this [link](http://katydid.github.io/relaxng/play/index.html?testsuite=049.1.v), which will load the first test, and then clicking the *NextTest* button.

### Example 1

The Simplified RelaxNG Grammar

```
<grammar>
    <start>
        <ref name="element1"></ref>
    </start>
    <define name="element1">
        <element>
            <name>foo</name>
            <empty></empty>
        </element>
    </define>
</grammar>
```

is translated to this Relapse Grammar

```
@element1
#element1={
    elem_foo:(<empty>|@ws);
    (@ws)*;
}
#ws=->whitespace($string)
#text=->anytext($string)
```


### Example 2

The Simplified RelaxNG Grammar

```
<grammar>
    <start>
        <ref name="element1"></ref>
    </start>
    <define name="element1">
        <element>
            <name>foo</name>
            <attribute>
                <name>bar</name>
                <text></text>
            </attribute>
        </element>
    </define>
</grammar>
```

is translated to this Relapse Grammar

```
@element1
#element1={
    elem_foo:attr_bar:(@text)*;
    (@ws)*;
}
#ws=->whitespace($string)
#text=->anytext($string)
```

### Example 3

The Simplified RelaxNG Grammar

```
<grammar>
    <start>
        <ref name="element1"></ref>
    </start>
    <define name="element1">
        <element>
            <name>foo</name>
            <choice>
                <group>
                    <attribute>
                        <name>bar</name>
                        <empty></empty>
                    </attribute>
                    <ref name="element2"></ref>
                </group>
                <group>
                    <attribute>
                        <name>bar</name>
                        <text></text>
                    </attribute>
                    <ref name="element3"></ref>
                </group>
            </choice>
        </element>
    </define>
    <define name="element2">
        <element>
            <name>baz1</name>
            <empty></empty>
        </element>
    </define>
    <define name="element3">
        <element>
            <name>baz2</name>
            <empty></empty>
        </element>
    </define>
</grammar>
```

is translated to this Relapse Grammar

```
@element1
#element1={
    elem_foo:(
        [
            attr_bar:@ws,
            @element2
        ] |
        [
            attr_bar:(@text)*,
            @element3
        ]
    );
    (@ws)*;
}
#element2={
    elem_baz1:(<empty>|@ws);
    (@ws)*;
}
#element3={
    elem_baz2:(<empty>|@ws);
    (@ws)*;
}
#ws=->whitespace($string)
#text=->anytext($string)
```

## Known Issues

There are quite a few known issues:
  - Only simplified grammars are supported.
  - namespaces are not supported.
  - datatypes: only string and token are currently supported.
  - datatypeLibraries are not supported.

I don't really intend to fix these, but you never know.

### Only handles simplified relaxng grammars.

http://www.kohsuke.org/relaxng/rng2srng/ seems to be quite effective at converting the full spectrum of what is possible within the relaxng grammar to the simplified grammar.

```
java -jar rng2srng.jar full.rng > simplified.rng
```
