package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Servasdfsdfnt sdfsdfmy geboku fssssucfsdfsdfkun sfsdfdsdfsfsdfdssfsdfsddsfdsdfsdfffsdfs\n")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
