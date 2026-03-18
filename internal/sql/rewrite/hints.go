package rewrite

import (
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

type Hint struct {
	NotNull  bool
	Node     ast.Node // The inner node
	Location int      // The location of the original sqlc.nullable/notnull call
	FuncName string   // "nullable" or "notnull"
}

type HintSet []*Hint

func (hs HintSet) ForNode(node ast.Node) (*Hint, bool) {
	for _, h := range hs {
		if h.Node == node {
			return h, true
		}
	}
	return nil, false
}

func Hints(raw *ast.RawStmt) (*ast.RawStmt, HintSet) {
	var hints []*Hint

	node := astutils.Apply(raw, func(cr *astutils.Cursor) bool {
		node := cr.Node()
		call, ok := node.(*ast.FuncCall)
		if !ok || call.Func == nil || call.Func.Schema != "sqlc" {
			return true
		}

		if call.Func.Name != "nullable" && call.Func.Name != "notnull" {
			return true
		}

		if len(call.Args.Items) == 0 {
			return true
		}

		inner := call.Args.Items[0]
		hints = append(hints, &Hint{
			NotNull:  call.Func.Name == "notnull",
			Node:     inner,
			Location: call.Location,
			FuncName: call.Func.Name,
		})

		cr.Replace(inner)
		return false
	}, nil)

	return node.(*ast.RawStmt), hints
}
