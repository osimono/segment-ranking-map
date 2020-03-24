package api

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func fromRequest(r *http.Request, v interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&v)
	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(decoder.Buffered())
		jsonBody := buf.String()
		logrus.Errorf("%s - failed to decode json: '%s'", err.Error(), jsonBody)
		return err
	}
	return
}

func respondDecodingError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("failed to decode json body"))
}

func respondInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("oops, something went wrong"))
}
