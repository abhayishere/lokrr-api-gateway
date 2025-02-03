package pkg

import (
	"log"
)

func LogInfo(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func LogError(err error) {
	log.Printf("[ERROR] %v\n", err)
}
