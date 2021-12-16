package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-toast/toast"
	"github.com/liyue201/goqr"
	"golang.design/x/clipboard"
	"image"
	"image/png"
	"log"
	"os"
)

var default_imagefile_name = "C:\\Users\\Dell\\go\\src\\SimpleQR\\icon\\simpleqr.png"

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

	file, _ := os.Create(default_imagefile_name)
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

		file, err := os.OpenFile("./tmp.txt", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		file.WriteString(string(qrMessage))
		file.Close()

		fmt.Println(qrMessage)
		message = message + string(qrMessage) + "\n"
		action = append(action, toast.Action{
			Type:      "protocol",
			Label:     "点击跳转" + string(qrMessage),
			Arguments: string(qrMessage),
		})
		action = append(action, toast.Action{
			Type:      "protocol",
			Label:     "复制内容",
			Arguments: "smartqr://C:\\Users\\DELL\\go\\src\\SimpleQR\\smartqr.ps1",
		})
	}

	fmt.Printf(message)
	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   "SimpleQR",
		Message: message,
		Icon:    default_imagefile_name,
		Actions: action,
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	ch := clipboard.Watch(context.TODO(), clipboard.FmtImage)
	for data := range ch {
		handleClipboardChange(data)
	}

}
