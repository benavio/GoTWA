package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slices"
)

const (
	timeFormat = ""
)

type Storage interface {
	CreateUser(id string, segments string, newAlbum album) album
	ReadUsers() []album
	ReadUser(id string) (album, error)
	UpdateUser(id string, a album) (album, error)
	DeleteUser(id string) error
	AddUserSegments(id string, segments string, a album) (album, error)
	DeleteUserSegments(id string, segments string, a album) error
	UserContains(id string) (string, error)
	CreateSegment(arr string, segments segmentslist) segmentslist
	ReadSegments() []string
	DeleteSegment(segment string) error
	DeletedList(id string, segments string, date string) error
	AddedList(id string, segments string, date string) error
	AddSegmentsToRandomUsers(percent string, segments string)
}

type PostgresStorage struct {
	db *sql.DB
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
	err = storage.CreateSchema2()
	if err != nil {
		log.Fatal(err)
	}
	return storage
}

func (p PostgresStorage) CreateSchema() error {
	_, err := p.db.Exec("create table if not exists albums (ID char(16) primary key, Segments text[], LogChanges text[])")
	return err
}

func (p PostgresStorage) CreateUser(id string, segments string, newAlbum album) album {
	var album album
	segment := strings.Split(segments, ",")
	for _, v := range segment {
		if len(v) != 0 {
			req := fmt.Sprintf("Сегмент %s добавлен: %s", v, time.Now().Format("2006-1-2 15:4:5"))
			album.Segments = append(album.Segments, v)
			if len(v) != 0 {
				album.LogChanges = append(album.LogChanges, req)
			}
		}
	}
	_, err := p.db.Exec("INSERT INTO albums (ID, Segments, LogChanges) VALUES($1,$2,$3)", id, (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges))
	if err != nil {
		log.Fatal("Create err ->  ", err)
	}
	return album
}

func (p PostgresStorage) ReadUsers() []album {
	rows, err := p.db.Query("select * from  albums")
	if err != nil {
		log.Fatal("Read error 1 -> ", err)
	}
	defer rows.Close()

	var albums []album
	for rows.Next() {
		var a album
		err := rows.Scan(&a.ID, pq.Array(&a.Segments), pq.Array(&a.LogChanges))
		if err != nil {
			log.Fatal("Read error 2 -> ", err)
		}
		albums = append(albums, a)
	}

	return albums
}

func (p PostgresStorage) ReadUser(id string) (album, error) {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges))
	if err != nil {
		if err == sql.ErrNoRows {
			return album, errors.New("Not found")
		}
		return album, err
	}
	return album, nil
}

func (p PostgresStorage) UpdateUser(id string, a album) (album, error) {
	result, _ := p.db.Exec("update albums set segments=$1, logchanges=$2 where id=$3", (*pq.StringArray)(&a.Segments), (*pq.StringArray)(&a.LogChanges), a.ID)
	err := handlerNotFound(result)
	return a, err
}

func (p PostgresStorage) DeleteUser(id string) error {
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

func (p PostgresStorage) CreateSchema2() error {
	_, err := p.db.Exec("create table if not exists segments (segmentslist text)")
	return err
}

func (p PostgresStorage) CreateSegment(arr string, sg segmentslist) segmentslist {
	segment := strings.Split(arr, ",")
	for _, v := range segment {
		if arr != "" && !checkSegment(v) {
			_, err := p.db.Exec("insert into segments (segmentslist) VALUES($1)", strings.ToUpper(v))
			if err != nil {
				log.Fatal("Create Segment err ->  ", err)
			}
		}
	}
	return sg
}

func (p PostgresStorage) ReadSegments() []string {
	rows, err := p.db.Query("select * from  segments")
	if err != nil {
		log.Fatal("Read Segments error 1 -> ", err)
	}
	defer rows.Close()
	var list []string
	for rows.Next() {
		var l string
		err := rows.Scan(&l)
		if err != nil {
			log.Fatal("Read Segments error 2 -> ", err)
		}
		list = append(list, l)
	}

	return list
}

func (p PostgresStorage) DeleteSegment(segment string) error {
	segments := strings.Split(segment, ",")
	for _, v := range segments {
		result, err := p.db.Exec("delete from segments where segmentslist=$1", strings.ToUpper(v))
		// result, err := p.db.Exec("delete from segments where segmentslist=$1", strings.ToUpper(v))
		if err != nil {
			log.Fatal(err)
			err = handlerNotFound(result)
		}
	}
	return nil
}

// --------------------------  JSON RESPONSE  -----------------------
func (p PostgresStorage) UserContains(id string) (string, error) {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges))
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("Not found")
		}
		return "", err
	}
	length := 0
	row = p.db.QueryRow("select cardinality(segments) from albums where id = $1", id)
	err = row.Scan(&length)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}
	segments := strings.Join(album.Segments, ", ")
	var req string
	if length > 0 {
		req = fmt.Sprintf("Пользователь %s состоит в %d сегментах: %s\n", strings.Trim(album.ID, " "), length, segments)
	} else {
		req = fmt.Sprintf("Пользователь %s не состоит ни в одном из сегментов\n", strings.Trim(album.ID, " "))
	}
	return req, nil
}

