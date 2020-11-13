package radare2cvt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadRose convert rose print disassembly output into an int list contains insts offsets only
func ReadRose(file string) (offsets []int) {
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
		line := lines.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		if strings.HasPrefix(fields[0], ";") {
			// Comments
			continue
		}
		if !strings.HasSuffix(fields[0], ":") {
			continue
		}
		if strings.HasSuffix(line, "|") {
			// Data section
			continue
		}
		strNum := fields[0][2 : len(fields[0])-1]
		addr64, err := strconv.ParseInt(strNum, 16, 64)
		if err != nil {
			continue
		}
		offsetsMap[int(addr64)] = true
	}
	for i := range offsetsMap {
		offsets = append(offsets, i)
	}
	return
}
