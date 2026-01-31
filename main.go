package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"connect4/Position"
	"connect4/Solver"

	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Game struct {
	Position *Position.Position
	mu   sync.Mutex
}


type GameState struct {
	Board      [][]int `json:"board"`
	Winner     int     `json:"winner"`     // -1: no winner, 0: player, 1: bot, 2: tie
	GameOver   bool    `json:"gameOver"`
	LastMove   int     `json:"lastMove"`
	NumMoves   int     `json:"numMoves"`
	CurrentPlayer int  `json:"currentPlayer"` // 0: player, 1: bot
}

type MoveRequest struct {
	GameID string `json:"gameId"`
	Column int `json:"column"`
}

type MoveResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	GameState GameState `json:"gameState"`
	BotMove   int       `json:"botMove,omitempty"`
}

//var gamePosition *Position.Position

var games = make(map[string]*Game)
var gamesMutex sync.Mutex

func (g *Game) getGameState() GameState {
	board := g.Position.BoardState()
	
	// Flip board vertically for display (top row first)
	flippedBoard := make([][]int, len(board))
	for i := range board {
		flippedBoard[len(board)-1-i] = board[i]
	}
	
	winner := -1
	gameOver := false
	
	// Check for tie
	if g.Position.NumMoves == g.Position.BoardHeight*g.Position.BoardWidth {
		winner = 2
		gameOver = true
	}
	
	// Check for win
	if g.Position.NumMoves > 0 && g.Position.WinningBoardState() {
		// Winner is the player who made the last move
		lastPlayer := 1 - g.Position.GetCurrentPlayer()
		winner = lastPlayer
		gameOver = true
	}
	
	return GameState{
		Board:         flippedBoard,
		Winner:        winner,
		GameOver:      gameOver,
		LastMove:      g.Position.LastMove,
		NumMoves:      g.Position.NumMoves,
		CurrentPlayer: g.Position.GetCurrentPlayer(),
	}
}

func newGame(c *gin.Context) {
	gameId := uuid.New().String()

	game := &Game{
		Position: Position.NewPosition(),
		
	}

	gamesMutex.Lock()
	games[gameId] = game
	gamesMutex.Unlock()

	gameState := game.getGameState()
	
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "New game started",
		"gameState": gameState,
		"gameId": gameId,


	})
}

// makePlayerMove handles the player's move and validation
func (g *Game) makeMove(c *gin.Context, moveReq MoveRequest) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Validate column
	if moveReq.Column < 0 || moveReq.Column >= g.Position.BoardWidth {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid column",
		})
		return
	}
	
	// Check if column is full
	if !g.Position.CanPlay(moveReq.Column) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Column is full",
		})
		return
	}
	
	// Check if game is over
	currentState := g.getGameState()
	if currentState.GameOver {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Game is already over",
		})
		return
	}
	
	// Check if it's player's turn (player is 0)
	if g.Position.GetCurrentPlayer() != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Not player's turn",
		})
		return
	}
	
	// Make player move
	g.Position.Play(moveReq.Column)
	
	// Return current game state after player move
	gameState := g.getGameState()
	c.JSON(http.StatusOK, MoveResponse{
		Success:   true,
		Message:   "Player move made successfully",
		GameState: gameState,
	})
}

// getBotMove handles the bot's move logic
func (g *Game) makeBotMove(c *gin.Context) {
	// Check if game is over
	g.mu.Lock()
	defer g.mu.Unlock()
	currentState := g.getGameState()
	if currentState.GameOver {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Game is already over",
		})
		return
	}
	
	// Check if it's bot's turn (bot is 1)
	if g.Position.GetCurrentPlayer() != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Not bot's turn",
		})
		return
	}
	
	// Make bot move
	botMove := Solver.MakeBestMove(g.Position)
	if botMove == -1 || !g.Position.CanPlay(botMove) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Bot could not make a valid move",
		})
		return
	}
	
	g.Position.Play(botMove)
	
	// Get final game state after bot move
	finalGameState := g.getGameState()
	
	c.JSON(http.StatusOK, MoveResponse{
		Success:   true,
		Message:   "Bot move made successfully",
		GameState: finalGameState,
		BotMove:   botMove,
	})
}


func (g *Game) getStatus(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	gameState := g.getGameState()
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"gameState": gameState,
	})
}

func getGameByID(req MoveRequest) (*Game, error) {
	
	gamesMutex.Lock()
	game, ok := games[req.GameID]
	gamesMutex.Unlock()

	if !ok {
		
		return nil, errors.New("error reading the game from games slice")
	}

	return game, nil
}

func moveHandler(c *gin.Context) {

	var moveReq MoveRequest
	if err := c.ShouldBindJSON(&moveReq); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}
	game, err := getGameByID(moveReq)
	if  err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Could not retrieve Game",
		})
		return
	}
	game.makeMove(c, moveReq)
}

func botmoveHandler(c *gin.Context) {
	var moveReq MoveRequest
	if err := c.ShouldBindJSON(&moveReq); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}
	game, err := getGameByID(moveReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Could not retrieve Game",
		})
		return
	}
	game.makeBotMove(c)
}

func statusHandler(c *gin.Context) {
	gameID := c.Query("gameId")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	gamesMutex.Lock()
	game, ok := games[gameID]
	gamesMutex.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false})
		return
	}

	game.mu.Lock()
	defer game.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"gameState": game.getGameState(),
	})
}




func main() {
	
	
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
		api.POST("/move", moveHandler)
		api.POST("/bot", botmoveHandler)
		api.GET("/status", statusHandler)
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