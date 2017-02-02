package main

import (
	"flag"
	"golang.org/x/oauth2"
	"log"
	"github.com/google/go-github/github"
	"os"
	"path"
	"io"
	"net/http"
)

var (
	token = flag.String("token", "", "GitHub Personal Access Token")
	owner = flag.String("owner", "", "GitHub Owner")
	repo = flag.String("repo", "", "GitHub Repo")
	output = flag.String("output", ".", "The output directory")
)

func main()  {
	flag.Parse()
	// init github
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	client := github.NewClient(oauth2.NewClient(oauth2.NoContext, ts))
	// get latest release
	release, _, err := client.Repositories.GetLatestRelease(*owner, *repo)
	if err != nil {
		log.Fatalln(err)
	}
	// download all assets
	for _, asset := range release.Assets {
		rc, redirect, err := client.Repositories.DownloadReleaseAsset(*owner, *repo, *asset.ID)
		if err != nil {
			log.Fatalln(err)
		}
		// handle redirect
		if redirect != "" {
			// get file
			resp, err := http.DefaultClient.Get(redirect)
			if err != nil {
				log.Fatalln(err)
			}
			writeFile(resp.Body, *asset.Name)
		} else {
			writeFile(rc, *asset.Name)
		}

	}
}

func writeFile(r io.ReadCloser, name string)  {
	fPath := path.Join(*output, name);
	// create file
	file, err := os.Create(fPath)
	if err != nil {
		log.Fatalln(err)
	}
	n, err := io.Copy(file, r)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Created file %s (%d Bytes)", fPath, n)
	r.Close()
}