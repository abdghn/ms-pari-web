package enum

type PARI int

const (
	CreateProduct PARI = iota + 1
	DetailProduct
)

// String - Creating common behavior - give the type a String function
func (p PARI) String() string {
	return [...]string{"/product/create", "/product/detail"}[p-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (p PARI) EnumIndex() int {
	return int(p)
}
