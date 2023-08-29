package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slices"
)

type Storage interface {
	Create(album) album
	Read() []album
	ReadOne(id string) (album, error)
	UserContains(id string) (album, error)
	Update(id string, a album) (album, error)
	AddSegment(id string, segment string, a album) (album, error)
	DeleteSegment(id string, segment string, a album) error
	Delete(id string) error
}

type PostgresStorage struct {
	db *sql.DB
}

func (p PostgresStorage) CreateSchema() error {
	_, err := p.db.Exec("create table if not exists albums (ID char(16) primary key, Segments text[], LogChanges text)")
	return err
}

func NewStorage() Storage {
	return NewPostgresStorage()
}
func NewPostgresStorage() PostgresStorage {
	connStr := "user=user dbname=db password=pass sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	storage := PostgresStorage{db: db}
	err = storage.CreateSchema()
	if err != nil {
		log.Fatal(err)
	}
	return storage
}

func (p PostgresStorage) Create(am album) album {
	_, err := p.db.Exec("INSERT INTO albums (ID, Segments, LogChanges) VALUES($1,$2,$3)", am.ID, (*pq.StringArray)(&am.Segments), am.LogChanges)
	if err != nil {
		log.Fatal("Create err ->  ", err)
	}
	return am
}

func (p PostgresStorage) Read() []album {
	rows, err := p.db.Query("select * from  albums")
	if err != nil {
		log.Fatal("Read error 1 -> ", err)
	}
	defer rows.Close()

	var albums []album
	for rows.Next() {
		var a album
		err := rows.Scan(&a.ID, pq.Array(&a.Segments), &a.LogChanges)
		if err != nil {
			log.Fatal("Read error 2 -> ", err)
		}
		albums = append(albums, a)
	}

	return albums
}

func (p PostgresStorage) ReadOne(id string) (album, error) {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), &album.LogChanges)
	if err != nil {
		if err == sql.ErrNoRows {
			return album, errors.New("Not found")
		}
		return album, err
	}
	return album, nil
}

func (p PostgresStorage) UserContains(id string) (album, error) {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), &album.LogChanges)
	if err != nil {
		if err == sql.ErrNoRows {
			return album, errors.New("Not found")
		}
		return album, err
	}
	length := 0
	row = p.db.QueryRow("select cardinality(segments) from albums where id = $1", id)
	err = row.Scan(&length)
	if err != nil {
		if err == sql.ErrNoRows {
			return album, err
		}
		return album, err
	}
	segments := strings.Join(album.Segments, ", ")
	if length > 0 {
		fmt.Printf("Пользователь %s состоит в %d сегментах: %s\n", strings.Trim(album.ID, " "), length, segments)
	} else {
		fmt.Printf("Пользователь %s не состоит ни в одном из сегментов\n", strings.Trim(album.ID, " "))
	}
	return album, nil
}

func (p PostgresStorage) Update(id string, a album) (album, error) {
	result, _ := p.db.Exec("update albums set segments=$1, logchanges=$2 where id=$3", (*pq.StringArray)(&a.Segments), a.LogChanges, a.ID)
	err := handlerNotFound(result)
	return a, err
}

func (p PostgresStorage) Delete(id string) error {
	result, err := p.db.Exec("delete from albums where id=$1", id)
	if err != nil {
		log.Fatal(err)
	}
	err = handlerNotFound(result)
	return err
}
func handlerNotFound(result sql.Result) error {
	countAffected, _ := result.RowsAffected()
	if countAffected == 0 {
		return errors.New("Not found")
	}
	return nil
}

func (p PostgresStorage) AddSegment(id string, segment string, a album) (album, error) {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), &album.LogChanges)
	if err != nil {
		if err == sql.ErrNoRows {
			return album, errors.New("Not found")
		}
		return album, err
	}
	if contains(segment, album.Segments) {
		return album, fmt.Errorf("Segment already exists in segments")
	}
	album.Segments = append(album.Segments, segment)
	_, err = p.db.Exec("update albums set segments=$1, logchanges=$2 WHERE id=$3", (*pq.StringArray)(&album.Segments), album.LogChanges, album.ID)
	if err != nil {
		return album, err
	}
	return album, nil
}

func (p PostgresStorage) DeleteSegment(id string, segment string, a album) error {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), &album.LogChanges)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Not found")
		}
		return err
	}
	if !contains(segment, album.Segments) {
		return errors.New("Segment missing in Segments")
	}
	for i, v := range album.Segments {
		if v == segment {
			// album.Segments = append(album.Segments[i:], album.Segments[:i+1]...)
			album.Segments = slices.Delete(album.Segments, i, i+1)
		}
	}
	_, err = p.db.Exec("update albums set segments=$1, logchanges=$2 WHERE id=$3", (*pq.StringArray)(&album.Segments), album.LogChanges, album.ID)
	if err != nil {
		return err
	}
	return nil
}

func contains(segment string, segments []string) bool {
	for _, s := range segments {
		if s == segment {
			return true
		}
	}
	return false
}

const (
	voice       = "AVITO_VOICE_MESSAGES"
	performance = "AVITO_PERFORMANCE_VAS"
	discount30  = "AVITO_DISCOUNT_30"
	discount50  = "AVITO_DISCOUNT_50"
)
