// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

// +build functionaltests

package tests

import (
	"os"
	"syscall"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/security/policy"
)

func TestRename(t *testing.T) {
	rule := &policy.RuleDefinition{
		ID:         "test_rule",
		Expression: `rename.old_filename == "{{.Root}}/test-rename" && rename.new_filename == "{{.Root}}/test2-rename"`,
	}

	test, err := newTestModule(nil, []*policy.RuleDefinition{rule}, testOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	testOldFile, testOldFilePtr, err := test.Path("test-rename")
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(testOldFile)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	testNewFile, testNewFilePtr, err := test.Path("test2-rename")
	if err != nil {
		t.Fatal(err)
	}

	_, _, errno := syscall.Syscall(syscall.SYS_RENAME, uintptr(testOldFilePtr), uintptr(testNewFilePtr), 0)
	if errno != 0 {
		t.Fatal(err)
	}

	event, _, err := test.GetEvent()
	if err != nil {
		t.Error(err)
	} else {
		if event.GetType() != "rename" {
			t.Errorf("expected rename event, got %s", event.GetType())
		}
	}

	if err := os.Rename(testNewFile, testOldFile); err != nil {
		t.Fatal(err)
	}

	_, _, errno = syscall.Syscall6(syscall.SYS_RENAMEAT, 0, uintptr(testOldFilePtr), 0, uintptr(testNewFilePtr), 0, 0)
	if errno != 0 {
		t.Fatal(err)
	}

	event, _, err = test.GetEvent()
	if err != nil {
		t.Error(err)
	} else {
		if event.GetType() != "rename" {
			t.Errorf("expected rename event, got %s", event.GetType())
		}
	}

	if err := os.Rename(testNewFile, testOldFile); err != nil {
		t.Fatal(err)
	}

	_, _, errno = syscall.Syscall6(316 /* syscall.SYS_RENAMEAT2 */, 0, uintptr(testOldFilePtr), 0, uintptr(testNewFilePtr), 0, 0)
	if errno != 0 {
		t.Fatal(err)
	}

	event, _, err = test.GetEvent()
	if err != nil {
		t.Error(err)
	} else {
		if event.GetType() != "rename" {
			t.Errorf("expected rename event, got %s", event.GetType())
		}
	}
}