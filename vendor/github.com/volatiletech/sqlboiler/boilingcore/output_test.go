package boilingcore

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type NopWriteCloser struct {
	io.Writer
}

func (NopWriteCloser) Close() error {
	return nil
}

func nopCloser(w io.Writer) io.WriteCloser {
	return NopWriteCloser{w}
}

func TestWriteFile(t *testing.T) {
	// t.Parallel() cannot be used

	// set the function pointer back to its original value
	// after we modify it for the test
	saveTestHarnessWriteFile := testHarnessWriteFile
	defer func() {
		testHarnessWriteFile = saveTestHarnessWriteFile
	}()

	var output []byte
	testHarnessWriteFile = func(_ string, in []byte, _ os.FileMode) error {
		output = in
		return nil
	}

	buf := &bytes.Buffer{}
	writePackageName(buf, "pkg")
	fmt.Fprintf(buf, "func hello() {}\n\n\nfunc world() {\nreturn\n}\n\n\n\n")

	if err := writeFile("", "", buf, true); err != nil {
		t.Error(err)
	}

	if string(output) != "package pkg\n\nfunc hello() {}\n\nfunc world() {\n\treturn\n}\n" {
		t.Errorf("Wrong output: %q", output)
	}
}

func TestFormatBuffer(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "package pkg\n\nfunc() {a}\n")

	// Only test error case - happy case is taken care of by template test
	_, err := formatBuffer(buf)
	if err == nil {
		t.Error("want an error")
	}

	if txt := err.Error(); !strings.Contains(txt, ">>>> func() {a}") {
		t.Error("got:\n", txt)
	}
}

func TestOutputFilenameParts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Filename string

		FirstDir    string
		Normalized  string
		IsSingleton bool
		IsGo        bool
		UsePkg      bool
	}{
		{"templates/00_struct.go.tpl", "templates", "struct.go", false, true, true},
		{"templates/singleton/00_struct.go.tpl", "templates", "struct.go", true, true, true},
		{"templates/notpkg/00_struct.go.tpl", "templates", "notpkg/struct.go", false, true, false},
		{"templates/js/singleton/00_struct.js.tpl", "templates", "js/struct.js", true, false, false},
		{"templates/js/00_struct.js.tpl", "templates", "js/struct.js", false, false, false},
	}

	for i, test := range tests {
		normalized, isSingleton, isGo, usePkg := outputFilenameParts(test.Filename)

		if normalized != test.Normalized {
			t.Errorf("%d) normalized wrong, want: %s, got: %s", i, test.Normalized, normalized)
		}
		if isSingleton != test.IsSingleton {
			t.Errorf("%d) isSingleton wrong, want: %t, got: %t", i, test.IsSingleton, isSingleton)
		}
		if isGo != test.IsGo {
			t.Errorf("%d) isGo wrong, want: %t, got: %t", i, test.IsGo, isGo)
		}
		if usePkg != test.UsePkg {
			t.Errorf("%d) usePkg wrong, want: %t, got: %t", i, test.UsePkg, usePkg)
		}
	}
}
