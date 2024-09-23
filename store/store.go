package store

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	DriverName = "sqlite"
	DBName     = "scheduler.db"
	Table      = "scheduler"
)

type Store struct {
	db *sql.DB
}

var DB Store

func NewStore() {
	db, err := CheckDB()
	if err != nil {
		log.Fatal(err)
	}
	DB = Store{db: db}
}

func CheckDB() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbEnv := os.Getenv("TODO_DBFILE")
	if dbEnv == "" {
		dbEnv = DBName
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbEnv)
	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		_, err := os.Create(dbFile)
		if err != nil {
			return nil, err
		}
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open(DriverName, dbFile)
	if err != nil {
		log.Fatal(err)
	}

	if install {
		query := `CREATE TABLE IF NOT EXISTS ` + Table + ` (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title VARCHAR(128) NOT NULL DEFAULT "",
				comment VARCHAR(128) NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX scheduler_date ON ` + Table + ` (date);`
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
	return db, err
}

func (s Store) Close() {
	s.db.Close()
}

func (s Store) Add(p Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO "+Table+" (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", p.Date),
		sql.Named("title", p.Title),
		sql.Named("comment", p.Comment),
		sql.Named("repeat", p.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s Store) Get(id int) (Task, error) {
	p := Task{}
	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM "+Table+" WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&p.ID, &p.Date, &p.Title, &p.Comment, &p.Repeat)
	return p, err
}

func (s Store) GetTasks(search string) ([]Task, error) {
	var tasks []Task
	q := ""
	date := ""
	if search != "" {
		q = " WHERE LOWER(title) LIKE LOWER(:search) OR LOWER(comment) LIKE LOWER(:search)"
		arr := strings.Split(search, ".")
		if len(arr) == 3 && len(search) == 10 {
			d, err := time.Parse("02.01.2006", search)
			if err == nil {
				q = " WHERE date = :date"
				date = d.Format("20060102")
			}
		}
	}
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM "+Table+q+
		" ORDER BY date LIMIT :limit",
		sql.Named("search", "%"+search+"%"),
		sql.Named("date", date),
		sql.Named("limit", 30))
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Task{}
		err := rows.Scan(&p.ID, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, p)
	}

	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, err
}
