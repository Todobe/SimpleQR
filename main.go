package main

import (
	"bytes"
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

func decode_clipboard() [][]uint8 {
	fmt.Println("hello world!")
	imgContext := clipboard.Read(clipboard.FmtImage)
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

func main() {
	qrcodeMessages := decode_clipboard()
	if qrcodeMessages == nil {
		fmt.Println("noting")
		return
	}

	message := "检测到二维码内容:\n"

	var action []toast.Action

	for _, qrMessage := range qrcodeMessages {
		fmt.Println(qrMessage)
		message = message + string(qrMessage) + "\n"
		action = append(action, toast.Action{
			Type:      "protocol",
			Label:     "点击跳转连接" + string(qrMessage),
			Arguments: string(qrMessage),
		})
	}
	action = append(action, toast.Action{
		Type:      "action",
		Label:     "reply",
		Arguments: "reply",
	})
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
