package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	counterMu sync.Mutex
	counter   uint64
)

func GenerateID(requesterID string) string {
	// increment counter
	counterMu.Lock()
	counter++
	counterMu.Unlock()

	// add timestamp
	timestamp := time.Now().UnixMilli()

	// add UID
	ID := fmt.Sprintf("%s-%d-%d", requesterID, timestamp, counter)

	return ID
}

// func GenerateServerID() {

// }