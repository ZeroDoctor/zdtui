package main

import "github.com/zerodoctor/zdtui/data"

type ScreenManager struct {
	vm      *ViewManager
	screens []data.IView
}
