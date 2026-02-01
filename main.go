package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
)

const header = "## 📕 Latest Blog Posts"

func main() {
	rssURL := os.Getenv("RSS_URL")
	if rssURL == "" {
		log.Fatal("RSS_URL environment variable is not set")
	}

	// 1. Fetch RSS Feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		log.Fatalf("Failed to parse RSS feed: %v", err)
	}

	// 2. Build the new blog post list
	var sb strings.Builder
	sb.WriteString("<ul>\n")
	count := 0
	for _, item := range feed.Items {
		if count >= 10 {
			break
		}
		sb.WriteString(fmt.Sprintf("<li><a href='%s' target='_blank'>%s</a></li>\n", item.Link, item.Title))
		count++
	}
	sb.WriteString("</ul>\n")
	newContent := sb.String()

	// 3. Read existing README.md
	readmeData, err := os.ReadFile("README.md")
	if err != nil {
		// If README doesn't exist, just start with empty string
		if os.IsNotExist(err) {
			readmeData = []byte("")
		} else {
			log.Fatalf("Failed to read README.md: %v", err)
		}
	}
	readmeContent := string(readmeData)

	// 4. Replace or Append
	var finalContent string
	if strings.Contains(readmeContent, header) {
		// Split and keep the part before the header
		parts := strings.SplitN(readmeContent, header, 2)
		// Reconstruct: content before header + header + newline + new content
		// We avoid keeping old content after the header by discarding parts[1]
		finalContent = parts[0] + header + "\n\n" + newContent
	} else {
		// Append to the end
		if len(readmeContent) > 0 && !strings.HasSuffix(readmeContent, "\n") {
			readmeContent += "\n"
		}
		finalContent = readmeContent + "\n" + header + "\n\n" + newContent
	}

	// 5. Write back to README.md
	err = os.WriteFile("README.md", []byte(finalContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write README.md: %v", err)
	}

	fmt.Println("README.md updated successfully")
}
