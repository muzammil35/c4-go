## Connect 4 AI Game
A web-based Connect 4 game where you can play against an AI opponent. The game features a modern, responsive UI and a powerful AI solver using negamax algorithm with alpha-beta pruning.
## Project Structure


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

## API Endpoints

The server provides these REST API endpoints:

POST /api/new - Start a new game
POST /api/move - Make a player move
GET /api/status - Get current game state 


## Features

Modern UI: Responsive design with glassmorphism effects
Real-time gameplay: Immediate feedback and smooth animations
Powerful AI: Uses your existing negamax solver with alpha-beta pruning
Visual feedback: Highlights last moves and win conditions
Mobile friendly: Works on both desktop and mobile devices


