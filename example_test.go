package relaxng

import (
	"fmt"
)

//Parses the simplified RelaxNG Grammar,
//translates it to relapse and then validates the input xml.
func ExampleValidate() {
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
}
