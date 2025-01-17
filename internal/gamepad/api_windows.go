// Copyright 2022 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gamepad

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	_DI_OK           = 0
	_DI_NOEFFECT     = _SI_FALSE
	_DI_PROPNOEFFECT = _SI_FALSE

	_DI_DEGREES = 100

	_DI8DEVCLASS_GAMECTRL = 4

	_DIDFT_ABSAXIS     = 0x00000002
	_DIDFT_AXIS        = 0x00000003
	_DIDFT_BUTTON      = 0x0000000C
	_DIDFT_POV         = 0x00000010
	_DIDFT_OPTIONAL    = 0x80000000
	_DIDFT_ANYINSTANCE = 0x00FFFF00

	_DIDOI_ASPECTPOSITION = 0x00000100

	_DIEDFL_ALLDEVICES = 0x00000000

	_DIENUM_STOP     = 0
	_DIENUM_CONTINUE = 1

	_DIERR_INPUTLOST   = windows.SEVERITY_ERROR<<31 | windows.FACILITY_WIN32<<16 | windows.ERROR_READ_FAULT
	_DIERR_NOTACQUIRED = windows.SEVERITY_ERROR<<31 | windows.FACILITY_WIN32<<16 | windows.ERROR_INVALID_ACCESS

	_DIJOFS_X  = uint32(unsafe.Offsetof(diJoyState{}.lX))
	_DIJOFS_Y  = uint32(unsafe.Offsetof(diJoyState{}.lY))
	_DIJOFS_Z  = uint32(unsafe.Offsetof(diJoyState{}.lZ))
	_DIJOFS_RX = uint32(unsafe.Offsetof(diJoyState{}.lRx))
	_DIJOFS_RY = uint32(unsafe.Offsetof(diJoyState{}.lRy))
	_DIJOFS_RZ = uint32(unsafe.Offsetof(diJoyState{}.lRz))

	_DIPH_DEVICE = 0
	_DIPH_BYID   = 2

	_DIPROP_AXISMODE = 2
	_DIPROP_RANGE    = 4

	_DIPROPAXISMODE_ABS = 0

	_DIRECTINPUT_VERSION = 0x0800

	_GWL_WNDPROC = -4

	_MAX_PATH = 260

	_RIDI_DEVICEINFO = 0x2000000b
	_RIDI_DEVICENAME = 0x20000007

	_RIM_TYPEHID = 2

	_SI_FALSE = 1

	_WM_DEVICECHANGE = 0x0219

	_XINPUT_CAPS_WIRELESS = 0x0002

	_XINPUT_DEVSUBTYPE_GAMEPAD      = 0x01
	_XINPUT_DEVSUBTYPE_WHEEL        = 0x02
	_XINPUT_DEVSUBTYPE_ARCADE_STICK = 0x03
	_XINPUT_DEVSUBTYPE_FLIGHT_STICK = 0x04
	_XINPUT_DEVSUBTYPE_DANCE_PAD    = 0x05
	_XINPUT_DEVSUBTYPE_GUITAR       = 0x06
	_XINPUT_DEVSUBTYPE_DRUM_KIT     = 0x08

	_XINPUT_GAMEPAD_DPAD_UP        = 0x0001
	_XINPUT_GAMEPAD_DPAD_DOWN      = 0x0002
	_XINPUT_GAMEPAD_DPAD_LEFT      = 0x0004
	_XINPUT_GAMEPAD_DPAD_RIGHT     = 0x0008
	_XINPUT_GAMEPAD_START          = 0x0010
	_XINPUT_GAMEPAD_BACK           = 0x0020
	_XINPUT_GAMEPAD_LEFT_THUMB     = 0x0040
	_XINPUT_GAMEPAD_RIGHT_THUMB    = 0x0080
	_XINPUT_GAMEPAD_LEFT_SHOULDER  = 0x0100
	_XINPUT_GAMEPAD_RIGHT_SHOULDER = 0x0200
	_XINPUT_GAMEPAD_A              = 0x1000
	_XINPUT_GAMEPAD_B              = 0x2000
	_XINPUT_GAMEPAD_X              = 0x4000
	_XINPUT_GAMEPAD_Y              = 0x8000
)

func diDftGetType(n uint32) byte {
	return byte(n)
}

func diJofsSlider(n int) uint32 {
	return uint32(unsafe.Offsetof(diJoyState{}.rglSlider) + uintptr(n)*unsafe.Sizeof(int32(0)))
}

func diJofsPOV(n int) uint32 {
	return uint32(unsafe.Offsetof(diJoyState{}.rgdwPOV) + uintptr(n)*unsafe.Sizeof(uint32(0)))
}

