package cache

import (
	"context"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type packagePool map[string]*packages.Package

type PackageCache struct {
	mu   sync.RWMutex
	pool packagePool
}

const windowsOS = "windows"

func (c *PackageCache) Load(pkgDir string) (*packages.Package, error) {
	loadDir := getLoadDir(pkgDir)
	cacheKey := loadDir

	if runtime.GOOS == windowsOS {
		cacheKey = strings.Replace(loadDir, "\\", "/", -1)
	}

	log.Printf("load dir %s\n", loadDir)
	log.Printf("cache key %s\n", cacheKey)
	c.mu.RLock()

	pkg := c.pool[cacheKey]
	if pkg != nil {
		c.mu.RUnlock()
		return pkg, nil
	}

	c.mu.RUnlock()

	cfg := &packages.Config{Mode: packages.LoadAllSyntax, Context:context.Background(), Tests: true}
	err := os.Chdir(loadDir)
	if err != nil {
		return nil, err
	}

	pkgList, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, err
	}

	if len(pkgList) == 0 {
		return nil, nil
	}

	pkg = pkgList[0]

	go c.push(pkg)

	return pkg, nil
}

func (c *PackageCache) push(pkg *packages.Package) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache(pkg)
}

func (c *PackageCache) cache(pkg *packages.Package) {
	if len(pkg.CompiledGoFiles[0]) == 0 {
		return
	}

	cacheKey := getCacheKey(pkg.CompiledGoFiles[0])

	if _, ok := c.pool[cacheKey]; ok {
		return
	}

	c.pool[cacheKey] = pkg
	log.Printf("cached package %s\n", cacheKey)
	for _, importPkg := range pkg.Imports {
		c.cache(importPkg)
	}
}

func getLoadDir(dir string) string {
	if runtime.GOOS != windowsOS {
		return dir
	}

	return dir[1:]
}

func getCacheKey(filename string) string {
	dir := filepath.Dir(filename)
	if runtime.GOOS != windowsOS {
		return dir
	}

	dirs := strings.Split(dir, ":")
	if len(dirs) >= 2 {
		dirs[0] = strings.ToLower(dirs[0])
		dir = strings.Join(dirs, ":")
	}

	dir = strings.Replace(dir, "\\", "/", -1)
	return dir
}