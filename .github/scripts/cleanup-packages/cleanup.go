package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type (
	PackageVersion struct {
		Id             int       `json:"id"`
		Name           string    `json:"name"`
		Url            string    `json:"url"`
		PackageHtmlUrl string    `json:"package_html_url"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Description    string    `json:"description"`
		HtmlUrl        string    `json:"html_url"`
		Metadata       struct {
			PackageType string `json:"package_type"`
			Container   struct {
				Tags []string `json:"tags"`
			} `json:"container"`
		} `json:"metadata"`
	}
)

func main() {
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		panic("missing TOKEN")
	}
	packagesStr, ok := os.LookupEnv("PACKAGES")
	if !ok {
		panic("missing PACKAGES")
	}
	packages := strings.Split(packagesStr, ",")
	ticket, ok := os.LookupEnv("TICKET")
	if !ok {
		panic("missing TICKET")
	}

	for _, pkg := range packages {
		processPackage(pkg, token, ticket)
	}
}

func processPackage(pkg, token, ticket string) {
	client := http.DefaultClient
	pkg = url.PathEscape(pkg)
	reqUrl := fmt.Sprintf("https://api.github.com/orgs/azarc-io/packages/container/%s/versions?state=active&per_page=100", pkg)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		panic(fmt.Errorf("invalid request: %s", err))
	}
	req.Header.Set("authorization", "Bearer "+token)

	rsp, err := client.Do(req)
	if err != nil {
		panic(fmt.Errorf("request failed: %w", err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("could not close reader: %s", err.Error()))
		}
	}(rsp.Body)

	if rsp.StatusCode == 404 {
		fmt.Println("package not found, skipping: " + reqUrl)
		return
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		panic(fmt.Errorf("could not read response: %w", err))
	}

	var versions []*PackageVersion
	if err := json.Unmarshal(body, &versions); err != nil {
		fmt.Printf("error when handling response from: %s\n", reqUrl)
		fmt.Printf("error could not unmarshal: %s\n", string(body))
		panic(fmt.Errorf("could not unmarshal response: %w", err))
	}

	var toDelete []string
	for _, version := range versions {
		if len(version.Metadata.Container.Tags) > 1 {
			// TODO check if tags actually contains a valid semver
			continue
		}
		for _, tag := range version.Metadata.Container.Tags {
			if strings.HasPrefix(strings.ToLower(tag), strings.ToLower(ticket)) {
				toDelete = append(toDelete, fmt.Sprintf("%d", version.Id))
			}
		}
	}

	for _, vid := range toDelete {
		reqUrl = fmt.Sprintf("https://api.github.com/orgs/azarc-io/packages/container/%s/versions/%s", pkg, vid)
		req, err := http.NewRequest("DELETE", reqUrl, nil)
		if err != nil {
			panic(fmt.Errorf("invalid request: %s", err))
		}
		req.Header.Set("authorization", "Bearer "+token)
		rsp, err := client.Do(req)
		if err != nil {
			panic(fmt.Errorf("delete request failed: %w", err))
		}

		if rsp.StatusCode != 204 {
			fmt.Println(fmt.Sprintf("delete was not succesful for package: %s version: %s", pkg, vid))
		} else {
			fmt.Println(fmt.Sprintf("deleted: %v in package: %s", vid, pkg))
		}
	}
}
