package data

type LoginPageData struct {
	Title   string
	Message string
}

type LoginForm struct {
	Title   string
	Message string
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
