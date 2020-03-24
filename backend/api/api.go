package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strava-segemnt-ranking/backend/strava"
	"time"
)

type Server struct {
	router *mux.Router
}

func NewServer() Server {
	return Server{
		router: mux.NewRouter(),
	}
}

func (s Server) Start() {
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("all is not lost!"))
	})
	apiRouter.HandleFunc("/rankings", SegmentHandler).Methods(http.MethodPost)

	srv := &http.Server{
		Handler: s.router,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		log.Infof("route: %s", t)
		return nil
	})

	log.Infof("starting server on address: %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func SegmentHandler(w http.ResponseWriter, r *http.Request) {
	type corner struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	type bounds struct {
		SouthWestCorner corner `json:"sw"`
		NorthEastCorner corner `json:"ne"`
	}

	var b bounds
	err := fromRequest(r, &b)
	if err != nil {
		respondDecodingError(w)
		return
	}

	log.Infof("received ranking request for bounds: %#v", b)

	token := os.Getenv("STRAVA_TOKEN")
	stravaApi := strava.NewApi(token)

	swCorner := strava.Corner{
		Latitude:  b.SouthWestCorner.Lat,
		Longitude: b.SouthWestCorner.Lng,
	}

	neCorner := strava.Corner{
		Latitude:  b.NorthEastCorner.Lat,
		Longitude: b.NorthEastCorner.Lng,
	}

	segments, err := stravaApi.Segments(swCorner, neCorner)
	if err != nil {
		log.Error(err.Error())
	}

	type segmentRanking struct {
		Start       corner `json:"start"`
		SegmentName string `json:"segmentName"`
		All         int    `json:"all"`
		Position    int    `json:"position"`
	}
	var rankings []segmentRanking
	for _, segment := range segments {
		ranking, err := stravaApi.SegmentRanking(segment.Id)
		if err != nil {
			log.Errorf("failed to read segment ranking on segment %d - %s", segment.Id, err.Error())
			continue
		}
		log.Debugf("%d/%d on Segment: %s (%d)", ranking.Position, ranking.All, segment.Name, segment.Id)
		rankings = append(rankings, segmentRanking{
			Start: corner{
				Lat: segment.Start[0],
				Lng: segment.Start[1],
			},
			SegmentName: segment.Name,
			All:         ranking.All,
			Position:    ranking.Position,
		})
	}
	err = json.NewEncoder(w).Encode(rankings)
	if err != nil {
		respondInternalServerError(w)
	}
}
