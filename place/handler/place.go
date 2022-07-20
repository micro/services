package handler

type Place struct {
	*Google
}

func New() *Place {
	g := NewGoogle()
	return &Place{g}
}
