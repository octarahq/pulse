package display

import "fmt"

func Draw(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}
