package main

import (
	"embed"
	_ "embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var trayIcon []byte

func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	// Check if we are in dev mode. Wails v3 sets different flags or you can check for common dev indicators.
	// For Wails v3, we can check for the build tag or just use application.Options.
	options := application.Options{
		Name:        "fakeRing",
		Services: []application.Service{
			application.NewService(&RingLightService{}),
			application.NewService(&App{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	}
	// SingleInstance with a fresh ID to ensure it's not blocked by old PIDs
	if os.Getenv("WAILS_VITE_PORT") == "" {
		options.SingleInstance = &application.SingleInstanceOptions{
			UniqueID: "com.fakering.pro.v1",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				windows := application.Get().Window.GetAll()
				if len(windows) > 0 {
					w := windows[0]
					w.Show()
					w.UnMinimise()
					w.Focus()
				}
			},
		}
	}

	app := application.New(options)
	StartOverlay()

	// Setup Tray
	systray := app.SystemTray.New()
	systray.SetIcon(trayIcon)
	systray.SetLabel("FakeRing")

	menu := app.NewMenu()
	menu.Add("Open Dashboard").OnClick(func(ctx *application.Context) {
		windows := app.Window.GetAll()
		if len(windows) > 0 {
			w := windows[0]
			w.Show()
			w.UnMinimise()
			w.Focus()
		}
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(menu)

	app.OnShutdown(func() {
		systray.Destroy()
	})

	// Create Window
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "FakeRing Controller",
		Width:  450,
		Height: 600,
		BackgroundColour: application.NewRGB(15, 23, 42),
		URL:              "/",
		DisableResize:    true,
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
