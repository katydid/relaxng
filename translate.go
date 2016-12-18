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
	"strings"

	"github.com/katydid/katydid/relapse/ast"
)

func translate(g *Grammar) (*ast.Grammar, error) {
	refs := make(ast.RefLookup)
	refs["main"] = translatePattern(g.Start, false)
	for _, d := range g.Define {
		pattern := translatePattern(d.Element.Right, false)
		if !hasNsName(d.Element.Left) {
			pattern = addXmlns(pattern)
		}
		pattern = newTreeNode(d.Element.Left, pattern)
		pattern = ast.NewInterleave(pattern,
			ast.NewZeroOrMore(newWhitespace()),
		)
		refs[d.Name] = pattern

	}
	gg := ast.NewGrammar(refs)
	gg.Format()
	return gg, nil
}

func addXmlns(p *ast.Pattern) *ast.Pattern {
	return ast.NewConcat(
		ast.NewOptional(
			ast.NewTreeNode(ast.NewStringName("attr_xmlns"), ast.NewZAny()),
		),
		p,
	)
}

func hasNsName(n *NameOrPattern) bool {
	if n.Choice != nil {
		return hasNsName(n.Choice.Left) || hasNsName(n.Choice.Right)
	}
	if n.NsName != nil {
		return true
	}
	if n.Name != nil {
		return true
	}
	return false
}

func hasAttr(p *NameOrPattern) bool {
	if p.NotAllowed != nil ||
		p.Empty != nil ||
		p.Text != nil ||
		p.Data != nil ||
		p.Value != nil ||
		p.List != nil {
		return false
	}
	if p.Attribute != nil {
		return true
	}
	if p.Ref != nil {
		return false //TODO
	}
	if p.OneOrMore != nil {
		return hasAttr(p.OneOrMore.NameOrPattern)
	}
	if p.Choice != nil {
		return hasAttr(p.Choice.Left) || hasAttr(p.Choice.Right)
	}
	if p.Group != nil {
		return hasAttr(p.Group.Left) || hasAttr(p.Group.Right)
	}
	if p.Interleave != nil {
		return hasAttr(p.Interleave.Left) || hasAttr(p.Interleave.Right)
	}
	panic(fmt.Sprintf("unreachable pattern %v", p))
}

func translatePattern(p *NameOrPattern, attr bool) *ast.Pattern {
	if p.NotAllowed != nil {
		return ast.NewNot(ast.NewZAny())
	}
	if p.Empty != nil {
		if attr {
			return newEmptyValue()
		}
		return ast.NewOr(
			ast.NewEmpty(),
			newEmptyValue(),
		)
	}
	if p.Text != nil {
		return ast.NewZeroOrMore(newTextValue())
	}
	if p.Data != nil {
		if len(p.Data.DatatypeLibrary) > 0 {
			panic("data datatypeLibrary not supported")
		}
		if p.Data.Except == nil {
			return ast.NewOr(newTextValue(), ast.NewEmpty())
		}
		expr, nullable := translateLeaf(p.Data.Except)
		v := ast.NewAnd(
			newTextValue(),
			ast.NewNot(expr),
		)
		if nullable {
			return ast.NewAnd(v, ast.NewNot(ast.NewEmpty()))
		}
		return ast.NewOr(v, ast.NewEmpty())
	}
	if p.Value != nil {
		v, nullable := translateLeaf(p)
		if !nullable {
			return v
		}
		return ast.NewOr(v, ast.NewEmpty())
	}
	if p.List != nil {
		return newList(p.List.NameOrPattern)
	}
	if p.Attribute != nil {
		nameExpr := translateNameClass(p.Attribute.Left, true)
		pattern := translatePattern(p.Attribute.Right, true)
		return ast.NewTreeNode(nameExpr, pattern)
	}
	if p.Ref != nil {
		return ast.NewReference(p.Ref.Name)
	}
	if p.OneOrMore != nil {
		inside := translatePattern(p.OneOrMore.NameOrPattern, attr)
		return ast.NewConcat(inside, ast.NewZeroOrMore(inside))
	}
	if p.Choice != nil {
		left := translatePattern(p.Choice.Left, attr)
		right := translatePattern(p.Choice.Right, attr)
		return ast.NewOr(left, right)
	}
	if p.Group != nil {
		left := translatePattern(p.Group.Left, attr)
		right := translatePattern(p.Group.Right, attr)
		if attr {
			return ast.NewInterleave(left, right)
		}
		if hasAttr(p.Group.Right) {
			if hasAttr(p.Group.Left) {
				return ast.NewInterleave(left, right)
			}
			return ast.NewConcat(right, left)
		}
		return ast.NewConcat(left, right)
	}
	if p.Interleave != nil {
		left := translatePattern(p.Interleave.Left, attr)
		right := translatePattern(p.Interleave.Right, attr)
		return ast.NewInterleave(left, right)
	}
	panic(fmt.Sprintf("unreachable pattern %v", p))
}

