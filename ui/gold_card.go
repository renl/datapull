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
	openLabel := widget.NewLabel(formatCommodityPriceLabel("Open", data.Open))
	highLabel := widget.NewLabel(formatCommodityPriceLabel("High", data.High))
	lowLabel := widget.NewLabel(formatCommodityPriceLabel("Low", data.Low))
	volumeLabel := widget.NewLabel(formatCommodityVolumeLabel("Volume", data.Volume))

	return widget.NewCard("Gold Price Data", "", container.NewVBox(closeLabel, openLabel, highLabel, lowLabel, volumeLabel))
}

func formatCommodityPriceLabel(label string, value *float64) string {
	if value == nil {
		return fmt.Sprintf("%s: unavailable", label)
	}

	return fmt.Sprintf("%s: $%.2f", label, *value)
}

func formatCommodityVolumeLabel(label string, value *int64) string {
	if value == nil {
		return fmt.Sprintf("%s: unavailable", label)
	}

	return fmt.Sprintf("%s: %d", label, *value)
}
