package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	genutils "github.com/pangine/pangineDSM-utils/general"
)

func r2Analysis(binFile, outFile string) bool {
	r2 := exec.Command("r2", "-Aqc", "pdr @@f > "+outFile, binFile)
	r2.Stderr = os.Stderr
	r2.Stdout = os.Stdout
	err := r2.Start()
	if err != nil {
		fmt.Println("\tr2 start failed")
		return false
	}
	err = r2.Wait()
	if err != nil {
		fmt.Println("\tr2 execution failed")
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

	r2Root := filepath.Join(inputDir, "radare2")
	binRoot := filepath.Join(inputDir, "bin")
	_ = os.Mkdir(r2Root, os.ModePerm)

	var dirList []string
	if singleDir != "" {
		dirList = []string{singleDir}
	} else {
		dirList = genutils.GetDirs(binRoot)
	}
	for _, dir := range dirList {
		r2Dir := filepath.Join(r2Root, dir)
		binDir := filepath.Join(binRoot, dir)
		_ = os.Mkdir(r2Dir, os.ModePerm)

		var fileList []string
		if singleTarget != "" && singleDir != "" {
			fileList = []string{singleTarget}
		} else {
			fileList = genutils.GetFiles(binDir, "")
		}

		for _, file := range fileList {
			fmt.Printf("%-15s:\n", file)
			binFile := filepath.Join(binDir, file)
			inFile := filepath.Join(r2Dir, file+".lst")
			if !r2Analysis(binFile, inFile) {
				fmt.Println("failed")
			} else {
				fmt.Println("done")
			}
		}
	}
}
