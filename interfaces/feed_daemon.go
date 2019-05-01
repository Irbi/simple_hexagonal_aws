package interfaces

import (
	"fmt"
	"test_aws/usecases"
	"time"
)

type FeedInteractor interface {
	SyncItems() ([]usecases.File, error)
}

type FeedHandler struct {
	FilesInteractor FeedInteractor
}

func (handler FeedHandler) Feed(tick time.Duration) {

	//items, _ := handler.FilesInteractor.SyncItems()
	//for _, item := range items {
	//	fmt.Printf("Sync file v.%s: %s", item.Version, item.Name)
	//}
	ticker := time.NewTicker(tick * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <- ticker.C:
				fmt.Println("Ticker started")
				items, _ := handler.FilesInteractor.SyncItems()
				for _, item := range items {
					fmt.Printf("Sync file v.%d: %s\n\r", item.Version, item.Name)
				}
			case <- quit:
				fmt.Println("Ticker stopped")
				ticker.Stop()
				return
			}
		}
	}()
}
