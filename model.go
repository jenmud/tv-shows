package main

import (
	"encoding/json"
)

// ShowImage is the location to the shows image.
type ShowImage struct {
	URL string `json:"showImage,omitempty"`
}

// Episode contains information about the next episode
// of a particular show.
type Episode struct {
	Channel     string `json:"channel,omitempty"`
	ChannelLogo string `json:"channelLogo,omitempty"`
	Date        string `json:"date,omitempty"`
	HTML        string `json:"html,omitempty"`
	URL         string `json:"url,omitempty"`
}

// Slug is the part of a URL that identifies a page in
// human-readable keywords.
type Slug struct {
	Name string `json:"slug,omitempty"`
}

// TVShow contains all the information related to a tv show.
type TVShow struct {
	Country       string     `json:"country,omitempty"`
	Description   string     `json:"description,omitempty"`
	DRM           bool       `json:"drm,omitempty"`
	EpisodeCount  int        `json:"episodeCount,omitempty"`
	Genre         string     `json:"genre,omitempty"`
	Image         *ShowImage `json:"image,omitempty"`
	Language      string     `json:"language,omitempty"`
	NextEpisode   *Episode   `json:"nextEpisode,omitempty"`
	PrimaryColour string     `json:"primaryColour,omitempty"`
	Seasons       []*Slug    `json:"seasons,omitempty"`
	Slug          string     `json:"slug,omitempty"`
	Title         string     `json:"title,omitempty"`
	TVChannel     string     `json:"tvChannel,omitempty"`
}

// RequestJSON is the base structure of a requests body in JSON format.
type RequestJSON struct {
	// Payload is a slice of tv shows.
	Payload []*TVShow `json:"payload,omitempty"`
}

// ResponseItem contains information about a particular show.
type ResponseItem struct {
	Image string `json:"image,omitempty"`
	Slug  string `json:"slug,omitempty"`
	Title string `json:"title,omitempty"`
}

// ResponseJSON contains a payload of information to be send as a
// JSON response.
type ResponseJSON struct {
	Payload []*ResponseItem `json:"response,omitempty"`
}

// ParseRequestJSON takes a JSON formatted byte array.
// JSON data should contain a `payload` field which maps to a array of
// tv shows.
func ParseRequestJSON(data []byte) ([]*TVShow, error) {
	reqData := &RequestJSON{}

	err := json.Unmarshal(data, reqData)
	if err != nil {
		return nil, err
	}

	return reqData.Payload, nil
}

// FilterTVShowsForDRM filters out all tv shows that do not have a their
// DRM enabled.
func FilterTVShowsForDRM(shows []*TVShow) []*TVShow {
	filtered := []*TVShow{}
	for _, tvShow := range shows {
		if tvShow.DRM {
			filtered = append(filtered, tvShow)
		}
	}
	return filtered
}

// FilterTVShowsWithEpisodes filters out all the tv shows that do
// not have any expisodes.
func FilterTVShowsWithEpisodes(shows []*TVShow) []*TVShow {
	filtered := []*TVShow{}
	for _, tvShow := range shows {
		if tvShow.EpisodeCount > 0 {
			filtered = append(filtered, tvShow)
		}
	}
	return filtered
}

// MakeResponseJSON takes tv shows and returns a response
// in the form of a JSON formatted byte array.
func MakeResponseJSON(shows []*TVShow) ([]byte, error) {
	payload := make([]*ResponseItem, len(shows))

	for index, show := range shows {
		payload[index] = &ResponseItem{
			Image: show.Image.URL,
			Slug:  show.Slug,
			Title: show.Title,
		}
	}

	return json.Marshal(
		&ResponseJSON{
			Payload: payload,
		},
	)
}
