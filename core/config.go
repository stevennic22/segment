package core

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"time"

	"github.com/joho/godotenv"
)

// Storage for
//
//	Application FilePath
//	Signature Secret
type Settings struct {
	BaseFilePath string
	Secret       []byte
	ServerPort   string
	DotEnv
}

var Config Settings = Settings{
	BaseFilePath: "/usr/src/app/content/",
}

func (s *Settings) FilePath(dirStr string) string {
	return path.Join(s.BaseFilePath, dirStr)
}

func (s *Settings) FetchServerPort() {
	sp := s.RetrieveValue("listen_port")
	if sp != ":" {
		s.ServerPort = ":" + sp
	} else {
		s.ServerPort = ":8080"
	}
}

type DotEnv struct {
	Location   string
	LastLoaded struct {
		Key       string
		Value     string
		Retrieved time.Time
	}
}

func (d *DotEnv) loadFile() {
	loadErr := godotenv.Load(d.Location)

	CheckError(loadErr, true)
}

func (d *DotEnv) InitEnvFile() {
	firstAttempt := godotenv.Load(d.Location)
	if firstAttempt != nil {
		log.Printf("Error loading .env file (first attempt): %s", firstAttempt)
		if d.Location != "./content/static/.env" {
			log.Println("Updating .env file and settings to use local path: ./content/static/")

			d.Location = "./content/static/.env"
			Config.BaseFilePath = "./content"

			secondAttempt := godotenv.Load(d.Location)
			if secondAttempt != nil {
				log.Fatalf("Error loading .env file (second attempt): %s", secondAttempt)
			}
		}
	}
}

func (d *DotEnv) fileModMoreRecent() bool {
	file, err := os.Stat(d.Location)

	if err != nil {
		log.Fatalf("Error loading file stats: %s", err)
	}

	modifiedtime := file.ModTime()

	return (modifiedtime.After(d.LastLoaded.Retrieved))
}

func (d *DotEnv) RetrieveValue(key string) string {

	// Check if:
	//  We've already retrieved that value
	//  The .env file has been updated
	if d.LastLoaded.Key == key && !d.fileModMoreRecent() {
		return (d.LastLoaded.Value)
	}

	d.loadFile()

	d.LastLoaded.Key = key
	d.LastLoaded.Value = os.Getenv(key)
	d.LastLoaded.Retrieved = time.Now()

	return d.LastLoaded.Value
}

// Event represents the structure of events to be broadcast
type Event struct {
	Type        string                 `json:"type"`
	Description string                 `json:"event"`
	Data        map[string]interface{} `json:"raw_data"`
	MessageID   string                 `json:"messageId"`
	Timestamp   time.Time              `json:"timestamp"`
}

// In-Memory list of events
var Events []Event

func (e *Event) UnmarshalJSON(data []byte) error {
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	e.Data = rawMap

	type TempRequestData Event
	temp := (*TempRequestData)(e)
	return json.Unmarshal(data, temp)
}

func FindEventIndexByMessageID(messageID string) int {
	for i, event := range Events {
		if event.MessageID == messageID {
			return i
		}
	}
	return -1 // Not found
}
