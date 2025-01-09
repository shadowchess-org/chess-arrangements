package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	EXPECTED_ARRANGEMENTS = 216481
	POSITIONS_PER_FILE    = 1000000
)

type Piece struct {
	name      string
	fenSymbol string
}

func main() {
	// Check if "full" argument is present
	isFull := len(os.Args) > 1 && os.Args[1] == "full"

	// Generate white arrangements
	whiteArrangements := generateForColor("white", []Piece{
		{"rook", "R"}, {"rook", "R"},
		{"knight", "N"}, {"knight", "N"},
		{"bishop", "B"}, {"bishop", "B"},
		{"queen", "Q"}, {"king", "K"},
		{"pawn", "P"}, {"pawn", "P"},
		{"pawn", "P"}, {"pawn", "P"},
		{"pawn", "P"}, {"pawn", "P"},
		{"pawn", "P"}, {"pawn", "P"},
	})

	// Generate black arrangements
	blackArrangements := generateForColor("black", []Piece{
		{"rook", "r"}, {"rook", "r"},
		{"knight", "n"}, {"knight", "n"},
		{"bishop", "b"}, {"bishop", "b"},
		{"queen", "q"}, {"king", "k"},
		{"pawn", "p"}, {"pawn", "p"},
		{"pawn", "p"}, {"pawn", "p"},
		{"pawn", "p"}, {"pawn", "p"},
		{"pawn", "p"}, {"pawn", "p"},
	})

	// If full argument is present, generate combined positions
	if isFull {
		generateCombinedPositions(whiteArrangements, blackArrangements)
	}
}

func generateForColor(color string, pieces []Piece) [][]Piece {
	arrangements := generateTwoRankArrangements(pieces)

	filename := fmt.Sprintf("chess_arrangements_%s.txt", color)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file for %s pieces: %v\n", color, err)
		return arrangements
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, arr := range arrangements {
		fen := arrangementToFEN(arr)
		_, err := fmt.Fprintf(writer, "Arrangement %d: %s\n", i+1, fen)
		if err != nil {
			fmt.Printf("Error writing to file for %s pieces: %v\n", color, err)
			break
		}
	}
	writer.Flush()

	fmt.Printf("Generated %d arrangements for %s pieces and saved to %s\n",
		len(arrangements), color, filename)
	return arrangements
}

func generateCombinedPositions(whiteArrangements, blackArrangements [][]Piece) {
	// Create full directory if it doesn't exist
	err := os.MkdirAll("full", 0755)
	if err != nil {
		fmt.Printf("Error creating full directory: %v\n", err)
		return
	}

	totalCombinations := uint64(len(whiteArrangements)) * uint64(len(blackArrangements))
	currentFile := 0
	currentPositionInFile := 0
	var writer *bufio.Writer
	var file *os.File
	count := uint64(0)

	// Write summary file
	summaryFile, err := os.Create(filepath.Join("full", "summary.txt"))
	if err == nil {
		fmt.Fprintf(summaryFile, "Total combinations: %d\nPositions per file: %d\nTotal files: %d\n",
			totalCombinations, POSITIONS_PER_FILE, (totalCombinations+uint64(POSITIONS_PER_FILE)-1)/uint64(POSITIONS_PER_FILE))
		summaryFile.Close()
	}

	openNewFile := func() error {
		if file != nil {
			writer.Flush()
			file.Close()
		}

		filename := filepath.Join("full", fmt.Sprintf("positions_%d.txt", currentFile+1))
		var err error
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
		writer = bufio.NewWriter(file)
		currentPositionInFile = 0
		return nil
	}

	// Open first file
	err = openNewFile()
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer func() {
		if file != nil {
			writer.Flush()
			file.Close()
		}
	}()

	for _, whiteArr := range whiteArrangements {
		whiteFen := arrangementToFEN(whiteArr)

		for _, blackArr := range blackArrangements {
			blackFen := arrangementToFEN(blackArr)
			combinedFen := whiteFen + "-" + blackFen

			count++
			_, err := fmt.Fprintf(writer, "Position %d: %s\n", count, combinedFen)
			if err != nil {
				fmt.Printf("Error writing combined position: %v\n", err)
				return
			}

			currentPositionInFile++
			if currentPositionInFile >= POSITIONS_PER_FILE {
				currentFile++
				err = openNewFile()
				if err != nil {
					fmt.Printf("Error creating new output file: %v\n", err)
					return
				}
			}

			// Print progress every 5 million positions
			if count%5000000 == 0 {
				fmt.Printf("Progress: %.2f%% (%d/%d)\n",
					float64(count)/float64(totalCombinations)*100,
					count, totalCombinations)
			}
		}
	}

	fmt.Printf("Generated %d combined positions across %d files\n", count, currentFile+1)
}

func generateTwoRankArrangements(pieces []Piece) [][]Piece {
	var result [][]Piece
	var firstRank []Piece
	generateCombinations(pieces, 8, firstRank, &result)
	return result
}

func generateCombinations(remainingPieces []Piece, firstRankSpaces int, currentFirstRank []Piece, result *[][]Piece) {
	if firstRankSpaces == 0 {
		arrangement := make([]Piece, 16)
		copy(arrangement[:8], currentFirstRank)
		copy(arrangement[8:], remainingPieces)
		if isValidArrangement() {
			*result = append(*result, arrangement)
		}
		return
	}

	if len(remainingPieces) < firstRankSpaces {
		return
	}

	used := make(map[string]bool)
	for i := 0; i < len(remainingPieces); i++ {
		if used[remainingPieces[i].name] {
			continue
		}
		used[remainingPieces[i].name] = true

		newRemaining := make([]Piece, 0, len(remainingPieces)-1)
		newRemaining = append(newRemaining, remainingPieces[:i]...)
		newRemaining = append(newRemaining, remainingPieces[i+1:]...)

		newFirstRank := append([]Piece{}, currentFirstRank...)
		newFirstRank = append(newFirstRank, remainingPieces[i])

		generateCombinations(newRemaining, firstRankSpaces-1, newFirstRank, result)
	}
}

func isValidArrangement() bool {
	return true
}

func arrangementToFEN(pieces []Piece) string {
	var firstRank []string
	for i := 0; i < 8; i++ {
		firstRank = append(firstRank, pieces[i].fenSymbol)
	}

	var secondRank []string
	for i := 8; i < 16; i++ {
		secondRank = append(secondRank, pieces[i].fenSymbol)
	}

	return strings.Join(firstRank, "") + "/" + strings.Join(secondRank, "")
}
