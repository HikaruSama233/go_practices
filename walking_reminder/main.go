package main
import (
	"gopkg.in/toast.v1"
	"log"
	"github.com/carlescere/scheduler"
	"runtime"
)

func main() {
	job := func() {
		notification := toast.Notification{
			AppID: "Microsoft.Windows.Shell.RunDialog",
			Title: "Walking Reminder",
			Message: "Time to Walk. >.<",
			Actions: []toast.Action{
				{"protocal", "OK", ""},
				{"protocal", "Will do!", ""},
			},
		}
		err := notification.Push()
		if err != nil {
			log.Fatalln(err)
		}
	}
	scheduler.Every(45).Minutes().Run(job)
	runtime.Goexit()
}
