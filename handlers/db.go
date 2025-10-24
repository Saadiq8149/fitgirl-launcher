package handlers

import (
	"encoding/json"
	"fitgirl-launcher/models"
	"fitgirl-launcher/utils"
	"fmt"
	"os"
	"regexp"
)

type DatabaseHandler struct {
	DatabasePath string
}

func NewDatabaseHandler() *DatabaseHandler {
	path, err := os.UserHomeDir()

	fmt.Println("home: ", path)

	if err != nil {
		fmt.Print("error getting home directory")
		panic(err)
	}

	path = path + string(os.PathSeparator) + "fitgirl-launcher"

	if _, err := os.Stat(path + string(os.PathSeparator) + "database.json"); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)

		if err != nil {
			fmt.Println("error creating parent dirs")
			panic(err)
		}

		file, err := os.Create(path + string(os.PathSeparator) + "database.json")

		fmt.Println("db path: ", path+string(os.PathSeparator)+"database.json")

		if err != nil {
			fmt.Println("error creating database file")
			panic(err)
		}

		defer file.Close()

		db := models.Database{
			Games: []models.Game{},
		}

		dbJson, err := json.Marshal(db)
		if err != nil {
			fmt.Println("error marshalling database")
			panic(err)
		}

		if _, err := file.Write(dbJson); err != nil {
			fmt.Println("error writing initial database")
			panic(err)
		}
	}

	return &DatabaseHandler{
		DatabasePath: path,
	}
}

func (dh *DatabaseHandler) LoadDatabase() (*models.Database, error) {
	file, err := os.Open(dh.DatabasePath + string(os.PathSeparator) + "database.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var db models.Database
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&db); err != nil {
		return nil, err
	}

	return &db, nil
}

func (dh *DatabaseHandler) SaveDatabase(db *models.Database) error {
	file, err := os.Create(dh.DatabasePath + string(os.PathSeparator) + "database.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(db); err != nil {
		return err
	}

	return nil
}

func sanitizePath(path string) string {
	illegalChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	return illegalChars.ReplaceAllString(path, "_")
}

func (dh *DatabaseHandler) AddGameToDatabase(game models.Game) error {
	db, err := dh.LoadDatabase()
	if err != nil {
		return err
	}

	game.InstallPath = sanitizePath(game.InstallPath)

	for _, g := range db.Games {
		if g.Url == game.Url {
			return nil
		}
	}

	db.Games = append(db.Games, game)

	if err := dh.SaveDatabase(db); err != nil {
		return err
	}

	return nil
}

func (dh *DatabaseHandler) RemoveGameFromDatabase(url string) error {
	db, err := dh.LoadDatabase()
	if err != nil {
		return err
	}

	for i, game := range db.Games {
		if game.Url == url {
			db.Games = append(db.Games[:i], db.Games[i+1:]...)
			break
		}
	}

	if err := dh.SaveDatabase(db); err != nil {
		return err
	}

	return nil
}

func (dh *DatabaseHandler) GetGameFromDatabase(url string) (models.Game, error) {
	db, err := dh.LoadDatabase()
	if err != nil {
		return models.Game{}, err
	}

	for _, game := range db.Games {
		if game.Url == url {
			return game, nil
		}
	}

	return models.Game{}, nil
}

func (db *DatabaseHandler) GetGameFromDatabaseByHash(hash string) (models.Game, error) {
	database, err := db.LoadDatabase()
	if err != nil {
		return models.Game{}, err
	}

	for _, game := range database.Games {
		if utils.InfoHashFromMagnet(game.Magnet) == hash {
			return game, nil
		}
	}

	return models.Game{}, nil
}

func (dh *DatabaseHandler) UpdateGameStatusDownloaded(url string) error {
	db, err := dh.LoadDatabase()
	if err != nil {
		return err
	}

	for i, game := range db.Games {
		if game.Url == url {
			db.Games[i].Status = utils.DB_DOWNLOADED
			break
		}
	}

	if err := dh.SaveDatabase(db); err != nil {
		return err
	}

	return nil
}

func (dh *DatabaseHandler) UpdateGameStatusInstalling(url string) error {
	db, err := dh.LoadDatabase()
	if err != nil {
		return err
	}

	for i, game := range db.Games {
		if game.Url == url {
			db.Games[i].Status = utils.DB_INSTALLING
			break
		}
	}

	if err := dh.SaveDatabase(db); err != nil {
		return err
	}

	return nil
}

func (dh *DatabaseHandler) UpdateGameStatusInstalled(url string) error {
	db, err := dh.LoadDatabase()
	if err != nil {
		return err
	}

	for i, game := range db.Games {
		if game.Url == url {
			db.Games[i].Status = utils.DB_INSTALLED
			break
		}
	}

	if err := dh.SaveDatabase(db); err != nil {
		return err
	}

	return nil
}
