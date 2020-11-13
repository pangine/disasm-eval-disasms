package ddisasmcvt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadDdisasm convert ddisasm output asm file into an int list contains insts offsets only
func ReadDdisasm(file string) (offsets []int) {
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
	// ddisasm devide a long nop inst into several one byte nops
	// such feature would cause errors in recognition without merging back
	prevNop := false
	for lines.Scan() {
		fields := strings.Fields(lines.Text())
		if len(fields) < 2 {
			continue
		}
		if !strings.HasSuffix(fields[0], ":") || strings.HasPrefix(fields[1], ".") {
			continue
		}
		offset, err := strconv.ParseInt(fields[0][:len(fields[0])-1], 16, 64)
		if err != nil {
			continue
		}
		if fields[1] == "nop" {
			if prevNop {
				continue
			}
			prevNop = true
		} else {
			prevNop = false
		}
		offsetsMap[int(offset)] = true
	}
	for i := range offsetsMap {
		offsets = append(offsets, i)
	}
	return
}
