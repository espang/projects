package boolParser

import (
	"errors"
	"fmt"
	"strings"
)

const (
	AND   = "and"
	OR    = "or"
	OPEN  = "("
	CLOSE = ")"
	NOT   = "not"
)

type Node interface {
	evaluate(map[string]bool) bool
}

type And struct {
	left, right Node
}

func (a *And) evaluate(m map[string]bool) bool {
	return a.left.evaluate(m) && a.right.evaluate(m)
}

type Or struct {
	left, right Node
}

func (o *Or) evaluate(m map[string]bool) bool {
	return o.left.evaluate(m) || o.right.evaluate(m)
}

type Not struct {
	child Node
}

func (n *Not) evaluate(m map[string]bool) bool {
	return !n.child.evaluate(m)
}

type Value struct {
	label string
}

func (v *Value) evaluate(m map[string]bool) bool {
	if value, ok := m[v.label]; ok {
		return value
	}
	panic(fmt.Sprintf(
		"No value with label '%s' in %v",
		v.label,
		m,
	))
}

type Expression struct {
	root Node
}

func (e *Expression) evaluate(valmap map[string]bool) bool {
	return e.root.evaluate(valmap)
}

func createSubTree(tokens []string) (Node, error) {
	//fmt.Println("Tokens: ", tokens)
	numberOfTokens := len(tokens)
	if numberOfTokens == 0 {
		return nil, errors.New("Cannot create tree with zero tokens")
	}
	if numberOfTokens == 1 {
		token := tokens[0]
		switch strings.ToLower(token) {
		case AND:
			return nil, errors.New("Got AND as a leaf")
		case OR:
			return nil, errors.New("Got OR as a leaf")
		default:
			return &Value{label: token}, nil
		}
	}

	var n Node
	var lvl int
	//search ors
	for i, token := range tokens {
		switch strings.ToLower(token) {
		case OR:
			if lvl != 0 {
				continue
			}
			left, err := createSubTree(tokens[:i])
			if err != nil {
				return nil, err
			}
			right, err := createSubTree(tokens[i+1:])
			if err != nil {
				return nil, err
			}
			n = &Or{left: left, right: right}
			return n, nil
		case OPEN:
			lvl--
		case CLOSE:
			lvl++
		}
	}

	if lvl != 0 {
		if lvl > 0 {
			return nil, errors.New("More closing than opening brackets")
		}
		return nil, errors.New("More opening than closing brackets")
	}
	//search ands
	for i, token := range tokens {
		switch strings.ToLower(token) {
		case AND:
			if lvl != 0 {
				continue
			}
			left, err := createSubTree(tokens[:i])
			if err != nil {
				return nil, err
			}
			right, err := createSubTree(tokens[i+1:])
			if err != nil {
				return nil, err
			}
			n = &And{left: left, right: right}
			return n, nil
		case OPEN:
			lvl--
		case CLOSE:
			lvl++
		}
		if strings.ToLower(token) == AND {

		}
	}

	if strings.ToLower(tokens[0]) == NOT {
		child, err := createSubTree(tokens[1:])
		if err != nil {
			return nil, err
		}
		n = &Not{child: child}
		return n, nil
	}
	start, end := tokens[0], tokens[numberOfTokens-1]

	if start == OPEN && end == CLOSE {
		return createSubTree(tokens[1 : numberOfTokens-1])
	}

	return nil, errors.New("could not handle expression!")
}

func NewExpression(txt string) (*Expression, error) {

	tokens := strings.Split(txt, " ")
	n, err := createSubTree(tokens)
	if err != nil {
		return nil, err
	}
	return &Expression{root: n}, nil
}
