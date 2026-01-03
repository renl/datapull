# DataPull

A desktop application built with Go and Fyne that fetches and displays various financial and property market data.

## Features

- **Property Data**: Singapore URA residential transaction data for tracked properties
- **US 30-Year Treasury**: Current US 30-year treasury yield
- **Pimco GIS Income Fund**: Fund price data
- **S&P 500**: Index price data
- **Crude Oil**: Oil futures price (CL=F)
- **Gold**: Gold futures price (GC=F)
- **Bitcoin**: BTC-USD price

## Screenshot

The application displays price cards for financial instruments and a table for property transactions.

## Requirements

- Go 1.23+
- GCC (for CGO, required by Fyne)

## Build

```bash
go build -ldflags="-s -w" -o datapull.exe .
```

## Run

```bash
./datapull.exe
```

## Data Sources

| Data | Source |
|------|--------|
| Property Transactions | Singapore URA |
| US 30Y Treasury | CNBC |
| Pimco GIS Income | FT Markets |
| S&P 500 | Yahoo Finance |
| Crude Oil | Yahoo Finance |
| Gold | Yahoo Finance |
| Bitcoin | Yahoo Finance |

## License

MIT
