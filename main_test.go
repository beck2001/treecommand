package main

import (
	"bytes"
	"testing"
)

const testCurrentDirResult = `|____.idea
`

const testCurrentDirResultFiles = `|----.idea
|       |----.gitignore (256b)
|       |----modules.xml (281b)
|       |----treecommand.iml (330b)
|       |____workspace.xml (3023b)
|----main.go (2784b)
|____main_test.go (987b)`

func TestTreeCurrentDir(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, ".", false)
	if err != nil {
		t.Errorf("dirTree command test failed - %v", err)
	}
	result := out.String()
	if result != testCurrentDirResult {
		t.Errorf("dirTree command test failed - result mismatch %v", err)
	}
}

func TestTreeCurrentDirWithFiles(t *testing.T) {
	out := new(bytes.Buffer)
	err := dirTree(out, ".", true)
	if err != nil {
		t.Errorf("dirTree command test failed - %v", err)
	}
	result := out.String()
	if result != testCurrentDirResultFiles {
		t.Errorf("dirTree command failed - result mismatch\nExpected: %v\nGot: %v", testCurrentDirResultFiles, result)
	}
}
