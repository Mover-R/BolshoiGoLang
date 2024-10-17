package server

import (
	"BolshiGoLang/internal/pkg/storage"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value string `json:"value"`
}

func NewServer(host string, st *storage.Storage) *Server {
	s := &Server{
		host:    host,
		storage: st,
	}

	return s
}

func (r Server) newAPI() *gin.Engine {
	engine := gin.New()

	engine.PUT("scalar/set/:key", r.handlerSet)
	engine.GET("scalar/get/:key", r.handlerGet)
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	return engine
}

func (r Server) handlerSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.Set(key, v.Value)
}

func (r Server) handlerGet(ctx *gin.Context) {
	key := ctx.Param("key")

	v, err := r.storage.Get(key)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: v})
}

func (r *Server) Start() {
	r.newAPI().Run(r.host)
}
