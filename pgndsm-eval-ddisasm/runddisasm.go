package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	genutils "github.com/pangine/pangineDSM-utils/general"
)

func ddisasmAnalysis(binFile, outFile string) bool {
	ddisasm := exec.Command("ddisasm", binFile, "--debug", "--asm", outFile)
	ddisasm.Stderr = os.Stderr
	ddisasm.Stdout = os.Stdout
	err := ddisasm.Start()
	if err != nil {
		fmt.Println("\tddisasm start failed")
		return false
	}
	err = ddisasm.Wait()
	if err != nil {
		fmt.Println("\tddisasm execution failed")
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

	ddmRoot := filepath.Join(inputDir, "ddisasm")
	binRoot := filepath.Join(inputDir, "bin")
	_ = os.Mkdir(ddmRoot, os.ModePerm)

	var dirList []string
	if singleDir != "" {
		dirList = []string{singleDir}
	} else {
		dirList = genutils.GetDirs(binRoot)
	}
	for _, dir := range dirList {
		ddmDir := filepath.Join(ddmRoot, dir)
		binDir := filepath.Join(binRoot, dir)
		_ = os.Mkdir(ddmDir, os.ModePerm)

		var fileList []string
		if singleTarget != "" && singleDir != "" {
			fileList = []string{singleTarget}
		} else {
			fileList = genutils.GetFiles(binDir, "")
		}

		for _, file := range fileList {
			fmt.Printf("%-15s:\n", file)
			binFile := filepath.Join(binDir, file)
			inFile := filepath.Join(ddmDir, file+".lst")
			if !ddisasmAnalysis(binFile, inFile) {
				fmt.Println("failed")
			} else {
				fmt.Println("done")
			}
		}
	}
}
