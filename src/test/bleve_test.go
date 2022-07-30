package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestBleveSearch(t *testing.T)  {
	s := "111111\n"
	i := strings.Index(s, "\n")
	fmt.Printf("i:%v\n",i)
	s2 := s[:i]
	fmt.Printf("s2: %v\n",s2)
}