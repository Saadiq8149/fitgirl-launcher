package models

type Game struct {
	Title       string `json:"title"`
	InstallPath string `json:"install_path"`
	Size        string `json:"size"`
	Url         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
	Status      string `json:"status"`
	Magnet      string `json:"magnet"`
}

type Database struct {
	Games []Game `json:"games"`
}
