package handler

import (
	"fmt"
	"log"
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
	}
)

func NewGinRouter(netProtocol, address string, db redisdb.DBAction) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // for production
	r := gin.New()
	// r.Use(cors.Default())
	//r.Use(static.Serve("/", static.LocalFile("./web", true)))
	r.Use(static.Serve("/", static.LocalFile("./web", true)))
	//api := r.Group("/api")
	h := handler{netProtocol, address, db}
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
		zap.S().Info(res)
		context.JSON(code, resp{Data: res, Success: err == nil})
	}
}

func (h *handler) cut(ctx *gin.Context) (interface{}, int, error) {
	var body struct {
		URL string `json:"url" binding:"required"`
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		zap.S().Errorw("Couldnot get body from request.", "err", err)
	}
	zap.S().Infow("Link for processing.", "link", body.URL)

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

	zap.S().Infof("Short link: %v", newUrl.String())

	return newUrl.String(), http.StatusCreated, nil
}

func (h *handler) expand(ctx *gin.Context) (interface{}, int, error) {
	hash := ctx.Param("hash")
	zap.S().Infof("Hash info: %s", hash)

	res, err := h.db.GetInfo(hash)
	log.Printf("----...>>> %v", err)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
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
