package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var englishWords []string

func init() {
	englishWordsFile, err := os.Open("assets/words_alpha.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer englishWordsFile.Close()
	scanner := bufio.NewScanner(englishWordsFile)
	for scanner.Scan() {
		englishWords = append(englishWords, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func RandomSlug() string {
	return fmt.Sprintf("%s-%v%v%v%v", englishWords[rand.Intn(len(englishWords))], rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
}
