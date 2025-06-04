package web

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractZip(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	if err := os.Mkdir(src, 0o755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.Mkdir(filepath.Join(src, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "sub", "b.txt"), []byte("b"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	zipPath := filepath.Join(tmp, "test.zip")
	zf, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(zf)
	data, _ := os.ReadFile(filepath.Join(src, "a.txt"))
	w, _ := zw.Create("a.txt")
	w.Write(data)
	zw.Create("sub/")
	data, _ = os.ReadFile(filepath.Join(src, "sub", "b.txt"))
	w, _ = zw.Create("sub/b.txt")
	w.Write(data)
	if err := zw.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	if err := zf.Close(); err != nil {
		t.Fatalf("close file: %v", err)
	}

	dest := filepath.Join(tmp, "dst")
	if err := os.Mkdir(dest, 0o755); err != nil {
		t.Fatalf("mkdir dst: %v", err)
	}
	if err := extractZip(zipPath, dest); err != nil {
		t.Fatalf("extractZip err: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "a.txt")); err != nil {
		t.Fatalf("missing file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "sub", "b.txt")); err != nil {
		t.Fatalf("missing file: %v", err)
	}
}

func TestBuildFileTree(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "dir"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "dir", "a.txt"), []byte("a"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("b"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	node, err := buildFileTree(root, root)
	if err != nil {
		t.Fatalf("build tree err: %v", err)
	}
	if !node.IsDir || len(node.Children) != 2 {
		t.Fatalf("unexpected tree: %+v", node)
	}
}
