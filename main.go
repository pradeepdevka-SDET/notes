package main

import (
	"database/sql"
	"log"
	"net/http"
	"notesapi/store"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // driver : imported only for its side effect(registers "postgres")
)

// shared connnection POOL, used by every handler
type apiConfig struct {
	store *store.Store
}

func startBackgroundWorker(s *store.Store) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C { // ticker.C fires a value every 10 second
		notes, err := s.GetNotes()
		if err != nil {
			log.Println("[background] error:", err)
			continue
		}
		log.Printf("[background] there are currently %d notes", len(notes))
	}
}

func setupRouter(cfg *apiConfig) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/notes", cfg.getNotes)
	r.POST("/notes", cfg.createNote)
	r.GET("/notes/:id", cfg.getNoteByID)
	r.PUT("/notes/:id", cfg.updateNote)
	r.DELETE("/notes/:id", cfg.deleteNote)
	return r
}
func main() {
	// http.HandleFunc("/hello",
	// 	func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintln(w, "Hello from my Go server!")

	// 	})
	// fmt.Println("Server running on http://localhost:8080")
	// http.ListenAndServe(":8080", nil)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using real environment variables")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Open db:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("connect db:", err)
	}
	log.Println("Connected to databse!")

	cfg := apiConfig{store: store.New(db)}
	r := setupRouter(&cfg)
	go startBackgroundWorker(cfg.store) //runs concurently with the web server
	r.Run(":" + port)                   // net/http:http.ListernAndServe(":8080",r)

}

/*
// GET /notes -> send back the whole list as JSON
func getNotes(c *gin.Context) {
	// C (the "context") hold BOTH the request and the response in one object.
	rows, err := db.Query("SELECT id, title FROM notes ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close() // always close rows wehn the function ends

	result := []models.Note{}
	for rows.Next() {
		var n models.Note
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
	var newNote models.Note

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
	var n models.Note
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
	var updated models.Note
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
*/
