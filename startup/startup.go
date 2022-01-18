package startup

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	openStartUp  = "Automatically start up at boot"
	closeStartUp = "Cancel automatically start up at boot"
	linkPath     = "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\StartUp\\SimpleQR_Link.exe.lnk"
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func CheckAutoStart() (bool, string, error) {
	_, err := os.Stat(linkPath)
	if err == nil {
		return true, closeStartUp, nil
	}
	if os.IsNotExist(err) {
		return false, openStartUp, nil
	}
	return false, openStartUp, err
}

func MakeLink() error {
	path := getCurrentDirectory() + "/autostartup/SimpleQRAutoStartUp.exe"
	c := exec.Command("cmd", "/C", "call", path)
	err := c.Run()
	return err
}

func RemoveLink() error {
	err := os.Remove(linkPath)
	return err
}
