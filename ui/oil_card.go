package ui

import (
	"datapull/api"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func fetchAndDisplayOilData(content *fyne.Container, FetchOilPrice func() (*api.OilPriceResponse, error), placeholder fyne.CanvasObject) {
	go func() {
		oilData, err := FetchOilPrice()
		if err != nil {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel(fmt.Sprintf("Error fetching oil data: %v", err)))
		} else {
			oilCard := renderOilCard(oilData)
			content.Objects = replaceObject(content.Objects, placeholder, oilCard)
		}
		content.Refresh()
	}()
}

func renderOilCard(data *api.OilPriceResponse) *widget.Card {
	closeLabel := widget.NewLabel(fmt.Sprintf("Close: $%.2f", data.Close))
	openLabel := widget.NewLabel(formatCommodityPriceLabel("Open", data.Open))
	highLabel := widget.NewLabel(formatCommodityPriceLabel("High", data.High))
	lowLabel := widget.NewLabel(formatCommodityPriceLabel("Low", data.Low))
	volumeLabel := widget.NewLabel(formatCommodityVolumeLabel("Volume", data.Volume))

	return widget.NewCard("Oil Price Data", "", container.NewVBox(closeLabel, openLabel, highLabel, lowLabel, volumeLabel))
}
