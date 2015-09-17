package yfotki

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	kServiceDocumentUrl         = "http://api-fotki.yandex.ru/api/me/"
	kOauthTokenKey              = "oauth_token"
	kMainAlbumId                = "photo-list"
	kUploadFieldImageName       = "image"
	kUploadFieldAccessName      = "access"
	kUploadFieldAccessValue     = "public"
	kUploadFieldTitleName       = "title"
	kUploadFieldPubChannelName  = "pub_channel"
	kUploadFieldAppPlatformName = "app_platform"
	kSizeTagXXXS                = "XXXS"
	kSizeTagXXS                 = "XXS"
	kSizeTagXS                  = "XS"
	kSizeTagS                   = "S"
	kSizeTagM                   = "M"
	kSizeTagL                   = "L"
	kSizeTagXL                  = "XL"
	kSizeTagOrig                = "orig"
)

type UploadData struct {
	XxxSmallImageUrl, XxSmallImageUrl, XSmallImageUrl, SmallImageUrl, MediumImageUrl, LargeImageUrl, XLargeImageUrl, OrigImageUrl string
	MainAlbumUrl                                                                                                                  string
	Error                                                                                                                         error
}

func getMainAlbumUrl(token string) (mainAlbumUrl string, err error) {

	type Album struct {
		XMLName xml.Name `xml:"collection"`
		Href    string   `xml:"href,attr"`
		Id      string   `xml:"id,attr"`
	}

	type Service struct {
		XMLName xml.Name `xml:"service"`
		Albums  []Album  `xml:"workspace>collection"`
	}

	mainAlbumUrl = ""
	err = nil

	url := fmt.Sprintf("%v?%v=%v", kServiceDocumentUrl, kOauthTokenKey, token)

	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var service Service
	err = xml.Unmarshal(body, &service)
	if err != nil {
		return
	}

	for _, album := range service.Albums {
		if album.Id == kMainAlbumId {
			mainAlbumUrl = album.Href
			return
		}
	}

	return
}

func doUploadFile(token, filePath, mainAlbumUrl, appName, appPlatform string) (responseData []byte, err error) {
	var bodyBuffer bytes.Buffer
	writer := multipart.NewWriter(&bodyBuffer)
	formWriter, err := writer.CreateFormFile(kUploadFieldImageName, filePath)
	if err != nil {
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(formWriter, file)
	if err != nil {
		return
	}

	writer.WriteField(kUploadFieldAccessName, kUploadFieldAccessValue)

	fileName := path.Base(filePath)
	title := strings.TrimSuffix(fileName, path.Ext(fileName))
	writer.WriteField(kUploadFieldTitleName, title)

	writer.WriteField(kUploadFieldPubChannelName, appName)
	writer.WriteField(kUploadFieldAppPlatformName, appPlatform)

	err = writer.Close()
	if err != nil {
		return
	}

	url := fmt.Sprintf("%v?%v=%v", mainAlbumUrl, kOauthTokenKey, token)

	request, err := http.NewRequest("POST", url, &bodyBuffer)
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	responseBodyBuffer := &bytes.Buffer{}
	_, err = responseBodyBuffer.ReadFrom(response.Body)
	if err != nil {
		return
	}
	response.Body.Close()
	responseData = responseBodyBuffer.Bytes()
	return
}

func doParseUploadResponse(responseData []byte, fileExtension string, uploadData *UploadData) {

	type Image struct {
		XMLName xml.Name `xml:"img"`
		Href    string   `xml:"href,attr"`
		Size    string   `xml:"size,attr"`
		Height  string   `xml:"height,attr"`
		Width   string   `xml:"width,attr"`
	}

	type Result struct {
		XMLName xml.Name `xml:"entry"`
		Images  []Image  `xml:"img"`
	}

	var result Result
	uploadData.Error = xml.Unmarshal(responseData, &result)
	if uploadData.Error != nil {
		return
	}

	for _, image := range result.Images {
		switch image.Size {
		case kSizeTagXXXS:
			uploadData.XxxSmallImageUrl = image.Href + fileExtension
		case kSizeTagXXS:
			uploadData.XxSmallImageUrl = image.Href + fileExtension
		case kSizeTagXS:
			uploadData.XSmallImageUrl = image.Href + fileExtension
		case kSizeTagS:
			uploadData.SmallImageUrl = image.Href + fileExtension
		case kSizeTagM:
			uploadData.MediumImageUrl = image.Href + fileExtension
		case kSizeTagL:
			uploadData.LargeImageUrl = image.Href + fileExtension
		case kSizeTagXL:
			uploadData.XLargeImageUrl = image.Href + fileExtension
		case kSizeTagOrig:
			uploadData.OrigImageUrl = image.Href + fileExtension
		}
	}

	return
}

func uploadFile(token, filePath, mainAlbumUrl, appName, appPlatform string, uploadDataChan chan UploadData) {

	var uploadData UploadData

	if mainAlbumUrl == "" {
		var err error
		mainAlbumUrl, err = getMainAlbumUrl(token)
		if err != nil {
			errorText := fmt.Sprintf("Error getting main album url: %v", err)
			uploadData.Error = errors.New(errorText)
			uploadDataChan <- uploadData
			return
		}
		uploadData.MainAlbumUrl = mainAlbumUrl
	}

	responseData, err := doUploadFile(token, filePath, mainAlbumUrl, appName, appPlatform)
	if err != nil {
		uploadData.Error = err
		uploadDataChan <- uploadData
		return
	}

	fileExtension := strings.ToLower(path.Ext(filePath))

	doParseUploadResponse(responseData, fileExtension, &uploadData)

	uploadDataChan <- uploadData
	return
}

// cachedMainAlbumUrl may be empty
func UploadFile(token, filePath, cachedMainAlbumUrl, appName, appPlatform string, uploadDataChan chan UploadData) {
	go uploadFile(token, filePath, cachedMainAlbumUrl, appName, appPlatform, uploadDataChan)
}
