package entities

type Entity []KeyValue

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
