package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func timeHandler(w http.ResponseWriter, req *http.Request) {
	currentTime := time.Now()
	switch req.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Heure actuelle : %v%v%v", currentTime.Hour(), "h", currentTime.Minute())
	}
}

func addEntrieHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}
		newEntrie := make(map[string]string)
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
		for key, value := range req.PostForm {
			if key == "entry" {
				newEntrie["entry"] = fmt.Sprint(value[0])
			}
			if key == "author" {
				newEntrie["author"] = fmt.Sprint(value[0])
			}
		}
		saveEntrie(newEntrie)
	}
}

func saveEntrie(newEntrie map[string]string) {
	saveFile, err := os.OpenFile("./entrie.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	defer saveFile.Close()

	w := bufio.NewWriter(saveFile)
	if err == nil {
		fmt.Fprintf(w, "%s: %s\n", newEntrie["author"], newEntrie["entry"])
	}
	w.Flush()
}

func getEntriesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		file, err := os.Open("entrie.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			if len(strings.Split(sc.Text(), ":")) > 1 {
				fmt.Fprintf(w, "%s\n", strings.Split(sc.Text(), ":")[1])
			}
		}

		if err := sc.Err(); err != nil {
			log.Fatal(err)
		}

	}
}

func main() {
	http.HandleFunc("/", timeHandler)
	http.HandleFunc("/add", addEntrieHandler)
	http.HandleFunc("/entries", getEntriesHandler)
	http.ListenAndServe(":4567", nil)
}
