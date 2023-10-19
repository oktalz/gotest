package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/akamensky/argparse"
)

func main() {
	args := argparse.NewParser("gotest", "Runs Go tests")
	args.SetHelp("", "help")

	tag := args.String("t", "tags", &argparse.Options{Required: true, Help: "build tags"})
	runs := args.Int("n", "runs", &argparse.Options{Required: false, Default: 1, Help: "different runs to split tests"})
	shift := args.Int("s", "shift", &argparse.Options{Required: false, Default: 0, Help: "shift to start with different index"})
	mode := args.Selector("m", "mode", []string{"list", "dir", "run"}, &argparse.Options{Required: false, Default: "run", Help: "list - list tests, dir - list tests with dir, run - normal run]"})
	info := args.Flag("", "info", &argparse.Options{Required: false, Default: false, Help: "print commands that are executed"})
	parallel := args.Int("p", "parallel", &argparse.Options{Required: false, Help: "number of tests run in parallel"})
	verbose := args.Flag("v", "verbose", &argparse.Options{Required: false, Help: "verbose mode"})

	err := args.Parse(os.Args)
	if err != nil {
		fmt.Print(args.Usage(err)) //nolint:forbidigo
		os.Exit(1)
	}
	tags := strings.Split(*tag, ",")

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

	for i := *shift; i < len(result); i += *runs {
		switch *mode {
		case "list":
			fmt.Println(result[i].Name) //nolint:forbidigo
		case "dir":
			fmt.Printf("%s,%s\n", result[i].Dir, result[i].Name) //nolint:forbidigo
		case "run":
			tests = append(tests, result[i].Name)
			dirs = append(dirs, "./"+result[i].Dir)
		}
	}
	if *mode == "run" {
		// fmt.Printf(`go test -run '^(%s)$' --tags=%s %s`+"\n", strings.Join(tests, "|"), strings.Join(tags, ","), strings.Join(dirs, " ")) //nolint:forbidigo
		hasErr = execTest(tests, tags, dirs, *info, verbose, parallel)
	}
	if hasErr != nil {
		fmt.Println(hasErr)
		os.Exit(1)
	}
}
