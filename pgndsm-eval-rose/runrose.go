package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	genutils "github.com/pangine/pangineDSM-utils/general"
)

func roseAnalysis(binFile, outFile string) bool {
	rose := exec.Command("rose-recursive-disassemble", binFile)
	rose.Stderr = os.Stderr
	out, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("\tout file %s create failed\n", outFile)
		return false
	}
	defer out.Close()
	rose.Stdout = out
	err = rose.Start()
	if err != nil {
		fmt.Println("\trose start failed")
		return false
	}
	err = rose.Wait()
	if err != nil {
		fmt.Println("\trose execution failed")
		return false
	}
	return true
}

func main() {
	argNum := len(os.Args)
	inputDir := os.Args[argNum-1]

	singleDirFlag := flag.String("sd", "", "only operate on a single dir")
	singleTargetFlag := flag.String("sf", "", "only operate on a single file")
	flag.Parse()

	singleDir := *singleDirFlag
	singleTarget := *singleTargetFlag

	roseRoot := filepath.Join(inputDir, "rose")
	binRoot := filepath.Join(inputDir, "bin")
	_ = os.Mkdir(roseRoot, os.ModePerm)

	var dirList []string
	if singleDir != "" {
		dirList = []string{singleDir}
	} else {
		dirList = genutils.GetDirs(binRoot)
	}
	for _, dir := range dirList {
		roseDir := filepath.Join(roseRoot, dir)
		binDir := filepath.Join(binRoot, dir)
		_ = os.Mkdir(roseDir, os.ModePerm)

		var fileList []string
		if singleTarget != "" && singleDir != "" {
			fileList = []string{singleTarget}
		} else {
			fileList = genutils.GetFiles(binDir, "")
		}

		for _, file := range fileList {
			fmt.Printf("%-15s:\n", file)
			binFile := filepath.Join(binDir, file)
			inFile := filepath.Join(roseDir, file+".lst")
			if !roseAnalysis(binFile, inFile) {
				fmt.Println("failed")
			} else {
				fmt.Println("done")
			}
		}
	}
}
