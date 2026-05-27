package store

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Article struct {
	Title    string
	Path     string
	FullPath string
}

type Section struct {
	Name     string
	Articles []Article
	Sub      map[string]*Section
}

type Store struct {
	mu      sync.RWMutex
	root    map[string]*Section
	dataDir string
}

func New(dataDir string) *Store {
	return &Store{dataDir: dataDir, root: map[string]*Section{}}
}

func (s *Store) Reload() error {
	root := map[string]*Section{}

	err := filepath.WalkDir(s.dataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return err
		}

		rel, _ := filepath.Rel(s.dataDir, path)
		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) < 2 {
			return nil
		}

		top := parts[0]
		if root[top] == nil {
			root[top] = &Section{Name: top, Sub: map[string]*Section{}}
		}

		article := Article{
			Title:    titleFromFilename(parts[len(parts)-1]),
			Path:     "/" + filepath.ToSlash(rel[:len(rel)-3]),
			FullPath: path,
		}

		if len(parts) == 2 {
			root[top].Articles = append(root[top].Articles, article)
		} else {
			sub := parts[1]
			if root[top].Sub[sub] == nil {
				root[top].Sub[sub] = &Section{Name: sub, Sub: map[string]*Section{}}
			}
			root[top].Sub[sub].Articles = append(root[top].Sub[sub].Articles, article)
		}
		return nil
	})

	if err != nil {
		return err
	}

	s.mu.Lock()
	s.root = root
	s.mu.Unlock()
	return nil
}

func (s *Store) Root() map[string]*Section {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.root
}

func (s *Store) Resolve(urlPath string) (string, bool) {
	rel := filepath.FromSlash(strings.TrimPrefix(urlPath, "/")) + ".md"
	full := filepath.Join(s.dataDir, rel)
	data, err := os.ReadFile(full)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func titleFromFilename(name string) string {
	name = strings.TrimSuffix(name, ".md")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	if len(name) == 0 {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
