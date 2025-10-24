package handlers

import (
	"fitgirl-launcher/models"
	"fitgirl-launcher/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type GameHandler struct {
	runningGames map[string]*os.Process
}

func CreateGameHandler() *GameHandler {
	return &GameHandler{
		runningGames: make(map[string]*os.Process),
	}
}
func (gh *GameHandler) LaunchGame(game models.Game) error {
	fmt.Printf("DEBUG: Launching game: %s\n", game.Title)
	fmt.Printf("DEBUG: Install path: %s\n", game.InstallPath)

	if proc, exists := gh.runningGames[game.Url]; exists && proc != nil {
		fmt.Printf("DEBUG: Found existing process for game, checking if running (PID: %d)\n", proc.Pid)
		if utils.IsProcessRunning(proc.Pid) {
			fmt.Printf("DEBUG: Game is already running\n")
			return fmt.Errorf("game is already running")
		}
		fmt.Printf("DEBUG: Existing process is not running, proceeding\n")
	}

	files, err := os.ReadDir(game.InstallPath)
	if err != nil {
		fmt.Printf("DEBUG: Error reading install directory: %v\n", err)
		return err
	}

	var executables []string
	var bestMatch string
	bestScore := -1

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".exe") {
			executables = append(executables, file.Name())
			fmt.Printf("DEBUG: Found executable: %s\n", file.Name())
		}
	}

	if len(executables) == 0 {
		fmt.Printf("DEBUG: No executables found in install path\n")
		return fmt.Errorf("no executable found in install path")
	}

	fmt.Printf("DEBUG: Total executables found: %d\n", len(executables))
	gameTitle := strings.ToLower(game.Title)
	fmt.Printf("DEBUG: Normalized game title: %s\n", gameTitle)

	for _, exe := range executables {
		exeName := strings.ToLower(strings.TrimSuffix(exe, ".exe"))
		score := 0
		fmt.Printf("DEBUG: Evaluating executable: %s (normalized: %s)\n", exe, exeName)

		if strings.Contains(exeName, gameTitle) || strings.Contains(gameTitle, exeName) {
			score += 100
			fmt.Printf("DEBUG: Title match bonus: +100\n")
		}

		gameTitleWords := strings.Fields(gameTitle)
		for _, word := range gameTitleWords {
			if len(word) > 2 && strings.Contains(exeName, word) {
				score += 10
				fmt.Printf("DEBUG: Word match bonus for '%s': +10\n", word)
			}
		}

		if strings.Contains(exeName, "setup") || strings.Contains(exeName, "install") || strings.Contains(exeName, "unity") ||
			strings.Contains(exeName, "uninstall") || strings.Contains(exeName, "config") {
			score -= 50
			fmt.Printf("DEBUG: Utility/installer penalty: -50\n")
		}

		fmt.Printf("DEBUG: Final score for %s: %d\n", exe, score)

		if score > bestScore {
			bestScore = score
			bestMatch = exe
			fmt.Printf("DEBUG: New best match: %s (score: %d)\n", bestMatch, bestScore)
		}
	}

	fmt.Printf("DEBUG: Selected executable: %s\n", bestMatch)
	gameExecutablePath := game.InstallPath + "/" + bestMatch
	fmt.Printf("DEBUG: Full executable path: %s\n", gameExecutablePath)

	cmd := exec.Command(gameExecutablePath)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("DEBUG: Error starting game: %v\n", err)
		return err
	}

	fmt.Printf("DEBUG: Game started successfully with PID: %d\n", cmd.Process.Pid)
	gh.runningGames[game.Url] = cmd.Process
	return nil
}
