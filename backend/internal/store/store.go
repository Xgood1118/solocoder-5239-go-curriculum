package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	dataDir string
	mu      sync.RWMutex
}

func NewStore(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	return &Store{dataDir: dataDir}, nil
}

func (s *Store) saveData(filename string, data interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dataDir, filename)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func (s *Store) loadData(filename string, data interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.dataDir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if file == nil {
		return nil
	}

	stat, err := file.Stat()
	if err != nil || stat.Size() == 0 {
		return nil
	}

	if err := json.NewDecoder(file).Decode(data); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
