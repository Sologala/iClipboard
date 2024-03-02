package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/rs/zerolog/log"
	"golang.design/x/clipboard"
)

var ServiceFlag bool = true

func main() {
	InitLogger(true)
	Gconfig.LoadConfig()

	onExit := func() {
		// now := time.Now()
		log.Info().Msg(`on_exit_.txt`)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	var tray_name = "iClipboard"
	// var clipboard = ClipBoardBase{"linux"}
	// var clip_board_type = clipboard.info()

	err := clipboard.Init()
	if err != nil {
		fmt.Println("iClipboard init failid")
		panic(err)
	}

	systray.SetTemplateIcon(IconData, IconData)
	systray.SetTitle(tray_name)
	systray.SetTooltip("iClipboard")

	go func() {
		RunHTTPServer()
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(IconData, IconData)
		systray.SetTitle(tray_name)
		// systray.SetTooltip("Pretty awesome棒棒嗒")
		mChecked := systray.AddMenuItemCheckbox("Enabled", "Check Me", true)
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		// Sets the icon of a menu item. Only available on Mac.
		mQuit.SetIcon(IconData)
		systray.AddSeparator()
		for {
			select {
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Disabled")
					ServiceFlag = false
				} else {
					mChecked.Check()
					mChecked.SetTitle("Enabled")
					ServiceFlag = true
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}

// func write_2_file(file string) error {
// 	var b []byte
// 	var err error
// 	file = "hhh.txt"
// 	b = clipboard.Read(clipboard.FmtText)
// 	if b == nil {
// 		b = clipboard.Read(clipboard.FmtImage)
// 		file = "hhh.png"
// 	}
//
// 	if file != "" && b != nil {
// 		err = os.WriteFile(file, b, os.ModePerm)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "failed to write data to file %s: %v", file, err)
// 		}
// 		return err
// 	}
//
// 	for len(b) > 0 {
// 		n, err := os.Stdout.Write(b)
// 		if err != nil {
// 			return err
// 		}
// 		b = b[n:]
// 	}
// 	return nil
// }
