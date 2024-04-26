package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/YudhistiraTA/sharing-vision-be/internal/logger"
	"github.com/YudhistiraTA/sharing-vision-be/internal/middlewares"
	"github.com/YudhistiraTA/sharing-vision-be/internal/response"
	"github.com/YudhistiraTA/sharing-vision-be/internal/validation"
	"github.com/YudhistiraTA/sharing-vision-be/model"
	"github.com/YudhistiraTA/sharing-vision-be/service/article"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Handler struct {
	as  *article.ArticleService
	log *zap.Logger
}

func ArticleService(ar *article.ArticleService, log *zap.Logger) *Handler {
	return &Handler{
		as:  ar,
		log: log,
	}
}

func Article(r chi.Router, log *zap.Logger, as *article.ArticleService) {
	h := ArticleService(as, log)
	r.Use(middlewares.IgnoreRequest, middlewares.Timeout, middlewares.CORS, logger.NewLoggingMiddleware(log))
	r.Post("/articles", h.CreateArticle)
	r.Get("/articles/{id}", h.GetArticle)
	r.Put("/articles/{id}", h.UpdateArticle)
	r.Delete("/articles/{id}", h.DeleteArticle)
	r.Get("/articles/{limit}/{offset}", h.GetArticles)
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid request ID", nil)
		return
	}
	article, err := h.as.FindByID(r.Context(), articleID)
	if err != nil {
		response.WriteError(w, err, "Article not found", nil)
		return
	}

	response.WriteSuccess(w, article, http.StatusOK)
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid request ID", nil)
		return
	}
	_, err = h.as.FindByID(r.Context(), articleID)
	if err != nil {
		response.WriteError(w, err, "Article not found", nil)
		return
	}
	err = h.as.Delete(r.Context(), articleID)
	if err != nil {
		response.WriteError(w, err, "Failed to delete article", nil)
		return
	}

	response.WriteSuccess(w, nil, http.StatusOK)
}

func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(chi.URLParam(r, "limit"))
	if err != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid limit", nil)
		return
	}
	offset, err := strconv.Atoi(chi.URLParam(r, "offset"))
	if err != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid offset", nil)
		return
	}
	articles, err := h.as.FindAll(r.Context(), limit, offset)
	if err != nil {
		response.WriteError(w, err, "Failed to get articles", nil)
		return
	}

	response.WriteSuccess(w, articles, http.StatusOK)
}

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article model.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid request body", nil)
		return
	}
	validationError := validation.ValidateStruct(article)
	if validationError != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid request body", validationError)
		return
	}

	err = h.as.Create(r.Context(), &article)
	if err != nil {
		response.WriteError(w, err, "Failed to create article", nil)
		return
	}

	response.WriteSuccess(w, nil, http.StatusCreated)
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid request ID", nil)
		return
	}
	_, err = h.as.FindByID(r.Context(), articleID)
	if err != nil {
		response.WriteError(w, err, "Article not found", nil)
		return
	}
	var article model.Article
	err = json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid request body", nil)
		return
	}
	validationError := validation.ValidateStruct(article)
	if validationError != nil {
		response.WriteError(w, response.ErrInvalidRequest, "Invalid request body", validationError)
		return
	}

	err = h.as.Update(r.Context(), articleID, &article)
	if err != nil {
		response.WriteError(w, err, "Failed to update article", nil)
		return
	}

	response.WriteSuccess(w, nil, http.StatusOK)
}