func newTreeNode(n *NameOrPattern, pattern *ast.Pattern) *ast.Pattern {
	if n.Choice != nil {
		return ast.NewOr(
			newTreeNode(n.Choice.Left, pattern),
			newTreeNode(n.Choice.Right, pattern),
		)
	}
	if n.AnyName != nil {
		if n.AnyName.Except == nil {
			return ast.NewTreeNode(
				ast.NewAnyName(),
				pattern,
			)
		}
		except := translateNameClass(n.AnyName.Except, false)
		return ast.NewTreeNode(ast.NewAnyNameExcept(except), pattern)
	}
	if n.NsName != nil {
		// if len(n.NsName.Ns) == 0 {
		// 	return ast.NewTreeNode(ast.NewAnyName(), pattern)
		// }
		panic("nsName is not supported")
	}
	if n.Name != nil {
		if len(n.Name.Ns) > 0 {
			return ast.NewTreeNode(ast.NewStringName("elem_"+n.Name.Text), ast.NewConcat(
				ast.NewTreeNode(ast.NewStringName("attr_xmlns"),
					newValue(n.Name.Ns),
				),
				pattern,
			))
		}
		return ast.NewTreeNode(ast.NewStringName("elem_"+n.Name.Text), pattern)
	}
	panic(fmt.Sprintf("unreachable nameclass %v", n))
}

func translateNameClass(n *NameOrPattern, attr bool) *ast.NameExpr {
	if n.Choice != nil {
		return ast.NewNameChoice(
			translateNameClass(n.Choice.Left, attr),
			translateNameClass(n.Choice.Right, attr),
		)
	}
	if n.AnyName != nil {
		if n.AnyName.Except == nil {
			return ast.NewAnyName()
		}
		except := translateNameClass(n.AnyName.Except, attr)
		return ast.NewAnyNameExcept(except)
	}
	if n.NsName != nil {
		// if len(n.NsName.Ns) == 0 {
		// 	if n.NsName.Except != nil {
		// 		return ast.NewAnyNameExcept(translateNameClass(n.NsName.Except, attr))
		// 	} else {
		// 		return ast.NewAnyName()
		// 	}
		// }
		panic("nsName is not supported")
	}
	if n.Name != nil {
		if len(n.Name.Ns) > 0 {
			panic(fmt.Sprintf("name ns <%v> is not supported", n.Name.Ns))
		}
		if attr {
			return ast.NewStringName("attr_" + n.Name.Text)
		}
		return ast.NewStringName("elem_" + n.Name.Text)
	}
	panic(fmt.Sprintf("unreachable nameclass %v", n))
}

func translateLeaf(p *NameOrPattern) (*ast.Pattern, bool) {
	if p.Value != nil {
		if len(p.Value.Ns) > 0 {
			panic("value ns not supported")
		}
		text := p.Value.Text
		if p.Value.IsString() {
			return newValue(text), len(text) == 0
		}
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		text = strings.Replace(text, "\t", "", -1)
		text = strings.TrimSpace(text)
		return newToken(text), len(text) == 0
	}
	if p.Choice != nil {
		l, nl := translateLeaf(p.Choice.Left)
		r, nr := translateLeaf(p.Choice.Right)
		return ast.NewOr(l, r), nl || nr
	}
	panic(fmt.Sprintf("unsupported leaf %v", p))
}
