package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestStruct struct {
	Url       string `json:"url"`
	Operation string `json:"operation"`
}

type ResponseStruct struct {
	ProcessedUrl string `json:"processed_url"`
}

type ErrMessage struct {
	Msg string `json:"msg"`
}

// routes requests that have ID based on HTTP
func ProcessUrl(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
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

	fmt.Println("Serving on port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), middleWareMux)

	if err != nil {
		log.Fatal(err)
	}
}
