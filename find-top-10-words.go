package main

import (
  "fmt"
  "html/template"
  "net/http"
  "log"
  "strings"
  "sort"
  "strconv"
)

func submit(w http.ResponseWriter, r *http.Request){
    fmt.Println("method:", r.Method) //get request method
    if r.Method == "GET" {
        t, _ := template.ParseFiles("submit.html")
        t.Execute(w, nil)
    } else {
        // logic part of submit
        r.ParseForm()

        fmt.Println("text: ", r.Form["textbody"])
        var count int
        wordDict := make(map[string]int)

        unparsedText := r.Form["textbody"]
        newText := strings.NewReplacer(",", "", ".", "", ";", "") //remove punctuation
        unparsedText[0] = newText.Replace(unparsedText[0])
        words := strings.Fields(strings.ToLower(unparsedText[0])) //convert to lowercase, writing to words
        //parsing each word, adding each word to a map and count them
        for i, word := range words {
              fmt.Println(i, " -> ", word)
              wordDict[word] = count
              cWord := word
              for i := 0; i < len(words); i++ {
                if words[i] == cWord {
                    count++
                    wordDict[word] = count
                }
              }
              count = 0 //reset
        }
        fmt.Println(wordDict)
        sortedList := rankByWordCount(wordDict)
        fmt.Println(sortedList)
        fmt.Fprintf(w, "Results\n")
        // Write TOP 10 most frequent words on the webpage
        for i, word := range sortedList {
            fmt.Fprintf(w, "Number " + strconv.Itoa(i+1) + ":\t" + word.String() + "\n")
            if i == 9 {
              return
            }
        }


    }
}
// function to convert struct to string for output
func (this keyPair) String() string {
    return this.Key + "\t\t" + strconv.Itoa(this.Value)
}
type keyPair struct {
    Key   string
    Value int
}

type WordList []keyPair
func (w WordList) Len() int { return len(w) }
func (w WordList) Less(i, j int) bool { return w[i].Value < w[j].Value }
func (w WordList) Swap(i, j int){ w[i], w[j] = w[j], w[i] }

func rankByWordCount(wordFrequencies map[string]int) WordList{
    wl := make(WordList, len(wordFrequencies))
    i := 0
    for k, v := range wordFrequencies {
        wl[i] = keyPair{k, v}
        i++
    }
    sort.Sort(sort.Reverse(wl))
    return wl
}

func main() {

    http.HandleFunc("/submit", submit)
    err := http.ListenAndServe(":8080", nil) //setting listening port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
