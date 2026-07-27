package handler

import (
	"encoding/json"
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
	r.Put("/profile", h.UpdateProfile)
	r.Post("/cv", h.UploadCV)
	r.Post("/jobs/{jobID}/apply", h.ApplyJob)
	r.Get("/jobs/{jobID}/match", h.GetJobMatch)
	r.Get("/applications", h.GetApplications)
	r.Delete("/applications/{id}", h.CancelApplication)
}

func (h *CandidatePortalHandler) ApplyJob(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	jobID := chi.URLParam(r, "jobID")

	if appErr := parseUploadForm(w, r, h.fileSvc.MaxUploadBytes()); appErr != nil {
		response.Error(w, appErr, requestID)
		return
	}

	file, header, err := r.FormFile("cv_file")
	var cvFileID, cvOriginalName string
	if err == nil && file != nil {
		defer file.Close()
		fileRecord, upErr := h.fileSvc.ProcessUpload(r.Context(), file, header, userID, "", "cv")
		if upErr != nil {
			if appErr, ok := errors.IsAppError(upErr); ok {
				response.Error(w, appErr, requestID)
			} else {
				response.Error(w, errors.NewInternal("failed to process cv upload"), requestID)
			}
			return
		}
		cvFileID = fileRecord.ID
		cvOriginalName = header.Filename
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

// GetJobMatch returns an on-demand AI CV↔job match for the Apply page preview.
// When the user has no parsed CV, it returns 200 with has_cv=false so the
// frontend can show a "upload CV" CTA instead of a fake score.
func (h *CandidatePortalHandler) GetJobMatch(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	jobID := chi.URLParam(r, "jobID")

	result, err := h.svc.GetJobMatch(r.Context(), userID, jobID)
	if err != nil {
		if err == service.ErrNoCV {
			response.JSON(w, http.StatusOK, map[string]interface{}{"has_cv": false}, nil, requestID)
			return
		}
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, requestID)
		} else {
			response.Error(w, errors.NewInternal("failed to compute match"), requestID)
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"has_cv":         true,
		"fit_score":      result.FitScore,
		"matched_skills": result.MatchedSkills,
		"missing_skills": result.MissingSkills,
		"summary":        result.Summary,
		"recommendation": result.Recommendation,
	}, nil, requestID)
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

func (h *CandidatePortalHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	var body struct {
		FullName  string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, errors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}

	if err := h.svc.UpdateProfile(r.Context(), userID, body.FullName, body.AvatarURL); err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "updated"}, nil, requestID)
}

func (h *CandidatePortalHandler) UploadCV(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	// Parse multipart form with the configured size cap (413 on overflow)
	if appErr := parseUploadForm(w, r, h.fileSvc.MaxUploadBytes()); appErr != nil {
		response.Error(w, appErr, requestID)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, errors.NewValidation("file", []string{"file part is required"}), requestID)
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

	err = h.svc.UploadCV(r.Context(), userID, header.Filename, fileRecord.ID, fileRecord.StorageKey)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	cvUrl, _ := h.fileSvc.SignedURL(fileRecord.StorageKey)

	// Fetch the freshly parsed profile so the client can show real extracted data.
	profile, _ := h.svc.GetProfile(r.Context(), userID)
	var parsedData interface{}
	if profile != nil {
		parsedData = profile.ParsedData
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message":     "CV uploaded",
		"file_name":   fileRecord.OriginalName,
		"cv_url":      cvUrl,
		"cv_file_id":  fileRecord.ID,
		"parsed_data": parsedData,
	}, nil, requestID)
}

func (h *CandidatePortalHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}

	apps, err := h.svc.GetApplications(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, apps, nil, requestID)
}

func (h *CandidatePortalHandler) CancelApplication(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), requestID)
		return
	}
	appID := chi.URLParam(r, "id")

	err := h.svc.CancelApplication(r.Context(), userID, appID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "withdrawn successfully"}, nil, requestID)
}

