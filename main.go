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
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f9;
					margin: 0;
					padding: 0;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
				}
				.container {
					background: white;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
					width: 300px;
					text-align: center;
				}
				h1 {
					color: #333;
				}
				label {
					display: block;
					margin-top: 10px;
					color: #555;
				}
				input[type="text"] {
					width: 100%;
					padding: 10px;
					margin-top: 5px;
					border: 1px solid #ccc;
					border-radius: 4px;
					box-sizing: border-box;
				}
				input[type="submit"] {
					background-color: #4CAF50;
					color: white;
					padding: 10px 20px;
					margin-top: 20px;
					border: none;
					border-radius: 4px;
					cursor: pointer;
				}
				input[type="submit"]:hover {
					background-color: #45a049;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Word Filter</h1>
				<form action="/filter" method="post">
					<label for="correct">Correctly placed letters:</label>
					<input type="text" id="correct" name="correct"><br>
					<label for="guessed">Guessed letters:</label>
					<input type="text" id="guessed" name="guessed"><br>
					<label for="exhausted">Exhausted letters:</label>
					<input type="text" id="exhausted" name="exhausted"><br>
					<input type="submit" value="Filter">
				</form>
			</div>
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
