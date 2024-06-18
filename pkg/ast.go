package veritas

import (
	"bytes"
)

// Expression is the interface for all our AST nodes
type Expression interface {
	String() string
	Len(map[string]int) int
	len() int
}

// VarLiteralExpression represents a variable reference
type VarLiteralExpression struct {
	Token Token
	Value byte
}

func (rl *VarLiteralExpression) String() string { return rl.Token.Literal }
func (rl *VarLiteralExpression) Len(varCounts map[string]int) int {
	varCounts[rl.String()]++
	return varCounts[rl.String()]
}
func (rl *VarLiteralExpression) len() int { return 1 }

// PrefixExpression represents an expression prefixed with an operator (!A)
type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (rl *PrefixExpression) Len(a map[string]int) int {
	return rl.Right.Len(a)
}
func (rl *PrefixExpression) len() int { return rl.Right.len() }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents an infix expression (A & B)
type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (rl *InfixExpression) Len(a map[string]int) int {
	return rl.Right.Len(a) + rl.Left.Len(a)
}
func (rl *InfixExpression) len() int { return rl.Left.len() + rl.Right.len() }

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
