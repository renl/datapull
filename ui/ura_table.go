package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func renderURATable(data [][]string) *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return len(data), len(data[0]) },
		func() fyne.CanvasObject {
			return container.NewStack(canvas.NewRectangle(theme.Color(theme.ColorNameBackground)), widget.NewLabel(""))
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			stack := o.(*fyne.Container)
			label := stack.Objects[1].(*widget.Label)
			label.SetText(data[i.Row][i.Col])
			if i.Row > 0 && data[i.Row][2] == "914.94" { // Highlight cells where AREA (SQFT) is 914.94
				stack.Objects[0].(*canvas.Rectangle).FillColor = theme.Color(theme.ColorNameSelection)
			} else if i.Row > 0 && data[i.Row][2] == "365.98" {
				stack.Objects[0].(*canvas.Rectangle).FillColor = theme.Color(theme.ColorNameSuccess)
			} else {
				stack.Objects[0].(*canvas.Rectangle).FillColor = theme.Color(theme.ColorNameBackground)
			}
			stack.Objects[0].Refresh()
		},
	)

	// Set column widths dynamically
	for i := 0; i < len(data[0]); i++ {
		maxWidth := float32(0)
		for j := 0; j < len(data); j++ {
			width := fyne.MeasureText(data[j][i], theme.TextSize(), fyne.TextStyle{}).Width
			if width > maxWidth {
				maxWidth = width
			}
		}
		table.SetColumnWidth(i, maxWidth+20) // Add some padding
	}

	// Set row heights
	for i := 0; i < len(data); i++ {
		table.SetRowHeight(i, 30)
	}

	return table
}

func fetchAndDisplayURAData(content *fyne.Container, FetchDataURA func() ([][]string, error), placeholder fyne.CanvasObject) {
	go func() {
		uraData, err := FetchDataURA()
		if err != nil {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel("Error fetching URA data"))
		} else if len(uraData) <= 1 {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel("No URA data found"))
		} else {
			table := renderURATable(uraData)
			scrollContainer := container.NewScroll(table)
			scrollContainer.SetMinSize(fyne.NewSize(1600, 600)) // Set the minimum size to show at least 20 rows
			content.Objects = replaceObject(content.Objects, placeholder, scrollContainer)
		}
		content.Refresh()
	}()
}
