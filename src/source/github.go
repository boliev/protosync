package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v55/github"
)

// Github source struct implements domain.Source
type Github struct {
	Name     string
	User     string
	Repo     string
	Path     string
	Branch   string
	Tag      string
	SyncPath string

	client *github.Client
}

// NewGithub constructor
func NewGithub(name string, user string, repo string, path string, branch string, tag string, syncPath string) *Github {
	return &Github{
		Name:     name,
		User:     user,
		Repo:     repo,
		Path:     path,
		Branch:   branch,
		Tag:      tag,
		SyncPath: syncPath,
		client:   github.NewClient(nil),
	}
}

// SyncProtos syncs Github proto files
func (source *Github) SyncProtos() error {
	return nil
}

// DownloadProto download a proto
func (source *Github) DownloadProto(url string) error {
	fileContent, _, _, err := source.client.Repositories.GetContents(
		context.Background(), "coolnotes", "users", "protos/user.proto", nil,
	)
	if err != nil {
		return err
	}
	// fmt.Printf("File: %s\n", *fileContent.DownloadURL)
	err = source.copyFile(*fileContent.DownloadURL, "proto/user.proto")
	if err != nil {
		return err
	}

	return nil
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

	return nil
}
