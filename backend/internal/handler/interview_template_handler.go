package handler

import (
	"encoding/json"
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

type InterviewTemplateHandler struct {
	svc *service.InterviewTemplateService
}

func NewInterviewTemplateHandler(svc *service.InterviewTemplateService) *InterviewTemplateHandler {
	return &InterviewTemplateHandler{svc: svc}
}

func (h *InterviewTemplateHandler) Routes(r chi.Router) {
	r.Route("/templates", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Route("/{template_id}", func(r chi.Router) {
			r.Get("/", h.GetByID)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})
}

func (h *InterviewTemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	p := pagination.FromRequest(r)

	templates, total, err := h.svc.List(r.Context(), companyID, p)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	items := make([]response.InterviewTemplateItem, 0, len(templates))
	for _, t := range templates {
		items = append(items, toInterviewTemplateItem(t))
	}

	meta := &pkgresponse.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: pagination.CalcTotalPages(total, p.PageSize),
	}
	pkgresponse.JSON(w, http.StatusOK, items, meta, requestID)
}

func (h *InterviewTemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	templateID := chi.URLParam(r, "template_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	t, err := h.svc.GetByID(r.Context(), companyID, templateID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, toInterviewTemplateItem(*t), nil, requestID)
}

func (h *InterviewTemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.CreateInterviewTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	t, err := h.svc.Create(r.Context(), companyID, userID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusCreated, toInterviewTemplateItem(*t), nil, requestID)
}

func (h *InterviewTemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	templateID := chi.URLParam(r, "template_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.UpdateInterviewTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}

	if err := h.svc.Update(r.Context(), companyID, templateID, &req); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "updated successfully"}, nil, requestID)
}

func (h *InterviewTemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	templateID := chi.URLParam(r, "template_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	if err := h.svc.Delete(r.Context(), companyID, templateID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"}, nil, requestID)
}

// toInterviewTemplateItem maps models.InterviewTemplate → response.InterviewTemplateItem
func toInterviewTemplateItem(t models.InterviewTemplate) response.InterviewTemplateItem {
	item := response.InterviewTemplateItem{
		ID:              t.ID,
		Name:            t.Name,
		Type:            t.Type,
		DurationMinutes: t.DurationMinutes,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
	if t.CompanyID.Valid {
		item.CompanyID = &t.CompanyID.String
	}
	if t.CreatedBy.Valid {
		item.CreatedBy = &t.CreatedBy.String
	}
	if t.Description.Valid {
		item.Description = &t.Description.String
	}
	if len(t.ConfigJSON) > 0 {
		item.Config = json.RawMessage(t.ConfigJSON)
	}
	return item
}
