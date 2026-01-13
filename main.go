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

	log.Println(" Serveur lancé sur http://localhost:8080")
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
		"XXXTentacion": "15UsOTvKbsmSXD2v6vYdyO", "Mac Miller": "4LLpKhyESsy734mZ8iUv2v",
		"Joyner Lucas": "6i6foloYpC8fw8EM3pYVE2", "Kendrick Lamar": "2YZyLoL8N0Wb9xBt1NhCgG",
		"ACDC":                        "711ySwwmS7A6H0UefSTpDn",
		"Pearl Jam":                   "1w5vsy1BhDY9YM4oPdUOSI",
		"Katy Perry":                  "6jJ0s8O6EPuF9sF3rT6S31",
		"Rihanna":                     "5pKCCm2ICUXRM0vYm09vU0",
		"Genesis":                     "3vYduoqX8U97D9v999p9C4",
		"Phil Collins":                "44tN9s7YAh6us9uS7Xoa8T",
		"Led Zeppelin":                "36QSy2vCdbXv0vO9vZDXoX",
		"The Jimi Hendrix Experience": "1EALhRhlL9pApm0ZzbCOjM",
		"Bee Gees":                    "1LZEwsT5SIlyGcyB6v9S1l",
		"Deep Purple":                 "568Wfb9Mfs962o9j9Ws9u0",
		"Aerosmith":                   "7Ey4PDuEah5DMs9yt2O6An",
		"Dire Straits":                "0v1XpXWs34F9JpXU9vX9X9",
		"Mamonas Assassinas":          "2p9i6U9X9u6U9X9u6U9X9u",
		"Thirty Seconds to Mars":      "06HL4z0CvFAxy0p7sP0STu",
		"Imagine Dragons":             "53XhwfbYqKCa1cCjXmD569",
		"Juice Wrld":                  "4MCBfE4YhI0S8m5pB0xy9O",
		"Logic":                       "49X9u6U9X9u6U9X9u6U9X9",
		"Alec Benjamin":               "5i7puiovTMqXvH49I0S8m5",
		"Bobby McFerrins":             "3u6S9u6S9u6S9u6S9u6S9u",
		"R3HAB":                       "6f6U9X9u6U9X9u6U9X9u6U",
		"Post Malone":                 "246E6vSwsR9pIrtc99o9uO",
		"Travis Scott":                "0Y5tJX1vP1UMmc9Lp0699B",
		"J. Cole":                     "6l3HvQhPBMvM99P9P9P9P9",
		"Nickelback":                  "6Y73S9u6S9u6S9u6S9u6S9",
		"Mobb Deep":                   "6O9X9u6U9X9u6U9X9u6U9X",
		"Guns N' Roses":               "3dBVy97mB7whYQpScyqAhH",
		"NWA":                         "49X9u6U9X9u6U9X9u6U9X9",
		"U2":                          "51w9P9P9P9P9P9P9P9P9P9",
		"Arctic Monkeys":              "7Ln80S9u6S9u6S9u6S9u6S",
		"Fall Out Boy":                "49X9u6U9X9u6U9X9u6U9X9",
		"Gorillaz":                    "3u6S9u6S9u6S9u6S9u6S9u",
		"Eagles":                      "0v1XpXWs34F9JpXU9vX9X9",
		"Linkin Park":                 "6XyY86sb0U5ODiMRBKSTmG",
		"Red Hot Chili Peppers":       "0L89CY0r9kLVdtT9spt8ab",
		"Eminem":                      "7dG3Yv97mB7whYQpScyqAhH",
		"Green Day":                   "7oPftvl6af9Gqy9vX9u6U9",
		"Metallica":                   "2ye2Wgw4gimLv2eAKykI1B",
		"Coldplay":                    "4gzpq5YvC9S9u6S9u6S9u6",
		"Maroon 5":                    "04gDigrSndNo9RS9u6S9u6",
		"Twenty One Pilots":           "3YQKmS9u6S9u6S9u6S9u6S",
		"The Rolling Stones":          "22b9P9P9P9P9P9P9P9P9P9",
		"Muse":                        "12Ch79S9m7mU7777777777",
		"Foo Fighters":                "7jy3rLJdDQY2crqiIB7oX3",
		"The Chainsmokers":            "69S9m7mU77777777777777",
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
