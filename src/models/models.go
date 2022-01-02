package models

type LinkEntry struct {
	ID   int
	Link string
	Hash string
}

type Result struct {
	Link string `json:"link"`
}

type CompressRequest struct {
	Link string `json:"link"`
}
