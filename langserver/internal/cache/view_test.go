package cache

import (
	"fmt"
	"github.com/saibing/bingo/langserver/internal/source"
	"github.com/saibing/bingo/pkg/lsp"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestAppendNode(t *testing.T) {
	f, fset, _ := loadViewFile(t)

	// ast.Inspect(fast, func(n ast.Node) bool {
	// 	if n != nil {
	// 		fmt.Println("pos", fset.Position(n.Pos()).Line, n)
	// 	}
	// 	return true
	// })

	start, ok, why := offsetForPosition(f.content, lsp.Position{Line: 8, Character: 0})
	if !ok {
		t.Fatalf("could not get offset start position: %v", why)
	}

	end, ok, why := offsetForPosition(f.content, lsp.Position{Line: 8, Character: 18})
	if !ok {
		t.Fatalf("could not get offset end position: %v", why)
	}

	lineContent := strings.TrimSpace(string(f.content[start:end]))

	expr, err := parser.ParseExpr(lineContent)
	if err != nil {
		t.Fatalf("could not parse line expression: %v", err)
	}

	ast.Print(fset, expr)

	fmt.Println("line content", lineContent)

	// TODO: now as we have full line ast we replace this block inside actual file AST
}

// TODO: pass test file path
func loadViewFile(t *testing.T) (*File, *token.FileSet, *ast.File) {
	v := NewView()
	v.getLoadDir = func(_ string) string {
		return "/Users/anjmao/s/bingo/"
	}

	f := v.GetFile(source.ToURI("/Users/anjmao/s/bingo/testdata/astappend.go"))

	fset, err := f.GetFileSet()
	if err != nil {
		t.Fatal(err)
	}

	fast, err := f.GetAST()
	if err != nil {
		t.Fatal(err)
	}

	return f, fset, fast
}

// TODO(anjmao): remove this func from test
func offsetForPosition(contents []byte, p lsp.Position) (offset int, valid bool, whyInvalid string) {
	line := 0
	col := 0
	// TODO(sqs): count chars, not bytes, per LSP. does that mean we
	// need to maintain 2 separate counters since we still need to
	// return the offset as bytes?
	for _, b := range contents {
		// fmt.Println("line", line, "col", col, "char", string(b))
		if line == p.Line && col == p.Character {
			return offset, true, ""
		}
		if (line == p.Line && col > p.Character) || line > p.Line {
			return 0, false, fmt.Sprintf("character %d is beyond line %d boundary", p.Character, p.Line)
		}
		offset++
		if b == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	if line == p.Line && col == p.Character {
		return offset, true, ""
	}
	if line == 0 {
		return 0, false, fmt.Sprintf("character %d is beyond first line boundary", p.Character)
	}
	return 0, false, fmt.Sprintf("file only has %d lines", line+1)
}
