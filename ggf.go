package ggf

// GGF is one of filename extention
type GGF struct {
	Path string
	Data map[string][]string
}

// New return GGF object
func New(path string) (*GGF, error) {
	ggf := &GGF{
		Path: path,
	}
	return ggf, nil
}
