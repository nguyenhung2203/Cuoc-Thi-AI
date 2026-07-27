package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"backend/internal/middleware"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/response"
	"backend/internal/pkg/validator"
	"backend/internal/service"
)

type CompanyHandler struct {
	companyService *service.CompanyService
}

func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

// Routes mounts the collection-level endpoints at /companies.
// Per-company GET/PUT live in ScopedRoutes: registering them here is dead code —
// chi's param node (/companies/{company_id} mount) shadows this mount's
// catch-all for any path with an id, so requests never reached these handlers.
func (h *CompanyHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
}

// ScopedRoutes mounts per-company endpoints. MUST be registered inside the
// /companies/{company_id} route group, after CompanyScopeMiddleware.
func (h *CompanyHandler) ScopedRoutes(r chi.Router) {
	r.Get("/", h.Get) // any active member — CompanyScopeMiddleware is enough
	r.With(middleware.RequireCompanyRole("owner", "admin")).Put("/", h.Update)
}

type CreateCompanyReq struct {
	Name     string `json:"name" validate:"required,max=200"`
	Website  string `json:"website" validate:"omitempty,max=255"`
	Industry string `json:"industry" validate:"omitempty,max=100"`
	Size     string `json:"size" validate:"omitempty,max=50"`
}

func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok {
		response.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	companies, err := h.companyService.ListCompanies(r.Context(), userID)
	if err != nil {
		response.Error(w, errors.NewInternal("failed to list companies"), "")
		return
	}

	response.JSON(w, http.StatusOK, companies, nil, "")
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok {
		response.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	var req CreateCompanyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if msgs := validator.Validate(&req); msgs != nil {
		response.Error(w, errors.NewValidation("validation failed", msgs), "")
		return
	}

	company, err := h.companyService.CreateCompany(r.Context(), userID, req.Name, req.Website, req.Industry, req.Size)
	if err != nil {
		response.Error(w, errors.NewInternal("failed to create company"), "")
		return
	}

	// Audit log company creation
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("company:create", "company", company.ID, company.ID, nil, company)
	}

	response.JSON(w, http.StatusCreated, company, nil, "")
}

func (h *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	company, err := h.companyService.GetCompany(r.Context(), companyID)
	if err != nil {
		response.Error(w, errors.NewNotFound("company not found"), "")
		return
	}
	response.JSON(w, http.StatusOK, company, nil, "")
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")

	var req CreateCompanyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if msgs := validator.Validate(&req); msgs != nil {
		response.Error(w, errors.NewValidation("validation failed", msgs), "")
		return
	}

	company, err := h.companyService.UpdateCompany(r.Context(), companyID, req.Name, req.Website, req.Industry, req.Size)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, "")
			return
		}
		response.Error(w, errors.NewInternal("failed to update company"), "")
		return
	}

	// Audit log company update
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("company:update", "company", companyID, companyID, nil, company)
	}

	response.JSON(w, http.StatusOK, company, nil, "")
}
