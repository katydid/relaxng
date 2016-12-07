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
	"fmt"
	sdebug "github.com/katydid/katydid/parser/debug"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/interp"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
)

type testCase struct {
	Filename       string
	Content        []byte
	SimpleFilename string
	SimpleContent  []byte
	Xmls           []xmlCase
}

func (this testCase) expectError() bool {
	return strings.HasSuffix(this.Filename, "i.rng")
}

type xmlCase struct {
	Filename string
	Content  []byte
}

func (this xmlCase) expectError() bool {
	return strings.HasSuffix(this.Filename, "i.xml")
}

type testSuite []testCase

func scanFiles() testSuite {
	cases := make(map[int]testCase)
	if err := filepath.Walk("./RelaxTestSuite/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		extension := filepath.Ext(path)
		if !(extension == ".rng" || extension == ".xml") {
			return nil
		}
		number, err := strconv.Atoi(filepath.Base(filepath.Dir(path)))
		if err != nil {
			return err
		}
		c := cases[number]
		if extension == ".rng" {
			if strings.HasSuffix(path, "s.rng") {
				c.SimpleFilename = path
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				c.SimpleContent = data
			} else {
				c.Filename = path
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				c.Content = data
			}
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			c.Xmls = append(c.Xmls, xmlCase{
				Filename: path,
				Content:  data,
			})
		}
		cases[number] = c
		return nil
	}); err != nil {
		panic(err)
	}
	num := len(cases) + 1
	suite := make(testSuite, len(cases))
	for i := 1; i < num; i++ {
		if c, ok := cases[i]; ok {
			suite[i-1] = c
		} else {
			panic(fmt.Sprintf("missing test %d", i))
		}
	}
	return suite
}

func debugValidate(katydid *ast.Grammar, xmlContent []byte) error {
	p := NewXMLParser()
	if err := p.Init(xmlContent); err != nil {
		return err
	}
	d := sdebug.NewLogger(p, sdebug.NewLineLogger())
	valid, err := interp.Interpret(katydid, d)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("not valid")
	}
	return nil
}

func testSimple(t *testing.T, spec testCase, debugParser bool) bool {
	debugStr := fmt.Sprintf("Original:\n%s\n", string(spec.SimpleContent))
	defer func() {
		r := recover()
		if r != nil {
			t.Fatalf("%srecover for %s: %v: %s", debugStr, spec.SimpleFilename, r, debug.Stack())
		}
	}()
	g, err := ParseGrammar(spec.SimpleContent)
	if err != nil {
		t.Fatalf("%sunparsable %s", debugStr, spec.SimpleFilename)
	}
	RemoveTODOs(g)
	debugStr += fmt.Sprintf("Parsed:\n%s\n", g.String())
	katydid, err := Translate(g)
	if err != nil {
		t.Fatalf("%sunexpected error <%s> for %s", debugStr, err, spec.SimpleFilename)
	}
	debugStr += fmt.Sprintf("To:\n%s\n", katydid.String())
	preInputDebug := debugStr
	passed := true
	for _, xml := range spec.Xmls {
		if debugParser {
			fmt.Printf("--- PARSING %s\n", xml.Filename)
		}
		debugStr = preInputDebug + fmt.Sprintf("Input:\n%s\n", string(xml.Content))
		if debugParser {
			err = debugValidate(katydid, xml.Content)
		} else {
			err = Validate(katydid, xml.Content)
		}
		if xml.expectError() {
			if err == nil {
				t.Errorf("%sexpected error for %s", debugStr, xml.Filename)
				passed = false
				continue
			}
			return passed
		}
		if err != nil {
			t.Errorf("%sgot unexpected error <%s> for %s", debugStr, err, xml.Filename)
			passed = false
			continue
		}
	}
	return passed
}

var namespaces = map[string]bool{
	"122": true, //<name ns="http://www.example.com">
	"124": true, //<name ns="http://www.example.com">
	"128": true, //<name ns="http://www.example.com">
	"131": true, //<name ns="http://www.example.com">
	"132": true, //<name ns="http://www.w3.org/XML/1998/namespace">
	"217": true,
	"218": true,
	"219": true, //<nsName ns="http://www.example.com">
	"220": true, //<nsName ns="http://www.example.com">
	"221": true, //<nsName ns="http://www.example.com">
	"248": true,
	"353": true,
	"354": true, //<nsName ns="http://www.example.com/1">
}

var datatypeLibrary = map[string]bool{
	"261": true,
}

func testNumber(filename string) string {
	return filepath.Base(filepath.Dir(filename))
}

func TestSimpleSuite(t *testing.T) {
	suite := scanFiles()
	passed := 0
	incorrectSkipped := 0
	failed := 0
	namespacesSkipped := 0
	datatypeLibrarySkipped := 0
	for _, spec := range suite {
		num := testNumber(spec.Filename)
		t.Run(num, func(t *testing.T) {
			if len(spec.SimpleFilename) == 0 {
				incorrectSkipped++
				t.Skip("incorrect specification")
				return
			}
			if namespaces[num] {
				namespacesSkipped++
				t.Skip("namespaces not supported")
				return
			}
			if datatypeLibrary[num] {
				datatypeLibrarySkipped++
				t.Skip("datatypeLibrary not supported")
				return
			}
			if testSimple(t, spec, true) {
				passed++
			} else {
				failed++
				t.Fatal()
			}
		})
	}
	t.Logf("passed: %d, failed: %d, namespace tests skipped: %d, datatypeLibrary tests skipped: %d, incorrect grammars skipped: %d", passed, failed, namespacesSkipped, datatypeLibrarySkipped, incorrectSkipped)
}