func diJofsButton(n int) uint32 {
	return uint32(unsafe.Offsetof(diJoyState{}.rgbButtons) + uintptr(n))
}

var (
	iidIDirectInput8W = windows.GUID{0xbf798031, 0x483a, 0x4da2, [...]byte{0xaa, 0x99, 0x5d, 0x64, 0xed, 0x36, 0x97, 0x00}}
	guidXAxis         = windows.GUID{0xa36d02e0, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
	guidYAxis         = windows.GUID{0xa36d02e1, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
	guidZAxis         = windows.GUID{0xa36d02e2, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
	guidRxAxis        = windows.GUID{0xa36d02f4, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
	guidRyAxis        = windows.GUID{0xa36d02f5, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
	guidRzAxis        = windows.GUID{0xa36d02e3, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
	guidSlider        = windows.GUID{0xa36d02e4, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
	guidPOV           = windows.GUID{0xa36d02f2, 0xc9f3, 0x11cf, [...]byte{0xbf, 0xc7, 0x44, 0x45, 0x53, 0x54, 0x00, 0x00}}
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	user32   = windows.NewLazySystemDLL("user32.dll")

	procGetCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	procGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")

	procCallWindowProcW        = user32.NewProc("CallWindowProcW")
	procGetActiveWindow        = user32.NewProc("GetActiveWindow")
	procGetRawInputDeviceInfoW = user32.NewProc("GetRawInputDeviceInfoW")
	procGetRawInputDeviceList  = user32.NewProc("GetRawInputDeviceList")
	procSetWindowLongPtrW      = user32.NewProc("SetWindowLongPtrW")
)

func getCurrentThreadId() uint32 {
	t, _, _ := procGetCurrentThreadId.Call()
	return uint32(t)
}

func getModuleHandleW() (uintptr, error) {
	m, _, e := procGetModuleHandleW.Call(0)
	if m == 0 {
		if e != nil && e != windows.ERROR_SUCCESS {
			return 0, fmt.Errorf("gamepad: GetModuleHandleW failed: %w", e)
		}
		return 0, fmt.Errorf("gamepad: GetModuleHandleW returned 0")
	}
	return m, nil
}

func callWindowProcW(lpPrevWndFunc uintptr, hWnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	r, _, _ := procCallWindowProcW.Call(lpPrevWndFunc, hWnd, uintptr(msg), wParam, lParam)
	return r
}

func getActiveWindow() uintptr {
	h, _, _ := procGetActiveWindow.Call()
	return h
}

func getRawInputDeviceInfoW(hDevice windows.Handle, uiCommand uint32, pData unsafe.Pointer, pcb *uint32) (uint32, error) {
	r, _, e := procGetRawInputDeviceInfoW.Call(uintptr(hDevice), uintptr(uiCommand), uintptr(pData), uintptr(unsafe.Pointer(pcb)))
	if uint32(r) == ^uint32(0) {
		if e != nil && e != windows.ERROR_SUCCESS {
			return 0, fmt.Errorf("gamepad: GetRawInputDeviceInfoW failed: %w", e)
		}
		return 0, fmt.Errorf("gamepad: GetRawInputDeviceInfoW returned -1")
	}
	return uint32(r), nil
}

func getRawInputDeviceList(pRawInputDeviceList *rawInputDeviceList, puiNumDevices *uint32) (uint32, error) {
	r, _, e := procGetRawInputDeviceList.Call(uintptr(unsafe.Pointer(pRawInputDeviceList)), uintptr(unsafe.Pointer(puiNumDevices)), unsafe.Sizeof(rawInputDeviceList{}))
	if uint32(r) == ^uint32(0) {
		if e != nil && e != windows.ERROR_SUCCESS {
			return 0, fmt.Errorf("gamepad: GetRawInputDeviceList failed: %w", e)
		}
		return 0, fmt.Errorf("gamepad: GetRawInputDeviceList returned -1")
	}
	return uint32(r), nil
}

func setWindowLongPtrW(hWnd uintptr, nIndex int32, dwNewLong uintptr) (uintptr, error) {
	h, _, e := procSetWindowLongPtrW.Call(hWnd, uintptr(nIndex), dwNewLong)
	if h == 0 {
		if e != nil && e != windows.ERROR_SUCCESS {
			return 0, fmt.Errorf("gamepad: SetWindowLongPtrW failed: %w", e)
		}
		return 0, fmt.Errorf("gamepad: SetWindowLongPtrW returned 0")
	}
	return h, nil
}

type directInputError uint32

func (d directInputError) Error() string {
	return fmt.Sprintf("DirectInput error: %d", d)
}

type diDataFormat struct {
	dwSize     uint32
	dwObjSize  uint32
	dwFlags    uint32
	dwDataSize uint32
	dwNumObjs  uint32
	rgodf      *diObjectDataFormat
}

type diDevCaps struct {
	dwSize                uint32
	dwFlags               uint32
	dwDevType             uint32
	dwAxes                uint32
	dwButtons             uint32
	dwPOVs                uint32
	dwFFSamplePeriod      uint32
	dwFFMinTimeResolution uint32
	dwFirmwareRevision    uint32
	dwHardwareRevision    uint32
	dwFFDriverVersion     uint32
}

type diDeviceInstanceW struct {
	dwSize          uint32
	guidInstance    windows.GUID
	guidProduct     windows.GUID
	dwDevType       uint32
	tszInstanceName [_MAX_PATH]uint16
	tszProductName  [_MAX_PATH]uint16
	guidFFDriver    windows.GUID
	wUsagePage      uint16
	wUsage          uint16
}

type diDeviceObjectInstanceW struct {
	dwSize              uint32
	guidType            windows.GUID
	dwOfs               uint32
	dwType              uint32
	dwFlags             uint32
	tszName             [_MAX_PATH]uint16
	dwFFMaxForce        uint32
	dwFFForceResolution uint32
	wCollectionNumber   uint16
	wDesignatorIndex    uint16
	wUsagePage          uint16
	wUsage              uint16
	dwDimension         uint32
	wExponent           uint16
	wReserved           uint16
}

type diJoyState struct {
	lX         int32
	lY         int32
	lZ         int32
	lRx        int32
	lRy        int32
	lRz        int32
	rglSlider  [2]int32
	rgdwPOV    [4]uint32
	rgbButtons [32]byte
}

type diObjectDataFormat struct {
	pguid   *windows.GUID
	dwOfs   uint32
	dwType  uint32
	dwFlags uint32
}

type diPropDword struct {
	diph   diPropHeader
	dwData uint32
}

type diPropHeader struct {
	dwSize       uint32
	dwHeaderSize uint32
	dwObj        uint32
	dwHow        uint32
}

type diPropRange struct {
	diph diPropHeader
	lMin int32
	lMax int32
}

type iDirectInput8W struct {
	vtbl *iDirectInput8W_Vtbl
}

type iDirectInput8W_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	CreateDevice           uintptr
	EnumDevices            uintptr
	GetDeviceStatus        uintptr
	RunControlPanel        uintptr
	Initialize             uintptr
	FindDevice             uintptr
	EnumDevicesBySemantics uintptr
	ConfigureDevices       uintptr
}

func (d *iDirectInput8W) CreateDevice(rguid *windows.GUID, lplpDirectInputDevice **iDirectInputDevice8W, pUnkOuter unsafe.Pointer) error {
	r, _, _ := syscall.Syscall6(d.vtbl.CreateDevice, 4,
		uintptr(unsafe.Pointer(d)),
		uintptr(unsafe.Pointer(rguid)), uintptr(unsafe.Pointer(lplpDirectInputDevice)), uintptr(pUnkOuter),
		0, 0)
	if r != _DI_OK {
		return fmt.Errorf("gamepad: IDirectInput8::CreateDevice failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

func (d *iDirectInput8W) EnumDevices(dwDevType uint32, lpCallback uintptr, pvRef unsafe.Pointer, dwFlags uint32) error {
	r, _, _ := syscall.Syscall6(d.vtbl.EnumDevices, 5,
		uintptr(unsafe.Pointer(d)),
		uintptr(dwDevType), lpCallback, uintptr(pvRef), uintptr(dwFlags),
		0)
	if r != _DI_OK {
		return fmt.Errorf("gamepad: IDirectInput8::EnumDevices failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

type iDirectInputDevice8W struct {
	vtbl *iDirectInputDevice8W_Vtbl
}

type iDirectInputDevice8W_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetCapabilities          uintptr
	EnumObjects              uintptr
	GetProperty              uintptr
	SetProperty              uintptr
	Acquire                  uintptr
	Unacquire                uintptr
	GetDeviceState           uintptr
	GetDeviceData            uintptr
	SetDataFormat            uintptr
	SetEventNotification     uintptr
	SetCooperativeLevel      uintptr
	GetObjectInfo            uintptr
	GetDeviceInfo            uintptr
	RunControlPanel          uintptr
	Initialize               uintptr
	CreateEffect             uintptr
	EnumEffects              uintptr
	GetEffectInfo            uintptr
	GetForceFeedbackState    uintptr
	SendForceFeedbackCommand uintptr
	EnumCreatedEffectObjects uintptr
	Escape                   uintptr
	Poll                     uintptr
	SendDeviceData           uintptr
	EnumEffectsInFile        uintptr
	WriteEffectToFile        uintptr
	BuildActionMap           uintptr
	SetActionMap             uintptr
	GetImageInfo             uintptr
}

func (d *iDirectInputDevice8W) Acquire() error {
	r, _, _ := syscall.Syscall(d.vtbl.Acquire, 1, uintptr(unsafe.Pointer(d)), 0, 0)
	if r != _DI_OK && r != _SI_FALSE {
		return fmt.Errorf("gamepad: IDirectInputDevice8::Acquire failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

func (d *iDirectInputDevice8W) EnumObjects(lpCallback uintptr, pvRef unsafe.Pointer, dwFlags uint32) error {
	r, _, _ := syscall.Syscall6(d.vtbl.EnumObjects, 4,
		uintptr(unsafe.Pointer(d)),
		lpCallback, uintptr(pvRef), uintptr(dwFlags),
		0, 0)
	if r != _DI_OK {
		return fmt.Errorf("gamepad: IDirectInputDevice8::EnumObjects failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

func (d *iDirectInputDevice8W) GetCapabilities(lpDIDevCaps *diDevCaps) error {
	r, _, _ := syscall.Syscall(d.vtbl.GetCapabilities, 2, uintptr(unsafe.Pointer(d)), uintptr(unsafe.Pointer(lpDIDevCaps)), 0)
	if r != _DI_OK {
		return fmt.Errorf("gamepad: IDirectInputDevice8::GetCapabilities failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

func (d *iDirectInputDevice8W) GetDeviceState(cbData uint32, lpvData unsafe.Pointer) error {
	r, _, _ := syscall.Syscall(d.vtbl.GetDeviceState, 3, uintptr(unsafe.Pointer(d)), uintptr(cbData), uintptr(lpvData))
	if r != _DI_OK {
		return fmt.Errorf("gamepad: IDirectInputDevice8::GetDeviceState failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

func (d *iDirectInputDevice8W) Poll() error {
	r, _, _ := syscall.Syscall(d.vtbl.Poll, 1, uintptr(unsafe.Pointer(d)), 0, 0)
	if r != _DI_OK && r != _DI_NOEFFECT {
		return fmt.Errorf("gamepad: IDirectInputDevice8::Poll failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

func (d *iDirectInputDevice8W) Release() uint32 {
	r, _, _ := syscall.Syscall(d.vtbl.Release, 1, uintptr(unsafe.Pointer(d)), 0, 0)
	return uint32(r)
}

func (d *iDirectInputDevice8W) SetDataFormat(lpdf *diDataFormat) error {
	r, _, _ := syscall.Syscall(d.vtbl.SetDataFormat, 2, uintptr(unsafe.Pointer(d)), uintptr(unsafe.Pointer(lpdf)), 0)
	if r != _DI_OK {
		return fmt.Errorf("gamepad: IDirectInputDevice8::SetDataFormat failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

func (d *iDirectInputDevice8W) SetProperty(rguidProp uintptr, pdiph *diPropHeader) error {
	r, _, _ := syscall.Syscall(d.vtbl.SetProperty, 3, uintptr(unsafe.Pointer(d)), rguidProp, uintptr(unsafe.Pointer(pdiph)))
	if r != _DI_OK && r != _DI_PROPNOEFFECT {
		return fmt.Errorf("gamepad: IDirectInputDevice8::SetProperty failed: %w", directInputError(syscall.Errno(r)))
	}
	return nil
}

type ridDeviceInfo struct {
	cbSize uint32
	dwType uint32
	hid    ridDeviceInfoHID // Originally, this member is a union.
}

type ridDeviceInfoHID struct {
	dwVendorId      uint32
	dwProductId     uint32
	dwVersionNumber uint32
	usUsagePage     uint16
	usUsage         uint16
	_               uint32 // A padding adjusting with the size of RID_DEVICE_INFO_KEYBOARD
	_               uint32 // A padding adjusting with the size of RID_DEVICE_INFO_KEYBOARD
}

type rawInputDeviceList struct {
	hDevice windows.Handle
	dwType  uint32
}

type xinputCapabilities struct {
	typ       byte
	subType   byte
	flags     uint16
	gamepad   xinputGamepad
	vibration xinputVibration
}

type xinputGamepad struct {
	wButtons      uint16
	bLeftTrigger  byte
	bRightTrigger byte
	sThumbLX      int16
	sThumbLY      int16
	sThumbRX      int16
	sThumbRY      int16
}

type xinputState struct {
	dwPacketNumber uint32
	Gamepad        xinputGamepad
}

type xinputVibration struct {
	wLeftMotorSpeed  uint16
	wRightMotorSpeed uint16
}
