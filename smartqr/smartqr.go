package smartqr

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-toast/toast"
	"github.com/liyue201/goqr"
	qrencode "github.com/skip2/go-qrcode"
	"golang.design/x/clipboard"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

var imageFileName string
var absPath, defaultImageFileName string
var tempFileName, defaultTempFileName string
var ps1FileName, defaultPs1FileName string

func constructSmartQRPS1(path string) string {
	return "Get-Content " + path + " | clip\n"
}

func writeTempFile(fileName string, content string) {
	var file *os.File

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}

	file.WriteString(content)
	file.Close()
}

func decodeImage(imgContext []byte) [][]uint8 {

	img, _, err := image.Decode(bytes.NewReader(imgContext))
	if err != nil {
		fmt.Printf("image.Decode error: %v\n", err)
		return nil
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		fmt.Printf("Recognize failed: %v\n", err)
		return nil
	}
	var ret [][]uint8
	for _, qrCode := range qrCodes {
		fmt.Printf("qrCode text: %s\n", qrCode.Payload)
		ret = append(ret, qrCode.Payload)
	}

	file, _ := os.Create(defaultImageFileName)
	defer file.Close()
	png.Encode(file, img)

	return ret
}

func handleClipboardChange(data []byte) {
	qrcodeMessages := decodeImage(data)
	if qrcodeMessages == nil {
		fmt.Println("noting")
		return
	}

	message := "检测到二维码内容:\n"

	var action []toast.Action

	for _, qrMessage := range qrcodeMessages {

		writeTempFile(defaultTempFileName, string(qrMessage))
		writeTempFile(defaultPs1FileName, constructSmartQRPS1(defaultTempFileName))

		message = message + string(qrMessage) + "\n"
		action = append(action, toast.Action{
			Type:      "protocol",
			Label:     "点击跳转" + string(qrMessage),
			Arguments: string(qrMessage),
		})
		action = append(action, toast.Action{
			Type:      "protocol",
			Label:     "复制内容",
			Arguments: "smartqr://" + defaultPs1FileName,
		})

		break
	}

	fmt.Printf(message)
	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   "SimpleQR",
		Message: message,
		Icon:    defaultImageFileName,
		Actions: action,
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

func EncodeQR() {
	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		notification := toast.Notification{
			AppID:   "Microsoft.Windows.Shell.RunDialog",
			Title:   "SimpleQR",
			Message: "未检测到有效内容",
		}
		err := notification.Push()
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	imgContext, err := qrencode.Encode(string(data), qrencode.Medium, 256)
	if err != nil {
		log.Fatalln(err)
	}

	img, _, err := image.Decode(bytes.NewReader(imgContext))
	if err != nil {
		fmt.Printf("image.Decode error: %v\n", err)
		return
	}
	file, _ := os.Create(defaultImageFileName)
	defer file.Close()
	png.Encode(file, img)

	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   "SimpleQR",
		Message: "已复制到剪贴板，并保存图片",
		Actions: []toast.Action{
			{"protocol", "在文件夹中打开", filepath.Join(absPath, "temp")},
			{"protocol", "打开图片", defaultImageFileName},
		},
	}

	clipboard.Write(clipboard.FmtImage, imgContext)

	err = notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func Run() {
	imageFileName = "temp/simpleqr.png"
	tempFileName = "temp/tmp.txt"
	ps1FileName = "temp/smartqr.ps1"
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	absPath = filepath.Dir(path)
	os.MkdirAll(filepath.Join(absPath, "temp"), os.ModeDir)
	defaultImageFileName = filepath.Join(absPath, imageFileName)
	defaultTempFileName = filepath.Join(absPath, tempFileName)
	defaultPs1FileName = filepath.Join(absPath, ps1FileName)

	ch := clipboard.Watch(context.TODO(), clipboard.FmtImage)
	for data := range ch {
		handleClipboardChange(data)
	}

}
