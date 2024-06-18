package veritas

func Eval(exp Expression, vars map[byte]int) int {
	switch node := exp.(type) {
	case *VarLiteralExpression:
		return vars[node.Value]
	case *PrefixExpression:
		rightV := Eval(node.Right, vars)
		return evalPrefixExpression(node.Operator, rightV)
	case *InfixExpression:
		leftV := Eval(node.Left, vars)
		rightV := Eval(node.Right, vars)
		return evalInfixExpression(leftV, node.Operator, rightV)
	}
	return 0
}

func evalPrefixExpression(operator string, right int) int {
	if operator != "!" {
		return -1
	}
	return right ^ 1
}

func evalInfixExpression(left int, operator string, right int) int {

	switch operator {
	case "&":
		return left & right
	case "|":
		return left | right
	case "^":
		if left != right {
			return 1
		} else {
			return 0
		}
	case "@":
		return (left ^ 1) | right
	default:
		return -1
	}
}
