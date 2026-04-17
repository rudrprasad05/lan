package api

import (
	"encoding/json"
	"errors"
	"io"
	"lan-share/daemon/api/res"
	"lan-share/daemon/internal/filetransfer"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type FileResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	CreatedAt   int64  `json:"createdAt"`
	DownloadURL string `json:"downloadUrl"`
}

type FileListResponse struct {
	res.BaseResponse
	Files []FileResponse `json:"files"`
}

type FileUploadResponse struct {
	res.BaseResponse
	File FileResponse `json:"file"`
}

func (s *Server) filesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listFilesHandler(w, r)
	case http.MethodPost:
		s.uploadFileHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listFilesHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	files, err := s.fileStore.List()
	if err != nil {
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	response := FileListResponse{
		BaseResponse: res.NewBaseResponse(),
		Files:        make([]FileResponse, 0, len(files)),
	}

	for _, file := range files {
		response.Files = append(response.Files, toFileResponse(file))
	}

	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("uploading file")

	filename, contentType, body, closeFn, err := extractUpload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer closeFn()

	stored, err := s.fileStore.Save(filename, contentType, body)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("sending file")
	_ = json.NewEncoder(w).Encode(FileUploadResponse{
		BaseResponse: res.NewBaseResponse(),
		File:         toFileResponse(stored),
	})
}

func (s *Server) fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("downlaoding file")

	id := strings.TrimPrefix(r.URL.Path, "/api/files/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	meta, file, err := s.fileStore.Open(id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": meta.Name})
	if disposition == "" {
		disposition = "attachment"
	}

	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Length", strconv.FormatInt(meta.Size, 10))

	http.ServeContent(w, r, meta.Name, time.Unix(meta.CreatedAt, 0), file)
}

func toFileResponse(file filetransfer.StoredFile) FileResponse {
	return FileResponse{
		ID:          file.ID,
		Name:        file.Name,
		Size:        file.Size,
		ContentType: file.ContentType,
		CreatedAt:   file.CreatedAt,
		DownloadURL: "/api/files/" + file.ID,
	}
}

func extractUpload(r *http.Request) (string, string, io.Reader, func(), error) {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, _ := mime.ParseMediaType(contentType)

	if strings.HasPrefix(mediaType, "multipart/") {
		reader, err := r.MultipartReader()
		if err != nil {
			return "", "", nil, noopClose, err
		}

		for {
			part, err := reader.NextPart()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return "", "", nil, noopClose, err
			}

			if part.FileName() == "" {
				_ = part.Close()
				continue
			}

			return part.FileName(), part.Header.Get("Content-Type"), part, func() {
				_ = part.Close()
			}, nil
		}

		return "", "", nil, noopClose, errors.New("multipart request did not include a file")
	}

	filename := strings.TrimSpace(r.Header.Get("X-Filename"))
	if filename == "" {
		filename = strings.TrimSpace(r.URL.Query().Get("filename"))
	}
	if filename == "" {
		filename = path.Base(strings.TrimSpace(r.URL.Query().Get("name")))
	}
	if filename == "" || filename == "." || filename == "/" {
		return "", "", nil, noopClose, errors.New("missing filename; send multipart form-data or provide X-Filename")
	}

	return filename, contentType, r.Body, func() {
		_ = r.Body.Close()
	}, nil
}

func noopClose() {}
