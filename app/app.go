package app

import (
	"datapull/api"
	"datapull/ui"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func Run() {
	myApp := app.New()
	myWindow := ui.SetupUI(myApp, FetchDataURA, FetchDataUS30Y, FetchDataPimcoGISIncome, FetchDataSP500, FetchOilPrice, FetchGoldPrice, FetchBitcoinPrice)

	// Load the icon file
	iconPath := "app.ico" // Path to your .ico file
	iconFile, err := os.Open(iconPath)
	if err != nil {
		panic(err)
	}
	defer iconFile.Close()

	iconResource := fyne.NewStaticResource("app.ico", readFile(iconFile))
	myApp.SetIcon(iconResource)    // Set app-wide icon
	myWindow.SetIcon(iconResource) // Set the window-specific icon

	myWindow.ShowAndRun()
}

func FetchDataURA() ([][]string, error) {
	return api.FetchDataURA()
}

func FetchDataUS30Y() ([][]string, error) {
	return api.FetchDataUS30Y()
}

func FetchDataPimcoGISIncome() ([]string, error) {
	return api.FetchDataPimcoGISIncome()
}

func FetchDataSP500() ([]string, error) {
	return api.FetchDataSP500()
}

func FetchOilPrice() (*api.OilPriceResponse, error) {
	return api.FetchOilPrice()
}

func FetchGoldPrice() (*api.GoldPriceResponse, error) {
	return api.FetchGoldPrice()
}

func FetchBitcoinPrice() (*api.BitcoinPriceResponse, error) {
	return api.FetchBitcoinPrice()
}

func readFile(file *os.File) []byte {
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}

	return data
}
