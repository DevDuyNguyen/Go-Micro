package main

import (
	"html/template"
	"net/http"
)

type Person struct{
	Name string
	Age int
}

func templateCreateHandler(w http.ResponseWriter, r *http.Request){
	template,err:=template.ParseFiles("test1.gohtml", "base.gohtml")
	if err!=nil{
		panic(err.Error())
	}
	err=template.Execute(w, Person{
		Name:"Duy",
		Age:12,
	})
	if err!=nil{
		panic(err.Error())
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", templateCreateHandler)
	server:=&http.Server{
		Addr: ":3000",
		Handler:mux,
	}
	err:=server.ListenAndServe()
	if err!=nil{
		panic(err.Error())
	}
}