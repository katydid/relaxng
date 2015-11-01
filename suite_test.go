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
	"io/ioutil"
	"os"
	"path/filepath"
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

func testOneCase(t *testing.T, spec testCase) {
	katydid, err := Translate(spec.Content)
	if spec.expectError() {
		if err == nil {
			t.Errorf("expected error for %s", spec.Filename)
		}
		return
	}
	if err != nil {
		t.Errorf("unexpected error %s for %s", err, spec.Filename)
		return
	}
	for _, xml := range spec.Xmls {
		err = Validate(katydid, xml.Content)
		if xml.expectError() {
			if err == nil {
				t.Errorf("expected error for %s", xml.Filename)
			}
			return
		}
		if err != nil {
			t.Errorf("got unexpected error %s for %s", err, xml.Filename)
		}
	}
}

func testSimple(t *testing.T, spec testCase) {
	debug := fmt.Sprintf("Original:\n%s\n", string(spec.SimpleContent))
	defer func() {
		r := recover()
		if r != nil {
			t.Fatalf("%srecover for %s: %#v", debug, spec.SimpleFilename, r)
		}
	}()
	g, err := ParseGrammar(spec.SimpleContent)
	if err != nil {
		t.Fatalf("%sunparsable %s", debug, spec.SimpleFilename)
	}
	debug += fmt.Sprintf("Parsed:\n%s\n", g.String())
	katydid, err := Translate(spec.SimpleContent)
	if err != nil {
		t.Fatalf("%sunexpected error <%s> for %s", debug, err, spec.SimpleFilename)
	}
	debug += fmt.Sprintf("To:\n%s\n", katydid.String())
	preInputDebug := debug
	for _, xml := range spec.Xmls {
		debug = preInputDebug + fmt.Sprintf("Input:\n%s\n", string(xml.Content))
		err = Validate(katydid, xml.Content)
		if xml.expectError() {
			if err == nil {
				t.Fatalf("%sexpected error for %s", debug, xml.Filename)
			}
			return
		}
		if err != nil {
			t.Fatalf("%sgot unexpected error <%s> for %s", debug, err, xml.Filename)
		}
	}
}

var namespaces = map[string]bool{
	"050": true,
	"051": true,
	"052": true,
	"095": true,
	"099": true,
	"104": true,
	"110": true,
	"122": true,
	"123": true,
	"124": true,
	"125": true,
	"126": true,
	"127": true,
	"128": true,
	"130": true,
	"131": true,
	"132": true,
	"133": true,
	"176": true,
	"217": true,
	"218": true,
	"219": true,
	"220": true,
	"221": true,
	"222": true,
	"248": true,
	"254": true,
	"255": true,
	"256": true,
	"258": true,
	"259": true,
	"262": true,
	"263": true,
	"266": true,
	"267": true,
	"270": true,
	"271": true,
	"272": true,
	"273": true,
	"274": true,
	"275": true,
	"280": true,
	"353": true,
	"354": true,
}

var fixable = map[string]bool{
	"120": true,
	"139": true,
	"142": true,
	"146": true,
	"147": true,
	"151": true,
	"190": true,
	"191": true,
	"194": true,
	"195": true,
	"215": true,
	"225": true,
	"226": true,
	"228": true,
	"232": true,
	"234": true,
	"236": true,
	"237": true,
	"238": true,
	"244": true,
	"250": true,
	"251": true,
	"261": true,
	"264": true,
	"265": true,
	"268": true,
	"269": true,
	"284": true,
	"368": true,
	"369": true,
	"372": true,
}

func testNumber(filename string) string {
	return filepath.Base(filepath.Dir(filename))
}

func TestSimpleSuite(t *testing.T) {
	suite := scanFiles()
	passed := 0
	incorrect := 0
	for _, spec := range suite {
		num := testNumber(spec.Filename)
		if len(spec.SimpleFilename) == 0 {
			//skipping incorrect specifications
			incorrect++
			continue
		}
		if namespaces[num] {
			t.Logf("%s [SKIP] namespaces not supported", num)
			continue
		}
		if fixable[num] {
			t.Errorf("%s [FAIL]", num)
			continue
		}
		testSimple(t, spec)
		t.Logf("%s [PASS]", num)
		passed++
	}
	total := passed + len(fixable)
	t.Logf("passed: %d/%d, failed: %d/%d, namespace tests skipped: %d, incorrect grammars skipped: %d", passed, total, len(fixable), total, len(namespaces), incorrect)
}
