package main

import (
	"fmt"
	"strings"
)

func main() {
	var tableDeta []string
	tableDeta = append(tableDeta, fmt.Sprintf("%s %s %s", "u", "p", "o"), "pp")
	m := strings.Join(tableDeta, ",")
	z := fmt.Sprintf("CREATE TABLE %s (%s);", "pppp", m)
	fmt.Println(z)
}
