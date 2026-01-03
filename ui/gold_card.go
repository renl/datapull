package ui

import (
	"datapull/api"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func fetchAndDisplayGoldData(content *fyne.Container, FetchGoldPrice func() (*api.GoldPriceResponse, error), placeholder fyne.CanvasObject) {
	go func() {
		goldData, err := FetchGoldPrice()
		if err != nil {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel(fmt.Sprintf("Error fetching gold data: %v", err)))
		} else {
			goldCard := renderGoldCard(goldData)
			content.Objects = replaceObject(content.Objects, placeholder, goldCard)
		}
		content.Refresh()
	}()
}

func renderGoldCard(data *api.GoldPriceResponse) *widget.Card {
	closeLabel := widget.NewLabel(fmt.Sprintf("Close: $%.2f", data.Close))
	openLabel := widget.NewLabel(fmt.Sprintf("Open: $%.2f", data.Open))
	highLabel := widget.NewLabel(fmt.Sprintf("High: $%.2f", data.High))
	lowLabel := widget.NewLabel(fmt.Sprintf("Low: $%.2f", data.Low))
	volumeLabel := widget.NewLabel(fmt.Sprintf("Volume: %d", data.Volume))

	return widget.NewCard("Gold Price Data", "", container.NewVBox(closeLabel, openLabel, highLabel, lowLabel, volumeLabel))
}
