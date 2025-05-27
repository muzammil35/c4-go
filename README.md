## Connect 4 AI Game
A web-based Connect 4 game where you can play against an AI opponent. The game features a modern, responsive UI and a powerful AI solver using negamax algorithm with alpha-beta pruning.
## Project Structure


connect4/ <br>
├── main.go                 # Gin API server <br>
├── go.mod                  # Go module dependencies <br>
├── Position/ <br>
│   └── position.go         # Game position logic <br>
├── Solver/ <br>
│   └── solver.go           # AI solver logic <br>
├── Transposition/ <br>
│   └── transposition.go    # Transposition table (you'll need to create this) <br>
└── templates/ <br>
    └── index.html          # Frontend HTML <br>

## API Endpoints

The server provides these REST API endpoints:

POST /api/new - Start a new game <br>
POST /api/move - Make a player move <br>
GET /api/status - Get current game state  <br>


## Features

Modern UI: Responsive design with glassmorphism effects <br>
Real-time gameplay: Immediate feedback and smooth animations <br>
Powerful AI: Uses your existing negamax solver with alpha-beta pruning <br>
Visual feedback: Highlights last moves and win conditions <br>
Mobile friendly: Works on both desktop and mobile devices <br>


