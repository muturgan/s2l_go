package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/muturgan/s2l_go/src/config"
	"github.com/muturgan/s2l_go/src/dal"
	"github.com/muturgan/s2l_go/src/models"
)

type server struct {
	dal dal.IDal
}

func newServer(dal dal.IDal) *server {
	return &server{dal: dal}
}

func (s *server) hashHandler(ctx *gin.Context) {
	hash := ctx.Param("hash")

	link, err := s.dal.GetLinkByHash(hash)
	if err != nil {
		fmt.Println("GetLinkByHash error!")
		fmt.Println(err)
		ctx.String(http.StatusInternalServerError, "Server error")
		return
	}
	if link == nil {
		fmt.Println("Not found!")
		ctx.String(http.StatusNotFound, "Not found")
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, link.Link)
}

func faviconHandler(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func (s *server) compressHandler(ctx *gin.Context) {
	var cr models.CompressRequest
	err := json.NewDecoder(ctx.Request.Body).Decode(&cr)
	if err != nil {
		errMessage := "Incorrect request body. It should be a valid json-serialized object with a \"link\" field which is a valid url"
		ctx.String(http.StatusBadRequest, errMessage)
		return
	}

	_, err = url.ParseRequestURI(cr.Link)
	if err != nil {
		errMessage := "Incorrect request body. It should be a valid json-serialized object with a \"link\" field which is a valid url"
		ctx.String(http.StatusBadRequest, errMessage)
		return
	}

	newShortLink, err := s.dal.CreateNewLink(cr.Link)
	if err != nil {
		fmt.Println("createNewLink error!")
		fmt.Println(err)
		ctx.String(http.StatusInternalServerError, "Server error")
		return
	}

	ctx.JSON(http.StatusOK, newShortLink)
}

func Serve(config *config.Config, dal dal.IDal) {
	router := gin.Default()

	s := newServer(dal)

	router.GET("/favicon.ico", faviconHandler)
	router.POST("/compress", s.compressHandler)
	router.GET("/_/:hash", s.hashHandler)
	router.GET("/:hash", s.hashHandler)

	fmt.Println("ok let's try to start at http://localhost" + config.GetServingAddress())

	err := router.Run(config.GetServingAddress())
	if err != nil {
		log.Fatal(err)
	}
}
