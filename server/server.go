package server

import (
	"encoding/json"
	"fmt"
	_ "github.com/mimminou/UrlCleaner-ByFood/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @Description	Process URL
type RequestStruct struct {
	// @Property		url string true "URL to process"
	Url string `json:"url"`
	// @Property		operation string true "Operation to perform"
	// @Enum			canonical, redirection, all
	Operation string `json:"operation"`
}

type ResponseStruct struct {
	ProcessedUrl string `json:"processed_url"`
}

// @Description	ErrMessage
// @Property		msg string true "Error message"
type ErrMessage struct {
	Msg string `json:"msg"`
}

// routes requests that have ID based on HTTP
//
//	@Summary		Process URL
//	@Description	Processes URLs depending on the requested operation
//	@Tags			ProcessURL
//	@Accept			json
//	@Produce		json
//	@Param			RequestStruct	body	RequestStruct	true	"Request Body"
//	@Success		200 {object}	ResponseStruct
//	@Failure		400 {object}	ErrMessage
//	@Failure		405
//	@Router			/ [post]
func ProcessUrl(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// redirect
	if r.Method == "GET" {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Error : Request Body is empty"})
		w.Write(jsonResponse)
		return
	}

	var request RequestStruct
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	decodeErr := decoder.Decode(&request)
	if decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Invalid request format"})
		w.Write(jsonResponse)
		return
	}

	if IsUrl(request.Url) == false {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Url format invalid"})
		w.Write(jsonResponse)
		return
	}

	// We are sure it's a URL from this point

	//check operation type and run it, if not valid, return BAD REQUEST
	switch request.Operation {
	case "canonical":
		processedUrl, err := GetCanonicalUrl(request.Url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(ResponseStruct{ProcessedUrl: processedUrl})
		w.Write(jsonResponse)
		return
	case "redirection":
		processedUrl, err := GetRedirectionUrl(request.Url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(ResponseStruct{ProcessedUrl: processedUrl})
		w.Write(jsonResponse)
		return

	case "all":
		processedUrl, err := GetCanonicalUrl(request.Url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}

		processedUrl, err = GetRedirectionUrl(processedUrl) //reprocess url
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(ErrMessage{Msg: err.Error()})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(ResponseStruct{ProcessedUrl: processedUrl})
		w.Write(jsonResponse)
		return

	default:
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(ErrMessage{Msg: "Invalid operation"})
		w.Write(jsonResponse)
		return
	}
}

func Serve(port uint16) {

	serverMux := http.NewServeMux()
	middleWareMux := Cors(ResponseLogging(Logging((serverMux))))
	serverMux.HandleFunc("/", ProcessUrl)

	serverMux.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	httpSwagger.URL("docs/swagger.yaml")
	serverMux.Handle("/docs/", httpSwagger.WrapHandler)

	fmt.Println("Serving on port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), middleWareMux)

	if err != nil {
		log.Fatal(err)
	}
}
