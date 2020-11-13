package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	genutils "github.com/pangine/pangineDSM-utils/general"
)

func bapAnalysis(binFile, outFile string) bool {
	bap := exec.Command("bap", binFile, "-d", "asm")
	bap.Stderr = os.Stderr
	out, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("\tout file %s create failed\n", outFile)
		return false
	}
	defer out.Close()
	bap.Stdout = out
	err = bap.Start()
	if err != nil {
		fmt.Println("\tbap start failed")
		return false
	}
	err = bap.Wait()
	if err != nil {
		fmt.Println("\tbap execution failed")
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

	bapRoot := filepath.Join(inputDir, "bap")
	binRoot := filepath.Join(inputDir, "bin")
	_ = os.Mkdir(bapRoot, os.ModePerm)

	var dirList []string
	if singleDir != "" {
		dirList = []string{singleDir}
	} else {
		dirList = genutils.GetDirs(binRoot)
	}
	for _, dir := range dirList {
		bapDir := filepath.Join(bapRoot, dir)
		binDir := filepath.Join(binRoot, dir)
		_ = os.Mkdir(bapDir, os.ModePerm)

		var fileList []string
		if singleTarget != "" && singleDir != "" {
			fileList = []string{singleTarget}
		} else {
			fileList = genutils.GetFiles(binDir, "")
		}

		for _, file := range fileList {
			fmt.Printf("%-15s:\n", file)
			binFile := filepath.Join(binDir, file)
			inFile := filepath.Join(bapDir, file+".lst")
			if !bapAnalysis(binFile, inFile) {
				fmt.Println("failed")
			} else {
				fmt.Println("done")
			}
		}
	}
}
