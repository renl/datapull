package ui

import "fyne.io/fyne/v2"

func replaceObject(objects []fyne.CanvasObject, oldObject, newObject fyne.CanvasObject) []fyne.CanvasObject {
	for i, obj := range objects {
		if obj == oldObject {
			objects[i] = newObject
			return objects
		}
	}
	return objects
}
