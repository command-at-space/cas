/* */

package chat

import (
	"log"
	"os"
)

func createCustomChatLogFile(f string) {
	chatLog, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("ERROR opening Chat log file %s\n", err)
	}
	cLog = log.New(chatLog, "CHAT:\t", log.Ldate|log.Ltime)
}
