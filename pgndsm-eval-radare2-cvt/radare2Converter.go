package radare2cvt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadRadare2 convert radare2 print disassembly output into an int list contains insts offsets only
func ReadRadare2(file string) (offsets []int) {
	offsets = make([]int, 0)
	offsetsMap := make(map[int]bool)
	fin, finerr := os.Open(file)
	if finerr != nil {
		//panic(file + " does not exist.")
		fmt.Println(file + " does not exist.")
		return
	}
	defer fin.Close()
	lines := bufio.NewScanner(fin)
	for lines.Scan() {
		var IsAddress bool
		line := lines.Text()
		fields := strings.Fields(line)
		// Check if the line start with an address
		for len(fields) > 0 {
			if strings.HasPrefix(fields[0], ";") {
				// Is comment
				break
			}
			if strings.HasSuffix(fields[0], ":") &&
				strings.HasSuffix(line, ";") {
				strNum := fields[0][:len(fields[0])-1]
				if _, err := strconv.ParseInt(strNum, 10, 64); err == nil {
					// A function define
					break
				}
			}
			if strings.HasPrefix(fields[0], "0x") {
				IsAddress = true
				break
			}
			fields = fields[1:]
		}
		if !IsAddress ||
			len(fields) < 3 {
			continue
		}
		MCodeIndex := 1
		if fields[1] == "~" {
			MCodeIndex = 2
		}
		if strings.HasPrefix(fields[MCodeIndex], ".") || //is data
			fields[MCodeIndex+1] == "invalid" { // invalid insn
			continue
		}
		addr64, err := strconv.ParseInt(fields[0][2:], 16, 64)
		if err != nil {
			fmt.Print("WARNING: ")
			fmt.Println(err)
			continue
		}
		// Use a map because there can be repeat instructions
		offsetsMap[int(addr64)] = true
	}
	for i := range offsetsMap {
		offsets = append(offsets, i)
	}
	return
}
