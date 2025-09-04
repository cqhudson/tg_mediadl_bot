package main

type YouTubeInformation struct {
	// The regex used to find YouTube links
	Regex string
}

func newYouTubeInformation(regex string) *YouTubeInformation {
	ytInfo := YouTubeInformation{Regex: regex}
	return &ytInfo
}
