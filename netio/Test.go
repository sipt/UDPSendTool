package netio

import (
	"fmt"
)

type Person struct {
	Name string
}

func (p *Person) P() {
	fmt.Println(p.Name)
}
