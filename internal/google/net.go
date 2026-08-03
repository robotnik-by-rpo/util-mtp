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
func getToken() *GoogleMetaData {
	return &GoogleMetaData{
		accessToken: os.Getenv("GOOGLE_TOKEN"),
		folder_id:   os.Getenv("GOOGLE_PATH_ID"),
	}
}

func validToken(tokens *GoogleMetaData) bool{
	return tokens.accessToken == "" && tokens.folder_id
}

// Function for sending to Google
func SendToGoogleDisk(ctx context.Context, filepath string) error {
	
	tokens := getToken()
	if validToken(tokens) {
		log.Printf("Empty Google's tokens")
		return fmt.Errorf("Empty tokens for Google")
	}

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
		return fmt.Errorf("Error getting information about file: %v", err)
	}

	driveFile := &drive.File{Name: fileInfo.Name(), Parents: []string{tokens.folder_id}}
	call := driveService.Files.Create(driveFile)
	call = call.Media(file, googleapi.ChunkSize(googleapi.DefaultChunkSize))

	res, err := call.Do()
	if err != nil{
		log.Printf("Error sendig into google disk: %v", err)
		return fmt.Errorf("Error sendig into google disk: %v", err)
	}

	fmt.Printf("File was uplouded successful %s", res.Id)
	return nil
}
