package evaluator

import (
	"github.com/jasonroelofs/late/template/ast"
	"github.com/jasonroelofs/late/template/object"
)

type Evaluator struct {
	template *ast.Template
}

func New(template *ast.Template) *Evaluator {
	return &Evaluator{template: template}
}

func (e *Evaluator) Run() []object.Object {
	var objects []object.Object

	for _, statement := range e.template.Statements {
		result := e.eval(statement)
		objects = append(objects, result)
	}

	return objects
}

func (e *Evaluator) eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Top-level Statements
	case *ast.RawStatement:
		return &object.String{Value: node.String()}

	case *ast.VariableStatement:
		return e.eval(node.Expression)

	// Expressions
	case *ast.InfixExpression:
		left := e.eval(node.Left)
		right := e.eval(node.Right)
		return e.evalInfix(node.Operator, left, right)

	case *ast.PrefixExpression:
		right := e.eval(node.Right)
		return e.evalPrefix(node.Operator, right)

	case *ast.FilterExpression:
		input := e.eval(node.Input)
		filter := e.eval(node.Filter)
		return e.evalFilter(input, filter)

	// Literals
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}

	case *ast.BooleanLiteral:
		return convertBoolean(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.FilterLiteral:
		return &object.Filter{Name: node.Name}

	default:
		return object.NULL
	}
}

func (e *Evaluator) evalInfix(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.OBJ_NUMBER && right.Type() == object.OBJ_NUMBER:
		return e.evalNumberOperation(operator, left, right)
	case operator == "==":
		return convertBoolean(left == right)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalNumberOperation(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value

	switch operator {
	case "+":
		return &object.Number{Value: leftVal + rightVal}
	case "-":
		return &object.Number{Value: leftVal - rightVal}
	case "*":
		return &object.Number{Value: leftVal * rightVal}
	case "/":
		return &object.Number{Value: leftVal / rightVal}
	case ">":
		return convertBoolean(leftVal > rightVal)
	case "<":
		return convertBoolean(leftVal < rightVal)
	case ">=":
		return convertBoolean(leftVal >= rightVal)
	case "<=":
		return convertBoolean(leftVal <= rightVal)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalPrefix(operator string, right object.Object) object.Object {
	switch {
	case right.Type() == object.OBJ_NUMBER:
		return e.evalNumberPrefix(operator, right)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalNumberPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return &object.Number{
			Value: right.(*object.Number).Value * -1,
		}
	default:
		return right
	}
}

func (e *Evaluator) evalFilter(input, filter object.Object) object.Object {
	//	filterFunc := e.Context.FindFilter(filter.(*object.Filter).Name)
	//
	//	if filterFunc == nil {
	//		return object.NULL
	//	}
	//
	//	return object.FromNativeType(filterFunc.Call(input))
	return object.NULL
}

func convertBoolean(value bool) object.Object {
	if value {
		return object.TRUE
	} else {
		return object.FALSE
	}
}
