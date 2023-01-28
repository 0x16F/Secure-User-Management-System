package validate

const (
	MinPasswordLength = 8
	MaxPasswordLength = 64
	MinLoginLength    = 3
	MaxLoginLength    = 24
	MinLimit          = 0
	MaxLimit          = 50
)

const (
	OrderAsc  = "ASC"
	OrderDesc = "DESC"
)

var Orders = []string{OrderAsc, OrderDesc}
