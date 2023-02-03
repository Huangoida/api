package Init

import (
	"api/service"
	"time"
)

func InitMetadata() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				service.MapInit()
			}
		}
	}()
	service.MapInit()
}
