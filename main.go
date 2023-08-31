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

type segmentslist struct {
	Segments string `json:"segmentlist"`
}

type album struct {
	ID         string   `json:"id"`
	Segments   []string `json:"segments"`
	LogChanges []string `json:"logchanges"`
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	id := c.Param("id")
	segments := c.PostForm("segments")
	c.BindJSON(&newAlbum)
	fmt.Println(id, segments)
	album := storage.CreateUser(id, segments, newAlbum)
	c.IndentedJSON(http.StatusCreated, album)
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, storage.ReadUsers())
}

func getUserById(c *gin.Context) {
	id := c.Param("id")
	album, err := storage.ReadUser(id)
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

func updateUserById(c *gin.Context) {
	id := c.Param("id")
	var newAlbum album
	c.BindJSON(&newAlbum)
	album, err := storage.UpdateUser(id, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func deleteUserById(c *gin.Context) {
	id := c.Param("id")
	err := storage.DeleteUser(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
	}
	c.IndentedJSON(http.StatusNoContent, album{})

}

func postSegments(c *gin.Context) {
	var segments segmentslist
	arr := c.Param("segmentlist")
	fmt.Println(arr)
	c.BindJSON(&segments)
	storage.CreateSegment(arr, segments)
	c.IndentedJSON(http.StatusCreated, segments)
}

func getSegments(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, storage.ReadSegments())
}

func deleteSegment(c *gin.Context) {
	segment := c.Param("segmentlist")
	err := storage.DeleteSegment(segment)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, HttpError{"not found"})
	}
	c.IndentedJSON(http.StatusNoContent, album{})

}

func addSegmetsToUserById(c *gin.Context) {
	id := c.Param("id")
	segments := c.Param("segments")
	var newAlbum album
	c.BindJSON(&newAlbum)
	album, err := storage.AddUserSegments(id, segments, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func removeSegmetsById(c *gin.Context) {
	id := c.Param("id")
	segments := c.Param("segments")
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

func getDeletedList(c *gin.Context) {
	id := c.Param("id")
	segments := c.Param("segments")
	date := c.Param("date")
	storage.DeletedList(id, segments, date)
	c.IndentedJSON(http.StatusOK, album{})
}

func getAddedList(c *gin.Context) {
	id := c.Param("id")
	segments := c.Param("segments")
	date := c.Param("date")
	storage.AddedList(id, segments, date)
	c.IndentedJSON(http.StatusOK, album{})
}

func AddSegmentsRandom(c *gin.Context) {
	segments := c.Param("segments")
	percent := c.Param("percent")
	storage.AddSegmentsToRandomUsers(percent, segments)
	c.IndentedJSON(http.StatusOK, album{})
}

func getRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.POST("/albums/add-user/:id", postAlbums)
	router.GET("/albums", getAlbums)
	// router.GET("/albums/:id", getAlbumById)
	router.GET("/albums/:id", getUserContainsById) //gj
	router.DELETE("/albums/:id", deleteUserById)
	router.PUT("/albums/:id/add-segments/:segments", addSegmetsToUserById)    //gj
	router.DELETE("/albums/:id/delete-segments/:segments", removeSegmetsById) //gj

	router.GET("/albums/:id/:segments/deleted-list/:date", getDeletedList) //gj
	router.GET("/albums/:id/:segments/added-list/:date", getAddedList)     //gj

	router.POST("/segments/add-new-segment/:segmentlist", postSegments)
	router.GET("/segments", getSegments)
	router.DELETE("/segments/delete-segment/:segmentlist", deleteSegment)

	router.PUT("/albums/automatic-add/:segments/:percent", AddSegmentsRandom)
	return router
}

func main() {
	router := getRouter()
	router.Run("localhost:8080")
}
