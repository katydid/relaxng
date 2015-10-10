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
	katydid, err := Translate(spec.SimpleContent)
	if err != nil {
		t.Errorf("unexpected error %s for %s", err, spec.SimpleFilename)
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

func TestSimpleSuite(t *testing.T) {
	suite := scanFiles()
	for _, spec := range suite {
		if len(spec.SimpleFilename) == 0 {
			//skipping incorrect specifications
			continue
		}
		testSimple(t, spec)
	}
}
