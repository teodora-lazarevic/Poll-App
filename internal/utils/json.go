package utils

import (
	"encoding/json"
	"net/http"
)

// Reads the JSON body and parses it into a target struct
func ReadJSON(writer http.ResponseWriter, request *http.Request, target interface{}) error {
	request.Body = http.MaxBytesReader(writer, request.Body, 1_048_576) // Limit the request body size to 1MB
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields in the JSON body

	return decoder.Decode(target) // Decode the JSON body into the target struct
}

// Sends the JSON response to the client
func WriteJSON(writer http.ResponseWriter, status int, data interface{}) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	return json.NewEncoder(writer).Encode(data) // Encode the data into JSON and write it to the response writer
}

// Sends the error response to the client with the given status code and message
func ErrorJSON(writer http.ResponseWriter, status int, message string) {
	WriteJSON(writer, status, map[string]string{"error": message})
}
