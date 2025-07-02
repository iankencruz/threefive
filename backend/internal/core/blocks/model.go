package blocks

type Block struct {
	Type  string         `json:"type"`
	Props map[string]any `json:"props"`
}
