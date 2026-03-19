package postgresql

import (
	"fmt"
	"strings"

	nodes "github.com/pganalyze/pg_query_go/v6"
)

func isArray(n *nodes.TypeName) bool {
	if n == nil {
		return false
	}
	return len(n.ArrayBounds) > 0
}

func isNotNull(n *nodes.ColumnDef) bool {
	if n.IsNotNull {
		return true
	}
	for _, c := range n.Constraints {
		switch inner := c.Node.(type) {
		case *nodes.Node_Constraint:
			if inner.Constraint.Contype == nodes.ConstrType_CONSTR_NOTNULL {
				return true
			}
			if inner.Constraint.Contype == nodes.ConstrType_CONSTR_PRIMARY {
				return true
			}
		}
	}
	return false
}

func isGenerated(n *nodes.ColumnDef) bool {
	for _, c := range n.Constraints {
		switch inner := c.Node.(type) {
		case *nodes.Node_Constraint:
			if inner.Constraint.Contype == nodes.ConstrType_CONSTR_GENERATED {
				return true
			}
		}
	}
	return false
}

func IsNamedParamFunc(node *nodes.Node) bool {
	fun, ok := node.Node.(*nodes.Node_FuncCall)
	return ok && joinNodes(fun.FuncCall.Funcname, ".") == "sqlc.arg"
}

func IsNamedParamSign(node *nodes.Node) bool {
	expr, ok := node.Node.(*nodes.Node_AExpr)
	return ok && joinNodes(expr.AExpr.Name, ".") == "@"
}

func makeByte(s string) byte {
	var b byte
	if s == "" {
		return b
	}
	return []byte(s)[0]
}

func makeUint32Slice(in []uint64) []uint32 {
	out := make([]uint32, len(in))
	for i, v := range in {
		out[i] = uint32(v)
	}
	return out
}

func makeString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func preProcessGenerated(sql string) string {
	var out strings.Builder
	idx := 0
	upper := strings.ToUpper(sql)
	for {
		start := strings.Index(upper[idx:], "GENERATED ALWAYS AS")
		if start == -1 {
			out.WriteString(sql[idx:])
			break
		}
		start += idx
		out.WriteString(sql[idx:start])

		// Find opening paren
		parenStart := strings.Index(sql[start:], "(")
		if parenStart == -1 {
			out.WriteString(sql[start : start+19]) // Skip "GENERATED ALWAYS AS"
			idx = start + 19
			continue
		}
		parenStart += start

		// Find matching closing paren
		depth := 0
		parenEnd := -1
		for i := parenStart; i < len(sql); i++ {
			if sql[i] == '(' {
				depth++
			} else if sql[i] == ')' {
				depth--
				if depth == 0 {
					parenEnd = i
					break
				}
			}
		}

		if parenEnd == -1 {
			out.WriteString(sql[start : parenStart+1])
			idx = parenStart + 1
			continue
		}

		expr := sql[parenStart+1 : parenEnd]
		out.WriteString(fmt.Sprintf("GENERATED ALWAYS AS (%s) STORED", expr))

		// Check if we need to skip VIRTUAL or STORED after the paren
		afterParen := sql[parenEnd+1:]
		spaceLen := len(afterParen) - len(strings.TrimLeft(afterParen, " \t\n\r"))
		afterParen = afterParen[spaceLen:]
		if strings.HasPrefix(strings.ToUpper(afterParen), "VIRTUAL") {
			idx = parenEnd + 1 + spaceLen + 7
		} else if strings.HasPrefix(strings.ToUpper(afterParen), "STORED") {
			idx = parenEnd + 1 + spaceLen + 6
		} else {
			idx = parenEnd + 1
		}
	}
	return out.String()
}
