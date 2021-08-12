package config

type Server struct {
	Folder Folder `yaml:"folder"`
}

type Folder struct {
	Path string `yaml:"path"`
}
