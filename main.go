package main

import (
	"fmt"
	"time"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"golang.design/x/clipboard"
)

type ClipBoardBase struct {
	platform_type string
}

func (d ClipBoardBase) info() string {
	return d.platform_type
}

func main() {
	onExit := func() {
		now := time.Now()
		println(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	var tray_name = "UniCLip"
	// var clipboard = ClipBoardBase{"linux"}
	// var clip_board_type = clipboard.info()

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle(tray_name)
	systray.SetTooltip("Lantern")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() { // 表示创建一个新的轻量级线程，异步执行一些函数。
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
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

		mQuit := systray.AddMenuItem("退出", "Quit the whole app")

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
                clipboard.Write(clipboard.FmtText, []byte("text data"))
			}
		}
	}()

}
