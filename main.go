package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Error string `json:"error"`
}

var storage = NewStorage()

func NewMemoryStorage() MemoryStorage {
	var albums = []album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jery", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan", Artist: "Sarah Vaughan", Price: 39.99},
	}
	return MemoryStorage{albums: albums}

}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func getAlbum(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, storage.Read())
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, HttpError{"bad_request"})
		return
	}
	storage.Create(newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
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

func getRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.GET("/albums", getAlbum)
	router.GET("/albums/:id", getAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.PUT("/albums/:id", updateAlbumsById)
	router.POST("/albums", postAlbums)
	return router
}

func main() {
	router := getRouter()
	router.Run("localhost:8080")
}
