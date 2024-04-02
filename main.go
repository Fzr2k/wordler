package main

import (
	"bufio"
	f "fmt"
	"os"
	"strings"
)

func main() {
	f.Println("Welcome to Wordler!")

	// Read the content of the words.txt file
	content, err := os.ReadFile("words.txt")
	if err != nil {
		f.Println("Error reading file:", err)
		return
	}

	// Convert the content to a string
	words := string(content)

	// Split the content by newline to get individual words
	filteredWords := strings.Split(words, "\n")

	var wordList []string
	for _, word := range filteredWords {
		if len(word) == 5 {
			wordList = append(wordList, word)
		}
	}

	// wordList := []string{"crane", "crank", "crone", "crony", "crown"}

	f.Println("Enter your correctly placed letters (Use _ for missing letters): ")
	correctLetters := getUserInput()

	f.Println("Enter your correctly guessed (but incorrectly placed) letters (Use _ for missing letters): ")
	guessedLetters := getUserInput()

	f.Println("Enter your exhausted letters: ")
	exhaustedLetters := getUserInput()

	resultingWords := filterWords(wordList, correctLetters, guessedLetters, exhaustedLetters)
	f.Println("Filtered Words:", resultingWords)
}

func getUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func filterWords(words []string, correctLetters, guessedLetters, exhaustedLetters string) []string {
	var filteredWords []string

	for _, word := range words {
		if strings.ContainsAny(word, exhaustedLetters) {
			continue
		}

		if !containsAtIndices(word, correctLetters) {
			continue
		}

		if !containsLetters(word, guessedLetters) {
			continue
		}

		filteredWords = append(filteredWords, word)
	}

	return filteredWords
}

func containsAtIndices(word, letters string) bool {
	for i, letter := range letters {
		if letter == '_' {
			continue // Skip empty slots
		}
		if i >= len(word) || word[i] != byte(letter) {
			return false
		}
	}
	return true
}

func containsLetters(word, letters string) bool {
	for _, letter := range letters {
		if !strings.ContainsRune(word, letter) {
			return false
		}
	}
	return true
}
