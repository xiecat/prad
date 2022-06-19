package interrupt

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

// HandleInterrupt handles interrupt signal using context.
// When the interrupt signal is received for the first time, the cancelFunc is called.
// When the interrupt signal is received for the second time, exit the program directly.
func HandleInterrupt(cancelFunc context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer close(signalChan)

	count := 0
	for {
		s, ok := <-signalChan
		if ok {
			if count == 0 {
				log.Infoln("Got signal 1st time:", s)
				count += 1
				cancelFunc()
			} else {
				log.Infoln("Got signal 2nd time:", s)
				os.Exit(1)
			}
		} else {
			return
		}
	}
}
