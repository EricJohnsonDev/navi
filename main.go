package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	var links []string
	// Read conf file for which file/directory to search on
	conf, err := ioutil.ReadFile("navi.conf")
	if err != nil {
		log.Fatal(err)
	}
	files := strings.Fields(strings.TrimSpace(string(conf)))
	fmt.Printf("Searching files: %s\n", files)

	// https://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go
	links = append(links, searchForLinks(files, []string{})...)

	// If there are no links, skip the rest for now
	if links != nil {
		fmt.Println("\nTotal links: ", links)
		links = trimLinks(links)

		// Execute GET on each link
		fmt.Print("\n---\n") // some silly formatting

		for _, link := range links {
			resp, _ := http.Get(link)
			fmt.Printf("%s - %s\n", link, resp.Status)
			resp.Body.Close()
		}
	}
}

// Trim links from regex to remove html and quotes
func trimLinks(links []string) []string {
	var trimmed []string

	// range on slices returns (index, value)
	for _, link := range links {
		link = strings.TrimPrefix(link, "<a href=\"")
		link = strings.TrimSuffix(link, "\"")
		trimmed = append(trimmed, link)
	}

	return trimmed

}

// Search an html file(s) for href links
func searchForLinks(files []string, matches []string) []string {
	fmt.Printf("searchForLinks(%s, %s)\n", files, matches)
	// var matches []string

	// Parse each file individually for links
	for _, filename := range files {
		fmt.Println("Filename: ", filename)
		content, err := ioutil.ReadFile(filename)

		if err != nil {
			// Couldn't read the file - let's see if it's a directory
			fmt.Println("(Directory) ", filename)
			dir, err := ioutil.ReadDir(filename)

			if err != nil {
				log.Fatal(err)
			}

			// Yep, it's a directory all right
			for _, dirFile := range dir {
				filepath := filename + "/" + dirFile.Name()
				matches = append(searchForLinks([]string{filepath}, matches), matches...)
			}
		}

		strContent := string(content)
		// Remove HTML comments
		reComment := regexp.MustCompile(`<!--.*-->`)
		strContent = reComment.ReplaceAllString(strContent, "")

		// Search for links
		re := regexp.MustCompile(`<a\shref=\".*\"`)

		// Add each link to our return variable individually
		for _, match := range re.FindAllString(strContent, -1) {
			fmt.Println("Found an html link: ", match)
			fmt.Println("Matches before appending: ", matches)
			matches = append(matches, match)
			fmt.Println("Matches after appending: ", matches)
		}
	}

	fmt.Println("Matches after filename loop: ", matches)

	if matches == nil {
		log.Println("No links found!")
		return nil
	} else {
		fmt.Printf("Found %d total links\n", len(matches))
		return matches
	}
}
