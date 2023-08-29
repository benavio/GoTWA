package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func handleRequest(w *httptest.ResponseRecorder, r *http.Request) {
	router := getRouter()
	router.ServeHTTP(w, r)
}

func createTestAlbum() album {
	testAlbum := album{ID: "2", Title: "TEST", Artist: "TESTOVICH", Price: 0.01}
	storage.Create(testAlbum)
	return testAlbum
}
func TestCreateAlbums(t *testing.T) {
	request, _ := http.NewRequest("POST", "/albums", strings.NewReader(`{"id": "4", "title": "Gib Beam", "artist": "John Coltrane", "price": 56.99}`))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusCreated {
		t.Fatal("status created 201", w.Code)
	}
}

func TestAlbumsList(t *testing.T) {
	request, _ := http.NewRequest("GET", "/albums", strings.NewReader(""))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusOK {
		t.Errorf("expected NotOK we got %d", w.Code)
	}
}

func TestAlbumDetail(t *testing.T) {
	testID := createTestAlbum().ID
	request, _ := http.NewRequest("GET", "/albums/"+testID, strings.NewReader(""))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusOK {
		t.Errorf("expected NotOK we got %d", w.Code)
	}
}

func TestAlbumNotFound(t *testing.T) {
	albumID := "9999"
	request, _ := http.NewRequest("GET", "/albums/"+albumID, strings.NewReader(""))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 we got %d", w.Code)
	}
}

func TestUpdateAlbum(t *testing.T) {
	testID := createTestAlbum().ID
	request, _ := http.NewRequest("PUT", "/albums/"+testID, strings.NewReader(`{"id": "4", "title": "TEST", "artist": "TEST", "price": 56.99}`))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 we got %d", w.Code)
	}
}

func TestUpdateAlbumNotFound(t *testing.T) {
	albumID := "9999"
	request, _ := http.NewRequest("PUT", "/albums/"+albumID, strings.NewReader(""))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 404 we got %d", w.Code)
	}
}

func TestUpdateBadStructure(t *testing.T) {
	albumID := "9999"
	request, _ := http.NewRequest("POST", "/albums/"+albumID, strings.NewReader(`{"title": "Karlos Makaroni"}`))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusNotFound {
		t.Fatal("status 404 but ")
	}
}

func TestCreateBadStructure(t *testing.T) {
	request, _ := http.NewRequest("POST", "/albums", strings.NewReader(""))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusBadRequest {
		t.Fatal("status 400", w.Code)
	}
}

func TestDeleteAlbum(t *testing.T) {
	testID := createTestAlbum().ID
	request, _ := http.NewRequest("DELETE", "/albums/"+testID, strings.NewReader(""))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusNoContent {
		t.Fatal("status 204", w.Code)
	}

}

func TestDeleteAlbumNotFound(t *testing.T) {
	albumID := "999999"
	request, _ := http.NewRequest("DELETE", "/albums/"+albumID, strings.NewReader(""))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusNotFound {
		t.Fatal("status 404", w.Code)
	}
}
