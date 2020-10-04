package main

import (
	"bufio"
	"fmt"
	"haruue.moe/x/cidr4"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	cs := cidr4.CIDRv4Set{}

	reader := bufio.NewReader(os.Stdin)

	for {
		bs, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[FATAL]: error when read from <stdin>")
			return
		}

		line := strings.TrimSpace(string(bs))
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "+") {
			line = line[1:]
		}
		minusMode := false
		if strings.HasPrefix(line, "-") {
			minusMode = true
			line = line[1:]
		}
		c, err := cidr4.ParseCIDRv4(line)
		if err != nil {
			log.Fatalf("[FATAL]: error when parse cidr: %s", err.Error())
			return
		}
		if minusMode {
			cs.Minus(c)
		} else {
			cs.Plus(c)
		}
	}

	for _, c := range cs {
		fmt.Println(c.String())
	}
}
