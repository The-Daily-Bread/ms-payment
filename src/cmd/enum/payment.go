package enum

type PaymentMethod string

const (
	CASH        PaymentMethod = "cash"
	CREDIT_CARD PaymentMethod = "credit_card"
	PIX         PaymentMethod = "pix"
)

type PaymentStatus int

const (
	DENIED PaymentStatus = iota
	APPROVED
)
