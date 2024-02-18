package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type TestOptions struct {
	Tags     *string
	Runs     *int
	Shift    *int
	Mode     *string
	Info     *bool
	Parallel *int
	Verbose  *bool
	Race     *bool
	Args     *string
	Timeout  *string
}

func parseTestOptions() (*TestOptions, error) {
	parser := argparse.NewParser("gotest", "Runs Go tests")
	parser.SetHelp("", "help")

	opts := &TestOptions{
		Tags:  parser.String("t", "tags", &argparse.Options{Required: true, Help: "build tags"}),
		Runs:  parser.Int("n", "runs", &argparse.Options{Required: false, Default: 1, Help: "different runs to split tests"}),
		Shift: parser.Int("s", "shift", &argparse.Options{Required: false, Default: 0, Help: "shift to start with different index"}),
		Mode:  parser.Selector("m", "mode", []string{"list", "dir", "run", "seq"}, &argparse.Options{Required: false, Default: "run", Help: "list - list tests, dir - list tests with dir, run - normal run]"}),
		Info:  parser.Flag("", "info", &argparse.Options{Required: false, Default: false, Help: "print commands that are executed"}),
		// Parallel: parser.Int("p", "parallel", &argparse.Options{Required: false, Help: "number of tests run in parallel"}),
		Verbose: parser.Flag("v", "verbose", &argparse.Options{Required: false, Help: "verbose mode"}),
		Race:    parser.Flag("", "race", &argparse.Options{Required: false, Help: "enable data race detection"}),
		Args:    parser.String("", "args", &argparse.Options{Required: false, Help: "pass command line arguments to test"}),
		Timeout: parser.String("", "timeout", &argparse.Options{Required: false, Help: "timeout to the test"}),
	}

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return nil, err
	}

	return opts, nil
}
