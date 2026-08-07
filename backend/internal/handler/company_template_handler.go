package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type CompanyTemplateHandler struct {
	svc *service.CompanyTemplateService
}

func NewCompanyTemplateHandler(svc *service.CompanyTemplateService) *CompanyTemplateHandler {
	return &CompanyTemplateHandler{svc: svc}
}

func (h *CompanyTemplateHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetByID)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
}

func (h *CompanyTemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)

	templates, err := h.svc.List(r.Context(), companyID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Prepare response that maps Tags JSON string back to arrays
	resp := make([]map[string]interface{}, 0, len(templates))
	for _, t := range templates {
		resp = append(resp, map[string]interface{}{
			"id":          t.ID,
			"title":       t.Title,
			"type":        t.Type,
			"description": t.Description,
			"tags":        t.GetTags(),
			"created_at":  t.CreatedAt,
			"updated_at":  t.UpdatedAt,
		})
	}

	pkgresponse.JSON(w, http.StatusOK, resp, nil, requestID)
}

func (h *CompanyTemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)

	var req service.CreateCompanyTemplateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json body", []string{err.Error()}), requestID)
		return
	}

	t, err := h.svc.Create(r.Context(), companyID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":          t.ID,
		"title":       t.Title,
		"type":        t.Type,
		"description": t.Description,
		"tags":        t.GetTags(),
		"created_at":  t.CreatedAt,
		"updated_at":  t.UpdatedAt,
	}, nil, requestID)
}

func (h *CompanyTemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	id := chi.URLParam(r, "id")

	t, err := h.svc.GetByID(r.Context(), companyID, id)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"id":          t.ID,
		"title":       t.Title,
		"type":        t.Type,
		"description": t.Description,
		"tags":        t.GetTags(),
		"created_at":  t.CreatedAt,
		"updated_at":  t.UpdatedAt,
	}, nil, requestID)
}

func (h *CompanyTemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	id := chi.URLParam(r, "id")

	var req service.UpdateCompanyTemplateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json body", []string{err.Error()}), requestID)
		return
	}

	t, err := h.svc.Update(r.Context(), companyID, id, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"id":          t.ID,
		"title":       t.Title,
		"type":        t.Type,
		"description": t.Description,
		"tags":        t.GetTags(),
		"created_at":  t.CreatedAt,
		"updated_at":  t.UpdatedAt,
	}, nil, requestID)
}

func (h *CompanyTemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), companyID, id); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "deleted"}, nil, requestID)
}
