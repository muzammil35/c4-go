package Position

import (
	"fmt"
	"sort"
)

// Position represents the Connect Four game state
type Position struct {
	BoardHeight   int
	BoardWidth    int
	NumMoves      int
	CurrentPositions [2]uint64
	LastMove      int
	BitShifts     []int
	ColumnOrder   []int
}

// NewPosition creates and initializes a new Position
func NewPosition() *Position {
	pos := &Position{
		BoardHeight:   6,
		BoardWidth:    7,
		NumMoves:      0,
		CurrentPositions: [2]uint64{0, 0},
		LastMove:      -1,
	}
	pos.BitShifts = pos.getBitShifts()
	pos.ColumnOrder = pos.getColumnOrder()
	return pos
}

// getBitShifts calculates bit shifts used for win detection
func (p *Position) getBitShifts() []int {
	return []int{
		1,              // vertical
		p.BoardHeight,  // diagonal \
		p.BoardHeight + 1, // horizontal
		p.BoardHeight + 2, // diagonal /
	}
}

// getColumnOrder generates the preferred column play order (middle first)
func (p *Position) getColumnOrder() []int {
	order := make([]int, p.BoardWidth)
	for i := 0; i < p.BoardWidth; i++ {
		order[i] = p.BoardWidth/2 + (1-2*(i%2))*(i+1)/2
	}
	return order
}

// GetMask returns the combined board state (all pieces)
func (p *Position) GetMask() uint64 {
	return p.CurrentPositions[0] | p.CurrentPositions[1]
}

// GetKey returns a unique game state identifier
func (p *Position) GetKey() uint64 {
	return p.GetMask() + p.CurrentPositions[p.GetCurrentPlayer()]
}

// TopMask returns a bit mask for the top position in a column
func (p *Position) TopMask(col int) uint64 {
	return uint64(1)<<(p.BoardHeight-1) << uint64(col*(p.BoardHeight+1))
}

// BottomMask returns a bit mask for the bottom position in a column
func (p *Position) BottomMask(col int) uint64 {
	return uint64(1) << uint64(col*(p.BoardHeight+1))
}

// CanPlay checks if a move in the given column is valid
func (p *Position) CanPlay(col int) bool {
	if p.NumMoves == p.BoardHeight*p.BoardWidth {
		return false
	}
	mask := p.GetMask()
	return (mask & p.TopMask(col)) == 0
}

// GetCurrentPlayer returns the index of the current player (0 or 1)
func (p *Position) GetCurrentPlayer() int {
	if p.NumMoves%2 == 0 {
		return 0
	}
	return 1
}

// Play makes a move in the specified column
func (p *Position) Play(col int) {
	currPlayer := p.GetCurrentPlayer()
	mask := p.GetMask()

	// Get opponent player index (1 if currPlayer is 0, otherwise 0)
	opponent := 1 - currPlayer

	// Update the current positions
	p.CurrentPositions[opponent] = p.CurrentPositions[currPlayer] ^ mask
	updatedMask := mask | (mask + p.BottomMask(col))
	p.CurrentPositions[currPlayer] = p.CurrentPositions[opponent] ^ updatedMask

	// Update move count and last move
	p.NumMoves++
	p.LastMove = col
}

// WinningBoardState checks if the last move created a winning alignment
func (p *Position) WinningBoardState() bool {
	opp := 1 - p.GetCurrentPlayer()
	for _, shift := range p.BitShifts {
		test := p.CurrentPositions[opp] & (p.CurrentPositions[opp] >> uint64(shift))
		if test&(test>>uint64(2*shift)) != 0 {
			return true
		}
	}
	return false
}

// GetScore returns the score of a complete game
func (p *Position) GetScore() int {
	return -((p.BoardWidth*p.BoardHeight + 1 - p.NumMoves) / 2)
}

// Helper for bitCount - counts set bits in a uint64
func bitCount(n uint64) int {
	count := 0
	for n > 0 {
		n &= n - 1
		count++
	}
	return count
}

// colSort provides a score for column preference (higher is better)
func (p *Position) colSort(col int) int {
	mask := p.GetMask()
	player := p.GetCurrentPlayer()
	position := p.CurrentPositions[player]
	oppPosition := position ^ mask
	newMask := mask | (mask + p.BottomMask(col))
	state := oppPosition ^ newMask

	count := 0
	for _, shift := range p.BitShifts {
		test := state & (state >> uint64(shift)) & (state >> uint64(2*shift))
		if test != 0 {
			count += bitCount(test)
		}
	}
	return count
}

// GetSearchOrder returns columns in preferred search order
func (p *Position) GetSearchOrder() []int {
	var validColumns []int
	for _, col := range p.ColumnOrder {
		if p.CanPlay(col) {
			validColumns = append(validColumns, col)
		}
	}

	// Sort by colSort score
	sort.Slice(validColumns, func(i, j int) bool {
		return p.colSort(validColumns[i]) > p.colSort(validColumns[j])
	})

	return validColumns
}

// IsWinningMove checks if playing in a column would create a win
func (p *Position) IsWinningMove(col int, position uint64) bool {
	mask := p.GetMask()
	oppPosition := position ^ mask
	newMask := mask | (mask + p.BottomMask(col))
	candidatePosition := oppPosition ^ newMask
	return p.ConnectedFour(candidatePosition)
}

// ConnectedFour checks if a position has four connected pieces
func (p *Position) ConnectedFour(position uint64) bool {
	// Horizontal check
	m := position & (position >> 7)
	if m&(m>>14) != 0 {
		return true
	}
	// Diagonal \
	m = position & (position >> 6)
	if m&(m>>12) != 0 {
		return true
	}
	// Diagonal /
	m = position & (position >> 8)
	if m&(m>>16) != 0 {
		return true
	}
	// Vertical
	m = position & (position >> 1)
	if m&(m>>2) != 0 {
		return true
	}
	return false
}

// BoardState returns a 2D representation of the board
func (p *Position) BoardState() [][]int {
	board := make([][]int, p.BoardHeight)
	for i := range board {
		board[i] = make([]int, p.BoardWidth)
		for j := range board[i] {
			board[i][j] = -1
		}
	}

	// Convert bit representation to 2D array
	for col := 0; col < p.BoardWidth; col++ {
		for row := 0; row < p.BoardHeight; row++ {
			pos := uint64(1) << uint64(col*(p.BoardHeight+1)+row)
			if p.CurrentPositions[0]&pos != 0 {
				board[row][col] = 0
			} else if p.CurrentPositions[1]&pos != 0 {
				board[row][col] = 1
			}
		}
	}

	return board
}

// PrintBoard displays the current board state
func (p *Position) PrintBoard() {
	board := p.BoardState()
	for i := 0; i < p.BoardHeight; i++ {
		for j := 0; j < p.BoardWidth; j++ {
			switch board[i][j] {
			case -1:
				fmt.Print(". ")
			case 0:
				fmt.Print("X ")
			case 1:
				fmt.Print("O ")
			}
		}
		fmt.Println()
	}
	fmt.Println("---------------")
	fmt.Println("0 1 2 3 4 5 6")
}

func main() {
	// Example usage
	pos := NewPosition()
	
	// Play some moves (example)
	columns := []int{3, 2, 4, 1}
	for _, col := range columns {
		fmt.Println("playing col: ", col)
		pos.Play(col)
		fmt.Println("current positions: ", pos.CurrentPositions)
		fmt.Printf("%b\n", pos.CurrentPositions[0])
		fmt.Printf("%b\n", pos.CurrentPositions[1])
		
	}
	
	// Show search order
	fmt.Println("Preferred columns:", pos.GetSearchOrder())
}