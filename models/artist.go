package models

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	SpotifyID    string   // Champ pour ton lecteur de musique
	Relations    Relation // Champ pour stocker les concerts
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
