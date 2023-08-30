package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Error string `json:"error"`
}

var storage = NewStorage()

type album struct {
	ID         string   `json:"id"`
	Segments   []string `json:"segments"`
	LogChanges string   `json:"logchanges"`
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpError{"bad_request"})
		fmt.Println(err)
		return
	}
	storage.Create(newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, storage.Read())
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")
	album, err := storage.ReadOne(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}
func getUserContainsById(c *gin.Context) {
	id := c.Param("id")
	album, err := storage.UserContains(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func updateAlbumsById(c *gin.Context) {
	id := c.Param("id")
	var newAlbum album
	c.BindJSON(&newAlbum)
	album, err := storage.Update(id, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func deleteAlbumById(c *gin.Context) {
	id := c.Param("id")
	err := storage.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
	}
	c.IndentedJSON(http.StatusNoContent, album{})

}

func addSegmetsToUserById(c *gin.Context) {
	id := c.Param("id")
	segments := c.PostForm("segments")
	var newAlbum album
	c.BindJSON(&newAlbum)
	fmt.Println(newAlbum.LogChanges)
	album, err := storage.AddSegmentsToUser(id, segments, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func removeSegmetsById(c *gin.Context) {
	id := c.Param("id")
	// segments := c.Param("segment")
	segments := c.Query("segments")
	fmt.Println(id, segments)
	var newAlbum album
	c.BindJSON(&newAlbum)
	err := storage.DeleteUserSegments(id, segments, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"HTTP_ERROR:not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album{})
}

func getRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.POST("/albums", postAlbums)
	router.GET("/albums", getAlbums)
	// router.GET("/albums/:id", getAlbumById)
	// router.PUT("/albums/:id", updateAlbumsById)
	router.GET("/albums/:id", getUserContainsById) //gj
	// router.DELETE("/albums/:id", deleteAlbumById)
	router.PUT("/albums/:id", addSegmetsToUserById) //gj
	router.DELETE("/albums/:id", removeSegmetsById) //gj
	return router
}

func main() {
	router := getRouter()
	router.Run("localhost:8080")
}
