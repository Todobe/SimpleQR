package main

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	linkPath = "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\StartUp\\SimpleQR_Link.exe.lnk"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func main() {
	dir := getCurrentDirectory()
	dir = getParentDirectory(dir)
	path := dir + "/SimpleQR.exe"
	src := path
	dst := linkPath
	fmt.Println(src, dst)
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println(err)
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println(err)
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
	if err != nil {
		fmt.Println(err)
	}
	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", src)
	oleutil.CallMethod(idispatch, "Save")
	idispatch = cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "WorkingDirectory", dir)
	oleutil.CallMethod(idispatch, "Save")
}
