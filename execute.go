package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func execTest(tests []string, tags []string, dirs []string, info bool, verbose *bool, parralel *int) error {
	// command1 := fmt.Sprintf(`go test%s -run '^(%s)$' --tags=%s %s`+"\n", parallelStr, strings.Join(tests, "|"), strings.Join(tags, ","), strings.Join(dirs, " "))
	// command1 := fmt.Sprintf(` '^(%s)$' --tags=%s %s`+"\n", strings.Join(tags, ","), strings.Join(dirs, " "))
	arg := []string{"test"}
	if verbose != nil && *verbose {
		arg = append(arg, "-v")
	}
	if parralel != nil && *parralel > 0 {
		arg = append(arg, "-p", strconv.Itoa(*parralel))
	}
	arg = append(arg, "-run")
	arg = append(arg, fmt.Sprintf(`^(%s)$`, strings.Join(tests, "|")))
	arg = append(arg, fmt.Sprintf(`--tags=%s`, strings.Join(tags, ",")))
	arg = append(arg, dirs...)

	if info {
		fmt.Println("go", strings.Join(arg, " ")) //nolint:forbidigo
	}
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", arg...)
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Dir = dir
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) //nolint:forbidigo
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	scanner = bufio.NewScanner(stdErr)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) //nolint:forbidigo
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
