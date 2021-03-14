package main

import (
	"lofty/cmd/input"
	"lofty/cmd/lofty"
	"os"
)

func main() {
	//                          111111111122222222223333333333444444444455555555556
	//                0123456789012345678901234567890123456789012345678901234567890
	//                ((100000+(2+3)+4)-5+((6+7)+8)+9)
	//                (((1+2-3+4)+1)-2)
	if input.IsPiping() && len(os.Args) > 1 {
		lofty.SearchPipeMode()
	}
}
