package main

import
(
	"log"
	"path/filepath"
	"fmt"
	"os"
	"net/http"
	"html/template"
	"sort"
)


type Page struct {
    Title string
    Body  []PageEntry
}

type PageEntry struct {
    Date string
    Images []string
    Caption []byte
}

func loadPage(title string) (*Page, error) {
    files, _ := filepath.Glob("data/"+title + "/*.txt")
	sort.Strings(files)
	entries := []PageEntry{}
	for i := 0; i < len(files); i++ {
		date := files[i][len("data/"+title + "/"):len(files[i])-4]
		images, _ := filepath.Glob("data/"+title + "/" + date + "*.jpg")
		caption, _ := os.ReadFile(files[i])
		entries = append(entries, PageEntry{Date: date, Images: images, Caption: caption})

	}
	return &Page{Title: title, Body: entries}, nil
}

func plantHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/plant/"):]
	p, err := loadPage(title)
    if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
        return
    }
    t, _ := template.ParseFiles("plant.html")
    t.Execute(w, p)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    files, _ := filepath.Glob("data/*")
	files2 := []string{}
	for i := 0; i < len(files); i++ {
		files2 = append(files2, files[i][len("data/"):])
	}
	sort.Strings(files2)
    t, _ := template.ParseFiles("index.html")
    t.Execute(w, files2)
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("about.html")
    t.Execute(w, nil)
}



func main() {
    http.HandleFunc("/plant/", plantHandler)
    http.HandleFunc("/home", homeHandler)
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusFound)
	})
	// files (images, css). this is a pretty bad way of handling it
	http.Handle("/x/", http.StripPrefix("/x/", http.FileServer(http.Dir("."))))
	fmt.Println("Running")
    log.Fatal(http.ListenAndServe(":4342", nil))
}

