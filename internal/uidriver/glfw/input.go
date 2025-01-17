// Copyright 2015 Hajime Hoshi
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

//go:build !android && !js && !ios
// +build !android,!js,!ios

package glfw

import (
	"math"
	"sync"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2/internal/driver"
	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
)

type gamepad struct {
	valid         bool
	guid          string
	name          string
	axisNum       int
	axes          [16]float64
	buttonNum     int
	buttonPressed [256]bool
	hatsNum       int
	hats          [16]int
}

type Input struct {
	keyPressed         map[glfw.Key]bool
	mouseButtonPressed map[glfw.MouseButton]bool
	onceCallback       sync.Once
	scrollX            float64
	scrollY            float64
	cursorX            int
	cursorY            int
	gamepads           [16]gamepad
	touches            map[driver.TouchID]pos // TODO: Implement this (#417)
	runeBuffer         []rune
	ui                 *UserInterface

	nativeGamepads
}

type pos struct {
	X int
	Y int
}

func (i *Input) CursorPosition() (x, y int) {
	if !i.ui.isRunning() {
		return 0, 0
	}

	i.ui.m.RLock()
	defer i.ui.m.RUnlock()
	return i.cursorX, i.cursorY
}

func (i *Input) AppendTouchIDs(touchIDs []driver.TouchID) []driver.TouchID {
	if !i.ui.isRunning() {
		return nil
	}

	i.ui.m.RLock()
	defer i.ui.m.RUnlock()
	for id := range i.touches {
		touchIDs = append(touchIDs, id)
	}
	return touchIDs
}

func (i *Input) TouchPosition(id driver.TouchID) (x, y int) {
	if !i.ui.isRunning() {
		return 0, 0
	}

	i.ui.m.RLock()
	defer i.ui.m.RUnlock()
	for tid, pos := range i.touches {
		if id == tid {
			return pos.X, pos.Y
		}
	}
	return 0, 0
}

func (i *Input) AppendInputChars(runes []rune) []rune {
	if !i.ui.isRunning() {
		return nil
	}

	i.ui.m.RLock()
	defer i.ui.m.RUnlock()
	return append(runes, i.runeBuffer...)
}

func (i *Input) resetForFrame() {
	if !i.ui.isRunning() {
		return
	}

	i.ui.m.Lock()
	defer i.ui.m.Unlock()
	i.runeBuffer = i.runeBuffer[:0]
	i.scrollX, i.scrollY = 0, 0
}

func (i *Input) IsKeyPressed(key driver.Key) bool {
	if !i.ui.isRunning() {
		return false
	}

	i.ui.m.Lock()
	defer i.ui.m.Unlock()
	if i.keyPressed == nil {
		i.keyPressed = map[glfw.Key]bool{}
	}
	gk, ok := driverKeyToGLFWKey[key]
	return ok && i.keyPressed[gk]
}

func (i *Input) IsMouseButtonPressed(button driver.MouseButton) bool {
	if !i.ui.isRunning() {
		return false
	}

	i.ui.m.Lock()
	defer i.ui.m.Unlock()
	if i.mouseButtonPressed == nil {
		i.mouseButtonPressed = map[glfw.MouseButton]bool{}
	}
	for gb, b := range glfwMouseButtonToMouseButton {
		if b != button {
			continue
		}
		if i.mouseButtonPressed[gb] {
			return true
		}
	}
	return false
}

func (i *Input) Wheel() (xoff, yoff float64) {
	if !i.ui.isRunning() {
		return 0, 0
	}

	i.ui.m.RLock()
	defer i.ui.m.RUnlock()
	return i.scrollX, i.scrollY
}

var glfwMouseButtonToMouseButton = map[glfw.MouseButton]driver.MouseButton{
	glfw.MouseButtonLeft:   driver.MouseButtonLeft,
	glfw.MouseButtonRight:  driver.MouseButtonRight,
	glfw.MouseButtonMiddle: driver.MouseButtonMiddle,
}

// update must be called from the main thread.
func (i *Input) update(window *glfw.Window, context driver.UIContext) error {
	i.ui.m.Lock()
	defer i.ui.m.Unlock()

	i.onceCallback.Do(func() {
		window.SetCharModsCallback(glfw.ToCharModsCallback(func(w *glfw.Window, char rune, mods glfw.ModifierKey) {
			// As this function is called from GLFW callbacks, the current thread is main.
			if !unicode.IsPrint(char) {
				return
			}

			i.ui.m.Lock()
			defer i.ui.m.Unlock()
			i.runeBuffer = append(i.runeBuffer, char)
		}))
		window.SetScrollCallback(glfw.ToScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
			// As this function is called from GLFW callbacks, the current thread is main.
			i.ui.m.Lock()
			defer i.ui.m.Unlock()
			i.scrollX = xoff
			i.scrollY = yoff
		}))
	})
	if i.keyPressed == nil {
		i.keyPressed = map[glfw.Key]bool{}
	}
	for gk := range glfwKeyToDriverKey {
		i.keyPressed[gk] = window.GetKey(gk) == glfw.Press
	}
	if i.mouseButtonPressed == nil {
		i.mouseButtonPressed = map[glfw.MouseButton]bool{}
	}
	for gb := range glfwMouseButtonToMouseButton {
		i.mouseButtonPressed[gb] = window.GetMouseButton(gb) == glfw.Press
	}
	cx, cy := window.GetCursorPos()
	// TODO: This is tricky. Rename the function?
	m := i.ui.currentMonitor()
	s := i.ui.deviceScaleFactor(m)
	cx = i.ui.dipFromGLFWPixel(cx, m)
	cy = i.ui.dipFromGLFWPixel(cy, m)
	cx, cy = context.AdjustPosition(cx, cy, s)

	// AdjustPosition can return NaN at the initialization.
	if !math.IsNaN(cx) && !math.IsNaN(cy) {
		i.cursorX, i.cursorY = int(cx), int(cy)
	}

	return i.updateGamepads()
}
