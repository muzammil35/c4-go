Connect 4 AI Game
A web-based Connect 4 game where you can play against an AI opponent. The game features a modern, responsive UI and a powerful AI solver using negamax algorithm with alpha-beta pruning.
Project Structure
connect4/
├── main.go                 # Gin API server
├── go.mod                  # Go module dependencies
├── Position/
│   └── position.go         # Game position logic
├── Solver/
│   └── solver.go           # AI solver logic
├── Transposition/
│   └── transposition.go    # Transposition table (you'll need to create this)
└── templates/
    └── index.html          # Frontend HTML
Setup Instructions
1. Prerequisites

Go 1.21 or later
Git

2. Project Setup

Create the project directory structure:

bashmkdir connect4
cd connect4

Initialize Go module:

bashgo mod init connect4

Create the package directories:

bashmkdir Position Solver Transposition templates

Add your existing code:

Place your position.go file in the Position/ directory
Place your solver.go file in the Solver/ directory
Create the main.go file in the root directory (provided above)
Create the templates/index.html file (provided above)


Create a simple Transposition table (if you don't have one):

Create Transposition/transposition.go:
gopackage Transposition

import "sync"

type TranspositionTable struct {
	table map[uint64]interface{}
	mutex sync.RWMutex
	size  int
}

func NewTranspositionTable(size int) *TranspositionTable {
	return &TranspositionTable{
		table: make(map[uint64]interface{}),
		size:  size,
	}
}

func (tt *TranspositionTable) Get(key uint64) (interface{}, bool) {
	tt.mutex.RLock()
	defer tt.mutex.RUnlock()
	val, exists := tt.table[key]
	return val, exists
}

func (tt *TranspositionTable) Put(key uint64, value interface{}) {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()
	
	// Simple size management - clear table if it gets too large
	if len(tt.table) >= tt.size {
		tt.table = make(map[uint64]interface{})
	}
	
	tt.table[key] = value
}

Install dependencies:

bashgo mod tidy
3. Running the Application

Start the server:

bashgo run main.go

Open your browser and navigate to:

http://localhost:8080
4. How to Play

Red pieces: Your moves
Yellow pieces: AI moves
Click on the column numbers (1-7) to drop your piece
The AI will automatically make its move after yours
First to connect 4 pieces wins!

5. API Endpoints
The server provides these REST API endpoints:

POST /api/new - Start a new game
POST /api/move - Make a player move
GET /api/status - Get current game state

6. Features

Modern UI: Responsive design with glassmorphism effects
Real-time gameplay: Immediate feedback and smooth animations
Powerful AI: Uses your existing negamax solver with alpha-beta pruning
Visual feedback: Highlights last moves and win conditions
Mobile friendly: Works on both desktop and mobile devices

7. Customization
You can customize the AI difficulty by modifying the searchDepth parameter in the Solve function call within MakeBestMove() in your solver code.
8. Troubleshooting

Port already in use: Change the port in main.go or stop other services using port 8080
CORS issues: The server is configured to allow all origins for development
Module not found: Make sure your go.mod file is in the root directory and run go mod tidy

Game Rules
Connect 4 is played on a 7x6 grid where players take turns dropping colored pieces from the top. The objective is to be the first to form a horizontal, vertical, or diagonal line of four of your pieces.
Enjoy playing against the AI! 🎮
