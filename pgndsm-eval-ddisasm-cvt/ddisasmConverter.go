package ddisasmcvt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func RunDdisasmConverter(gtirbFile, insnFile string) bool {
	ddisasm := exec.Command("ddisasmConverter.py", gtirbFile, outFile)
	ddisasm.Stderr = os.Stderr
	ddisasm.Stdout = os.Stdout
	err := ddisasm.Start()
	if err != nil {
		fmt.Println("\tddisasm converter start failed")
		return false
	}
	err = ddisasm.Wait()
	if err != nil {
		fmt.Println("\tddisasm converter execution failed")
		return false
	}
	return true
}

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
		offset, err := strconv.ParseInt(lines.Text(), 16, 64)
		offsets = append(offsets, offset)
	}
	return
}
