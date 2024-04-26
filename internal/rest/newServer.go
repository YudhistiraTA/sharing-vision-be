package rest

import (
	"net/http"
	"time"

	"github.com/YudhistiraTA/sharing-vision-be/model"
	"github.com/YudhistiraTA/sharing-vision-be/service/article"
	"github.com/go-chi/chi"
)

// NewServer creates a new HTTP server.
// It returns a pointer to the http.Server.
// Consumes the services to be started
// and the configuration for the server.
//
// It currently only accepts the user service.
func NewServer(conf model.ServerConfig, as *article.ArticleService) *http.Server {
	r := chi.NewRouter()

	Article(r, conf.Log, as)

	return &http.Server{
		Addr:         conf.Addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
