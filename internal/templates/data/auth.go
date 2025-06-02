package data

type LoginPageData struct {
	Title string
	Form  struct {
		Email string
	}
	Errors map[string]string
}

type LoginForm struct {
	Email    string
	Password string
}

type RegisterPageData struct {
	Title  string
	Form   RegisterForm
	Errors map[string]string
}

type RegisterForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}
