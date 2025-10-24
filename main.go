package main

import (
	"embed"
	"fitgirl-launcher/handlers"

	"github.com/superturkey650/go-qbittorrent/qbt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()

	qb := qbt.NewClient("http://localhost:8080")
	dh := handlers.NewDatabaseHandler()
	fs := handlers.CreateFitgirlScraperHandler()
	th := handlers.CreateTorrentHandler(qb, dh)
	ih := handlers.CreateInstallHandler(th)
	sh := handlers.CreateSyncHandler(dh, th, ih)
	gh := handlers.CreateGameHandler()

	go sh.Sync()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "wails-events",
		Width:            1024,
		Height:           768,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			fs,
			th,
			ih,
			dh,
			sh,
			gh,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
