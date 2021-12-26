package easyjson

//easyjson:json
type Person struct {
	Name    string   `json:"name"`
	Age     int32    `json:"age"`
	School  string   `json:"school"`
	Hobbies []string `json:"hobbies"`
}
