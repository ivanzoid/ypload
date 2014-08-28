package main

import (
	"./config"
	"./yfotki"
	"./ylogin"
	"flag"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"os"
	"path"
)

const (
	kAppId               = "e2b26273dab84121bf3f9c2be4bb8915"
	kLocalHttpServerPort = 30171
)

var (
	flagAllSizes bool
)

func openLoginPage(appId string) {
	urlString := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=token&client_id=%v", appId)
	err := open.Start(urlString)
	if err != nil {
		fmt.Printf("Can't open browser: %v\n", err)
		os.Exit(1)
	}
}

func getTokenData() ylogin.TokenData {
	tokenDataChan := make(chan ylogin.TokenData)
	ylogin.Login(kLocalHttpServerPort, tokenDataChan)

	openLoginPage(kAppId)

	tokenData := <-tokenDataChan

	return tokenData
}

func init() {
	const (
		allSizesDefault = false
		allSizesUsage   = "Print URLs for all image sizes"
	)
	flag.BoolVar(&flagAllSizes, "a", allSizesDefault, allSizesUsage)
}

func usage() {
	appName := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "%v - upload images to Yandex.Fotki\n", appName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Usage:\n\n")
	fmt.Fprintf(os.Stderr, "\t%v [options] <imageFile1> [<imageFile2> ...]\n", appName)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options are:\n")
	fmt.Fprintf(os.Stderr, "\n")

	flag.CommandLine.VisitAll(func(flag *flag.Flag) {
		defaultValue := ""
		if len(flag.DefValue) > 0 {
			defaultValue = fmt.Sprintf(" [%v]", flag.DefValue)
		}
		fmt.Fprintf(os.Stderr, "\t-%s:\t%s%s\n", flag.Name, flag.Usage, defaultValue)
	})
}

func main() {

	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		return
	}

	filePaths := flag.Args()

	cfg, _ := config.Load()
	needToSaveCfg := false

	if cfg == nil || cfg.TokenExpired() {
		fmt.Printf("Getting new OAuth token...\n")
		tokenData := getTokenData()
		cfg = &config.Config{}
		cfg.OauthToken = tokenData.Token
		cfg.UpdateExpirationDateTime(tokenData.ExpiresIn)
		needToSaveCfg = true
	}

	if needToSaveCfg {
		err := cfg.Save()
		if err != nil {
			fmt.Printf("Can't save config: %v\n", err)
		}
	}

	for i, filePath := range filePaths {

		fmt.Printf("%v", filePath)

		uploadDataChan := make(chan yfotki.UploadData)
		yfotki.UploadFile(cfg.OauthToken, filePath, cfg.MainAlbumUrl, uploadDataChan)
		uploadData := <-uploadDataChan
		if uploadData.Error != nil {
			fmt.Printf(": error uploading: %v\n\n", uploadData.Error)
			continue
		}
		if cfg.MainAlbumUrl == "" {
			cfg.MainAlbumUrl = uploadData.MainAlbumUrl
			needToSaveCfg = true
		}

		fmt.Printf(":\n")

		if flagAllSizes {
			fmt.Printf("Original:  %v\n", uploadData.OrigImageUrl)
			fmt.Printf("XXX-Small: %v\n", uploadData.XxxSmallImageUrl)
			fmt.Printf("XX-Small:  %v\n", uploadData.XxSmallImageUrl)
			fmt.Printf("X-Small:   %v\n", uploadData.XSmallImageUrl)
			fmt.Printf("Small:     %v\n", uploadData.SmallImageUrl)
			fmt.Printf("Medium:    %v\n", uploadData.MediumImageUrl)
			fmt.Printf("Large:     %v\n", uploadData.LargeImageUrl)
			fmt.Printf("X-Large:   %v\n", uploadData.XLargeImageUrl)
		} else {
			fmt.Printf("%v\n", uploadData.OrigImageUrl)
		}

		if i != len(filePaths)-1 {
			fmt.Printf("\n")
		}
	}

	if needToSaveCfg {
		err := cfg.Save()
		if err != nil {
			fmt.Printf("Can't save config: %v\n", err)
		}
	}
}
