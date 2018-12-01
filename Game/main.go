package main

import "net/http"

func HTML_MAIN(w http.ResponseWriter, r* http.Request)  {

}

func main() {
    http.HandleFunc("/", HTML_MAIN)
}
