package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/pkg/validator"
	"backend/internal/service"
)

type QuestionHandler struct {
	svc *service.QuestionService
}

func NewQuestionHandler(svc *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{svc: svc}
}

func (h *QuestionHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
}

func (h *QuestionHandler) List(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)

	q := r.URL.Query()
	jobID := q.Get("job_id")
	qType := q.Get("type")
	level := q.Get("level")
	keyword := q.Get("keyword")
	p := pagination.FromRequest(r)

	items, total, err := h.svc.List(r.Context(), companyID, jobID, qType, level, keyword, &p)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	meta := &pkgresponse.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      int(total),
		TotalPages: pagination.CalcTotalPages(int(total), p.PageSize),
	}
	pkgresponse.JSON(w, http.StatusOK, items, meta, requestID)
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)

	var req request.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	q, err := h.svc.Create(r.Context(), companyID, userID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusCreated, q, nil, requestID)
}

func (h *QuestionHandler) Update(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	id := chi.URLParam(r, "id")

	var req request.UpdateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	if err := h.svc.Update(r.Context(), companyID, id, &req); err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "updated successfully"}, nil, requestID)
}

func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), companyID, id); err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"}, nil, requestID)
}
