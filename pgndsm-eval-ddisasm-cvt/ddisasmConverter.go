package ddisasmcvt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)


// ReadDdisasm convert ddisasm output asm file into an int list contains insts offsets only
func ReadDdisasm(file string) (offsets []int) {
	offsets = make([]int, 0)
	fin, finerr := os.Open(file)
	if finerr != nil {
		fmt.Println(file + " does not exist.")
		return
	}
	defer fin.Close()

	lines := bufio.NewScanner(fin)
	for lines.Scan() {
		offset, _ := strconv.ParseInt(lines.Text(), 0, 64)
		offsets = append(offsets, int(offset))
	}
	return
}