func (p PostgresStorage) AddUserSegments(id string, segments string, a album) (album, error) {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges))
	if err != nil {
		if err == sql.ErrNoRows {
			return album, errors.New("Not found")
		}
		return album, err
	}
	segment := strings.Split(segments, ",")
	for _, v := range segment {
		if checkSegment(v) && !checkContains(v, album.Segments) {
			album.Segments = append(album.Segments, strings.ToUpper(v))
			if len(v) != 0 {
				req := fmt.Sprintf("Сегмент %s добавлен: %s", strings.ToUpper(v), time.Now().Format("2-1-2006"))
				album.LogChanges = append(album.LogChanges, req)
			}
		}
	}
	_, err = p.db.Exec("update albums set segments=$1, logchanges=$2 WHERE id=$3", (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges), album.ID)
	if err != nil {
		return album, err
	}
	return album, nil
}

func (p PostgresStorage) DeleteUserSegments(id string, segments string, a album) error {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges))
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Not found")
		}
		return err
	}
	segment := strings.Split(segments, ",")
	for _, val := range segment {
		if checkContains(val, album.Segments) {
			for i, v := range album.Segments {
				if slug.Make(v) == slug.Make(val) {
					req := fmt.Sprintf("Сегмент %s удалён: %s", strings.ToUpper(v), time.Now().Format("2-1-2006"))
					album.Segments = slices.Delete(album.Segments, i, i+1)
					album.LogChanges = append(album.LogChanges, req)
				}
			}
		}
	}
	_, err = p.db.Exec("update albums set segments=$1, logchanges=$2 WHERE id=$3", (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges), album.ID)
	if err != nil {
		return err
	}
	return nil
}
func checkContains(segment string, segments []string) bool {
	for _, s := range segments {
		if slug.Make(s) == slug.Make(segment) {
			return true
		}
	}
	return false
}

func checkSegment(segment string) bool {
	arr := storage.ReadSegments()
	for _, v := range arr {
		if slug.Make(v) == segment {
			return true
		}
	}
	return false
}

func (p PostgresStorage) DeletedList(id string, segments string, date string) error {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges))
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	file, err := os.Create(fmt.Sprintf("Deleted-%s-%s-%s.csv", id, segments, date))
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	for _, v := range album.LogChanges {
		logchanges := strings.Split(v, " ")
		if logchanges[1] == segments && logchanges[2] == "удалён:" && logchanges[3] == date {
			logchanges[0] = id
			if err := writer.Write(logchanges); err != nil {
				return err
			}
		}
	}
	writer.Flush()
	return nil
}
func (p PostgresStorage) AddedList(id string, segments string, date string) error {
	var album album
	row := p.db.QueryRow("select * from albums where id = $1", id)
	err := row.Scan(&album.ID, (*pq.StringArray)(&album.Segments), (*pq.StringArray)(&album.LogChanges))
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	file, err := os.Create(fmt.Sprintf("Added-%s-%s-%s.csv", id, segments, date))
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	for _, v := range album.LogChanges {
		logchanges := strings.Split(v, " ")
		if logchanges[1] == segments && logchanges[2] == "добавлен:" && logchanges[3] == date {
			logchanges[0] = id
			if err := writer.Write(logchanges); err != nil {
				return err
			}
		}
	}
	writer.Flush()
	return nil
}
func (p PostgresStorage) AddSegmentsToRandomUsers(percent string, segments string) {
	album := storage.ReadUsers()
	fmt.Println(album[0].ID)
	if len(album) > 2 {
		percent, err := strconv.ParseFloat(percent, 8)
		if err != nil {
			log.Fatal("NaN")
		}
		fmt.Println(len(album), (percent / 100))
		var leng float64 = float64(len(album))
		floatCurrentMembers := leng * (percent / 100)
		currentMembers := int(math.Round(floatCurrentMembers))
		fmt.Println("------------------------", currentMembers)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for currentMembers != 0 {
			q := r.Intn(len(album))
			storage.AddUserSegments(album[q].ID, segments, album[q])
			currentMembers--
		}
	}

}

// func segmentsArray() []string {
// 	segments := []string{
// 		"AVITO_VOICE_MESSAGES",
// 		"AVITO_PERFORMANCE_VAS",
// 		"AVITO_DISCOUNT_30",
// 		"AVITO_DISCOUNT_50",
// 	}
// 	return segments
// }
