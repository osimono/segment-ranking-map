package main

import (
	"github.com/sirupsen/logrus"
	"strava-segemnt-ranking/backend/api"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	server := api.NewServer()
	server.Start()
	/*athlete, err := stravaApi.Athlete(os.Getenv("STRAVA_ATHLETE_ID"))
	if err != nil {
		logrus.Errorf("oops, %s", err.Error())
	}
	logrus.Infof("using athlete: %s", athlete.Name)*/

}
