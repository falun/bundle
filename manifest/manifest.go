package manifest

type Manifest struct {
	Root  string `json:"root"`
	Files []File `json:"files"`
}

type File struct {
	Path     string `json:"path"`
	Compress bool   `json:"compress"`
}
