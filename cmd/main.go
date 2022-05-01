package main

import "lofty/cmd/lofty"

func main() {
	//lofty.SearchString("lazy & !dog", "The lazy dog jumped over the fence")
	//////////////////////////////012345678901234567890
	//lofty.SearchString("lazy or (the or jumped)", "The lazy dog jumped over the fence")
	/////////////////////////////012345678901234567890
	lofty.SearchString("0123456789012345678901234567890", "The lazy dog jumped over the fence")
	lofty.SearchString("lazy and not dog", "The lazy dog jumped over the fence")
	lofty.SearchString("not (((lazy and not dog))) or ((((lazy and dog)))) or not((((lazy or dog)))) or ((((not lazy and dog))))", "The lazy dog jumped over the fence")
	lofty.SearchString("lazy or not (the or not jumped)", "The lazy dog jumped over the fence")
}
