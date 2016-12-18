package relaxng

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenJSSuite(t *testing.T) {
	tests := scanFiles()
	fmt.Fprintf(os.Stderr, "var tests = [\n")
	for _, test := range tests {
		tbase := filepath.Base(filepath.Dir(test.Filename))
		for _, c := range test.Xmls {
			cbases := strings.Split(filepath.Base(c.Filename), ".")
			cbase := strings.Join(cbases[:len(cbases)-1], ".")
			fmt.Fprintf(os.Stderr, "'%s.%s',\n", tbase, cbase)
		}
	}
	fmt.Fprintf(os.Stderr, "]\n")
}
