package Solver

import (
	"connect4/Position"
	"connect4/Transposition"
	"math"
)

// TTEntry represents an entry in the transposition table
type TTEntry struct {
	Value int
	Col   int
	UB    bool // Upper bound flag
	LB    bool // Lower bound flag
}

type Result struct {
	Score int
	Col   int
}

// GetWinScore calculates the win score based on the current position
func GetWinScore(position *Position.Position) int {
	return ((position.BoardWidth*position.BoardHeight + 1) - position.NumMoves) / 2
}

// TieGame checks if the game is a tie
func TieGame(position *Position.Position) bool {
	return position.NumMoves == position.BoardHeight*position.BoardWidth
}

// Negamax implements the negamax algorithm with alpha-beta pruning and transposition table
func Negamax(position *Position.Position, alpha, beta int, transpositionTable *Transposition.TranspositionTable, maxDepth int) (int, int) {
	if maxDepth == 0 {
		return 0, 0
	}

	
	key := position.GetKey()

	// Check transposition table
	if entry, exists := transpositionTable.Get(key); exists {
		cachedEntry := entry.(TTEntry)
		if cachedEntry.LB {
			if cachedEntry.Value > alpha {
				alpha = cachedEntry.Value
			}
			if alpha >= beta {
				return cachedEntry.Value, cachedEntry.Col
			}
		} else {
			return cachedEntry.Value, cachedEntry.Col
		}
	}

	// Check for terminal states
	if TieGame(position) {
		return 0, position.LastMove
	}

	// Check if previous player won
	prevPlayer := 1
	if position.GetCurrentPlayer() == 1 {
		prevPlayer = 0
	}
	if position.ConnectedFour(position.CurrentPositions[prevPlayer]) {
		return -1 * GetWinScore(position), position.LastMove
	}

	// Look for immediate win
	for _, col := range position.GetSearchOrder() {
		if position.CanPlay(col) {
			if position.IsWinningMove(col, position.CurrentPositions[position.GetCurrentPlayer()]) {
				return GetWinScore(position), col
			}
		}
	}

	//if maxDepth > 3 {
		//return concurrentNegamax(position, alpha, beta, transpositionTable, maxDepth)
	//}

	bestCol := -1
	bestScore := math.MinInt32

	// Recursive search
	for _, col := range position.GetSearchOrder() {
		if position.CanPlay(col) {
			// Create a copy of the position and make the move
			newPosition := Position.NewPosition()
			newPosition.BoardHeight = position.BoardHeight
			newPosition.BoardWidth = position.BoardWidth
			newPosition.NumMoves = position.NumMoves
			newPosition.CurrentPositions = [2]uint64{position.CurrentPositions[0], position.CurrentPositions[1]}
			newPosition.BitShifts = position.BitShifts
			newPosition.ColumnOrder = position.ColumnOrder
			newPosition.Play(col)

			// Recursive call with negated alpha/beta
			score, _ := Negamax(newPosition, -beta, -alpha, transpositionTable, maxDepth-1)
			score = -score

			// Beta cutoff
			if score >= beta {
				transpositionTable.Put(key, TTEntry{
					Value: score,
					Col:   col,
					LB:    true,
				})
				return score, col
			}

			// Update alpha
			if score > alpha {
				alpha = score
				bestCol = col
				transpositionTable.Put(key, TTEntry{
					Value: score,
					Col:   bestCol,
				})
			} else if score > bestScore {
				bestScore = score
				bestCol = col
			}
		}
	}

	return alpha, bestCol
}



// Solve uses iterative deepening with Negamax to find the best move
func Solve(position *Position.Position, weak bool, loopIters int, searchDepth int) (int, int) {
	minVal := -(position.BoardWidth*position.BoardHeight - position.NumMoves) / 2
	maxVal := (position.BoardWidth*position.BoardHeight + 1 - position.NumMoves) / 2
	bestMove := -1

	if weak {
		minVal = -3
		maxVal = 3
	}

	tt := Transposition.NewTranspositionTable(1000000) // One million entries

	for minVal < maxVal {
		mid := minVal + (maxVal-minVal)/2
		
		if mid <= 0 && minVal/2 < mid {
			mid = minVal / 2
		} else if mid >= 0 && maxVal/2 > mid {
			mid = maxVal / 2
		}

		// Use a null window search
		r, move := Negamax(position, mid, mid+1, tt, searchDepth)

		if r <= mid {
			maxVal = r
		} else {
			minVal = r
		}
		
		bestMove = move
	}

	return minVal, bestMove
}

// MakeBestMove analyzes the position and returns the best move
func MakeBestMove(position *Position.Position) int {
	_, bestMove := Solve(position, false, 5, 10)
	return bestMove
}