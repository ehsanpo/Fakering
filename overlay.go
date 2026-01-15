package main

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	gdi32            = syscall.NewLazyDLL("gdi32.dll")
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procShowWindow          = user32.NewProc("ShowWindow")
	procUpdateWindow        = user32.NewProc("UpdateWindow")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procBeginPaint          = user32.NewProc("BeginPaint")
	procEndPaint            = user32.NewProc("EndPaint")
	procFillRect            = user32.NewProc("FillRect")
	procCreateSolidBrush    = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject        = gdi32.NewProc("DeleteObject")
	procGetClientRect       = user32.NewProc("GetClientRect")
	procInvalidateRect      = user32.NewProc("InvalidateRect")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
)

const (
	WS_POPUP             = 0x80000000
	WS_EX_TOPMOST        = 0x00000008
	WS_EX_TOOLWINDOW     = 0x00000080
	WS_EX_LAYERED        = 0x00080000
	WS_EX_TRANSPARENT    = 0x00000020
	WS_EX_NOACTIVATE     = 0x08000000

	LWA_COLORKEY = 0x00000001
	LWA_ALPHA    = 0x00000002

	WM_PAINT   = 0x000F
	WM_DESTROY = 0x0002

	SW_SHOW = 5
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type MONITORINFOEXW struct {
	Size    uint32
	Monitor RECT
	Work    RECT
	Flags   uint32
	Device  [32]uint16
}

type WNDCLASSEXW struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   syscall.Handle
	Icon       syscall.Handle
	Cursor     syscall.Handle
	Background syscall.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     syscall.Handle
}

type PAINTSTRUCT struct {
	Hdc         syscall.Handle
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	HWnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Point   POINT
}

type MonitorSettings struct {
	Color   uint32 // BGR
	Enabled bool
	Alpha   byte
	Width   int32
}

type OverlayManager struct {
	windows  map[string]syscall.Handle
	settings map[string]*MonitorSettings
	enabled  bool
}

var manager = &OverlayManager{
	windows:  make(map[string]syscall.Handle),
	settings: make(map[string]*MonitorSettings),
	enabled:  true,
}

func (m *OverlayManager) SetColor(monitorName string, r, g, b uint8) {
	color := uint32(b)<<16 | uint32(g)<<8 | uint32(r)
	if s, ok := m.settings[monitorName]; ok {
		s.Color = color
		m.RedrawMonitor(monitorName)
	}
}

func (m *OverlayManager) SetAlpha(monitorName string, alpha uint8) {
	if s, ok := m.settings[monitorName]; ok {
		s.Alpha = alpha
		if hwnd, ok := m.windows[monitorName]; ok {
			procSetLayeredWindowAttributes.Call(uintptr(hwnd), 0xFF00FF, uintptr(alpha), LWA_COLORKEY|LWA_ALPHA)
		}
	}
}

func (m *OverlayManager) SetWidth(monitorName string, width int32) {
	if s, ok := m.settings[monitorName]; ok {
		s.Width = width
		m.RedrawMonitor(monitorName)
	}
}

func (m *OverlayManager) SetEnabled(enabled bool) {
	m.enabled = enabled
	for name, hwnd := range m.windows {
		if enabled && m.settings[name].Enabled {
			procShowWindow.Call(uintptr(hwnd), SW_SHOW)
		} else {
			procShowWindow.Call(uintptr(hwnd), 0)
		}
	}
}

func (m *OverlayManager) SetMonitorEnabled(name string, enabled bool) {
	if s, ok := m.settings[name]; ok {
		s.Enabled = enabled
		if hwnd, ok := m.windows[name]; ok {
			if enabled && m.enabled {
				procShowWindow.Call(uintptr(hwnd), SW_SHOW)
			} else {
				procShowWindow.Call(uintptr(hwnd), 0)
			}
		}
	}
}

func (m *OverlayManager) RedrawMonitor(name string) {
	if hwnd, ok := m.windows[name]; ok {
		procInvalidateRect.Call(uintptr(hwnd), 0, 1)
	}
}



func WndProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc, _, _ := procBeginPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&ps)))
		if hdc != 0 {
			var rect RECT
			procGetClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))

			// Find monitor color for this HWND
			var s *MonitorSettings
			for name, h := range manager.windows {
				if h == hwnd {
					s = manager.settings[name]
					break
				}
			}

			if s == nil {
				procEndPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&ps)))
				return 0
			}

			// Fill background with colorkey (Magenta 0xFF00FF)
			bgBrush, _, _ := procCreateSolidBrush.Call(0xFF00FF)
			procFillRect.Call(hdc, uintptr(unsafe.Pointer(&rect)), bgBrush)
			procDeleteObject.Call(bgBrush)

			if manager.enabled && s.Enabled {
				ringColor := s.Color
				coreBrush, _, _ := procCreateSolidBrush.Call(uintptr(ringColor))
				
				// Helper to draw border rect
				drawBorder := func(brush uintptr, inset int32) {
					// Top
					r := RECT{inset, inset, rect.Right - inset, inset + 1}
					procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), brush)
					// Bottom
					r = RECT{inset, rect.Bottom - inset - 1, rect.Right - inset, rect.Bottom - inset}
					procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), brush)
					// Left
					r = RECT{inset, inset, inset + 1, rect.Bottom - inset}
					procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), brush)
					// Right
					r = RECT{rect.Right - inset - 1, inset, rect.Right - inset, rect.Bottom - inset}
					procFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), brush)
				}

				// 1. Draw solid core
				for i := int32(0); i < s.Width; i++ {
					drawBorder(coreBrush, i)
				}
				procDeleteObject.Call(coreBrush)
			}




			procEndPaint.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&ps)))
		}
		return 0
	case WM_DESTROY:
		return 0
	default:
		ret, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
		return ret
	}
}


func registerClass() {
	className, _ := syscall.UTF16PtrFromString("OverlayWindowClass")
	
	// Load standard arrow cursor
	hCursor, _, _ := user32.NewProc("LoadCursorW").Call(0, uintptr(32512)) // IDC_ARROW

	var wc WNDCLASSEXW
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.Style = 0x0003 // CS_HREDRAW | CS_VREDRAW
	wc.WndProc = syscall.NewCallback(WndProc)
	wc.ClassName = className
	wc.Instance = 0
	wc.Cursor = syscall.Handle(hCursor)

	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
}


func createOverlayForMonitor(rect RECT, deviceName string, alpha byte) {
	className, _ := syscall.UTF16PtrFromString("OverlayWindowClass")
	title, _ := syscall.UTF16PtrFromString("RingLightOverlay")

	hwnd, _, _ := procCreateWindowExW.Call(
		WS_EX_TOPMOST|WS_EX_LAYERED|WS_EX_TRANSPARENT|WS_EX_NOACTIVATE|WS_EX_TOOLWINDOW,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(title)),
		WS_POPUP,
		uintptr(rect.Left), uintptr(rect.Top),
		uintptr(rect.Right-rect.Left), uintptr(rect.Bottom-rect.Top),
		0, 0, 0, 0,
	)

	if hwnd != 0 {
		procSetLayeredWindowAttributes.Call(hwnd, 0xFF00FF, uintptr(alpha), LWA_COLORKEY|LWA_ALPHA)
		procShowWindow.Call(hwnd, SW_SHOW)
		procUpdateWindow.Call(hwnd)
		manager.windows[deviceName] = syscall.Handle(hwnd)
	}
}



func StartOverlay() {
	go func() {
		runtime.LockOSThread()
		registerClass()

		callback := syscall.NewCallback(func(hMonitor syscall.Handle, hdcMonitor syscall.Handle, lprcMonitor *RECT, dwData uintptr) uintptr {
			var mi MONITORINFOEXW
			mi.Size = uint32(unsafe.Sizeof(mi))
			procGetMonitorInfoW.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(&mi)))
			
			deviceName := syscall.UTF16ToString(mi.Device[:])
			// Clean up device name (e.g., \\.\DISPLAY1 -> DISPLAY1)
			if len(deviceName) > 4 && deviceName[:4] == `\\.\` {
				deviceName = deviceName[4:]
			}
			
			// Initialize settings for this monitor if new
			if _, ok := manager.settings[deviceName]; !ok {
				manager.settings[deviceName] = &MonitorSettings{
					Color:   0x08EEFF, // Default Yellow (BGR)
					Enabled: false,
					Alpha:   200,
					Width:   30,
				}
			}

			s := manager.settings[deviceName]
			createOverlayForMonitor(mi.Work, deviceName, s.Alpha)

			return 1

		})

		procEnumDisplayMonitors.Call(0, 0, callback, 0)


		var msg MSG
		for {
			ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if int32(ret) <= 0 {
				break
			}
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}()
}

