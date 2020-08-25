package gee

import (
	"fmt"
	"testing"
)
func TestmatchChild(t *testing.T){
	nod:=&node{}
	pr:=nod.matchChild("").pattern
	fmt.Println(pr)
}