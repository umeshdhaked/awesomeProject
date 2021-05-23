package api

import (
	"fmt"
	"github.com/umeshdhaked/awesomeProject/packages/pubsub"
	"log"
	"net/http"
	"strings"
)

func HandleRequest() {
	http.HandleFunc("/publish/", handlePublish)
	log.Fatal( http.ListenAndServe(":8080",nil) )
}

func handlePublish( w http.ResponseWriter, r *http.Request ) {

	data := strings.TrimPrefix(r.URL.Path, "/publish/")

	pubsub.Publish("topic1", "apiData:"+data)

	fmt.Fprintln(w,"response:"+data)
}