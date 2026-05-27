package fetch

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Pull(url, srcType, destDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch: status %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "red-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err = io.Copy(tmp, resp.Body); err != nil {
		return err
	}
	tmp.Close()

	switch srcType {
	case "tar.gz", "tgz":
		return extractTarGz(tmp.Name(), destDir)
	case "zip":
		return extractZip(tmp.Name(), destDir)
	default:
		return fmt.Errorf("fetch: unsupported type %q", srcType)
	}
}

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err = writeEntry(dest, hdr.Name, hdr.FileInfo().IsDir(), tr); err != nil {
			return err
		}
	}
	return nil
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		err = writeEntry(dest, f.Name, f.FileInfo().IsDir(), rc)
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func writeEntry(dest, name string, isDir bool, r io.Reader) error {
	parts := strings.SplitN(filepath.ToSlash(name), "/", 2)
	if len(parts) < 2 {
		return nil
	}
	rel := parts[1]
	if rel == "" {
		return nil
	}

	target := filepath.Join(dest, filepath.FromSlash(rel))
	if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("fetch: path traversal blocked: %s", name)
	}

	if isDir {
		return os.MkdirAll(target, 0755)
	}

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	return err
}
