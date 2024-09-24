package util

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
)

func IsSupprotedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP:
		return true
	}

	return false
}
