package tag

import (
	"testing"

	"github.com/jasonroelofs/late/context"
	"github.com/jasonroelofs/late/object"
	s "github.com/jasonroelofs/late/template/statement"
)

type TestStatement struct {
	Out string
}

func (t TestStatement) String() string {
	return t.Out
}

type TestEval struct {
	StatementsRan []s.Statement
}

func (t *TestEval) EvalAll(stmts []s.Statement) object.Object {
	results := &object.Array{}

	for _, stmt := range stmts {
		results.Elements = append(results.Elements, t.Eval(stmt))
	}

	return results
}

func (t *TestEval) Eval(stmt s.Statement) object.Object {
	t.StatementsRan = append(t.StatementsRan, stmt)
	return object.New(stmt.String())
}
func (t *TestEval) Interrupt() string { return "" }
func (t *TestEval) ClearInterrupt()   {}

func TestExpressionsAreTruthy(t *testing.T) {
	tag := new(If)
	eval := new(TestEval)
	ctx := context.New()
	ctx.SetEvaluator(eval)

	results := &ParseResult{
		TagName:    "if",
		Nodes:      []object.Object{object.New("Value")},
		Statements: []s.Statement{&TestStatement{Out: "Statement 1"}},
	}

	result := tag.Eval(ctx, results).(*object.Array)
	if result.Get(0).Value() != "Statement 1" {
		t.Fatalf("Did not execute the success block, got %v", result)
	}
}

func TestElseIfIsTruthy(t *testing.T) {
	tag := new(If)
	eval := new(TestEval)
	ctx := context.New()
	ctx.SetEvaluator(eval)

	results := &ParseResult{
		TagName:    "if",
		Nodes:      []object.Object{object.FALSE},
		Statements: []s.Statement{&TestStatement{Out: "Statement 1"}},
		SubTagResults: []*ParseResult{
			&ParseResult{
				TagName:    "elsif",
				Nodes:      []object.Object{object.New(123)},
				Statements: []s.Statement{&TestStatement{Out: "Statement 2"}},
			},
		},
	}

	result := tag.Eval(ctx, results).(*object.Array)
	if result.Get(0).Value() != "Statement 2" {
		t.Fatalf("Did not execute the elsif block, got %v", result)
	}
}
