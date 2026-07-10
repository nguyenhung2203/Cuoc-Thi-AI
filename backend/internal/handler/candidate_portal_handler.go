package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/response"
	"backend/internal/service"
)

type CandidatePortalHandler struct {
	svc     *service.CandidatePortalService
	fileSvc *service.FileService
}

func NewCandidatePortalHandler(svc *service.CandidatePortalService, fileSvc *service.FileService) *CandidatePortalHandler {
	return &CandidatePortalHandler{svc: svc, fileSvc: fileSvc}
}

func (h *CandidatePortalHandler) Routes(r chi.Router) {
	r.Get("/dashboard", h.GetDashboardStats)
	r.Get("/interviews", h.GetInterviews)
	r.Get("/profile", h.GetProfile)
	r.Post("/cv", h.UploadCV)
	r.Post("/jobs/{jobID}/apply", h.ApplyJob)
	r.Get("/jobs/{jobID}/check-applied", h.CheckApplied)
	r.Get("/applications", h.GetApplications)
}

func (h *CandidatePortalHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	applications, err := h.svc.GetApplications(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, applications, nil, requestID)
}

func (h *CandidatePortalHandler) ApplyJob(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	jobID := chi.URLParam(r, "jobID")

	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		response.Error(w, errors.NewBadRequest("failed to parse form data"), requestID)
		return
	}

	file, header, err := r.FormFile("cv_file")
	var cvFileID, cvOriginalName string
	if err == nil && file != nil {
		defer file.Close()
		fileRecord, uploadErr := h.fileSvc.ProcessUpload(r.Context(), file, header, userID, "", "cv")
		if uploadErr == nil {
			cvFileID = fileRecord.ID
			cvOriginalName = fileRecord.OriginalName
		} else {
			// fallback if upload fails but we still want to apply
			cvFileID = "local-" + header.Filename
			cvOriginalName = header.Filename
		}
	}

	err = h.svc.ApplyForJob(r.Context(), userID, jobID, cvFileID, cvOriginalName)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, requestID)
		} else {
			response.Error(w, errors.NewInternal("failed to apply for job"), requestID)
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "applied successfully"}, nil, requestID)
}

func (h *CandidatePortalHandler) CheckApplied(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}
	jobID := chi.URLParam(r, "jobID")

	applied, err := h.svc.CheckApplied(r.Context(), userID, jobID)
	if err != nil {
		response.Error(w, errors.NewInternal("failed to check application status"), requestID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"has_applied": applied}, nil, requestID)
}

func (h *CandidatePortalHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	stats, err := h.svc.GetDashboardStats(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, stats, nil, requestID)
}

func (h *CandidatePortalHandler) GetInterviews(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	interviews, err := h.svc.GetInterviews(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, interviews, nil, requestID)
}

func (h *CandidatePortalHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	profile, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, profile, nil, requestID)
}

func (h *CandidatePortalHandler) UploadCV(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		response.Error(w, errors.NewBadRequest("failed to parse form data"), requestID)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, errors.NewBadRequest("file is required"), requestID)
		return
	}
	defer file.Close()

	fileRecord, err := h.fileSvc.ProcessUpload(r.Context(), file, header, userID, "", "cv")
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, requestID)
		} else {
			response.Error(w, errors.NewInternal("failed to process file upload"), requestID)
		}
		return
	}

	err = h.svc.UploadCV(r.Context(), userID, header.Filename, fileRecord.ID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	cvUrl := "http://localhost:18080/uploads/" + fileRecord.StorageKey

	response.JSON(w, http.StatusOK, map[string]string{
		"message":   "CV uploaded",
		"file_name": fileRecord.OriginalName,
		"cv_url":    cvUrl,
	}, nil, requestID)
}

