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
	exprparser "github.com/katydid/katydid/expr/parser"
	"github.com/katydid/katydid/funcs"
	"github.com/katydid/katydid/relapse/ast"
)

func translate(g *Grammar) (*relapse.Grammar, error) {
	refs := make(relapse.RefLookup)
	refs["main"] = translatePattern(g.Start, false)
	for _, d := range g.Define {
		nameExpr := translateNameClass(d.Element.Left, false)
		pattern := translatePattern(d.Element.Right, false)
		refs[d.Name] = relapse.NewTreeNode(nameExpr, pattern)
	}
	return relapse.NewGrammar(refs), nil
}

func newLeaf(exprStr string) *relapse.Pattern {
	expr, err := exprparser.NewParser().ParseExpr(exprStr)
	if err != nil {
		panic(err)
	}
	return &relapse.Pattern{LeafNode: &relapse.LeafNode{
		Expr: expr,
	}}
}

func translatePattern(p *NameOrPattern, attr bool) *relapse.Pattern {
	if p.NotAllowed != nil {
		return relapse.NewEmptySet()
	}
	if p.Empty != nil {
		return relapse.NewEmpty()
	}
	if p.Text != nil {
		return relapse.NewZeroOrMore(newLeaf(funcs.Sprint(funcs.TypeString())))
	}
	if p.Data != nil {
		if p.Data.Except == nil {
			return newLeaf(funcs.Sprint(funcs.TypeString()))
		}
		expr := translateLeaf(p.Data.Except, funcs.StringVar())
		return newLeaf(funcs.Sprint(funcs.And(funcs.TypeString(), funcs.Not(expr))))
	}
	if p.Value != nil {
		return newLeaf(funcs.Sprint(translateLeaf(p, funcs.StringVar())))
	}
	if p.List != nil {
		panic("todo: list")
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
			return relapse.NewOr(
				relapse.NewConcat(left, right),
				relapse.NewConcat(right, left),
			)
		}
		return relapse.NewConcat(left, right)
	}
	if p.Interleave != nil {
		left := translatePattern(p.Interleave.Left, attr)
		right := translatePattern(p.Interleave.Right, attr)
		return relapse.NewOr(
			relapse.NewConcat(left, right),
			relapse.NewConcat(right, left),
		)
	}
	panic(fmt.Sprintf("unreachable pattern %v", p))
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
		panic("nsName is not supported")
	}
	if n.Name != nil {
		return relapse.NewName(n.Name.Text)
	}
	panic(fmt.Sprintf("unreachable nameclass %v", n))
}

func translateLeaf(p *NameOrPattern, v funcs.String) funcs.Bool {
	if p.Value != nil {
		if p.Value.IsString() {
			return funcs.StringEq(funcs.StringConst(p.Value.Text), v)
		}
		return funcs.StringEq(funcs.StringConst(p.Value.Text), Token(v))
	}
	if p.Choice != nil {
		return funcs.Or(translateLeaf(p.Choice.Left, v), translateLeaf(p.Choice.Right, v))
	}
	panic(fmt.Sprintf("unsupported leaf %v", p))
}
