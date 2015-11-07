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
	"github.com/katydid/katydid/funcs"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/combinator"
	"strings"
)

func translate(g *Grammar) (*relapse.Grammar, error) {
	refs := make(relapse.RefLookup)
	refs["main"] = translatePattern(g.Start, false)
	for _, d := range g.Define {
		pattern := translatePattern(d.Element.Right, false)
		if !hasNsName(d.Element.Left) {
			pattern = addXmlns(pattern)
		}
		pattern = newTreeNode(d.Element.Left, pattern)
		refs[d.Name] = pattern

	}
	return relapse.NewGrammar(refs), nil
}

func addXmlns(p *relapse.Pattern) *relapse.Pattern {
	return relapse.NewConcat(
		relapse.NewOptional(
			relapse.NewTreeNode(relapse.NewStringName("@xmlns"), relapse.NewZAny()),
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

func translatePattern(p *NameOrPattern, attr bool) *relapse.Pattern {
	if p.NotAllowed != nil {
		return relapse.NewEmptySet()
	}
	if p.Empty != nil {
		if attr {
			return combinator.Value(funcs.StringEq(Token(funcs.StringVar()), funcs.StringConst("")))
		}
		return relapse.NewOr(
			relapse.NewEmpty(),
			combinator.Value(funcs.StringEq(Token(funcs.StringVar()), funcs.StringConst(""))),
		)
	}
	if p.Text != nil {
		return relapse.NewZeroOrMore(combinator.Value(funcs.TypeString(funcs.StringVar())))
	}
	if p.Data != nil {
		if len(p.Data.DatatypeLibrary) > 0 {
			panic("data datatypeLibrary not supported")
		}
		if p.Data.Except == nil {
			return relapse.NewOr(combinator.Value(funcs.TypeString(funcs.StringVar())), relapse.NewEmpty())
		}
		expr, nullable := translateLeaf(p.Data.Except, funcs.StringVar())
		v := combinator.Value(funcs.And(
			funcs.TypeString(funcs.StringVar()),
			funcs.Not(expr),
		))
		if nullable {
			return relapse.NewAnd(v, relapse.NewNot(relapse.NewEmpty()))
		}
		return relapse.NewOr(v, relapse.NewEmpty())
	}
	if p.Value != nil {
		v, nullable := translateLeaf(p, funcs.StringVar())
		if !nullable {
			return combinator.Value(v)
		}
		return relapse.NewOr(combinator.Value(v), relapse.NewEmpty())
	}
	if p.List != nil {
		regexStr, nullable := listToRegex(p.List.NameOrPattern)
		val := combinator.Value(funcs.Regex(funcs.StringConst("^"+regexStr+"$"), Token(funcs.StringVar())))
		if !nullable {
			return val
		}
		return relapse.NewOr(val, relapse.NewEmpty())
	}
	if p.Attribute != nil {
		nameExpr := translateNameClass(p.Attribute.Left, true)
		pattern := translatePattern(p.Attribute.Right, true)
		return relapse.NewTreeNode(nameExpr, pattern)
	}
	if p.Ref != nil {
		return relapse.NewReference(p.Ref.Name)
	}
	if p.OneOrMore != nil {
		inside := translatePattern(p.OneOrMore.NameOrPattern, attr)
		return relapse.NewConcat(inside, relapse.NewZeroOrMore(inside))
	}
	if p.Choice != nil {
		left := translatePattern(p.Choice.Left, attr)
		right := translatePattern(p.Choice.Right, attr)
		return relapse.NewOr(left, right)
	}
	if p.Group != nil {
		left := translatePattern(p.Group.Left, attr)
		right := translatePattern(p.Group.Right, attr)
		if attr {
			return relapse.NewInterleave(left, right)
		}
		if hasAttr(p.Group.Right) {
			if hasAttr(p.Group.Left) {
				return relapse.NewInterleave(left, right)
			}
			return relapse.NewConcat(right, left)
		}
		return relapse.NewConcat(left, right)
	}
	if p.Interleave != nil {
		left := translatePattern(p.Interleave.Left, attr)
		right := translatePattern(p.Interleave.Right, attr)
		return relapse.NewInterleave(left, right)
	}
	panic(fmt.Sprintf("unreachable pattern %v", p))
}

func newTreeNode(n *NameOrPattern, pattern *relapse.Pattern) *relapse.Pattern {
	if n.Choice != nil {
		return relapse.NewOr(
			newTreeNode(n.Choice.Left, pattern),
			newTreeNode(n.Choice.Right, pattern),
		)
	}
	if n.AnyName != nil {
		if n.AnyName.Except == nil {
			return relapse.NewTreeNode(
				relapse.NewAnyName(),
				pattern,
			)
		}
		except := translateNameClass(n.AnyName.Except, false)
		return relapse.NewTreeNode(relapse.NewAnyNameExcept(except), pattern)
	}
	if n.NsName != nil {
		// if len(n.NsName.Ns) == 0 {
		// 	return relapse.NewTreeNode(relapse.NewAnyName(), pattern)
		// }
		panic("nsName is not supported")
	}
	if n.Name != nil {
		if len(n.Name.Ns) > 0 {
			return relapse.NewTreeNode(relapse.NewStringName(n.Name.Text), relapse.NewConcat(
				relapse.NewTreeNode(relapse.NewStringName("@xmlns"),
					combinator.Value(
						funcs.StringEq(
							funcs.StringVar(),
							funcs.StringConst(n.Name.Ns),
						),
					),
				),
				pattern,
			))
		}
		return relapse.NewTreeNode(relapse.NewStringName(n.Name.Text), pattern)
	}
	panic(fmt.Sprintf("unreachable nameclass %v", n))
}

func translateNameClass(n *NameOrPattern, attr bool) *relapse.NameExpr {
	if n.Choice != nil {
		return relapse.NewNameChoice(
			translateNameClass(n.Choice.Left, attr),
			translateNameClass(n.Choice.Right, attr),
		)
	}
	if n.AnyName != nil {
		if n.AnyName.Except == nil {
			return relapse.NewAnyName()
		}
		except := translateNameClass(n.AnyName.Except, attr)
		return relapse.NewAnyNameExcept(except)
	}
	if n.NsName != nil {
		// if len(n.NsName.Ns) == 0 {
		// 	if n.NsName.Except != nil {
		// 		return relapse.NewAnyNameExcept(translateNameClass(n.NsName.Except, attr))
		// 	} else {
		// 		return relapse.NewAnyName()
		// 	}
		// }
		panic("nsName is not supported")
	}
	if n.Name != nil {
		if len(n.Name.Ns) > 0 {
			panic(fmt.Sprintf("name ns <%v> is not supported", n.Name.Ns))
		}
		if attr {
			return relapse.NewStringName("@" + n.Name.Text)
		}
		return relapse.NewStringName(n.Name.Text)
	}
	panic(fmt.Sprintf("unreachable nameclass %v", n))
}

func translateLeaf(p *NameOrPattern, v funcs.String) (funcs.Bool, bool) {
	if p.Value != nil {
		if len(p.Value.Ns) > 0 {
			panic("value ns not supported")
		}
		text := p.Value.Text
		if p.Value.IsString() {
			return funcs.StringEq(funcs.StringConst(text), v), len(text) == 0
		}
		text = strings.Replace(text, "\n", "", -1)
		text = strings.Replace(text, "\r", "", -1)
		text = strings.Replace(text, "\t", "", -1)
		text = strings.TrimSpace(text)
		return funcs.StringEq(funcs.StringConst(text), Token(v)), len(text) == 0
	}
	if p.Choice != nil {
		l, nl := translateLeaf(p.Choice.Left, v)
		r, nr := translateLeaf(p.Choice.Right, v)
		return funcs.Or(l, r), nl || nr
	}
	panic(fmt.Sprintf("unsupported leaf %v", p))
}

func listToRegex(p *NameOrPattern) (string, bool) {
	if p.Empty != nil {
		return "", true
	}
	if p.Data != nil {
		if p.Data.Except == nil {
			return `(\S)*`, false
		}
	}
	if p.Value != nil {
		if len(p.Value.Ns) > 0 {
			panic("list value ns not supported")
		}
		return p.Value.Text, len(p.Value.Text) == 0
	}
	if p.OneOrMore != nil {
		s, nullable := listToRegex(p.OneOrMore.NameOrPattern)
		return `(\s)?` + s + "(\\s" + s + ")*", nullable
	}
	if p.Choice != nil {
		l, nl := listToRegex(p.Choice.Left)
		r, nr := listToRegex(p.Choice.Right)
		return "(" + l + "|" + r + ")", nl || nr
	}
	if p.Group != nil {
		l, nl := listToRegex(p.Group.Left)
		r, nr := listToRegex(p.Group.Right)
		return l + `\s` + r, nl && nr
	}
	panic(fmt.Sprintf("unsupported list %v", p))
}
