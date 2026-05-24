package ui

import (
	"datapull/api"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SetupUI(myApp fyne.App, FetchDataURA func() ([][]string, error), FetchDataUS30Y func() ([][]string, error), FetchDataPimcoGISIncome func() ([]string, error), FetchDataSP500 func() ([]string, error), FetchOilPrice func() (*api.OilPriceResponse, error), FetchGoldPrice func() (*api.GoldPriceResponse, error), FetchBitcoinPrice func() (*api.BitcoinPriceResponse, error)) fyne.Window {
	myWindow := myApp.NewWindow("Data Pull")

	content := container.NewVBox()

	us30yCardPlaceholder := widget.NewLabel("Loading US30Y data...")
	pimcoCardPlaceholder := widget.NewLabel("Loading Pimco data...")
	sp500CardPlaceholder := widget.NewLabel("Loading S&P 500 data...")
	oilCardPlaceholder := widget.NewLabel("Loading Oil data...")
	goldCardPlaceholder := widget.NewLabel("Loading Gold data...")
	bitcoinCardPlaceholder := widget.NewLabel("Loading Bitcoin data...")
	uraTablePlaceholder := widget.NewLabel("Loading URA data...")

	cardContainer := container.NewHBox(us30yCardPlaceholder, pimcoCardPlaceholder, sp500CardPlaceholder, oilCardPlaceholder, goldCardPlaceholder, bitcoinCardPlaceholder)
	content.Objects = []fyne.CanvasObject{
		cardContainer,
		uraTablePlaceholder,
	}

	refreshButton := widget.NewButton("Refresh", func() {
		cardContainer.Objects = []fyne.CanvasObject{us30yCardPlaceholder, pimcoCardPlaceholder, sp500CardPlaceholder, oilCardPlaceholder, goldCardPlaceholder, bitcoinCardPlaceholder}
		content.Objects = []fyne.CanvasObject{
			cardContainer,
			uraTablePlaceholder,
		}
		cardContainer.Refresh()
		content.Refresh()

		fetchAndDisplayUS30YData(cardContainer, FetchDataUS30Y, us30yCardPlaceholder)
		fetchAndDisplayPimcoData(cardContainer, FetchDataPimcoGISIncome, pimcoCardPlaceholder)
		fetchAndDisplaySP500Data(cardContainer, FetchDataSP500, sp500CardPlaceholder)
		fetchAndDisplayOilData(cardContainer, FetchOilPrice, oilCardPlaceholder)
		fetchAndDisplayGoldData(cardContainer, FetchGoldPrice, goldCardPlaceholder)
		fetchAndDisplayBitcoinData(cardContainer, FetchBitcoinPrice, bitcoinCardPlaceholder)
		fetchAndDisplayURAData(content, FetchDataURA, uraTablePlaceholder)
	})

	myWindow.SetContent(container.NewVBox(
		refreshButton,
		content,
	))

	myWindow.Resize(fyne.NewSize(1600, 900))

	fetchAndDisplayUS30YData(cardContainer, FetchDataUS30Y, us30yCardPlaceholder)
	fetchAndDisplayPimcoData(cardContainer, FetchDataPimcoGISIncome, pimcoCardPlaceholder)
	fetchAndDisplaySP500Data(cardContainer, FetchDataSP500, sp500CardPlaceholder)
	fetchAndDisplayOilData(cardContainer, FetchOilPrice, oilCardPlaceholder)
	fetchAndDisplayGoldData(cardContainer, FetchGoldPrice, goldCardPlaceholder)
	fetchAndDisplayBitcoinData(cardContainer, FetchBitcoinPrice, bitcoinCardPlaceholder)
	fetchAndDisplayURAData(content, FetchDataURA, uraTablePlaceholder)

	return myWindow
}
