package gs

import (
	"fmt"
	"os"
)

// Print banner.txt if it exists.
func PrintBanner() {
	if data, err := os.ReadFile("banner.txt"); err == nil {
		fmt.Println(string(data))
	}
}
