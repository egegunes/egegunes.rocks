package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type comment struct {
	Name      string
	Message   string
	CreatedAt string
}

func (c comment) DisplayDate() string {
	createdAt, err := time.Parse(time.RFC3339, c.CreatedAt)
	if err != nil {
		return c.CreatedAt
	}
	return createdAt.Format("2006-01-02")
}

func getComments(db *sql.DB) ([]comment, error) {
	var comments []comment

	rows, err := db.Query("SELECT name, message, created_at FROM comments ORDER BY id DESC")
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		var c comment
		err = rows.Scan(&c.Name, &c.Message, &c.CreatedAt)
		if err != nil {
			return comments, err
		}
		comments = append(comments, c)
	}

	return comments, err
}

func putComment(db *sql.DB, c comment) error {
	_, err := db.Exec("INSERT INTO comments(name, message, created_at) VALUES(?, ?, ?)", c.Name, c.Message, c.CreatedAt)

	return err
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS
	    comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		message TEXT,
		created_at TEXT
	    );
	`)
	if err != nil {
		log.Fatalf("can't create table: %v", err)
	}

	router := gin.Default()
	router.LoadHTMLFiles("index.tmpl")
	router.GET("/", func(c *gin.Context) {
		comments, err := getComments(db)
		if err != nil {
			log.Printf("can't get comments: %s", err)
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"comments": comments})
	})
	router.POST("/", func(c *gin.Context) {
		name := c.PostForm("name")
		message := c.PostForm("message")
		err := putComment(db, comment{name, message, time.Now().Format(time.RFC3339)})
		if err != nil {
			log.Printf("can't put comment: %s", err)
		}
		comments, err := getComments(db)
		if err != nil {
			log.Printf("can't get comments: %s", err)
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"comments": comments})
	})

	router.Run()
}
