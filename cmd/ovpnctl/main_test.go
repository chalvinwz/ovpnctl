package main

import (
	"errors"
	"testing"
)

func TestMainSuccessDoesNotExit(t *testing.T) {
	origExec, origExit := execute, exitFn
	defer func() { execute, exitFn = origExec, origExit }()

	execute = func() error { return nil }
	exitCalled := false
	exitFn = func(code int) { exitCalled = true }

	main()

	if exitCalled {
		t.Fatalf("did not expect exit on success")
	}
}

func TestMainErrorExitsWithCodeOne(t *testing.T) {
	origExec, origExit := execute, exitFn
	defer func() { execute, exitFn = origExec, origExit }()

	execute = func() error { return errors.New("boom") }
	code := 0
	exitFn = func(c int) { code = c }

	main()

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}
