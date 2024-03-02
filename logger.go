package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(stdio bool) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	if !stdio {
		f, err := os.OpenFile(GLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err)
		}
		currentTime := time.Now()
		GLogPath = currentTime.Format("2006-01-02_15-04-05") + "_file.txt"
		GLogPath = filepath.Join(LogFolder, GLogPath)
		createFolderIfNotExists(LogFolder)
		fmt.Println("path is ", GLogPath)
		log.Logger = log.Output(f)
	} else{
        log.Info().Msgf("will log to stdio")
    }
    
}
