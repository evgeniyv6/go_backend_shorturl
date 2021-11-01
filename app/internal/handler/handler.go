package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/static"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/redisdb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	handler struct {
		netProtocol string
		host        string
		db          redisdb.DBAction
		logger      *zap.SugaredLogger
	}
)

func NewGinRouter(netProtocol, address string, db redisdb.DBAction, logger *zap.SugaredLogger) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode) // for production
	r := gin.New()
	r.Use(static.Serve("/", static.LocalFile("./web", true)))
	h := handler{netProtocol, address, db, logger}
	r.POST("/cut", respHandler(h.cut))
	r.GET("/:hash/info", respHandler(h.expand))
	r.GET("/:hash", h.redirect)
	return r
}

func respHandler(h func(ctx *gin.Context) (interface{}, int, error)) gin.HandlerFunc {
	return func(context *gin.Context) {
		res, code, err := h(context)
		if err != nil {
			res = err.Error()
		}
		context.JSON(code, resp{Data: res, Success: err == nil})
	}
}

func (h *handler) cut(ctx *gin.Context) (interface{}, int, error) {
	var body struct {
		URL string `json:"url" binding:"required"`
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.logger.Errorw("Couldnot get body from request.", "err", err)
	}
	h.logger.Infow("Link for processing.", "link", body.URL)

	uri, err := url.ParseRequestURI(body.URL)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	hashString, err := h.db.Save(uri.String())
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Couldnot save to database: %v", err)
	}

	newUrl := url.URL{
		Scheme: h.netProtocol,
		Host:   h.host,
		Path:   hashString,
	}

	h.logger.Infof("Short link: %v", newUrl.String())

	return newUrl.String(), http.StatusCreated, nil
}

func (h *handler) expand(ctx *gin.Context) (interface{}, int, error) {
	hash := ctx.Param("hash")
	h.logger.Infof("Hash info: %s", hash)

	res, err := h.db.GetInfo(hash)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}

	return res, http.StatusOK, nil
}

func (h *handler) redirect(ctx *gin.Context) {
	hash := ctx.Param("hash")

	link, err := h.db.GetLink(hash)

	h.logger.Infof("Redirect to %s", link)

	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, link)
}
