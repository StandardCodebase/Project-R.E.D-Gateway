package fs

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

func SecureJoin(baseDir, requestedPath string) (string, error) {
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}

	targetAbs, err := filepath.Abs(filepath.Join(baseDir, requestedPath))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(targetAbs, baseAbs+string(os.PathSeparator)) && targetAbs != baseAbs {
		return "", os.ErrPermission
	}

	return targetAbs, nil
}

func ReadFileWithContext(ctx context.Context, path string) ([]byte, error) {
	type result struct {
		data []byte
		err  error
	}

	ch := make(chan result, 1)
	go func() {
		data, err := os.ReadFile(path)
		ch <- result{data, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-ch:
		return r.data, r.err
	}
}