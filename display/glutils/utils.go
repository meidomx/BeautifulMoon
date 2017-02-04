package glutils

type Bindable interface {
	Bind()
	ResetBinding()
}

func NewBindableFunc(bindable Bindable, f func()) func() {
	return func() {
		bindable.Bind()
		f()
		bindable.ResetBinding()
	}
}
