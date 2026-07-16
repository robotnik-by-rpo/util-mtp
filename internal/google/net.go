package google

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"log"
	"os"
)

type GoogleMetaData struct {
	accessToken string
	folder_id   string
}

// Function for getting for Google
func getTokken() *GoogleMetaData {
	return &GoogleMetaData{
		accessToken: os.Getenv("GOOGLE_TOKEN"),
		folder_id:   os.Getenv("GOOGLE_PATH_ID"),
	}
}

// Function for sending to Google
func SendToGoogleDisk(ctx *context.Context, filepath string) error {
	
	tokens := getTokken()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokens.accessToken})
	driveService, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil{
		log.Printf("Error initing google service: %v", err)
		return fmt.Errorf("Error initing google service: %v", err)
	}

	file, err := os.Open(filepath)
	if err != nil{
		log.Printf("Error opening file: %v", err)
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil{
		log.Printf("Error getting inforamation about file: %v", err)
		return fmt.Errorf("Error getting inforamation about file: %v", err)
	}

	driveFile := &drive.File{Name: file}
	call := client.Files.Create(driveFile)
	call = call.Media(file, googleapi.ChunkSize(0))

	res, err := call.Do()
	if err != nil{
		log.Printf("Error sendig into google disk: %v", err)
		return fmt.Errorf("Error sendig into google disk: %v", err)
	}

	fmt.Printf("File was uplouded successful %s", res.Id)

}
