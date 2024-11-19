package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Orders struct {
	OBDs []Order `json:"OBDs"`
}

type Order struct {
	OBD   string `json:"OBD"`
	Lines []Line `json:"Lines"`
}

type Line struct {
	PartNumber   string      `json:"Part Number"`
	Batch        string      `json:"Batch"`
	SerialNumber SerialField `json:"Serial Number"`
	HandlingUnit string      `json:"Handling Unit"`
	Quantity     int         `json:"Quantity"`
	Completed    bool        `json:"Completed"`
}

type SerialField []string

func (sf *SerialField) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*sf = []string{single}
		return nil
	}
	var array []string
	if err := json.Unmarshal(data, &array); err == nil {
		*sf = array
		return nil
	}

	return fmt.Errorf("invalid Serial Number Field: %s", string(data))
}

func main() {
	fmt.Println("App Starting.....")
	data, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println(err)
	}

	var orders Orders

	err = json.Unmarshal(data, &orders)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orders)

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, orders)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		input := r.PostFormValue("handlingUnit")
		if err != nil {
			fmt.Println("Failed to convert String")
		}
		var linesToPick []Line
		for _, v := range orders.OBDs {
			for _, w := range v.Lines {
				if !w.Completed {
					if w.HandlingUnit == input {
						linesToPick = append(linesToPick, w)
						w.Completed = true
						break
					}
				}
			}
		}
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "Display Lines", linesToPick)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/show-line/", h2)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
