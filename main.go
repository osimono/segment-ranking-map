package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"strava-segemnt-ranking/strava"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	token := os.Getenv("STRAVA_TOKEN")
	stravaApi := strava.NewApi(token)

	athlete, err := stravaApi.Athlete(os.Getenv("STRAVA_ATHLETE_ID"))
	if err != nil {
		logrus.Errorf("oops, %s", err.Error())
	}
	logrus.Infof("using athlete: %s", athlete.Name)

	swCorner := strava.Corner{
		Latitude:  48.454595,
		Longitude: 10.050606,
	}

	neCorner := strava.Corner{
		Latitude:  48.486576,
		Longitude: 10.098846,
	}

	segments, err := stravaApi.Segments(swCorner, neCorner)
	if err != nil {
		logrus.Error(err.Error())
	}

	for _, segment := range segments {
		ranking, err := stravaApi.SegmentRanking(segment.Id)
		if err != nil {
			logrus.Errorf("failed to read segment ranking on segment %d - %s", segment.Id, err.Error())
			continue
		}
		logrus.Infof("%d/%d on Segment: %s (%d)", ranking.Position, ranking.All, segment.Name, segment.Id)
	}
	logrus.Infof("found segments: %s", segments)
}
