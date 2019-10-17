package signingin

type (
	LoginForm struct {
		Login    string
		Password string
	}
	Account struct {
		Login       string
		AccessToken string
	}
	Value struct {
	}
)

func New() *Value {
	return &Value{}
}
