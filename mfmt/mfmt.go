package mfmt

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var stdlib = make(map[string]bool)

var module string

func DoFmt() {
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || strings.HasPrefix(path, ".") || strings.HasPrefix(path, "vendor") || strings.HasSuffix(path, ".pb.go") || filepath.Ext(path) != ".go" {
			return nil
		}

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "read file %s err", path)
		}

		digest0 := sha256.Sum256(raw)
		if raw, err = format.Source(raw); err != nil {
			return errors.Wrapf(err, "format file %s err", path)
		}

		file, err := parser.ParseFile(token.NewFileSet(), "", raw, 0)
		if err != nil {
			return errors.Wrapf(err, "parse file %s err", path)
		}

		var first, last int
		var imports []*ast.ImportSpec
		comments := make(map[string]string)

		ast.Inspect(file, func(n ast.Node) bool {
			switch spec := n.(type) {
			case *ast.ImportSpec:
				if first == 0 {
					first = int(spec.Pos())
				}
				last = int(spec.End())

				imports = append(imports, spec)

				k := last - 1
				for ; k < len(raw); k++ {
					if raw[k] == '\r' || raw[k] == '\n' {
						break
					}
				}

				comment := string(raw[last-1 : k])
				if index := strings.Index(comment, "//"); index != -1 {
					comments[spec.Path.Value] = strings.TrimSpace(comment[index+2:])
				}
			}
			return true
		})

		if imports != nil {
			buf := bytes.NewBuffer(nil)
			buf.Write(raw[:first-2])
			buf.WriteString(sort(imports, comments))
			buf.Write(raw[last-1:])

			if raw, err = format.Source(buf.Bytes()); err != nil {
				return errors.Wrapf(err, "double format file %s err", path)
			}
		}

		digest1 := sha256.Sum256(raw)
		if !bytes.Equal(digest0[:], digest1[:]) {
			fmt.Println("格式化：", path)
		}

		if err = ioutil.WriteFile(path, raw, info.Mode()); err != nil {
			return errors.Wrapf(err, "write file %s err", path)
		}

		return nil
	})
	if err != nil {
		log.Fatal("scan project err", zap.Error(err))
	}
}

func sort(imports []*ast.ImportSpec, comments map[string]string) string {
	system := bytes.NewBuffer(nil)
	group := bytes.NewBuffer(nil)
	others := bytes.NewBuffer(nil)

	for _, pkg := range imports {
		value := strings.Trim(pkg.Path.Value, `"`)
		switch {
		case stdlib[value]:
			if pkg.Name != nil {
				system.WriteString(pkg.Name.String())
				system.WriteString(" ")
			}

			system.WriteString(pkg.Path.Value)
			if comment, ok := comments[pkg.Path.Value]; ok {
				system.WriteString(" ")
				system.WriteString("// ")
				system.WriteString(comment)
			}
			system.WriteString("\n")

		case strings.HasPrefix(value, module):
			if pkg.Name != nil {
				group.WriteString(pkg.Name.String())
				group.WriteString(" ")
			}

			group.WriteString(pkg.Path.Value)
			if comment, ok := comments[pkg.Path.Value]; ok {
				group.WriteString(" ")
				group.WriteString("// ")
				group.WriteString(comment)
			}
			group.WriteString("\n")

		default:
			if pkg.Name != nil {
				others.WriteString(pkg.Name.String())
				others.WriteString(" ")
			}

			others.WriteString(pkg.Path.Value)
			if comment, ok := comments[pkg.Path.Value]; ok {
				others.WriteString(" ")
				others.WriteString("// ")
				others.WriteString(comment)
			}
			others.WriteString("\n")
		}
	}

	return fmt.Sprintf("%s\n%s\n%s", system.String(), group.String(), others.String())
}
