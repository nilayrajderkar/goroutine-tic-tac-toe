package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var Board struct {
	currBoard          [3][3]int
	insertedValueCount int
	mu                 sync.Mutex // To ensure thread-safe operations
}

func isGameWon(board [3][3]int, value int) bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if board[i][0] == value && board[i][1] == value && board[i][2] == value {
			fmt.Println("Game won by completing a row by:", value)
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if board[0][i] == value && board[1][i] == value && board[2][i] == value {
			fmt.Println("Game won by completing a column by:", value)
			return true
		}
	}

	// Check diagonals
	if board[0][0] == value && board[1][1] == value && board[2][2] == value {
		fmt.Println("Game won by completing a diagonal by:", value)
		return true
	}

	if board[0][2] == value && board[1][1] == value && board[2][0] == value {
		fmt.Println("Game won by completing a diagonal by:", value)
		return true
	}

	return false
}
	

func fillBoard(value int, ch chan int) {
	Board.mu.Lock() // Lock the Board to prevent concurrent modifications
	defer Board.mu.Unlock() // Ensure the mutex is unlocked if any return occurs
	for {
		i := rand.Intn(3)
		j := rand.Intn(3)
		if Board.currBoard[i][j] == 0 {
			Board.currBoard[i][j] = value
			Board.insertedValueCount++
			fmt.Println("Board after inserting value:", value)
			printBoard()
			// Reutrn if the board is filled or someone has won
			if isGameWon(Board.currBoard, value){
				ch <- 2 // Send the game status to end the game
				return
			}

			if Board.insertedValueCount == 9 {
				fmt.Println("Game ends in a draw. Board filled")
				ch <- 2 // Send the game status to end the game
				return
			}

			ch <- (value) % 2 // Send the next value to be inserted
			return
		}
	}
}

func printBoard() {

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch Board.currBoard[i][j] {
			case 1:
				fmt.Print("X ")
			case 2:
				fmt.Print("O ")
			default:
				fmt.Print(". ") // Use '.' for empty spaces
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	fmt.Println("Game started")
	done := make(chan bool)
	ch := make(chan int)

	// Launch the goroutines to fill the board with 1 and 2
	go func() {
		for {
			gameStatus := <-ch
			if gameStatus == 2 {
				done <- true
				return
			} else if gameStatus == 0 {
				go fillBoard(1, ch)
			} else if gameStatus == 1 {
				go fillBoard(2, ch)
			}
		}
	}()
	
	gameStatus := rand.Intn(2) 
	ch <- gameStatus // Start the process by sending the initial value
	
	<-done
	fmt.Println("Game over")
}
