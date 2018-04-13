package evaluator

import (
	"github.com/jasonroelofs/late"
	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
	"github.com/jasonroelofs/late/tag"
	"github.com/jasonroelofs/late/template/ast"
)

type Evaluator struct {
	context  *context.Context
	template *ast.Template
}

func New(template *ast.Template, context *context.Context) *Evaluator {
	return &Evaluator{
		template: template,
		context:  context,
	}
}

func (e *Evaluator) Run() []object.Object {
	var objects []object.Object

	for _, statement := range e.template.Statements {
		result := e.eval(statement)
		objects = append(objects, result)
	}

	return objects
}

func (e *Evaluator) Set(variable string, value interface{}) {
	e.context.Set(variable, value)
}

func (e *Evaluator) Get(variable string) object.Object {
	return e.context.Get(variable)
}

func (e *Evaluator) Eval(node tag.Statement) object.Object {
	return e.eval(node.(ast.Node))
}

func (e *Evaluator) eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Top-level Statements
	case *ast.RawStatement:
		return object.New(node.String())

	case *ast.VariableStatement:
		return e.eval(node.Expression)

	case *ast.TagStatement:
		return e.evalTagStatement(node)

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

	case *ast.IndexExpression:
		left := e.eval(node.Left)
		index := e.eval(node.Index)
		return e.evalIndex(left, index)

	// Literals
	case *ast.NumberLiteral:
		return object.New(node.Value)

	case *ast.BooleanLiteral:
		return convertBoolean(node.Value)

	case *ast.StringLiteral:
		return object.New(node.Value)

	case *ast.Identifier:
		return e.evalIdentifier(node.Value)

	case *ast.FilterLiteral:
		return e.evalFilterLiteral(node)

	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node)

	default:
		return object.NULL
	}
}

func (e *Evaluator) evalTagStatement(node *ast.TagStatement) object.Object {
	return node.Tag.Eval(e, e.prepareTagResults(node))
}

func (e *Evaluator) prepareTagResults(node *ast.TagStatement) *tag.ParseResult {
	var results []object.Object

	for _, node := range node.Nodes {
		switch node := node.(type) {
		case *ast.Identifier:
			results = append(results, object.New(node.Value))
		default:
			results = append(results, e.eval(node))
		}
	}

	parseResults := &tag.ParseResult{
		TagName: node.TagName,
		Nodes:   results,
	}

	if node.BlockStatement != nil {
		for _, stmt := range node.BlockStatement.Statements {
			parseResults.Statements = append(parseResults.Statements, stmt.(tag.Statement))
		}

		for _, subTag := range node.SubTags {
			parseResults.SubTagResults = append(
				parseResults.SubTagResults,
				e.prepareTagResults(subTag),
			)
		}
	}

	return parseResults
}

func (e *Evaluator) evalInfix(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.TYPE_NUMBER && right.Type() == object.TYPE_NUMBER:
		return e.evalNumberOperation(operator, left, right)
	case operator == "==":
		return convertBoolean(left.Value() == right.Value())
	case operator == "!=":
		return convertBoolean(left.Value() != right.Value())
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalNumberOperation(operator string, left, right object.Object) object.Object {
	leftVal := left.Value().(float64)
	rightVal := right.Value().(float64)

	switch operator {
	case "+":
		return object.New(leftVal + rightVal)
	case "-":
		return object.New(leftVal - rightVal)
	case "*":
		return object.New(leftVal * rightVal)
	case "/":
		return object.New(leftVal / rightVal)
	case ">":
		return convertBoolean(leftVal > rightVal)
	case "<":
		return convertBoolean(leftVal < rightVal)
	case ">=":
		return convertBoolean(leftVal >= rightVal)
	case "<=":
		return convertBoolean(leftVal <= rightVal)
	case "==":
		return convertBoolean(leftVal == rightVal)
	case "!=":
		return convertBoolean(leftVal != rightVal)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalPrefix(operator string, right object.Object) object.Object {
	switch {
	case right.Type() == object.TYPE_NUMBER:
		return e.evalNumberPrefix(operator, right)
	default:
		return object.NULL
	}
}

func (e *Evaluator) evalNumberPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return object.New(right.Value().(float64) * -1)
	default:
		return right
	}
}

func (e *Evaluator) evalIdentifier(name string) object.Object {
	return e.context.Get(name)
}

func (e *Evaluator) evalFilterLiteral(node *ast.FilterLiteral) object.Object {
	filterObj := &object.Filter{
		Name:       node.Name,
		Parameters: make(map[string]object.Object),
	}

	for paramName, paramExp := range node.Parameters {
		filterObj.Parameters[paramName] = e.eval(paramExp)
	}

	return filterObj
}

func (e *Evaluator) evalFilter(input, filter object.Object) object.Object {
	filterName := filter.(*object.Filter).Name
	filterParams := filter.(*object.Filter).Parameters

	filterFunc := late.FindFilter(filterName)

	if filterFunc == nil {
		return object.NULL
	}

	return filterFunc.Call(input, filterParams)
}

func (e *Evaluator) evalIndex(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.TYPE_ARRAY && index.Type() == object.TYPE_NUMBER:
		return e.evalArrayAccess(left, index)
	case left.Type() == object.TYPE_HASH:
		return e.evalHashAccess(left, index)
	default:
		// Unknown action "index" on this object
		return object.NULL
	}
}

func (e *Evaluator) evalArrayAccess(left, index object.Object) object.Object {
	array := left.(*object.Array)
	idx := int(index.Value().(float64))

	if idx < 0 || len(array.Elements) <= idx {
		return object.NULL
	}

	return array.Elements[idx]
}

func (e *Evaluator) evalHashAccess(left, index object.Object) object.Object {
	hash := left.(*object.Hash)

	return hash.Get(index)
}

func (e *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral) object.Object {
	array := &object.Array{}

	for _, expr := range node.Expressions {
		array.Elements = append(array.Elements, e.eval(expr))
	}

	return array
}

func convertBoolean(value bool) object.Object {
	if value {
		return object.TRUE
	} else {
		return object.FALSE
	}
}
