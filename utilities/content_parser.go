package utilities

import (
	"regexp"
	"strings"
)

var (
	hashtagRegex = regexp.MustCompile(`#(\w+)`)
	mentionRegex = regexp.MustCompile(`@(\w+)`)
)

func ParseHashtags(content string) []string {
	matches := hashtagRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var tags []string
	for _, match := range matches {
		tag := strings.ToLower(match[1])
		if !seen[tag] {
			seen[tag] = true
			tags = append(tags, tag)
		}
	}
	return tags
}

func ParseMentions(content string) []string {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var usernames []string
	for _, match := range matches {
		username := strings.ToLower(match[1])
		if !seen[username] {
			seen[username] = true
			usernames = append(usernames, username)
		}
	}
	return usernames
}
