package main

import (
	"fmt"
	"lofty/cmd/lofty"
)

func main() {
	//lofty.SearchString("lazy & !dog", "The lazy dog jumped over the fence")
	//////////////////////////////012345678901234567890
	//lofty.SearchString("lazy or (the or jumped)", "The lazy dog jumped over the fence")
	/////////////////////////////012345678901234567890
	//lofty.SearchString("0123456789012345678901234567890", "The lazy dog jumped over the fence")
	fmt.Println(lofty.SearchString("ladπle|(the|jumped)", "The lazy dog jumped over the fence"))
	fmt.Println(lofty.SearchString(" lazy  &  dog ", "The lazy dog jumped over the fence"))
	fmt.Println(lofty.SearchString("\"lazy\" & \"dog\"", "The lazy dog jumped over the fence"))
}
