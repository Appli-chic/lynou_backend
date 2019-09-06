package validator

type SignUpUserForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=100"`
	Name     string `validate:"required"`
}

type LoginUserForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=100"`
}

type RefreshingTokenForm struct {
	RefreshToken string `validate:"required"`
}
