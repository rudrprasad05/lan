package filetransfer

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	metaSuffix = ".meta.json"
)

type Service struct {
	dir string
}

type StoredFile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	CreatedAt   int64  `json:"createdAt"`
}

func NewService(dir string) *Service {
	return &Service{dir: dir}
}

func (s *Service) Save(name, contentType string, src io.Reader) (StoredFile, error) {
	if err := s.ensureDir(); err != nil {
		return StoredFile{}, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = "file"
	}

	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(name))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	id, err := randomID()
	if err != nil {
		return StoredFile{}, err
	}

	storedName := filepath.Base(name)
	finalDataPath := s.dataPath(id, storedName)
	tmpDataPath := finalDataPath + ".tmp"

	dst, err := os.Create(tmpDataPath)
	if err != nil {
		return StoredFile{}, err
	}

	size, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil {
		_ = os.Remove(tmpDataPath)
		return StoredFile{}, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmpDataPath)
		return StoredFile{}, closeErr
	}
	if err := os.Rename(tmpDataPath, finalDataPath); err != nil {
		_ = os.Remove(tmpDataPath)
		return StoredFile{}, err
	}

	file := StoredFile{
		ID:          id,
		Name:        storedName,
		Size:        size,
		ContentType: contentType,
		CreatedAt:   time.Now().Unix(),
	}

	if err := writeJSONAtomically(s.metaPath(id), file); err != nil {
		_ = os.Remove(finalDataPath)
		return StoredFile{}, err
	}

	return file, nil
}

func (s *Service) List() ([]StoredFile, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	files := make([]StoredFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), metaSuffix) {
			continue
		}

		meta, err := s.readMeta(strings.TrimSuffix(entry.Name(), metaSuffix))
		if err != nil {
			continue
		}

		files = append(files, meta)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].CreatedAt > files[j].CreatedAt
	})

	return files, nil
}

func (s *Service) Open(id string) (StoredFile, *os.File, error) {
	if !validID(id) {
		return StoredFile{}, nil, os.ErrNotExist
	}

	meta, err := s.readMeta(id)
	if err != nil {
		return StoredFile{}, nil, err
	}

	f, err := os.Open(s.resolveDataPath(id, meta.Name))
	if err != nil {
		return StoredFile{}, nil, err
	}

	return meta, f, nil
}

func (s *Service) ensureDir() error {
	return os.MkdirAll(s.dir, 0o755)
}

func (s *Service) readMeta(id string) (StoredFile, error) {
	var meta StoredFile

	data, err := os.ReadFile(s.metaPath(id))
	if err != nil {
		return StoredFile{}, err
	}

	if err := json.Unmarshal(data, &meta); err != nil {
		return StoredFile{}, err
	}

	return meta, nil
}

func (s *Service) dataPath(id, originalName string) string {
	return filepath.Join(s.dir, id+safeExt(originalName))
}

func (s *Service) metaPath(id string) string {
	return filepath.Join(s.dir, id+metaSuffix)
}

func (s *Service) resolveDataPath(id, originalName string) string {
	currentPath := s.dataPath(id, originalName)
	if _, err := os.Stat(currentPath); err == nil {
		return currentPath
	}

	legacyPath := filepath.Join(s.dir, id+".bin")
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath
	}

	return currentPath
}

func ResolveStorageDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.Join("services", "daemon", "tmp")
	}

	candidates := []string{
		filepath.Join(cwd, "services", "daemon", "tmp"),
		filepath.Join(cwd, "tmp"),
		filepath.Join(cwd, "..", "tmp"),
		filepath.Join(cwd, "..", "..", "tmp"),
	}

	for _, candidate := range candidates {
		parent := filepath.Dir(candidate)
		info, err := os.Stat(parent)
		if err == nil && info.IsDir() {
			return candidate
		}
	}

	return filepath.Join(cwd, "tmp")
}

func writeJSONAtomically(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	return nil
}

func randomID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func validID(id string) bool {
	if len(id) == 0 {
		return false
	}

	if strings.Contains(id, "/") || strings.Contains(id, string(filepath.Separator)) {
		return false
	}

	_, err := hex.DecodeString(id)
	return err == nil && len(id)%2 == 0
}

func safeExt(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" || ext == "." {
		return ""
	}

	if len(ext) > 16 {
		return ""
	}

	for _, r := range ext[1:] {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return ""
		}
	}

	return ext
}
