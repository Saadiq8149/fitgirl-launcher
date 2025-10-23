package handlers

import (
	"fitgirl-launcher/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type InstallHandler struct {
	TorrentHandler *TorrentHandler
}

func CreateInstallHandler(th *TorrentHandler) *InstallHandler {
	return &InstallHandler{
		TorrentHandler: th,
	}
}

func (ih *InstallHandler) InstallRepack(magnet string) error {
	torrent, err := ih.TorrentHandler.GetTorrent(magnet)
	if err != nil {
		return err
	}

	if torrent == nil {
		return fmt.Errorf("torrent not found")
	}

	if torrent.State != utils.TORRENTUPLOADING && torrent.State != utils.TORRENTSTALLEDUPLOAD {
		return fmt.Errorf("torrent still downlaoding")
	}

	installerFilePath := fmt.Sprintf("%s/%s", torrent.ContentPath, "setup.exe")

	fmt.Println("Path: ", installerFilePath)

	gameName := strings.Split(torrent.ContentPath, `\`)[len(strings.Split(torrent.ContentPath, `\`))-1]

	err = exec.Command(installerFilePath, "/VERYSILENT", "/SP-", "/NOCANCEL", "/NORESTART", "/SUPPRESSMSGBOXES", "/COMPONENTS=text", `/DIR=D:\Games\Fitgirl-Launcher\`+gameName).Start()

	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	err = exec.Command("./SoundVolumeView.exe", "/Mute", "setup.tmp").Start()

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	err = os.RemoveAll(`D:\Games\Fitgirl-Launcher\` + gameName + `\_Redist`)

	if err != nil {
		return err
	}

	return nil
}
