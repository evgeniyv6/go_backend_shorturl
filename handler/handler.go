package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_backend_shorturl/redisdb"
	"net/http"
	"net/url"
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
	}
)

func NewGinRouter(netProtocol, address string, db redisdb.DBAction) *gin.Engine {
	r := gin.Default()
	h := handler{netProtocol, address, db}
	r.POST("/cut", respHandler(h.cut))
	r.GET("/:hash/info", respHandler(h.expand))
	r.GET("/:hash", h.redirect)
	return r
}

func respHandler(h func(ctx *gin.Context) (interface{}, int, error)) gin.HandlerFunc {
	return func(context *gin.Context) {
		res, errCode, err := h(context)
		if err != nil {
			res = err.Error()
		}
		context.Status(errCode)
		zap.S().Info(res)
		context.JSON(http.StatusOK, resp{Data: res, Success: err == nil})
	}
}

func (h *handler) cut(ctx *gin.Context) (interface{}, int, error) {
	var body struct {
		URL string `json:"url" binding:"required"`
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		zap.S().Errorf("Couldnot get body from request, error: %s", err)
	}
	zap.S().Infof("link for processing: %s",body.URL)

	uri, err := url.ParseRequestURI(body.URL)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	hashString, err := h.db.Save(uri.String())
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("Couldnot save to database: %v", err)
	}

	newUrl := url.URL{
		Scheme: h.netProtocol,
		Host:   h.host,
		Path:   hashString,
	}

	zap.S().Infof("Short link: %v", newUrl.String())

	return newUrl.String(), http.StatusCreated, nil
}

func (h *handler) expand(ctx *gin.Context) (interface{}, int, error) {
	hash := ctx.Param("hash")
	zap.S().Infof("hash info: %s", hash)

	res, err := h.db.GetInfo(hash)
	if err != nil {
		return "", http.StatusNotFound, fmt.Errorf("URL not found")
	}

	return res, http.StatusOK, nil
}

func (h *handler) redirect(ctx *gin.Context) {
	hash := ctx.Param("hash")

	link, err := h.db.GetLink(hash)

	zap.S().Infof("Redirect to %s", link)

	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, link)
}
