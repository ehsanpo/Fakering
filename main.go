package main

import (
	"embed"
	_ "embed"
	"log"
	"os"
	"time"

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

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
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
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(&RingLightService{}),
			application.NewService(&App{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	}

	// SingleInstance is only desired in production.
	// Wails v3 dev mode typically sets WAILS_VITE_PORT or similar environment variables.
	if os.Getenv("WAILS_VITE_PORT") == "" {
		options.SingleInstance = &application.SingleInstanceOptions{
			UniqueID: "com.fakering.elite",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				// Activate the main window when a second instance is launched
				windows := application.Get().Window.GetAll()
				if len(windows) > 0 {
					windows[0].Show()
					windows[0].UnMinimise()
					windows[0].Focus()
				}
			},
		}
	}

	app := application.New(options)

	StartOverlay()


	// Create system tray
	systray := app.SystemTray.New()
	systray.SetIcon(trayIcon)
	systray.SetLabel("fakeRing")

	// Add system tray menu
	menu := app.NewMenu()
	menu.Add("Show Window").OnClick(func(ctx *application.Context) {
		// Show and unminimize the window
		windows := app.Window.GetAll()
		if len(windows) > 0 {
			windows[0].Show()
			windows[0].UnMinimise()
		}
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(menu)

	// Create a new window for the controller UI
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "FakeRing Controller",
		Width: 450,
		Height: 600,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(15, 23, 42),
		URL:              "/",
		DisableResize:    true,
	})


	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
