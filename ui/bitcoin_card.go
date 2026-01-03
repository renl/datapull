package ui

import (
	"datapull/api"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func fetchAndDisplayBitcoinData(content *fyne.Container, FetchBitcoinPrice func() (*api.BitcoinPriceResponse, error), placeholder fyne.CanvasObject) {
	go func() {
		bitcoinData, err := FetchBitcoinPrice()
		if err != nil {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel(fmt.Sprintf("Error fetching bitcoin data: %v", err)))
		} else {
			bitcoinCard := renderBitcoinCard(bitcoinData)
			content.Objects = replaceObject(content.Objects, placeholder, bitcoinCard)
		}
		content.Refresh()
	}()
}

func renderBitcoinCard(data *api.BitcoinPriceResponse) *widget.Card {
	closeLabel := widget.NewLabel(fmt.Sprintf("Close: $%.2f", data.Close))
	openLabel := widget.NewLabel(fmt.Sprintf("Open: $%.2f", data.Open))
	highLabel := widget.NewLabel(fmt.Sprintf("High: $%.2f", data.High))
	lowLabel := widget.NewLabel(fmt.Sprintf("Low: $%.2f", data.Low))
	volumeLabel := widget.NewLabel(fmt.Sprintf("Volume: %d", data.Volume))

	return widget.NewCard("Bitcoin Price Data", "", container.NewVBox(closeLabel, openLabel, highLabel, lowLabel, volumeLabel))
}
