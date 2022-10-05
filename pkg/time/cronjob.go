package time

import (
	"fmt"
	"time"

	"github.com/danyouknowme/awayfromus/pkg/api"
	"github.com/go-co-op/gocron"
)

func RunCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Days().Do(func() {
		err := api.UpdateUserResourceExpiredDate()
		if err != nil {
			fmt.Println("error to update resource dayleft: ", err.Error())
		}
		fmt.Println("update dayleft")
	})

	s.StartBlocking()
}
