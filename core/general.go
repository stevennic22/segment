package core

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func CheckError(err error, hardStop bool) bool {
	if err != nil {
		if hardStop {
			log.Fatalf("Fatal error within application: %s", err.Error())
		}
		log.Printf("Recoverable error within application: %s", err.Error())
		return (true)
	}

	return (false)
}

func DelayDuration(step int, randomizationFactor float64, maximumDelayTime time.Duration) time.Duration {
	// Step should be at least 1
	if step == 0 {
		step += 1
	}

	// Add randomization factor to step
	// Divide by 1.5 for funsies
	// Could also multiple by 1/3
	// I'm just a code comment
	// Do what you want
	power := (randomizationFactor + float64(step)) / 1.5

	// Turn 2^power into a time duration in seconds
	calculatedTime := time.Duration(math.Pow(2, power)) * time.Second

	// Return calculated time if shorter than maximum delay time
	if calculatedTime < maximumDelayTime {
		return calculatedTime
	} else {
		return maximumDelayTime
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

var CurrentEnv = DotEnv{}

func (d DotEnv) fileModMoreRecent() bool {
	file, err := os.Stat(d.Location)

	if err != nil {
		log.Fatalf("Error loading file stats: %s", err)
	}

	modifiedtime := file.ModTime()

	return (modifiedtime.After(d.LastLoaded.Retrieved))
}

func (d DotEnv) RetrieveValue(key string) string {

	if d.LastLoaded.Key == key && d.fileModMoreRecent() {
		return (d.LastLoaded.Value)
	}

	err := godotenv.Load(d.Location)

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	d.LastLoaded.Key = key
	d.LastLoaded.Value = os.Getenv(key)
	d.LastLoaded.Retrieved = time.Now()

	return d.LastLoaded.Value
}
