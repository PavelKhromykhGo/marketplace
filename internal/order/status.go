package order

const (
	StatusNew             = "new"
	StatusAwaitingPayment = "awaiting_payment"
	StatusPaid            = "paid"
	StatusShipped         = "shipped"
	StatusDelivered       = "delivered"
	StatusCancelled       = "cancelled"
)

var allowedStatusTransitions = map[string]map[string]struct{}{
	StatusNew: {
		StatusAwaitingPayment: {},
		StatusCancelled:       {},
	},
	StatusAwaitingPayment: {
		StatusPaid:      {},
		StatusCancelled: {},
	},
	StatusPaid: {
		StatusShipped:   {},
		StatusCancelled: {},
	},
	StatusShipped: {
		StatusDelivered: {},
	},
}

func IsValidStatusTransition(from, to string) bool {
	if from == to {
		return true
	}
	nexts, ok := allowedStatusTransitions[from]
	if !ok {
		return false
	}
	_, ok = nexts[to]
	return ok
}
