package cache

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/saibing/bingo/langserver/internal/source"
)

func TestAppendNode(t *testing.T) {
	f, fset, _ := loadViewFile(t)

	// ast.Inspect(fast, func(n ast.Node) bool {
	// 	if n != nil {
	// 		fmt.Println("pos", fset.Position(n.Pos()).Line, n)
	// 	}
	// 	return true
	// })

	start, end, err := getLineOffsetRange(f.content, 8)
	if err != nil {
		t.Fatalf("could not get line offset range: %v", err)
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

// TODO(anjmao): move this func to view.go onces it is done
func getLineOffsetRange(contents []byte, posLine int) (int, int, error) {
	line := 0
	col := 0
	start := 0
	end := 0
	offset := 0
	for _, b := range contents {
		// fmt.Println("line", line, "col", col, "char", string(b))
		if line == posLine && start == 0 {
			start = offset
		}

		if line != posLine && start != 0 {
			end = offset
			return start, end, nil
		}

		if line > posLine {
			return 0, 0, fmt.Errorf("character is beyond line %d boundary", posLine)
		}

		offset++
		if b == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}

	return 0, 0, fmt.Errorf("file only has %d lines", line+1)
}
