package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"groupie/models"
)

// Client avec Timeout pour éviter le chargement infini
var client = &http.Client{Timeout: 3 * time.Second}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist", ArtistHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🚀 Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "API indisponible", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	var artists []models.Artist
	json.NewDecoder(resp.Body).Decode(&artists)

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, artists)
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 1. Récupérer l'artiste (Directement par ID)
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'artiste", 500)
		return
	}
	defer resp.Body.Close()

	var artist models.Artist
	json.NewDecoder(resp.Body).Decode(&artist)

	// 2. Récupérer les concerts (Relations)
	relResp, err := client.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)
	if err == nil {
		defer relResp.Body.Close()
		var rel models.Relation
		json.NewDecoder(relResp.Body).Decode(&rel)
		artist.Relations = rel
	}

	// 3. Spotify Map
	spotifyMap := map[string]string{
		"Queen": "1dfeR4HaWDbWqFHLkxsg1d", "Pink Floyd": "0k17h0D3J5VfsdmQ1iZtE9",
		"SOJA": "6Fx1cjY6uJqB3FqomkLzXU", "Scorpions": "27T030eWyCQRmDyuvr1kxY",
	}
	artist.SpotifyID = spotifyMap[artist.Name]

	// 4. Exécuter le template
	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Template artist.html introuvable", 500)
		return
	}
	tmpl.Execute(w, artist)
}
