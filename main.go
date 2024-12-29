package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type URL struct {
MainURL string `json:"main_url"`
SortURL string `json:"sort_url"`
}

var urlDB = make(map[string]URL)

func generateShortURL(url string) string {
 hash := md5.New()
 hash.Write([]byte(url))
 data := hash.Sum(nil)
 hased := hex.EncodeToString(data)
 return hased[:8]

}

func addURL(url string) string{
	SortURL := generateShortURL(url) 
	urlDB[SortURL] =URL{
		MainURL: url,
        SortURL: SortURL,
	}
	return SortURL
}

func deleteURL(id string) {
	delete(urlDB, id)
}

func updateURL(id, newURL string) {
	urlDB[id] = URL{
		MainURL: newURL,
        SortURL: id,
	}
}

func getURl(id string) (URL,error) {
	if url, ok := urlDB[id]; ok {
        return url, nil
    }
	return URL{}, fmt.Errorf("URL not found")
}

func handler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w , "Server started...")
}

func ShortUrlHandler(w http.ResponseWriter, r *http.Request){
	var data struct {
		URL_ string `json:"url"`	
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	sortUrl := addURL(data.URL_)
	 var res struct {
		URL string `json:"resp_url"`
	}
	res.URL = sortUrl
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)


}
func RedirectUrlHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURl(id)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
	http.Redirect(w, r, url.MainURL, http.StatusFound)
}

func DeleteUrlHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Path[len("/delete/"):]
	fmt.Println("Deleted id", id)
    deleteURL(id)
	fmt.Println("After Deleted id", id)
	fmt.Println(urlDB)
    fmt.Fprintf(w, "URL deleted successfully")
}

func UpdateUrlHandler(w http.ResponseWriter, r *http.Request){
	var data struct {
		ID_ string `json:"id"`
		URL_ string `json:"url"`	
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err!= nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	updateURL(data.ID_, data.URL_)
	fmt.Fprintf(w, "URL updated successfully")

}
func main() {
fmt.Println("Developing  URL sortner")
http.HandleFunc("/", handler)
http.HandleFunc("/sortner", ShortUrlHandler)
http.HandleFunc("/redirect/", RedirectUrlHandler)
http.HandleFunc("/delete/", DeleteUrlHandler)
http.HandleFunc("/update/", UpdateUrlHandler)


err := http.ListenAndServe(":8000", nil)
if err!= nil {
    panic(err)
}


}