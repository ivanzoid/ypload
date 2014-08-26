package yfotki

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	kServiceDocumentUrl     = "http://api-fotki.yandex.ru/api/me/"
	kOauthTokenKey          = "oauth_token"
	kMainAlbumId            = "photo-list"
	kUploadFieldImageName   = "image"
	kUploadFieldAccessName  = "access"
	kUploadFieldAccessValue = "public"
)

type UploadData struct {
	OrigImageUrl, SmallImageUrl, LargeImageUrl, XLargeImageUrl, XxLargeImageUrl string
	MainAlbumUrl                                                                string
	Error                                                                       error
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

	log.Printf("Url: %v\n", url)

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

func doUploadFile(token, filePath, mainAlbumUrl string) (responseString string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formWriter, err := writer.CreateFormFile(kUploadFieldImageName, filePath)
	if err != nil {
		return
	}
	_, err = io.Copy(formWriter, file)
	if err != nil {
		return
	}
	writer.WriteField(kUploadFieldAccessName, kUploadFieldAccessValue)
	err = writer.Close()
	if err != nil {
		return
	}

	url := fmt.Sprintf("%v?%v=%v", mainAlbumUrl, kOauthTokenKey, token)
	request := http.NewRequest("POST", url, body)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	responseBodyBuffer := &bytes.Buffer{}
	_, err = responseBodyBuffer.ReadFrom(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	responseString = string(responseBodyBuffer)
	return
}

func doParseUploadResponse(responseString string) (uploadData UploadData) {
	log.Printf("Response: %v\n", responseString)
}

func uploadFile(token, filePath, mainAlbumUrl string, uploadDataChan chan UploadData) {
	var uploadData UploadData
	var err error
	if mainAlbumUrl == "" {
		mainAlbumUrl, err = getMainAlbumUrl(token)
		if err != nil {
			goto fail
		}
		uploadData.MainAlbumUrl = mainAlbumUrl
	}

	responseString, err := doUploadFile(token, filePath, mainAlbumUrl)
	if err != nil {
		goto fail
	}

	uploadData = doParseUploadResponse(responseString)
	uploadDataChan <- uploadData
	return

fail:
	uploadData.Error = err
	uploadDataChan <- uploadData
	return
}

// cachedMainAlbumUrl may be empty
func UploadFile(token, filePath, cachedMainAlbumUrl string, uploadDataChan chan UploadData) {
	go uploadFile(token, filePath, errChan)
}
