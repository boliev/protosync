package source

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v55/github"
)

// Github source struct implements domain.Source
type Github struct {
	Name     string
	User     string
	Repo     string
	Path     string
	Ref      string
	SyncPath string

	client *github.Client
}

// NewGithub constructor
func NewGithub(name string, user string, repo string, path string, ref string, syncPath string) *Github {
	return &Github{
		Name:     name,
		User:     user,
		Repo:     repo,
		Path:     path,
		Ref:      ref,
		SyncPath: syncPath,
		client:   github.NewClient(nil),
	}
}

// SyncProtos syncs Github proto files
func (source *Github) SyncProtos() error {
	opts := source.getOpts()
	fileContent, dirContent, _, err := source.client.Repositories.GetContents(
		context.Background(), source.User, source.Repo, source.Path, opts,
	)
	if err != nil {
		return err
	}
	if fileContent != nil {
		err = source.syncFile(*fileContent.DownloadURL, source.SyncPath)
	}

	if len(dirContent) > 0 {
		hasher := sha1.New()
		hasher.Write([]byte("main"))
		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		tree, resp, err := source.client.Git.GetTree(context.Background(), source.User, source.Repo, sha, true)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", tree)
		fmt.Printf("%s\n", resp)
		return fmt.Errorf("sorry directory content does not supported yet")
	}

	if err != nil {
		return err
	}

	// fmt.Printf("%s\n", source.Name)
	// fmt.Printf("-- %s\n", fileContent)
	// fmt.Printf("-- %s\n", repoContent)
	// fmt.Printf("-- %s\n", dirContent)
	// fmt.Printf("%s\n", "###############")
	return nil
}

func (source *Github) getOpts() *github.RepositoryContentGetOptions {
	if source.Ref != "" {
		return &github.RepositoryContentGetOptions{
			Ref: source.Ref,
		}
	}

	return nil
}

func (source *Github) syncFile(downloadUrl string, dst string) error {
	log.Printf("source %s: syncing file %s to %s\n", source.Name, downloadUrl, dst)
	err := source.copyFile(downloadUrl, dst)
	if err != nil {
		return err
	}
	log.Printf("- done\n")

	return nil
}

func (source *Github) copyFile(src string, dst string) error {
	resp, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("error while making GET request: %w", err)
	}
	defer resp.Body.Close()

	dst = source.prepareDst(src, dst)
	if err = source.createDir(dst); err != nil {
		return err
	}

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

func (source *Github) prepareDst(src string, dst string) string {
	isdir := source.isPathADir(dst)
	if isdir {
		if !strings.HasSuffix(dst, "/") {
			dst = dst + "/"
		}
		dst = dst + filepath.Base(src)
	}

	return dst
}

func (source *Github) createDir(dst string) error {
	dir := filepath.Dir(dst)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("no directory %s found. Creating.\n", dir)
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (source *Github) isPathADir(dst string) bool {
	return !strings.HasSuffix(dst, ".proto")
}
