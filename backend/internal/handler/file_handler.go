package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type FileHandler struct {
	svc *service.FileService
}

func NewFileHandler(svc *service.FileService) *FileHandler {
	return &FileHandler{svc: svc}
}

func (h *FileHandler) Routes(r chi.Router) {
	r.Route("/files", func(r chi.Router) {
		r.Post("/", h.UploadFile)
		r.Get("/{file_id}/signed-url", h.GetSignedURL)
	})
	
	// Company scoped routes for candidate CV
	r.Post("/companies/{company_id}/candidates/{candidate_id}/cv", h.UploadCV)
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	h.handleUpload(w, r, "generic")
}

func (h *FileHandler) UploadCV(w http.ResponseWriter, r *http.Request) {
	h.handleUpload(w, r, "cv")
}



func (h *FileHandler) handleUpload(w http.ResponseWriter, r *http.Request, fileType string) {
	requestID := getRequestID(r)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	
	companyID := chi.URLParam(r, "company_id") // Might be empty for generic files

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		pkgresponse.Error(w, apierrors.NewBadRequest("file too large or invalid form"), requestID)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		pkgresponse.Error(w, apierrors.NewBadRequest("missing file"), requestID)
		return
	}
	defer file.Close()

	fileRecord, err := h.svc.ProcessUpload(r.Context(), file, header, userID, companyID, fileType)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusCreated, fileRecord, nil, requestID)
}

func (h *FileHandler) GetSignedURL(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	fileID := chi.URLParam(r, "file_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	userRole, _ := r.Context().Value(middleware.CtxUserRole).(string)

	isRecruiter := userRole == "recruiter" || userRole == "admin"
	companyID := r.URL.Query().Get("company_id")
	
	if isRecruiter && companyID == "" {
		pkgresponse.Error(w, apierrors.NewValidation("company_id query param is required for recruiter", nil), requestID)
		return
	}

	url, expiresAt, err := h.svc.GetSignedURL(r.Context(), fileID, userID, companyID, isRecruiter)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, signedURLResponse{
		URL:       url,
		ExpiresAt: expiresAt,
	}, nil, requestID)
}

type signedURLResponse struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}
