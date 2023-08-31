package main

import (
	"errors"
)

type MemoryStorage struct {
	albums []album
}

func NewMemoryStorage() MemoryStorage {
	var albums = []album{
		{ID: "1000", Segments: []string{}, LogChanges: []string{}},
		{ID: "1002", Segments: []string{}, LogChanges: []string{}},
		{ID: "1004", Segments: []string{}, LogChanges: []string{}},
	}
	return MemoryStorage{albums: albums}

}

func (s MemoryStorage) Create(am album) album {
	s.albums = append(s.albums, am)
	return am
}
func (s MemoryStorage) ReadOne(id string) (album, error) {
	for _, a := range s.albums {
		if a.ID == id {
			return a, nil
		}
	}
	return album{}, errors.New("not found")
}

func (s MemoryStorage) Read() []album {
	return s.albums
}

func (s MemoryStorage) Update(id string, newAlbum album) (album, error) {
	for i := range s.albums {
		if s.albums[i].ID == id {
			s.albums[i] = newAlbum
			return s.albums[i], nil
		}
	}
	return album{}, errors.New("not found")
}

func (s MemoryStorage) UpdateSegment(id string, newAlbum album) (album, error) {
	segment := "Voice"
	for i, obj := range s.albums {
		if obj.ID == id {
			s.albums[i].Segments = append(s.albums[i].Segments, segment)
			break
		}
	}
	return album{}, errors.New("not found")
}

func (s MemoryStorage) Delete(id string) error {
	for i, a := range s.albums {
		if a.ID == id {
			s.albums = append(s.albums[:i], s.albums[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
