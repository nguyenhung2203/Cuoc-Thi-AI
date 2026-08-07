package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) Routes(r chi.Router) {
	r.Post("/upload", h.UploadFile)
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)

	// Giới hạn file tối đa 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		pkgresponse.Error(w, errors.NewValidation("file is too large (max 10MB)", nil), requestID)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		pkgresponse.Error(w, errors.NewValidation("missing file in request", nil), requestID)
		return
	}
	defer file.Close()

	// Tạo tên file an toàn để tránh trùng lặp
	ext := filepath.Ext(header.Filename)
	originalName := strings.TrimSuffix(header.Filename, ext)
	// Thay khoảng trắng và các ký tự đặc biệt nếu cần
	safeName := fmt.Sprintf("%s-%d-%s%s", strings.ReplaceAll(originalName, " ", "_"), time.Now().Unix(), uuid.New().String()[:8], ext)

	// Tạo thư mục uploads nếu chưa có
	workDir, _ := os.Getwd()
	uploadDir := filepath.Join(workDir, "uploads")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to create upload directory"), requestID)
		return
	}

	dstPath := filepath.Join(uploadDir, safeName)
	dst, err := os.Create(dstPath)
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to save file"), requestID)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to write file"), requestID)
		return
	}

	// Trả về URL đường dẫn tương đối để phía Frontend tự ghép domain (hoặc trả đường dẫn absolute)
	fileURL := fmt.Sprintf("/uploads/%s", safeName)

	respData := map[string]interface{}{
		"url":       fileURL,
		"filename":  header.Filename,
		"size":      header.Size,
		"mime_type": header.Header.Get("Content-Type"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": respData,
	})
}
