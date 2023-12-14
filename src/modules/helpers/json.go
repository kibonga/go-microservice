package helpers

import "fmt"

type JsonResp struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Maximum payload size - 1MB
const maxBytes = 1048576

func (config *Config123) ReadJSON() {
	fmt.Println("Hello from ReadJson")
}
