package libs

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func EchoSnippet(filepath string, startln int, endln int) {
	if startln > endln || startln < 1 {
		fmt.Println("Invalid line height")
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentln := 1

	for scanner.Scan() {
		if currentln >= startln && currentln <= endln {
			fmt.Printf("%d | %s\n", currentln, scanner.Text())
		}
		if currentln >= endln {
			return
		}
		currentln++
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scan file %s", err)
	}
}
