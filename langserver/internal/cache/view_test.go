package cache

import (
	"fmt"
	"github.com/saibing/bingo/langserver/internal/source"
	"go/ast"
	"testing"
)

func TestAppendNode(t *testing.T) {
	v := NewView()
	v.getLoadDir = func(_ string) string {
		return "/Users/anjmao/s/bingo/"
	}
	// TODO: use test path
	f := v.GetFile(source.ToURI("/Users/anjmao/s/bingo/langserver/internal/cache/view.go"))

	fset, err := f.GetFileSet()
	if err != nil {
		t.Fatal(err)
	}

	fast, err := f.GetAST()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fset)
	ast.Print(fset, fast)
}
