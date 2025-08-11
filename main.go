package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

const filePath = "./gopher.json"

var err error

type arcModel struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []optionsModel
}

type optionsModel struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handlerStruct struct {
	Data     map[string]arcModel
	Template *template.Template
}

func (h *handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arcValue := r.URL.Path
	arcValue, _ = strings.CutPrefix(arcValue, "/")
	if arcValue == "" {
		arcValue = "intro"
	}
	if err = h.Template.ExecuteTemplate(w, "Index", h.Data[arcValue]); err != nil {
		err = fmt.Errorf("erro ao executar o template: %w", err)
	}
}

func extractData() map[string]arcModel {
	var fileBytes []byte
	if fileBytes, err = os.ReadFile(filePath); err != nil {
		err = fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	var fileData map[string]arcModel
	if err = json.Unmarshal(fileBytes, &fileData); err != nil {
		err = fmt.Errorf("erro ao ler o JSON: %w", err)
	}
	return fileData
}

func main() {
	templ := template.Must(template.ParseGlob("templates/*.html"))
	pageHandle := handlerStruct{
		Data:     extractData(),
		Template: templ,
	}
	http.Handle("/", &pageHandle)
	fmt.Println("Servidor Iniciado")
	http.ListenAndServe(":8080", nil)
}
