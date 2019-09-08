package validator

type PostForm struct {
	Text string `validate:"required"`
}
