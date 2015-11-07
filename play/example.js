var relaxngDefault = `<grammar>
    <start>
        <ref name="element1"></ref>
    </start>
    <define name="element1">
        <element>
            <name>Whats</name>
            <attribute>
                <name>up</name>
                <text></text>
            </attribute>
        </element>
    </define>
</grammar>`;

var xmlDefault = `<Whats up="E"/>`;

var defaults = {
	"relaxng": relaxngDefault,
	"xml": xmlDefault
}
