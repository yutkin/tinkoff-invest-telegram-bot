package currency

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCurrencyConverter_getCurrencyConvertRate(t *testing.T) {
	conv := New(os.Getenv("CURRENCY_API_TOKEN"), 3*time.Second)
	fmt.Println(conv.ConvertRate("USD", "RUB"))
}
