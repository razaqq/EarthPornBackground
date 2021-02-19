package main

import (
	"context"
	"fmt"
	"github.com/reujab/wallpaper"
	"github.com/rubenfonseca/fastimage"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	_ "image/jpeg"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var (
	user32           = syscall.NewLazyDLL("User32.dll")
	getSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func GetSystemMetrics(nIndex int) uint32 {
	index := uintptr(nIndex)
	ret, _, _ := getSystemMetrics.Call(index)
	return uint32(ret)
}

func main() {
	screenWidth := GetSystemMetrics(SM_CXSCREEN)
	screenHeight := GetSystemMetrics(SM_CYSCREEN)
	screenRatio := float32(screenWidth) / float32(screenHeight)

	// get top posts from today
	client, _ := reddit.NewReadonlyClient()
	posts, _, err := client.Subreddit.TopPosts(
		context.Background(),
		"EarthPorn",
		&reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: 100,
			},
			Time: "today",
		})

	if err != nil {
		log.Fatalln(err.Error())
	}

	for i, p := range posts {
		fmt.Println(i)

		_, size, err := fastimage.DetectImageType(p.URL)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if size == nil {
			fmt.Println("Not an image.")
			continue
		}

		// make sure we got minimum resolution
		if size.Height < screenHeight || size.Width < screenWidth {
			fmt.Printf("Minimum resolution not met: %d:%d\n", size.Width, size.Height)
			continue
		}

		// check if the ratio is too far off
		imageRatio := float32(size.Width) / float32(size.Height)
		if imageRatio < 0.9*screenRatio || imageRatio > 1.1*screenRatio {
			fmt.Printf("Bad ratio for image: %f\n", imageRatio)
			continue
		}

		// set as background
		fmt.Println(p.URL)
		if err := wallpaper.SetFromURL(p.URL); err != nil {
			fmt.Println(err.Error())
			continue
		}

		// force desktop refresh
		exec.Command("RUNDLL32.EXE USER32.DLL, UpdatePerUserSystemParameters 1, True")
		os.Exit(0)
	}
}
