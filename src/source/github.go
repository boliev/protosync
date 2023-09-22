package source

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/boliev/protosync/src/domain"
	"github.com/google/go-github/v55/github"
)

// Github source struct implements domain.Source
type Github struct {
}

// NewGithub constructor
func NewGithub() domain.Source {
	return &Github{}
}

// GetAllProtos returns all protos
func (source *Github) GetAllProtos() ([]domain.Proto, error) {
	return []domain.Proto{
		{URL: "https://github.com/coolnotes/users/blob/main/protos/user.proto"},
	}, nil
}

// DownloadProto download a proto
func (source *Github) DownloadProto(url string) (string, error) {
	client := github.NewClient(nil)
	// fmt.Println("\n\n")

	fileContent, _, _, err := client.Repositories.GetContents(
		context.Background(), "coolnotes", "users", "protos/user.proto", nil,
	)
	if err != nil {
		// fmt.Println(err)
		return "", nil
	}
	// fmt.Printf("File: %s\n", *fileContent.DownloadURL)
	err = source.copyFile(*fileContent.DownloadURL, "proto/user.proto")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Directory: %#v\n", directoryContent)
	// fmt.Printf("Response: %#v\n", resp)

	return "", nil
}

func (source *Github) copyFile(src string, dst string) error {
	resp, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("error while making GET request: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error while creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error while copying data: %w", err)
	}
	fmt.Printf("Copied %s to %s\n", src, dst)

	return nil
}
