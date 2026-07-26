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
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	h.handleUpload(w, r, "generic")
}

func (h *FileHandler) handleUpload(w http.ResponseWriter, r *http.Request, fileType string) {
	requestID := getRequestID(r)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)

	companyID := chi.URLParam(r, "company_id") // Might be empty for generic files

	if appErr := parseUploadForm(w, r, h.svc.MaxUploadBytes()); appErr != nil {
		pkgresponse.Error(w, appErr, requestID)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("file", []string{"file part is required"}), requestID)
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

	isAdmin := userRole == "admin"
	isRecruiter := userRole == "recruiter"
	companyID := r.URL.Query().Get("company_id")

	if isRecruiter && companyID == "" {
		pkgresponse.Error(w, apierrors.NewValidation("company_id query param is required for recruiter", nil), requestID)
		return
	}

	url, expiresAt, err := h.svc.GetSignedURL(r.Context(), fileID, userID, companyID, isRecruiter, isAdmin)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Admin can read any file bypassing company scope, so leave an audit trail.
	if isAdmin {
		if ah := middleware.GetAuditHelper(r); ah != nil {
			ah.Log("file:admin_access", "file", fileID, companyID, nil, map[string]string{"file_id": fileID})
		}
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
