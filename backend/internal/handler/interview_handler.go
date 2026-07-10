package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"backend/internal/livekit"
	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/response"
	"backend/internal/service"
)

type InterviewHandler struct {
	svc *service.InterviewService
}

func NewInterviewHandler(svc *service.InterviewService) *InterviewHandler {
	return &InterviewHandler{svc: svc}
}

// Routes sets up public routes (e.g. joining via token)
func (h *InterviewHandler) Routes(r chi.Router) {
	r.Get("/join/{invite_token}", h.JoinByToken)
}

// ProtectedRoutes sets up company-scoped routes
func (h *InterviewHandler) ProtectedRoutes(r chi.Router) {
	r.With(middleware.RequirePermission("interview:create")).Post("/", h.CreateInterview)
	r.With(middleware.RequirePermission("interview:read")).Get("/", h.ListInterviews)
	r.With(middleware.RequirePermission("interview:read")).Get("/{interview_id}", h.GetInterview)
	r.With(middleware.RequirePermission("interview:update")).Post("/{interview_id}/start", h.StartInterview)
	r.With(middleware.RequirePermission("interview:update")).Post("/{interview_id}/end", h.EndInterview)
	r.With(middleware.RequirePermission("interview:update")).Post("/{interview_id}/cancel", h.CancelInterview)
	r.With(middleware.RequirePermission("interview:read")).Get("/{interview_id}/room", h.GetRoom)
	r.Get("/{interview_id}/room/access-token", h.GetRoomAccessToken)
	r.Post("/{interview_id}/room/token", h.GetRecruiterRoomToken)
}

func (h *InterviewHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	room, err := h.svc.GetRoom(r.Context(), interviewID, companyID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	response.JSON(w, http.StatusOK, room, nil, requestID)
}

func (h *InterviewHandler) CancelInterview(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	if err := h.svc.CancelInterview(r.Context(), interviewID, companyID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "cancelled"}, nil, requestID)
}

func (h *InterviewHandler) ListInterviews(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	items, err := h.svc.ListInterviews(r.Context(), companyID, 100, 0) // stub pagination
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	response.JSON(w, http.StatusOK, items, nil, requestID)
}

func (h *InterviewHandler) GetInterview(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	item, err := h.svc.GetInterview(r.Context(), interviewID, companyID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	response.JSON(w, http.StatusOK, item, nil, requestID)
}

func (h *InterviewHandler) JoinByToken(w http.ResponseWriter, r *http.Request) {
	inviteToken := chi.URLParam(r, "invite_token")
	requestID := r.Context().Value(middleware.CtxRequestID)
	if reqID, ok := requestID.(string); ok {
		requestID = reqID
	} else {
		requestID = ""
	}

	room, err := h.svc.JoinByToken(r.Context(), inviteToken)
	if err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			response.Error(w, appErr, requestID.(string))
		} else {
			response.Error(w, apierrors.NewInternal("failed to join room"), requestID.(string))
		}
		return
	}

	response.JSON(w, http.StatusOK, room, nil, requestID.(string))
}

func (h *InterviewHandler) CreateInterview(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req service.CreateInterviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierrors.NewBadRequest("invalid request body"), requestID)
		return
	}

	req.CompanyID = companyID

	res, err := h.svc.CreateInterview(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log interview creation
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("interview:create", "interview", res.InterviewID, companyID, nil, map[string]interface{}{
			"job_id":       req.JobID,
			"candidate_id": req.CandidateID,
		})
	}

	response.JSON(w, http.StatusCreated, res, nil, requestID)
}

func (h *InterviewHandler) StartInterview(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	err := h.svc.StartInterview(r.Context(), interviewID, companyID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log start
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("interview:start", "interview", interviewID, companyID, map[string]string{"status": "scheduled"}, map[string]string{"status": "active"})
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "interview started"}, nil, requestID)
}

func (h *InterviewHandler) EndInterview(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	err := h.svc.EndInterview(r.Context(), interviewID, companyID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log end
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("interview:end", "interview", interviewID, companyID, map[string]string{"status": "active"}, map[string]string{"status": "completed"})
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "interview ended"}, nil, requestID)
}

func (h *InterviewHandler) GetRoomAccessToken(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	token, err := h.svc.GenerateRoomAccessToken(r.Context(), interviewID, companyID, userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"access_token": token}, nil, requestID)
}

func (h *InterviewHandler) GetRecruiterRoomToken(w http.ResponseWriter, r *http.Request) {
	interviewID := chi.URLParam(r, "interview_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	roomID := "room-" + interviewID

	livekitSecret := os.Getenv("LIVEKIT_API_SECRET")
	if livekitSecret == "" {
		livekitSecret = "devsecret"
	}
	livekitKey := os.Getenv("LIVEKIT_API_KEY")
	if livekitKey == "" {
		livekitKey = "devkey"
	}

	tokenString, err := livekit.GenerateToken(
		livekitKey,
		livekitSecret,
		roomID,
		userID,
		"Recruiter",
		"recruiter",
		interviewID,
	)
	if err != nil {
		response.Error(w, apierrors.NewInternal("Failed to generate token"), requestID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"token":             tokenString,
		"livekit_token":     tokenString,
		"room_access_token": tokenString,
	}, nil, requestID)
}
