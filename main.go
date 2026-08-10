package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // driver : imported only for its side effect(registers "postgres")
)

var db *sql.DB //shared connnection POOL, used by every handler

type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	// http.HandleFunc("/hello",
	// 	func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintln(w, "Hello from my Go server!")

	// 	})
	// fmt.Println("Server running on http://localhost:8080")
	// http.ListenAndServe(":8080", nil)

	connStr := "postgres://spurge:1411@localhost:5432/notesdb?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Open db:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("connect db:", err)
	}
	log.Println("Connected to databse!")

	r := gin.Default() //the router (net/http: http.NewServeMux,plus free logging )
	r.GET("/notes", getNotes)
	r.POST("/notes", createNote)
	r.GET("/notes/:id", getNoteByID)
	r.PUT("/notes/:id", updateNote)
	r.DELETE("/notes/:id", deleteNote)
	r.Run(":8080") // net/http:http.ListernAndServe(":8080",r)
}

// GET /notes -> send back the whole list as JSON
func getNotes(c *gin.Context) {
	// C (the "context") hold BOTH the request and the response in one object.
	rows, err := db.Query("SELECT id, title FROM notes ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close() // always close rows wehn the function ends
	result := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result = append(result, n)
	}
	// FINAL check after rows.Next() finishes
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST /notes -> read a note from the request body , add it to the list
func createNote(c *gin.Context) {
	var newNote Note

	// Read the JSON body into newNote. if the body isnt valid reply 400
	if err := c.ShouldBindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := db.QueryRow("INSERT INTO notes (title) VALUES ($1) RETURNING id",
		newNote.Title).Scan(&newNote.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newNote) //201 create
}
func getNoteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var n Note
	err = db.QueryRow("SELECT id, title FROM notes where id=($1)", id).Scan(&n.ID, &n.Title)
	if err != nil && err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, n)
}

func updateNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var updated Note
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated.ID = id
	result, err := db.Exec("UPDATE notes SET title=$1 WHERE id =$2", updated.Title, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	counts, _ := result.RowsAffected()
	if counts == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "nothing matched"})
		return
	}
	c.JSON(http.StatusOK, updated)
}
func deleteNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	result, err := db.Exec("DELETE FROM notes WHERE id =$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	counts, _ := result.RowsAffected()
	if counts == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "nothing matched"})
		return
	}
	c.Status(http.StatusNoContent)
}
