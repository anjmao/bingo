package cache

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/saibing/bingo/langserver/internal/source"
	"golang.org/x/tools/go/ast/astutil"
)

func TestAppendNode(t *testing.T) {
	testLineNum := 8
	f, fset, fast := loadViewFile(t, "/Users/anjmao/s/bingo/", "/Users/anjmao/s/bingo/testdata/astappend.go")

	fmt.Println("-----original file ast-----")
	ast.Print(fset, fast)

	p := &nodeLineParser{fset, fast}

	newNode, err := p.parseLineExpr(f.content, testLineNum)
	if err != nil {
		t.Fatalf("could not parse line expression: %v", err)
	}

	fmt.Println("-----new line node ast-----")
	ast.Print(fset, newNode)

	p.replaceNode(testLineNum+1, newNode)

	fmt.Println("-----file ast after replace-----")
	ast.Print(fset, fast)
}

func loadViewFile(t *testing.T, rootDir, sourceFile string) (*File, *token.FileSet, *ast.File) {
	v := NewView()
	v.getLoadDir = func(_ string) string {
		return rootDir
	}

	f := v.GetFile(source.ToURI(sourceFile))

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

type nodeLineParser struct {
	fset *token.FileSet
	ast  *ast.File
}

func (p *nodeLineParser) replaceNode(lineNum int, newNode ast.Node) {
	astutil.Apply(p.ast, func(c *astutil.Cursor) bool {
		n := c.Node()
		if n == nil {
			return false
		}
		line := p.fset.Position(n.Pos()).Line
		if line != lineNum {
			return true
		}
		if _, ok := n.(*ast.CallExpr); ok {
			// ast.Print(p.fset, ex)
			c.Replace(newNode)
			return false
		}
		return true
	}, nil)
}

func (p *nodeLineParser) parseLineExpr(contents []byte, posLine int) (ast.Node, error) {
	start, end, err := p.getLineOffsetRange(contents, posLine)
	if err != nil {
		return nil, err
	}
	lineContent := strings.TrimSpace(string(contents[start:end]))
	return parser.ParseExpr(lineContent)
}

func (p *nodeLineParser) getLineOffsetRange(contents []byte, posLine int) (int, int, error) {
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

