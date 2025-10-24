package handlers

import (
	"fitgirl-launcher/models"
	"fitgirl-launcher/utils"
	"fmt"
	"time"
)

type SyncHandler struct {
	dh *DatabaseHandler
	th *TorrentHandler
	ih *InstallHandler
}

func CreateSyncHandler(dh *DatabaseHandler, th *TorrentHandler, ih *InstallHandler) *SyncHandler {
	return &SyncHandler{
		dh: dh,
		th: th,
		ih: ih,
	}
}

func (sh *SyncHandler) Sync() {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Sync panic: %v\n", r)
				}
			}()
			sh.syncOnce()
		}()
		time.Sleep(5 * time.Second)
	}
}

func (sh *SyncHandler) syncOnce() {
	games, err := sh.dh.LoadDatabase()
	if err != nil {
		return
	}

	for _, game := range games.Games {
		sh.processGame(game)
	}
}

func (sh *SyncHandler) processGame(game models.Game) {
	torrent, err := sh.th.GetTorrent(game.Magnet)
	if err != nil {
		return
	}

	if torrent == nil {
		return
	}

	// if sh.ih.IsInstallationError(game) {
	// 	fmt.Printf("Installation error detected for game: %s\n", game.Url)
	// 	sh.dh.UpdateGameStatusDownloaded(game.Url)
	// 	return
	// }

	if sh.ih.IsInstallationCompleted(game) {
		fmt.Printf("Installation completed for game: %s\n", game.Url)
		err := sh.dh.UpdateGameStatusInstalled(game.Url)
		if err != nil {
			fmt.Printf("Failed to update status: %v\n", err)
		}
		return
	}

	if game.Status == utils.DB_DOWNLOADING && (torrent.State == utils.TORRENT_STALLED_UPLOAD || torrent.State == utils.TORRENT_UPLOADING) {
		err := sh.dh.UpdateGameStatusDownloaded(game.Url)
		if err != nil {
			fmt.Printf("Failed to update status: %v\n", err)
		}
		return
	}

	if game.Status == utils.DB_DOWNLOADED && game.InstallPath != "" {
		err := sh.dh.UpdateGameStatusInstalling(game.Url)
		if err != nil {
			return
		}
		err = sh.ih.InstallRepack(game.Magnet, game.InstallPath)
		if err != nil {
			return
		}
	}
}
