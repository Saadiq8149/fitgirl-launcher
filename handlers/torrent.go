package handlers

import (
	"fitgirl-launcher/models"
	"fitgirl-launcher/utils"
	"fmt"
	"os"

	"github.com/superturkey650/go-qbittorrent/qbt"
)

type TorrentHandler struct {
	Client          *qbt.Client
	DatabaseHandler *DatabaseHandler
}

func CreateTorrentHandler(qb *qbt.Client, dh *DatabaseHandler) *TorrentHandler {
	return &TorrentHandler{
		Client:          qb,
		DatabaseHandler: dh,
	}
}

func (th *TorrentHandler) AddTorrent(magnet string, repack models.FitgirlRepack, url string) error {
	if err := th.Client.DownloadLinks([]string{magnet}, qbt.DownloadOptions{
		Category: &utils.FITGIRLCATEGORY,
	}); err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	game := models.Game{
		Title:       repack.Name,
		InstallPath: homeDir + string(os.PathSeparator) + "fitgirl-repacks" + string(os.PathSeparator) + "games" + string(os.PathSeparator) + repack.Name,
		Thumbnail:   repack.CoverImage,
		Size:        repack.RepackSize,
		Url:         url,
		Magnet:      magnet,
		Status:      utils.DB_DOWNLOADING,
	}

	if err := th.DatabaseHandler.AddGameToDatabase(game); err != nil {
		return fmt.Errorf("failed to add game to database: %v", err)
	}

	return nil
}

func (th *TorrentHandler) GetTorrents() ([]qbt.TorrentInfo, error) {
	torrents, err := th.Client.Torrents(qbt.TorrentsOptions{
		Category: &utils.FITGIRLCATEGORY,
	})
	if err != nil {
		return nil, err
	}

	return torrents, nil
}

func (th *TorrentHandler) GetTorrent(magnet string) (*qbt.TorrentInfo, error) {
	infoHash := utils.InfoHashFromMagnet(magnet)

	if infoHash == "" {
		return nil, fmt.Errorf("invalid magnet link")
	}

	torrents, err := th.Client.Torrents(qbt.TorrentsOptions{
		Category: &utils.FITGIRLCATEGORY,
	})

	if err != nil {
		return nil, err
	}

	for _, torrent := range torrents {
		if torrent.Hash == infoHash {
			return &torrent, nil
		}
	}

	return nil, nil
}

func (th *TorrentHandler) RemoveTorrent(magnet string) error {
	infoHash := utils.InfoHashFromMagnet(magnet)

	if infoHash == "" {
		return fmt.Errorf("invalid magnet link")
	}

	if err := th.Client.Delete([]string{infoHash}, true); err != nil {
		return err
	}
	return nil
}

func (th *TorrentHandler) GetTorrentFiles(magnet string) ([]qbt.TorrentFile, error) {
	infoHash := utils.InfoHashFromMagnet(magnet)

	if infoHash == "" {
		return nil, fmt.Errorf("invalid magnet link")
	}

	files, err := th.Client.TorrentFiles(infoHash)
	if err != nil {
		return nil, err
	}

	return files, nil
}
