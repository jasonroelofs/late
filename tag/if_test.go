package tag

import (
	"testing"

	"github.com/jasonroelofs/late/object"
)

type TestEnv struct {
	StatementsRan []Statement
}

func (t *TestEnv) Eval(stmt Statement) object.Object {
	t.StatementsRan = append(t.StatementsRan, stmt)
	return object.New(stmt.String())
}

func (t *TestEnv) Set(_ string, _ interface{}) {}
func (t *TestEnv) Get(_ string) object.Object {
	return object.NULL
}

type TestStatement struct {
	Out string
}

func (t TestStatement) String() string {
	return t.Out
}

func TestExpressionsAreTruthy(t *testing.T) {
	tag := new(If)

	env := new(TestEnv)
	results := &ParseResult{
		TagName:    "if",
		Nodes:      []object.Object{object.New("Value")},
		Statements: []Statement{&TestStatement{Out: "Statement 1"}},
	}

	result := tag.Eval(env, results)
	if result.Value() != "Statement 1" {
		t.Fatalf("Did not execute the success block, got %v", result)
	}
}

func TestElseIfIsTruthy(t *testing.T) {
	tag := new(If)
	env := new(TestEnv)

	results := &ParseResult{
		TagName:    "if",
		Nodes:      []object.Object{object.FALSE},
		Statements: []Statement{&TestStatement{Out: "Statement 1"}},
		SubTagResults: []*ParseResult{
			&ParseResult{
				TagName:    "elsif",
				Nodes:      []object.Object{object.New(123)},
				Statements: []Statement{&TestStatement{Out: "Statement 2"}},
			},
		},
	}

	result := tag.Eval(env, results)
	if result.Value() != "Statement 2" {
		t.Fatalf("Did not execute the elsif block, got %v", result)
	}
}
