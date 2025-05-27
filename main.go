package main

import (
	"net/http"
	"os"

	"connect4/Position"
	"connect4/Solver"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type GameState struct {
	Board      [][]int `json:"board"`
	Winner     int     `json:"winner"`     // -1: no winner, 0: player, 1: bot, 2: tie
	GameOver   bool    `json:"gameOver"`
	LastMove   int     `json:"lastMove"`
	NumMoves   int     `json:"numMoves"`
	CurrentPlayer int  `json:"currentPlayer"` // 0: player, 1: bot
}

type MoveRequest struct {
	Column int `json:"column"`
}

type MoveResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	GameState GameState `json:"gameState"`
	BotMove   int       `json:"botMove,omitempty"`
}

var gamePosition *Position.Position

func initGame() {
	gamePosition = Position.NewPosition()
}

func getGameState() GameState {
	board := gamePosition.BoardState()
	
	// Flip board vertically for display (top row first)
	flippedBoard := make([][]int, len(board))
	for i := range board {
		flippedBoard[len(board)-1-i] = board[i]
	}
	
	winner := -1
	gameOver := false
	
	// Check for tie
	if gamePosition.NumMoves == gamePosition.BoardHeight*gamePosition.BoardWidth {
		winner = 2
		gameOver = true
	}
	
	// Check for win
	if gamePosition.NumMoves > 0 && gamePosition.WinningBoardState() {
		// Winner is the player who made the last move
		lastPlayer := 1 - gamePosition.GetCurrentPlayer()
		winner = lastPlayer
		gameOver = true
	}
	
	return GameState{
		Board:         flippedBoard,
		Winner:        winner,
		GameOver:      gameOver,
		LastMove:      gamePosition.LastMove,
		NumMoves:      gamePosition.NumMoves,
		CurrentPlayer: gamePosition.GetCurrentPlayer(),
	}
}

func newGame(c *gin.Context) {
	initGame()
	gameState := getGameState()
	
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "New game started",
		"gameState": gameState,
	})
}

// makePlayerMove handles the player's move and validation
func makeMove(c *gin.Context) {
	var moveReq MoveRequest
	if err := c.ShouldBindJSON(&moveReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}
	
	// Validate column
	if moveReq.Column < 0 || moveReq.Column >= gamePosition.BoardWidth {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid column",
		})
		return
	}
	
	// Check if column is full
	if !gamePosition.CanPlay(moveReq.Column) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Column is full",
		})
		return
	}
	
	// Check if game is over
	currentState := getGameState()
	if currentState.GameOver {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Game is already over",
		})
		return
	}
	
	// Check if it's player's turn (player is 0)
	if gamePosition.GetCurrentPlayer() != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Not player's turn",
		})
		return
	}
	
	// Make player move
	gamePosition.Play(moveReq.Column)
	
	// Return current game state after player move
	gameState := getGameState()
	c.JSON(http.StatusOK, MoveResponse{
		Success:   true,
		Message:   "Player move made successfully",
		GameState: gameState,
	})
}

// getBotMove handles the bot's move logic
func makeBotMove(c *gin.Context) {
	// Check if game is over
	currentState := getGameState()
	if currentState.GameOver {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Game is already over",
		})
		return
	}
	
	// Check if it's bot's turn (bot is 1)
	if gamePosition.GetCurrentPlayer() != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Not bot's turn",
		})
		return
	}
	
	// Make bot move
	botMove := Solver.MakeBestMove(gamePosition)
	if botMove == -1 || !gamePosition.CanPlay(botMove) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Bot could not make a valid move",
		})
		return
	}
	
	gamePosition.Play(botMove)
	
	// Get final game state after bot move
	finalGameState := getGameState()
	
	c.JSON(http.StatusOK, MoveResponse{
		Success:   true,
		Message:   "Bot move made successfully",
		GameState: finalGameState,
		BotMove:   botMove,
	})
}


func getStatus(c *gin.Context) {
	gameState := getGameState()
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"gameState": gameState,
	})
}

func main() {
	// Initialize the game
	initGame()
	
	// Create Gin router
	r := gin.Default()
	
	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))
	
	// API routes
	api := r.Group("/api")
	{
		api.POST("/new", newGame)
		api.POST("/move", makeMove)
		api.POST("/bot", makeBotMove)
		api.GET("/status", getStatus)
	}
	
	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	
	// Serve the main page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Connect 4",
		})
	})
	
	port := "8080"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}
	
	r.Run(":" + port)
}