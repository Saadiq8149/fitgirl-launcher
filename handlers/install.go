package handlers

import (
	"fitgirl-launcher/models"
	"fitgirl-launcher/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type InstallHandler struct {
	TorrentHandler         *TorrentHandler
	CurrentInstallationPID int
}

func CreateInstallHandler(th *TorrentHandler) *InstallHandler {
	return &InstallHandler{
		TorrentHandler: th,
	}
}

func (ih *InstallHandler) InstallRepack(magnet string, outputDir string) error {

	if ih.CurrentInstallationPID != 0 {
		return fmt.Errorf("another installation is in progress")
	}

	torrent, err := ih.TorrentHandler.GetTorrent(magnet)
	if err != nil {
		return err
	}

	if torrent == nil {
		return fmt.Errorf("torrent not found")
	}

	if torrent.State != utils.TORRENT_UPLOADING && torrent.State != utils.TORRENT_STALLED_UPLOAD {
		return fmt.Errorf("torrent still downlaoding")
	}

	cmd := exec.Command(
		filepath.Join(torrent.ContentPath, "setup.exe"),
		"/VERYSILENT", "/SP-", "/NOCANCEL", "/NORESTART", "/SUPPRESSMSGBOXES",
		"/DIR="+outputDir, "COMPONENTS=text",
	)

	// cmd := exec.Command(
	// 	filepath.Join(torrent.ContentPath, "setup.exe"), "/SUPPRESSMSGBOXES",
	// )
	cmd.Dir = torrent.ContentPath // crucial

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start installer: %w", err)
	}

	ih.CurrentInstallationPID = cmd.Process.Pid

	time.Sleep(5 * time.Second)

	err = exec.Command(`C:\Users\shaik\fitgirl-launcher\SoundVolumeView.exe`, "/Mute", "setup.tmp").Start()

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	err = os.RemoveAll(outputDir + string(os.PathSeparator) + "_Redist")

	if err != nil {
		return err
	}

	return nil
}

func (ih *InstallHandler) IsInstallationError(game models.Game) bool {
	// If game status says it's installing but we have no installation PID, that's an error
	if game.Status == utils.DB_INSTALLING && ih.CurrentInstallationPID == 0 {
		return true
	}

	if game.Status == utils.DB_INSTALLING && ih.CurrentInstallationPID != 0 && !utils.IsProcessRunning(ih.CurrentInstallationPID) {
		return true
	}

	return false
}

func (ih *InstallHandler) IsInstallationCompleted(game models.Game) bool {
	if ih.CurrentInstallationPID == 0 {
		return false
	}

	if game.Status != utils.DB_INSTALLING {
		fmt.Printf("DEBUG: Game status is not DB_INSTALLING (status: %v), returning false\n", game.Status)
		return false
	}

	processRunning := utils.IsProcessRunning(ih.CurrentInstallationPID)

	if !processRunning && ih.CurrentInstallationPID != 0 {
		ih.CurrentInstallationPID = 0
		return true
	}

	return false
}

func (ih *InstallHandler) IsInstallationInProgress() bool {
	return ih.CurrentInstallationPID != 0
}
