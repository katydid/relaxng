package relaxng

import (
	"strings"
	"testing"
)

func TestSimpleParse(t *testing.T) {
	//http://relaxng.org/spec-20011203.html
	example5p1 := `<?xml version="1.0"?>
<grammar xmlns="http://relaxng.org/ns/structure/1.0">
  <start>
    <ref name="foo.element"/>
  </start>

  <define name="foo.element">
    <element>
      <name ns="">foo</name>
      <group>
        <ref name="bar1.element"/>
        <ref name="bar2.element"/>
      </group>
    </element>
  </define>

  <define name="bar1.element">
    <element>
      <name ns="http://www.example.com/n1">bar1</name>
      <empty/>
    </element>
  </define>

  <define name="bar2.element">
    <element>
      <name ns="http://www.example.com/n2">bar2</name>
      <empty/>
    </element>
  </define>
</grammar>`
	g, err := ParseGrammar([]byte(example5p1))
	if err != nil {
		panic(err)
	}
	s := g.String()
	t.Logf("%s", s)
	if !strings.Contains(s, `ref name="bar2.element"`) {
		t.Fatalf("expected ref name bar2.element in group")
	}
	if !strings.Contains(s, `ref name="bar1.element"`) {
		t.Fatalf("expected ref name bar1.element in group")
	}
	if !strings.Contains(s, `<define name="bar1.element">`) {
		t.Fatalf("expected define bar1")
	}
}
