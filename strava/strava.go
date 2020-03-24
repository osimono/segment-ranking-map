package strava

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Api struct {
	client *http.Client
	base   string
	token  string
}

func NewApi(token string) Api {
	return Api{
		client: http.DefaultClient,
		base:   "https://www.strava.com/api/v3",
		token:  token,
	}
}

type Athlete struct {
	Name string `json:"username"`
}

func (a Api) Athlete(id string) (athlete Athlete, err error) {
	resp, err := a.executeRequest("athletes/"+id, nil)
	if err != nil {
		log.Error(err.Error())
		return athlete, err
	}

	defer resp.Body.Close()
	json.NewDecoder(bufio.NewReader(resp.Body)).Decode(&athlete)
	return
}

type Corner struct {
	Latitude  float64
	Longitude float64
}

func (c Corner) bounds() string {
	return fmt.Sprintf("%f,%f", c.Latitude, c.Longitude)
}

type Segment struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (a Api) Segments(southwestCorner, northEastCorner Corner) (segments []Segment, err error) {
	resp, err := a.executeRequest("segments/explore", []param{
		{name: "bounds", value: fmt.Sprintf("%s,%s", southwestCorner.bounds(), northEastCorner.bounds())},
		{name: "activity_type", value: "riding"},
		{name: "min_cat", value: "0"},
		{name: "max_cat", value: "0"},
	})
	if err != nil {
		log.Error(err.Error())
		return segments, err
	}

	defer resp.Body.Close()
	type stravaSegmentsResponse struct {
		Segments []Segment `json:"segments"`
	}
	r := stravaSegmentsResponse{}
	json.NewDecoder(bufio.NewReader(resp.Body)).Decode(&r)
	return r.Segments, nil
}

type Ranking struct {
	Position    int
	All         int
	EffortCount int
}

func (a Api) SegmentRanking(segmentId int64) (ranking Ranking, err error) {
	resp, err := a.executeRequest(fmt.Sprintf("segments/%d/leaderboard", segmentId), nil)
	if err != nil {
		log.Error(err.Error())
		return ranking, err
	}

	type stravaLeaderboardEntry struct {
		AthleteName string `json:"athlete_name"`
		Rank        int    `json:"rank"`
	}

	type stravaLeaderboard struct {
		Entries    []stravaLeaderboardEntry `json:"entries"`
		EntryCount int                      `json:"entry_count"`
	}

	r := stravaLeaderboard{}
	json.NewDecoder(bufio.NewReader(resp.Body)).Decode(&r)

	ranking.All = r.EntryCount
	for _, entry := range r.Entries {
		if entry.AthleteName == "Osimon O." {
			ranking.Position = entry.Rank
			break
		}
	}
	return ranking, nil
}

type param struct {
	name, value string
}

func (a Api) executeRequest(urlPart string, params []param) (resp *http.Response, err error) {
	url := fmt.Sprintf("%s/%s?access_token=%s", a.base, urlPart, a.token)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	if len(params) > 0 {
		values := request.URL.Query()
		for _, p := range params {
			values.Add(p.name, p.value)
		}
		request.URL.RawQuery = values.Encode()
	}
	request.Header.Add("Accept", "application/json")

	resp, err = a.client.Do(request)
	if err != nil {
		return resp, err
	}
	log.Debugf("requesting url: %s - response-status: %s", request.URL.String(), resp.Status)
	return resp, err
}
