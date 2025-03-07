package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/furkandogmus/HaberlerPlus/pkg/sources"
	"github.com/furkandogmus/HaberlerPlus/pkg/utils"
)

func main() {
	showVersion := flag.Bool("v", false, "Versiyon bilgisini göster!")
	showHelp := flag.Bool("h", false, "Yardım bilgisini göster!")
	flag.Parse()

	if *showVersion {
		fmt.Printf("HaberlerPlus Versiyon %s\n", utils.Version)
		return
	}
	if *showHelp {
		fmt.Println("CLI Haber Bülteni Plus")
		fmt.Println("Çeşitli haber kaynaklarından haber başlıklarını ve linklerini gösterir.")
		fmt.Println()
		fmt.Println("Seçenekler:")
		fmt.Println("-h  yardım bilgisini verir.")
		fmt.Println("-v  versiyon bilgisini verir.")
		return
	}

	// Get all available news sources
	allSources := sources.GetAllSources()

	// Display available news sources
	fmt.Printf("%sHaber Kaynakları:\n", utils.Cyan)
	for i, source := range allSources {
		fmt.Printf("%s%d. %s%s\n", utils.Yellow, i+1, source.Name(), utils.Reset)
	}

	// Get user input for news source
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%sHaber kaynağı numarası girin: %s", utils.Cyan, utils.Reset)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input = strings.TrimSpace(input)
	sourceNum, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal("Lütfen geçerli bir sayı girin.")
	}

	if sourceNum < 1 || sourceNum > len(allSources) {
		fmt.Println("Geçersiz haber kaynağı numarası.")
		return
	}

	// Get the selected news source
	selectedSource := allSources[sourceNum-1]

	// Display categories for the selected source
	categories := selectedSource.Categories()
	fmt.Printf("%s%s Kategorileri:\n", utils.Cyan, selectedSource.Name())
	for i, cat := range categories {
		fmt.Printf("%s%d. %s%s\n", utils.Yellow, i+1, cat, utils.Reset)
	}

	// Get user input for category
	fmt.Printf("%sKategori numarası girin: %s", utils.Cyan, utils.Reset)
	input, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	input = strings.TrimSpace(input)
	categoryNum, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal("Lütfen geçerli bir sayı girin.")
	}

	if categoryNum < 1 || categoryNum > len(categories) {
		fmt.Println("Geçersiz kategori numarası.")
		return
	}

	// Fetch news for the selected category
	newsItems, err := selectedSource.FetchNews(categoryNum - 1)
	if err != nil {
		log.Fatal(err)
	}

	// Display the news
	fmt.Printf("%s%s - %s kategorisinden haberler:%s\n", utils.Green, selectedSource.Name(), categories[categoryNum-1], utils.Reset)

	if len(newsItems) == 0 {
		fmt.Println("Bu kategoride haber bulunamadı.")
		return
	}

	for _, item := range newsItems {
		
			fmt.Printf("%s%s: %s", utils.Red, item.Title, utils.Reset)
			fmt.Printf("%s%s%s\n", utils.Gray, item.URL, utils.Reset)
		
	}
}