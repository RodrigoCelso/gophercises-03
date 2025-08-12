package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

const filePath = "./gopher.json"

var err error

type arcModel struct {
	Title   string         `json:"title"`
	Story   []string       `json:"story"`
	Options []optionsModel `json:"options"`
}

type optionsModel struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handlerStruct struct {
	Story        map[string]arcModel
	Template     *template.Template
	Introduction string
}

func newHandler(story map[string]arcModel, intro string, opts ...func(*handlerStruct)) *handlerStruct {
	handler := &handlerStruct{Story: story, Introduction: intro}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}

func withTemplate(templ *template.Template) func(*handlerStruct) {
	return func(h *handlerStruct) {
		h.Template = templ
	}
}

func (h *handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arcValue := r.URL.Path
	arcValue, _ = strings.CutPrefix(arcValue, "/")
	if arcValue == "" {
		arcValue = h.Introduction
	}
	if err = h.Template.Execute(w, h.Story[arcValue]); err != nil {
		err = fmt.Errorf("erro ao executar o template: %w", err)
	}
}

func extractData() (map[string]arcModel, error) {
	var fileBytes []byte
	if fileBytes, err = os.ReadFile(filePath); err != nil {
		err = fmt.Errorf("erro ao abrir o arquivo: %w", err)
		return nil, err
	}
	var fileData map[string]arcModel
	if err = json.Unmarshal(fileBytes, &fileData); err != nil {
		err = fmt.Errorf("erro ao ler o JSON: %w", err)
		return nil, err
	}
	return fileData, nil
}

func writeArc(currentArc arcModel) []optionsModel {
	fmt.Println(currentArc.Title)
	fmt.Print("---------------------------\n\n")
	for _, s := range currentArc.Story {
		fmt.Println(s + "\n")
	}
	fmt.Print("---------------------------\n")
	for i, o := range currentArc.Options {
		fmt.Printf("%d - %s\n\n", i+1, o.Text)
	}
	fmt.Println("0 - Stop reading")
	fmt.Println("\nChoose your own path, type the number to choose where to go: ")
	return currentArc.Options
}

func main() {
	port := flag.String("port", "8080", "Port of web application")
	intro := flag.String("intro", "intro", "Introduction arc")
	clVersion := flag.Bool("cl", false, "Command line version")
	flag.Parse()

	data, err := extractData()
	if err != nil {
		panic(error.Error(err))
	}

	if *clVersion {
		for {
			options := writeArc(data[*intro])
			var response int
			if _, err = fmt.Scanln(&response); err != nil {
				err = fmt.Errorf("valor digitado não confere. %v", err)
				panic(err)
			}
			if response == 0 {
				os.Exit(0)
			}
			*intro = options[response-1].Arc
		}
	}

	templ := template.Must(template.ParseGlob("templates/index.html"))
	pageHandle := newHandler(data, *intro, withTemplate(templ))
	fmt.Println("Servidor Iniciado na porta", *port)
	http.ListenAndServe(":"+*port, pageHandle)
}
