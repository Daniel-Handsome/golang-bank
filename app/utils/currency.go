package utils



// const (
// 	USD = "USD"
// 	EUR = "EUR"
// 	GBP = "GBP"
// )

// var currencies = map[string]string {
// 	"USD" : "USD",
// 	"EUR" : "EUR",
// 	"GBP" : "GBP",
// }

var currencies = []string {
	"USD",
	"EUR",
	"GBP",
}

func IsSupporedCurrency(currency string) bool {
	for _, v := range currencies {
		if v == currency {
            return true
        }
	}
	return false
}