package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	testOptions, err := parseTestOptions()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	tags := strings.Split(*testOptions.Tags, ",")

	result := []Test{}
	fset := token.NewFileSet()

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			isMatch := false

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				for _, tag := range tags {
					if strings.Contains(line, tag) {
						isMatch = true
						break
					}
				}
				if isMatch {
					break
				}
			}

			if !isMatch {
				return nil
			}

			for _, decl := range f.Decls {
				if fn, isFn := decl.(*ast.FuncDecl); isFn && fn.Name.IsExported() && fn.Recv == nil {
					if !strings.HasPrefix(fn.Name.Name, "Test") {
						continue
					}

					dir := filepath.Dir(path)
					result = append(result, Test{
						Name: fn.Name.Name,
						Dir:  dir,
					})
				}
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err) //nolint:forbidigo
		os.Exit(1)
	}

	tests := []string{}
	dirs := []string{}
	var hasErr error = nil

	for i := *testOptions.Shift; i < len(result); i += *testOptions.Runs {
		switch *testOptions.Mode {
		case "list":
			fmt.Println(result[i].Name) //nolint:forbidigo
		case "dir":
			fmt.Printf("%s,%s\n", result[i].Dir, result[i].Name) //nolint:forbidigo
		case "seq":
			hasErrSeq := execTest([]string{result[i].Name}, tags, []string{"./" + result[i].Dir}, *testOptions)
			if hasErrSeq != nil {
				fmt.Println(hasErrSeq)
				hasErr = errors.New("error")
			}
		case "run":
			tests = append(tests, result[i].Name)
			dirs = append(dirs, "./"+result[i].Dir)
		}
	}
	if *testOptions.Mode == "run" {
		// fmt.Printf(`go test -run '^(%s)$' --tags=%s %s`+"\n", strings.Join(tests, "|"), strings.Join(tags, ","), strings.Join(dirs, " ")) //nolint:forbidigo
		hasErr = execTest(tests, tags, dirs, *testOptions)
	}
	if hasErr != nil {
		fmt.Println(hasErr)
		os.Exit(1)
	}
}
