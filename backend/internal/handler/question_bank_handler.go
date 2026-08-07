package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/dto/response"
	"backend/internal/middleware"
	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/pkg/validator"
	"backend/internal/service"
)

type QuestionBankHandler struct {
	svc *service.QuestionBankService
}

func NewQuestionBankHandler(svc *service.QuestionBankService) *QuestionBankHandler {
	return &QuestionBankHandler{svc: svc}
}

func (h *QuestionBankHandler) Routes(r chi.Router) {
	r.Route("/question-bank", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Route("/{question_id}", func(r chi.Router) {
			r.Get("/", h.GetByID)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})
}

func (h *QuestionBankHandler) List(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	q := r.URL.Query()
	jobID := q.Get("job_id")
	questionType := q.Get("question_type")
	level := q.Get("level")
	keyword := q.Get("keyword")
	p := pagination.FromRequest(r)

	questions, total, err := h.svc.List(r.Context(), companyID, jobID, questionType, level, keyword, p)
	if err != nil {
		log.Printf("question bank list failed company_id=%s: %v", companyID, err)
		writeServiceError(w, err, requestID)
		return
	}

	items := make([]response.QuestionBankItem, 0, len(questions))
	for _, q := range questions {
		items = append(items, toQuestionBankItem(q))
	}

	meta := &pkgresponse.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: pagination.CalcTotalPages(total, p.PageSize),
	}
	pkgresponse.JSON(w, http.StatusOK, items, meta, requestID)
}

func (h *QuestionBankHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	questionID := chi.URLParam(r, "question_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	qb, err := h.svc.GetByID(r.Context(), companyID, questionID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, toQuestionBankItem(*qb), nil, requestID)
}

func (h *QuestionBankHandler) Create(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.CreateQuestionBankRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	qb, err := h.svc.Create(r.Context(), companyID, userID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusCreated, toQuestionBankItem(*qb), nil, requestID)
}

func (h *QuestionBankHandler) Update(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	questionID := chi.URLParam(r, "question_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.UpdateQuestionBankRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}

	if err := h.svc.Update(r.Context(), companyID, questionID, &req); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "updated successfully"}, nil, requestID)
}

func (h *QuestionBankHandler) Delete(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	questionID := chi.URLParam(r, "question_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	if err := h.svc.Delete(r.Context(), companyID, questionID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"}, nil, requestID)
}

// toQuestionBankItem maps models.QuestionBank → response.QuestionBankItem
func toQuestionBankItem(q models.QuestionBank) response.QuestionBankItem {
	item := response.QuestionBankItem{
		ID:            q.ID,
		QuestionText:  q.QuestionText,
		QuestionType:  q.QuestionType,
		IsAIGenerated: q.IsAIGenerated,
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
	}
	if q.CompanyID.Valid {
		item.CompanyID = &q.CompanyID.String
	}
	if q.JobID.Valid {
		item.JobID = &q.JobID.String
	}
	if q.CreatedBy.Valid {
		item.CreatedBy = &q.CreatedBy.String
	}
	if q.Level.Valid {
		item.Level = &q.Level.String
	}
	if len(q.SkillTags) > 0 {
		var tags []string
		json.Unmarshal(q.SkillTags, &tags)
		item.SkillTags = tags
	}
	if len(q.ExpectedSignals) > 0 {
		var signals []string
		json.Unmarshal(q.ExpectedSignals, &signals)
		item.ExpectedSignals = signals
	}
	return item
}
