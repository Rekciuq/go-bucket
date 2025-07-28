package uploadhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeResponse(urlId string, w http.ResponseWriter, filePath string) {
	w.WriteHeader(http.StatusCreated)
	res := postResponse{FileURL: fmt.Sprintf("%s/%s", filePath, urlId)}
	json.NewEncoder(w).Encode(res)
}
