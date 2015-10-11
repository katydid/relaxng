#Converts RelaxNG to Katydid

Katydid still requires some work to fully support XML.
When this is done testing of the translations can start.

## Example

```
#element1

@element1 = foo: (
    [
        "@bar":<empty>,
        #element2
    ] |
    [
        "@bar":(->type($string))*,
        #element3
    ]
)

@element2 = baz1: <empty>

@element3 = baz2: <empty>
```

is translated from:

```
<Grammar>
    <start>
        <ref name="element1"></ref>
    </start>
    <define name="element1">
        <element>
            <name ns="">foo</name>
            <choice>
                <group>
                    <attribute>
                        <name ns="">bar</name>
                        <empty></empty>
                    </attribute>
                    <ref name="element2"></ref>
                </group>
                <group>
                    <attribute>
                        <name ns="">bar</name>
                        <text></text>
                    </attribute>
                    <ref name="element3"></ref>
                </group>
            </choice>
        </element>
    </define>
    <define name="element2">
        <element>
            <name ns="">baz1</name>
            <empty></empty>
        </element>
    </define>
    <define name="element3">
        <element>
            <name ns="">baz2</name>
            <empty></empty>
        </element>
    </define>
</Grammar>
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
