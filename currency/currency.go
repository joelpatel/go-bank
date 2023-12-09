package currency

const (
	INR = "INR"
	USD = "USD"
)

// Return true if the currency is supported, else returns false.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case INR, USD:
		return true
	default:
		return false
	}
}
