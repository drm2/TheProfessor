package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

var local *leveldb.DB

func main() {

	sent := findWords("David is a banana butt!")
	for i := 0; i < len(sent); i++ {
		fmt.Print(Get(sent[i]) + " ")
	}

	// fmt.Println(Get(""))

	// err = db.Delete([]byte("key"), nil)
}

func Get(word string) string {

	if v := check(word); v != "" && v != "[error]" { //Word is in DB
		return v + "[db]"
	} else if v := searchWeb(word); v != "" && v != "[error]" { //Word is not in DB, Get from web
		if v != "" && v != "[error]" && v != "[nf]" {
			save(word, v)
		}

		return v
		// v := searchWeb(word)
		// if v != "" || v != "[error]" {
		// 	return v
		// }
	}
	return ""
}

//THIS IS CAUSING ISSUES!!!!
//Consider having multiple source options
//this may be a better api
//http://developer.pearson.com/content-apis/get-started
func searchWeb(word string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://wordsapiv1.p.mashape.com/words/"+word+"/", nil)
	if err != nil {
		return "[error]"
	}
	// ...
	req.Header.Add(`X-Mashape-Key`, `cSGPQwkQYNmshsigy0rUnRztVmj7p1FNHsBjsnHTpXlA4cE1RE`)
	//req.Header.Add("Accept","application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "[error]"
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()

	var v string
	if len(parse(s).Results) > 0 {
		v = parse(s).Results[0].PartOfSpeech
	} else {
		v = "[nf]"
	}
	return v
}

func save(word, pos string) {
	fmt.Println("[save{" + word + ":" + pos + "}]")
	db.Put([]byte(word), []byte(pos), nil)
}

func check(word string) string {
	val, err := db.Get([]byte(word), nil)
	if err != nil {
		return ""
	} else {
		return string(val)
	}
}

var db *leveldb.DB

func init() {
	Db, err := leveldb.OpenFile("dictionary", nil)
	db = Db
	if err != nil {
		panic(err)
	}
}

type Word struct {
	// Frequency     float64 `json:"frequency"`
	// Pronunciation struct {
	// 	All string `json:"all"`
	// } `json:"pronunciation"`
	Results []struct {
		// Definition   string   `json:"definition"`
		// MemberOf     []string `json:"memberOf"`
		PartOfSpeech string `json:"partOfSpeech"`
		// Synonyms     []string `json:"synonyms"`
		// TypeOf       []string `json:"typeOf"`
	} `json:"results"`
	// Syllables struct {
	// 	Count int      `json:"count"`
	// 	List  []string `json:"list"`
	// } `json:"syllables"`
	Word string `json:"word"`
}

func parse(jsonString string) Word {
	res := Word{}
	json.Unmarshal([]byte(jsonString), &res)
	return res
}

func removeChar(subject string, char string) string {
	var newString string = ""
	for i := 0; i < len(subject); i++ {
		if string(subject[i]) != char {
			newString += string(subject[i])
		}
	}
	return newString
}

func findWords(sentence string) []string {
	var lastIndex int = 0
	var numletters int8 = 0
	var numspaces int8 = 0
	var numwords int = 0
	var inWord bool = false

	for i := 0; i < len(sentence); i++ {
		if string(sentence[i]) == " " {
			if inWord {
				numspaces++
				numwords++
			}
			inWord = false
		} else {
			numletters++
			inWord = true
		}

	}

	words := make([]string, numwords+1)
	numwords = 0

	for i := 0; i < len(sentence); i++ {
		if string(sentence[i]) == " " {
			if inWord {
				numspaces++
				words[numwords] = sentence[lastIndex:i]
				numwords++
				lastIndex = i
			}
			inWord = false
		} else {
			numletters++
			inWord = true
		}

	}

	if lastIndex < len(sentence) {
		words[numwords] = sentence[lastIndex:len(sentence)]
	}

	for i := range words {
		if strings.Contains(words[i], ".") {
			words[i] = removeChar(words[i], ".")
		} else if strings.Contains(words[i], ",") {
			words[i] = removeChar(words[i], ",")
		} else if strings.Contains(words[i], "?") {
			words[i] = removeChar(words[i], "?")
		} else if strings.Contains(words[i], "!") {
			words[i] = removeChar(words[i], "!")
		}
	}

	return words
}