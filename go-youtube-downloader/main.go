package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
)

var ytClient youtube.Client

var (
	destinationPath *string
	links           *string
	linksFromFile   *string
)

func init() {

	destinationPath = flag.String("destination", "./downloads", "destination path")
	links = flag.String("links", "", "video links")
	linksFromFile = flag.String("linksFromFile", "", "path to file with links")
}

func downloadVideo(link string, path string) {

	videoID, err := youtube.ExtractVideoID(link)
	if err != nil {
		log.Fatal(err)
		return
	}

	video, err := ytClient.GetVideo(videoID)
	if err != nil {
		log.Fatal(err)
		return
	}

	formats := video.Formats.WithAudioChannels()
	title := video.Title + " | " + video.Author + ".mp4"

	stream, _, err := ytClient.GetStream(video, &formats[0])
	log.Printf("Downloading: %s", title)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer stream.Close()

	file, err := os.Create(path + "/" + title)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Downloaded: %s", title)

}

func main() {

	var wg sync.WaitGroup

	flag.Parse()

	ytClient = youtube.Client{}

	_ = os.Mkdir(*destinationPath, os.ModePerm)

	if *linksFromFile != "" {
		file, err := os.Open(*linksFromFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var i int = 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				downloadVideo(link, *destinationPath)
				i++
			}(scanner.Text())
		}
		wg.Wait()

		log.Printf("Done - Downloaded %d videos", i)
		os.Exit(0)
	}

	if *links == "" {
		log.Fatal("No links provided")
		return
	}

	parsedLinks := strings.Split(strings.ReplaceAll(*links, " ", ""), ",")

	for _, link := range parsedLinks {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			downloadVideo(link, *destinationPath)
		}(link)
	}

	wg.Wait()

	log.Printf("Done - Downloaded %d videos", len(parsedLinks))
	os.Exit(0)

}
