package helpers

import (
	"encoding/json"
	"io"
	"net/http"
)

type JsonResp struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Maximum payload size - 1MB
const maxBytes int64 = 1048576

func ReadSingleJson(w http.ResponseWriter, r *http.Request, res any) error {
	// Set the maximum size limit (bytes per request)
	// Request Body is a stream of data, so we can set maximum size limit
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	// Create a JSON decoder that will read the JSON data directly from the incoming HTTP request
	decoder := json.NewDecoder(r.Body)
	// Decode the data from the incoming stream, into the result data variable
	err := decoder.Decode(res)
	if err != nil {
		return err
	}

	// Try to decode once more, to ensure only one json object is sent
	// This is a dummy conversion
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return err
	}

	return nil
}

func ReadAllJsonMutable(w http.ResponseWriter, r *http.Request, res *[]any) error {
	decoder := json.NewDecoder(r.Body)

	for {
		var t any

		if err := decoder.Decode(&t); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		*res = append(*res, t)
	}

	return nil
}

// This is not a production ready function, it is just to illustrate how can multiple value decoding be done
// Returns a pointer to slice which under the hood contains pointer to array, which is redundant
func ReadAllJsonImmutable(w http.ResponseWriter, r *http.Request) (*[]any, error) {
	decoder := json.NewDecoder(r.Body)

	res := []any{}
	for {
		var t any

		if err := decoder.Decode(&t); err == io.EOF {
			break
		} else if err != nil {
			return &res, err
		}

		res = append(res, t)
	}

	return &res, nil
}

func WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	jsonEncoding, err := json.Marshal(data)
	if err != nil {
		return err
	}

	setHeaders(w, status, headers...)

	_, err = w.Write(jsonEncoding)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	// Init status code with default value
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	resp := &JsonResp{
		Error:   true,
		Message: err.Error(),
	}

	return WriteJson(w, statusCode, resp)
}

func setHeaders(w http.ResponseWriter, status int, headers ...http.Header) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if len(headers) > 0 {
		for _, header := range headers {
			for k, v := range header {
				w.Header()[k] = v
			}
		}
	}

}
