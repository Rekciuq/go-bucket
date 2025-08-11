package videohandler

import (
	"fmt"
	"net/http"
)

func VideoHandler() http.Handler {
	router := http.NewServeMux()
	fileServerRoot := http.Dir("./uploads/videos")
	fileServerHandler := http.FileServer(fileServerRoot)
	fmt.Println("File server opened!")
	router.Handle("/", fileServerHandler)

	return router
}
