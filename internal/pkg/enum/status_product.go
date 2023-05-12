package enum

type StatusProduct string

const (
	Processing StatusProduct = "processing"
	Approved                 = "approved"
	Rejected                 = "rejected"
)
