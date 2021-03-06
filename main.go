package main

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Data struct {
	Key   string
	Value float64
}

type DataList []Data

func (d DataList) Len() int {
	return len(d)
}

func (d DataList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DataList) Less(i, j int) bool {
	return d[i].Value < d[j].Value
}

func main() {
	// CSV Reader
	file, err := os.Open("games.csv")
	var gameNames []string
	var sales []float64
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	gameNames, sales = formatRecords(records)
	sortedData := mapData(gameNames, sales)
	createChart(sortedData)
}

func formatRecords(records [][]string) ([]string, []float64) {
	gameNames := []string{}
	sales := []float64{}
	for _, r := range records[1:] {
		gameNames = append(gameNames, r[0])
		val, _ := strconv.ParseFloat(r[1], 64)
		sales = append(sales, val)
	}
	return gameNames, sales
}

func mapData(gameNames []string, sales []float64) DataList {
	dataMap := map[string]float64{}
	for index, value := range gameNames {
		dataMap[value] = sales[index]
	}
	data := make(DataList, len(dataMap))
	iterator := 0
	for k, v := range dataMap {
		data[iterator] = Data{k, v}
		iterator++
	}
	sort.Sort(sort.Reverse(data))
	return data
}

func createChart(sortedData DataList) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "PC Games Sales",
		Subtitle: "Best selling PC games",
	}))
	bar.SetXAxis([]string{
		sortedData[0].Key[:4], sortedData[1].Key[:4], sortedData[2].Key[:4], sortedData[3].Key[:4],
		sortedData[4].Key[:4], sortedData[5].Key[:4], sortedData[6].Key[:4],
		sortedData[7].Key[:4], sortedData[8].Key[:4], sortedData[9].Key[:4],
	}).AddSeries("Values", generateBarItems(sortedData))
	f, _ := os.Create("games.html")
	bar.Render(f)
}

func generateBarItems(data DataList) []opts.BarData {
	barData := []float64{}
	items := make([]opts.BarData, 0)
	for i := 0; i < 10; i++ {
		barData = append(barData, data[i].Value)
	}
	for _, v := range barData {
		items = append(items, opts.BarData{Value: v})
	}
	return items
}
