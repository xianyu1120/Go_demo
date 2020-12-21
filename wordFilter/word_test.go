package wordFilter

import (
	"fmt"
	"testing"
)

func TestWord(t *testing.T) {
	str := "世界,我的,下贱,淫贱"
	all, isBool := CheckoutWordAll(str)
	fmt.Println(all, isBool)
}
