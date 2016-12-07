package relaxng

import (
	"testing"
)

func TestRemoveTODO(t *testing.T) {
	g := &Grammar{
		Start: &NameOrPattern{
			Value: &Value{
				Ns:   "TODO",
				Text: "Hello",
			},
		},
		Define: []Define{
			{
				Element: Pair{
					Right: &NameOrPattern{
						Value: &Value{
							Ns:   "TODO",
							Text: "Hello",
						},
					},
				},
			},
		},
	}
	RemoveTODOs(g)
	if g.Start.Value.Ns != "" {
		t.Fatalf("Ns not cleared")
	}
	if g.Start.Value.Text != "Hello" {
		t.Fatalf("Hello is where?")
	}
	if g.Define[0].Element.Right.Value.Ns != "" {
		t.Fatalf("Define Ns not cleared")
	}
}
