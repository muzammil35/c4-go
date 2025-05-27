package main

import (
	//"html/template"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Create a structure for the Connect 4 board
type Board struct {
	Cells [6][7]string // 6 rows, 7 columns (standard Connect 4 board)
}

func main() {
	// Create an instance of the Gin router
	r := gin.Default()

	// Serve static files like CSS, images, etc.
	r.Static("/assets", "./assets")

	// Route to serve the landing page
	r.GET("/", func(c *gin.Context) {
		board := Board{}
		// Initialize the board with empty strings
		for row := 0; row < 6; row++ {
			for col := 0; col < 7; col++ {
				board.Cells[row][col] = ""
			}
		}
		// Render the template and pass the board data
		c.HTML(http.StatusOK, "index.html", gin.H{
			"board": board,
		})
	})

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Run the server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
