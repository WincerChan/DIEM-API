package filefactory

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	var (
		in  = "file.go"
		out = "file_cp.go"
	)
	CopyFile(in, out)
	inInfo, _ := os.Stat(in)
	outInfo, _ := os.Stat(out)
	if inInfo.Mode() != outInfo.Mode() {
		t.Errorf("Input file mode is %s, output file mode is %s", inInfo.Mode(), outInfo.Mode())
	}
	if inInfo.Size() != outInfo.Size() {
		t.Errorf("Input file size is %d, output file mode is %d", inInfo.Size(), outInfo.Size())
	}
	os.Remove(out)
}

func TestNewFile(t *testing.T) {
	var (
		in = "fh/df.txt"
	)
	_ = NewFile(in)
	_, err := os.Stat(in)
	if err != nil {
		t.Errorf("New file error %v", err)
	}
	os.Remove(in)
	os.Remove(filepath.Dir(in))
}
