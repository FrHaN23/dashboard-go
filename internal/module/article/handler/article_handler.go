package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/frhan23/dashboard-go/internal/entity"
	"github.com/frhan23/dashboard-go/internal/module/article/usecase"
	"github.com/frhan23/dashboard-go/internal/transport"
)

type ArticleHandler struct {
	articleUC usecase.ArticleUsecase
}

func NewArticleHandler(uc usecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{articleUC: uc}
}

func (h *ArticleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /article/", h.CreateArticle)
	mux.HandleFunc("GET /article", h.GetArticles)
	mux.HandleFunc("GET /article/{id}", h.GetArticleByID)
	mux.HandleFunc("POST /article/{id}", h.UpdateArticle)
	mux.HandleFunc("PUT /article/{id}", h.UpdateArticle)
	mux.HandleFunc("PATCH /article/{id}", h.UpdateArticle)
	mux.HandleFunc("DELETE /article/{id}", h.DeleteArticle)
}

type ArticleRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

func (req *ArticleRequest) Validate() *entity.AppError {
	if req.Title == "" || len(req.Title) < 20 {
		return &entity.AppError{
			Code:    entity.ErrorCodeBadRequest,
			Message: "Title is required and must be at least 20 characters",
		}
	}
	if req.Content == "" || len(req.Content) < 200 {
		return &entity.AppError{
			Code:    entity.ErrorCodeBadRequest,
			Message: "Content is required and must be at least 200 characters",
		}
	}
	if req.Category == "" || len(req.Category) < 3 {
		return &entity.AppError{
			Code:    entity.ErrorCodeBadRequest,
			Message: "Category is required and must be at least 3 characters",
		}
	}
	status := strings.ToLower(req.Status)
	if status != "publish" && status != "draft" && status != "thrash" {
		return &entity.AppError{
			Code:    entity.ErrorCodeBadRequest,
			Message: "Status is required and must be either 'publish', 'draft', or 'thrash'",
		}
	}
	return nil
}

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, &entity.AppError{Code: entity.ErrorCodeBadRequest, Message: "Invalid JSON payload"})
		return
	}

	if appErr := req.Validate(); appErr != nil {
		transport.WriteError(w, appErr)
		return
	}

	article := entity.Article{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}

	err := h.articleUC.Create(r.Context(), &article)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteJSON(w, http.StatusCreated, map[string]any{})
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitStr := query.Get("limit")
	pageStr := query.Get("page")
	statusStr := query.Get("status")

	if limitStr == "" {
		limitStr = "10"
	}

	if pageStr == "" || pageStr == "0" {
		pageStr = "1"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		transport.WriteError(w, &entity.AppError{Code: entity.ErrorCodeBadRequest, Message: "Invalid limit parameter"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		transport.WriteError(w, &entity.AppError{Code: entity.ErrorCodeBadRequest, Message: "Invalid page parameter"})
		return
	}

	offset := (page - 1) * limit

	articles, err := h.articleUC.Fetch(r.Context(), limit, offset, statusStr)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteJSON(w, http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		transport.WriteError(w, &entity.AppError{Code: entity.ErrorCodeBadRequest, Message: "Invalid ID parameter"})
		return
	}

	article, err := h.articleUC.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "article not found" {
			transport.WriteError(w, &entity.AppError{
				Code:    entity.ErrorCodeNotFound,
				Message: "Article not found",
			})
			return
		}

		transport.WriteError(w, err)
		return
	}

	transport.WriteJSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		transport.WriteError(w, &entity.AppError{Code: entity.ErrorCodeBadRequest, Message: "Invalid ID parameter"})
		return
	}

	var req ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, &entity.AppError{Code: entity.ErrorCodeBadRequest, Message: "Invalid JSON payload"})
		return
	}

	if appErr := req.Validate(); appErr != nil {
		transport.WriteError(w, appErr)
		return
	}

	article := entity.Article{
		ID:       id,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}

	err = h.articleUC.Update(r.Context(), &article)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteJSON(w, http.StatusOK, map[string]any{})
}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		transport.WriteError(w, &entity.AppError{Code: entity.ErrorCodeBadRequest, Message: "Invalid ID parameter"})
		return
	}

	err = h.articleUC.Delete(r.Context(), id)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteJSON(w, http.StatusOK, map[string]any{})
}
