package login

// Wrap a function closure with the Value interface
type wrapFunc func(optsKvList string) error

func (f wrapFunc) Set(value string) error {
	return f(value)
}

func (f wrapFunc) String() string {
	return "setFunc"
}
