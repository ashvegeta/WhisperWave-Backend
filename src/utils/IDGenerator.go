package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	MIDcounterMu sync.Mutex
	MIDcounter   int64
	UIDcounterMu sync.Mutex
	UIDcounter   int64
	GIDcounterMu sync.Mutex
	GIDcounter   int64
)

func GenerateMessageID(requesterID string) string {
	// increment counter
	MIDcounterMu.Lock()
	MIDcounter++
	MIDcounterMu.Unlock()

	// add timestamp
	timestamp := time.Now().UnixMicro()

	// add MID
	return fmt.Sprintf("M-%s-%d-%d", requesterID, timestamp, UIDcounter)
}

func GenerateUserID() string {
	// increment counter
	UIDcounterMu.Lock()
	UIDcounter++
	UIDcounterMu.Unlock()

	// add timestamp
	timestamp := time.Now().UnixMicro()

	return fmt.Sprintf("U-%d-%d", timestamp, UIDcounter)
}

func GenerateGroupID() string {
	// increment counter
	GIDcounterMu.Lock()
	GIDcounter++
	GIDcounterMu.Unlock()

	// add timestamp
	timestamp := time.Now().UnixMicro()

	return fmt.Sprintf("G-%d-%d", timestamp, GIDcounter)
}
