package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"connect4/Position"
	"connect4/Solver"
)

// TestSolver tests the Connect Four solver against a file of test positions
func TestSolver(testfile string) {
	failedTests := 0
	movesToResult := make(map[string]int)
	
	// Read test file
	file, err := os.Open(testfile)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", testfile, err)
		return
	}
	defer file.Close()

	// Parse the test cases
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 {
			moves := parts[0]
			result, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Error parsing result in line %s: %v\n", line, err)
				continue
			}
			movesToResult[moves] = result
		}
	}

	totalTests := len(movesToResult)
	fmt.Printf("Running %d tests from %s\n", totalTests, testfile)

	// Run each test case
	for moves, expectedResult := range movesToResult {
		// Initialize bit boards
		player1Bs := uint64(0)
		player2Bs := uint64(0)
		colHeights := [7]int{0, 0, 0, 0, 0, 0, 0}

		// Convert move sequence to bit board representation
		for i, moveChar := range moves {
			col, _ := strconv.Atoi(string(moveChar))
			col-- // Convert 1-7 to 0-6

			// Calculate bit position
			row := colHeights[col]
			bitPos := col*(6+1) + row

			// Set bit in appropriate player's bit board
			if i%2 == 0 {
				player1Bs |= uint64(1) << uint(bitPos)
			} else {
				player2Bs |= uint64(1) << uint(bitPos)
			}
			
			colHeights[col]++
		}

		// Create position
		pos := Position.NewPosition()
		pos.NumMoves = len(moves)
		pos.CurrentPositions = [2]uint64{player1Bs, player2Bs}

		// Run solver
		result, _ := Solver.Solve(pos, false, 8, 25)
		
		fmt.Printf("Position after moves %s:\n", moves)
		pos.PrintBoard()
		fmt.Printf("Result: %d, Expected: %d\n", result, expectedResult)
		
		if result != expectedResult {
			failedTests++
			fmt.Printf("❌ Test failed for moves: %s\n", moves)
		} else {
			fmt.Printf("✓ Test passed\n")
		}
		fmt.Println("---------------------")
	}

	fmt.Printf("%d/%d tests passed\n", totalTests-failedTests, totalTests)
}

func main() {
	// Run tests
	fmt.Println("Running mini tests...")
	TestSolver("Test/mini_test.txt")
	
	fmt.Println("\nRunning hard tests...")
	TestSolver("Test/hard_test.txt")
}