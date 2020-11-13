package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	disasmsutils "github.com/pangine/disasm-eval-disasms/disasmsutils"
	bapcvt "github.com/pangine/disasm-eval-disasms/pgndsm-eval-bap-cvt"
	rstapi "github.com/pangine/pangineDSM-import/rstAPI"
	objx86coff "github.com/pangine/pangineDSM-obj-x86-coff"
	objx86elf "github.com/pangine/pangineDSM-obj-x86-elf"
	genutils "github.com/pangine/pangineDSM-utils/general"
	objectapi "github.com/pangine/pangineDSM-utils/objectAPI"
)

func main() {
	argNum := len(os.Args)
	inputDir := os.Args[argNum-1]

	ltFlag := flag.String("l", "x86_64-PC-Linux-GNU-ELF", "the llvm triple for the target binaries")
	singleDirFlag := flag.String("sd", "", "only operate on a single dir")
	singleTargetFlag := flag.String("sf", "", "only operate on a single file")
	rvlISAFlag := flag.String("ra", "", "specify a ISA to start llvmmc-resolver (by default it will be auto detected according to input llvm triple)")
	printFlag := flag.Bool("print", false, "Print supported llvm triple types for this program")
	flag.Parse()

	llvmTriple := *ltFlag
	singleDir := *singleDirFlag
	singleTarget := *singleTargetFlag
	rvlISA := *rvlISAFlag
	printLLVM := *printFlag

	if printLLVM {
		genutils.PrintSupportLlvmTriple(disasmsutils.LLVMTriples)
		return
	}

	llvmTripleStruct := genutils.ParseLlvmTriple(genutils.CheckLlvmTriple(llvmTriple, disasmsutils.LLVMTriples))
	osEnvObj := llvmTripleStruct.OS + "-" + llvmTripleStruct.Env + "-" + llvmTripleStruct.Obj

	if rvlISA == "" {
		rvlISA = llvmTripleStruct.Arch
	}

	fmt.Println("Start llvmmc-resolver...")
	resolver := exec.Command("resolver", "-p", rvlISA)
	resolver.Start()
	time.Sleep(time.Second)
	defer resolver.Process.Kill()

	var object objectapi.Object
	switch osEnvObj {
	case "Linux-GNU-ELF":
		object = objx86elf.ObjectElf{}
	case "Win32-MSVC-COFF":
		object = objx86coff.ObjectCoff{}
	}

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
			fmt.Printf("%-15s: ", file)
			binFile := filepath.Join(binDir, file)
			inFile := filepath.Join(bapDir, file+".lst")
			if _, err := os.Stat(inFile); os.IsNotExist(err) {
				fmt.Println("original output does not exist")
				continue
			}
			outFile := filepath.Join(bapDir, file+"_bap.out")
			bi := object.ParseObj(binFile)
			que := bapcvt.ReadBap(inFile)
			ety := object.InstLstFixForPrefix(que, bi)
			bout, err := os.Create(outFile)
			if err != nil {
				fmt.Println("cannot create output capnp file")
				continue
			}
			rstapi.WriteRstFromList(ety, bout)
			bout.Close()
			fmt.Println("done")
		}
	}
}
