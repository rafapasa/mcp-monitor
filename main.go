//go:build windows

package main

import (
	"fmt"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/rafapasa/mcp-monitor/internal/captor"
	"github.com/rafapasa/mcp-monitor/internal/llm"
)

func main() {
	llm.Init()

	var lblResult *walk.Label
	var lblTimer *walk.Label
	var window *walk.MainWindow

	captor.Init(func(img []byte) {
		go func() {
			window.Synchronize(func() {
				lblResult.SetText("ANALISANDO...")
				color, _ := walk.NewSolidColorBrush(walk.RGB(255, 255, 0))
				lblResult.SetBackground(color)
			})
			result := llm.ProcessImageSync(img) // vamos criar essa função sync
			window.Synchronize(func() {
				lblResult.SetText(result.Acao)
				//PAGAR, PASSAR, CORRER, SUBIR.
				switch result.Acao {
				case "SUBIR":
					color, _ := walk.NewSolidColorBrush(walk.RGB(0, 200, 0))
					lblResult.SetBackground(color)
				case "PAGAR":
					color, _ := walk.NewSolidColorBrush(walk.RGB(200, 0, 0))
					lblResult.SetBackground(color)
				default:
					color, _ := walk.NewSolidColorBrush(walk.RGB(80, 80, 80))
					lblResult.SetBackground(color)
				}
				lblResult.SetTextColor(walk.RGB(255, 255, 255))
			})
			// countdown 10s
			for i := 10; i >= 0; i-- {
				window.Synchronize(func() {
					lblTimer.SetText(fmt.Sprintf("%ds - %s", i, time.Now().Format("15:04:05")))
				})
				time.Sleep(1 * time.Second)
			}
			window.Synchronize(func() {
				color, _ := walk.NewSolidColorBrush(walk.RGB(20, 20, 20))
				lblResult.SetBackground(color)
				lblResult.SetText("EXPIRADO")
			})
		}()
	})

	MainWindow{
		AssignTo: &window,
		Title:    "MCP - F9 CAPTURA",
		MinSize:  Size{320, 280},
		Size:     Size{340, 400},
		Layout:   VBox{},
		Children: []Widget{
			Label{AssignTo: &lblTimer, Text: "Aguardando F9...", MinSize: Size{0, 30}},
			Label{
				AssignTo:  &lblResult,
				Text:      "PRONTO",
				MinSize:   Size{0, 150},
				Font:      Font{Family: "Arial Black", PointSize: 28, Bold: true},
				Alignment: AlignHCenterVCenter,
			},
			PushButton{
				Text:      "F9 - CAPTURAR AGORA",
				OnClicked: func() { go captor.CaptureScreen() },
			},
		},
	}.Create()

	// Always on top
	win.SetWindowPos(window.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
	// Hotkey global F9
	// win.keyNewAction().Shortcut()KeyString(&walk.RegistryKey{}window.Handle(), 1, 0, win.VK_F9)

	window.Run()
}
