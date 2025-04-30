package scraper

import (
	"net/url"
	"regexp"
	"slices"
	"strings"
)

func isValidURL(rawURL string) (string, bool) {
	// Verify domain name
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}

	validDomains := []string{
		"www.instagram.com",
		"instagram.com",
	}

	domainValid := slices.Contains(validDomains, parsedURL.Host)
	if !domainValid {
		return "", false
	}

	// Check path starts with /reel/ or /p/ or /username/ or /reels/
	path := strings.Trim(parsedURL.Path, "/")
	pathParts := strings.Split(path, "/")

	var reelID string

	// Validate and extract from different URL formats
	switch len(pathParts) {
	case 2:
		// Format: /reel/{reelID}/ or /reels/{reelID}/ or /p/{reelID}/
		if !slices.Contains([]string{"reel", "reels", "p"}, pathParts[0]) {
			return "", false
		}
		reelID = pathParts[1]
	case 3:
		// Format: /{username}/reel/{reelID}/
		if pathParts[0] == "" {
			return "", false
		}
		if !validUsername(pathParts[0]) && pathParts[1] != "reel" {
			return "", false
		}
		reelID = pathParts[2]
	default:
		return "", false
	}

	// Validate Reel ID format
	if !validReelID(reelID) {
		return "", false
	}
	return reelID, true
}

func validUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[A-Za-z0-9._]{1,30}$`)
	if !usernameRegex.MatchString(username) {
		return false
	}
	return true
}

func validReelID(reelID string) bool {
	reelIDRegex := regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)
	if !reelIDRegex.MatchString(reelID) {
		return false
	}
	return true
}
