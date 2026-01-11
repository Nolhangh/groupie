package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"groupie/models"
)

func main() {
	// Routes
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist", ArtistHandler)

	// Fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🚀 Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HomeHandler : Récupère tous les artistes pour la page d'accueil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Erreur serveur API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var artists []models.Artist
	json.NewDecoder(resp.Body).Decode(&artists)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur template index", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, artists)
}

// ArtistHandler : Récupère UNIQUEMENT l'artiste sélectionné (Beaucoup plus rapide)
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 1. Récupérer les données de base de l'artiste par son ID
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + idStr)
	if err != nil {
		http.Error(w, "Artiste introuvable", http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	var artist models.Artist
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		http.Error(w, "Erreur de données", http.StatusInternalServerError)
		return
	}

	// 2. Récupérer les relations (concerts) pour cet artiste spécifique
	relResp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + idStr)
	if err == nil {
		defer relResp.Body.Close()
		var rel models.Relation
		json.NewDecoder(relResp.Body).Decode(&rel)
		artist.Relations = rel
	}

	// 3. Ajouter l'ID Spotify manuellement
	spotifyMap := map[string]string{
		"Queen":      "1dfeR4HaWDbWqFHLkxsg1d",
		"Pink Floyd": "0k17h0D3J5VfsdmQ1iZtE9",
		"SOJA":       "6Fx1cjY6uJqB3FqomkLzXU",
		"Scorpions":  "27T030eWyCQRmDyuvr1kxY",
	}
	artist.SpotifyID = spotifyMap[artist.Name]

	// 4. Envoyer au template
	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Fichier artist.html manquant", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, artist)
}
