package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	ghidraimport "github.com/pangine/pangineDSM-import/ghidraCall"
	rstapi "github.com/pangine/pangineDSM-import/rstAPI"
	genutils "github.com/pangine/pangineDSM-utils/general"
	objectapi "github.com/pangine/pangineDSM-utils/objectAPI"
)

func ghidraAnalysis(hlLoc, binFile, scriptDir, logFile string) bool {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	headless := exec.Command(hlLoc,
		"/tmp", "groundtruth"+strconv.Itoa(r1.Int()),
		"-import", binFile,
		"-postScript", "OffsetOutput.java", logFile,
		"-scriptPath", scriptDir,
		"-overwrite", "-readOnly")
	headless.Stderr = os.Stderr
	headless.Stdout = os.Stdout
	err := headless.Start()
	if err != nil {
		fmt.Println("\tghidra start failed")
		return false
	}
	err = headless.Wait()
	if err != nil {
		fmt.Println("\tghidra execution failed")
		return false
	}
	return true
}

func ghidraPrefixFix(object objectapi.Object, binFile, ghiIn, ghiOut string) {
	bi := object.ParseObj(binFile)
	que, _ := rstapi.ReadRst(ghiIn)
	bout, err := os.Create(ghiOut)
	if err != nil {
		fmt.Println("Prefix fix failed")
		return
	}
	defer bout.Close()
	ety := object.InstLstFixForPrefix(que, bi)
	rstapi.WriteRstFromList(ety, bout)
	fmt.Println("Prefix fix succeeded")
}

func main() {
	argNum := len(os.Args)
	inputDir := os.Args[argNum-1]

	singleDirFlag := flag.String("sd", "", "only operate on a single dir")
	singleTargetFlag := flag.String("sf", "", "only operate on a single file")
	flag.Parse()

	singleDir := *singleDirFlag
	singleTarget := *singleTargetFlag

	ghiRoot := filepath.Join(inputDir, "ghidra")
	binRoot := filepath.Join(inputDir, "bin")
	_ = os.Mkdir(ghiRoot, os.ModePerm)

	pdir := ghidraimport.GoDir()

	var dirList []string
	if singleDir != "" {
		dirList = []string{singleDir}
	} else {
		dirList = genutils.GetDirs(binRoot)
	}
	for _, dir := range dirList {
		binDir := filepath.Join(binRoot, dir)
		ghiDir := filepath.Join(ghiRoot, dir)
		_ = os.Mkdir(ghiDir, os.ModePerm)

		var fileList []string
		if singleTarget != "" && singleDir != "" {
			fileList = []string{singleTarget}
		} else {
			fileList = genutils.GetFiles(binDir, "")
		}

		for _, file := range fileList {
			fmt.Printf("%-15s:\n", file)
			binFile := filepath.Join(binDir, file)
			ghiLog := filepath.Join(ghiDir, file+"_capnp.out")
			if !ghidraAnalysis("analyzeHeadless", binFile, filepath.Join(pdir, "ghidraScript"), ghiLog) {
				fmt.Println("failed")
			} else {
				fmt.Println("done")
			}
		}
	}
}
