package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/rs/zerolog/log"
	"golang.design/x/clipboard"
	"os"
)

func main() {
	InitLogger(true)
	Gconfig.LoadConfig()

	onExit := func() {
		// now := time.Now()
		log.Info().Msg(`on_exit_.txt`)
	}

	systray.Run(onReady, onExit)
}

func write_2_file(file string) error {
	var b []byte
	var err error
	file = "hhh.txt"
	b = clipboard.Read(clipboard.FmtText)
	if b == nil {
		b = clipboard.Read(clipboard.FmtImage)
		file = "hhh.png"
	}

	if file != "" && b != nil {
		err = os.WriteFile(file, b, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write data to file %s: %v", file, err)
		}
		return err
	}

	for len(b) > 0 {
		n, err := os.Stdout.Write(b)
		if err != nil {
			return err
		}
		b = b[n:]
	}
	return nil
}

func onReady() {
	var tray_name = "UniCLip"
	// var clipboard = ClipBoardBase{"linux"}
	// var clip_board_type = clipboard.info()

	err := clipboard.Init()
	if err != nil {
		fmt.Println("clipboard init failid")
		panic(err)
	}

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle(tray_name)
	systray.SetTooltip("Lantern")
    
    go func(){
        RunHTTPServer()        
    }()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle(tray_name)
		systray.SetTooltip("Pretty awesome棒棒嗒")
		mChecked := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
		mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		// Sets the icon of a menu item. Only available on Mac.
		mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		b_write_sth_to_clip_board := systray.AddMenuItem("write some to clipboard", "---------")

		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		// Sets the icon of a menu item. Only available on Mac.
		mQuit.SetIcon(icon.Data)
		systray.AddSeparator()
		for {
			select {
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Unchecked")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Checked")
				}
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("Disabled")
				mEnabled.Disable()
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			case <-b_write_sth_to_clip_board.ClickedCh:
				// clipboard.Write(clipboard.FmtText, []byte("text data"))
				write_2_file("hh")
			}
		}
	}()

}
