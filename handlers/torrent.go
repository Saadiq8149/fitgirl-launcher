package handlers

import (
	"fitgirl-launcher/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/superturkey650/go-qbittorrent/qbt"
)

type TorrentHandler struct {
	Client *qbt.Client
}

func CreateTorrentHandler(qb *qbt.Client) *TorrentHandler {
	return &TorrentHandler{
		Client: qb,
	}
}

func (th *TorrentHandler) AddTorrent(magnet string) error {
	if err := th.Client.DownloadLinks([]string{magnet}, qbt.DownloadOptions{
		Category: &utils.FITGIRLCATEGORY,
	}); err != nil {
		return err
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
	infoHash := th.infoHashFromMagnet(magnet)

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
	infoHash := th.infoHashFromMagnet(magnet)

	if infoHash == "" {
		return fmt.Errorf("invalid magnet link")
	}

	if err := th.Client.Delete([]string{infoHash}, true); err != nil {
		return err
	}
	return nil
}

func (th *TorrentHandler) GetTorrentFiles(magnet string) ([]qbt.TorrentFile, error) {
	infoHash := th.infoHashFromMagnet(magnet)

	if infoHash == "" {
		return nil, fmt.Errorf("invalid magnet link")
	}

	files, err := th.Client.TorrentFiles(infoHash)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (th *TorrentHandler) infoHashFromMagnet(magnet string) string {
	r, _ := regexp.Compile(`urn:btih:([a-fA-F0-9]{40})`)
	matches := r.FindStringSubmatch(magnet)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return ""
}
