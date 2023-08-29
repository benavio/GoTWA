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
	testAlbum := album{
		ID:         "1000",
		Segments:   []string{"segment1", "segment2"},
		LogChanges: "log changes",
	}
	storage.Create(testAlbum)
	return testAlbum
}
func TestCreateAlbums(t *testing.T) {
	request, _ := http.NewRequest("POST", "/albums", strings.NewReader(`{"id": "1000", "Segments": {"1000"}, "LogChanges": "qwe"}`))
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
	request, _ := http.NewRequest("GET", "/albums/1000", strings.NewReader(""))
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
	request, _ := http.NewRequest("POST", "/albums/1000", strings.NewReader(`{"id": "1000", "segments": ["Gib_Beam"], "logchanges": ""`))
	w := httptest.NewRecorder()
	handleRequest(w, request)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 we got %d", w.Code)
	}
}

func TestUpdateAlbumNotFound(t *testing.T) {
	albumID := "9999"
	request, _ := http.NewRequest("POST", "/albums/"+albumID, strings.NewReader(""))
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
	request, _ := http.NewRequest("DELETE", "/albums/", strings.NewReader(""))
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
