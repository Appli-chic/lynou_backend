package validator

type PostFileForm struct {
	Name      string `validate:"required"`
	Thumbnail string
	Type      uint `validate:"required"`
}

type PostForm struct {
	Text  string         `validate:"required"`
	Files []PostFileForm `validate:"required"`
}
