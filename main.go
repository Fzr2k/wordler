package main

import (
	"fmt"
	"os"
	"html/template"
	"net/http"
	"strings"
)

var wordList []string

func init() {
	// Read the content of the words.txt file
	content, err := os.ReadFile("words.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		wordList = []string{"error"}
		return
	}

	// Convert the content to a string
	allWords := string(content)

	// Split the content by newline to get individual words
	fiveLetterWords := strings.Split(allWords, "\n")

	for _, word := range fiveLetterWords {
		if len(word) == 5 {
			wordList = append(wordList, word)
		}
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/filter", filterHandler)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Define HTML template
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Word Filter</title>
	</head>
	<body>
		<h1>Word Filter</h1>
		<form action="/filter" method="post">
			<label for="correct">Correctly placed letters:</label>
			<input type="text" id="correct" name="correct"><br><br>
			<label for="guessed">Guessed letters:</label>
			<input type="text" id="guessed" name="guessed"><br><br>
			<label for="exhausted">Exhausted letters:</label>
			<input type="text" id="exhausted" name="exhausted"><br><br>
			<input type="submit" value="Filter">
		</form>
	</body>
	</html>
	`

	// Execute template
	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	r.ParseForm()
	correct := r.Form.Get("correct")
	guessed := r.Form.Get("guessed")
	exhausted := r.Form.Get("exhausted")

	// Filter words
	filteredWords := filterWords(correct, guessed, exhausted)

	// Print filtered words as response
	fmt.Fprintf(w, "Filtered Words: %s", strings.Join(filteredWords, ", "))
}

func filterWords(correctLetters, guessedLetters, exhaustedLetters string) []string {
	var filteredWords []string

	for _, word := range wordList {
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
