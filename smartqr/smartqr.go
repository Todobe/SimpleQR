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

var CacheContent = []byte("Nothing here now.")

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

	message := ""

	var action []toast.Action

	for _, qrMessage := range qrcodeMessages {
		CacheContent = qrMessage
		message = message + string(qrMessage) + "\n"
		action = append(action, toast.Action{
			Type:      "protocol",
			Label:     "点击跳转" + string(qrMessage),
			Arguments: string(qrMessage),
		})
	}
	fmt.Printf(message)
	clipboard.Write(clipboard.FmtText, CacheContent)
	//notification := toast.Notification{
	//	AppID:   "SimpleQR",
	//	Title:   "检测到二维码内容：",
	//	Message: message,
	//	Icon:    defaultImageFileName,
	//	Actions: action,
	//}
	//err := notification.Push()
	//if err != nil {
	//	log.Fatalln(err)
	//}
}

func EncodeQR() {
	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		fmt.Println("No valid data in clipboard.")
		//notification := toast.Notification{
		//	AppID:   "SimpleQR",
		//	Title:   "Encode 失败",
		//	Message: "未检测到剪贴板内有效内容",
		//}
		//err := notification.Push()
		//if err != nil {
		//	log.Fatalln(err)
		//}
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

	//notification := toast.Notification{
	//	AppID:   "SimpleQR",
	//	Title:   "Encode成功",
	//	Message: "已复制到剪贴板，并保存图片",
	//	Actions: []toast.Action{
	//		{"protocol", "在文件夹中打开", filepath.Join(absPath, "temp")},
	//		{"protocol", "打开图片", defaultImageFileName},
	//	},
	//}

	clipboard.Write(clipboard.FmtImage, imgContext)

	//err = notification.Push()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	return
}

func CopyContent() {
	clipboard.Write(clipboard.FmtText, CacheContent)
}

func Run() {
	imageFileName = "temp/simpleqr.png"
	path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	absPath = filepath.Dir(path)
	os.MkdirAll(filepath.Join(absPath, "temp"), os.ModeDir)
	defaultImageFileName = filepath.Join(absPath, imageFileName)

	ch := clipboard.Watch(context.TODO(), clipboard.FmtImage)
	for data := range ch {
		handleClipboardChange(data)
	}

}
