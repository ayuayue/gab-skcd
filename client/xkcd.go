package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"skcd/model"
	"time"
)

// ComicNumber 返回信息条数
type ComicNumber int

const (
	// BaseURL api
	BaseURL string = "https://xkcd.com"
	// DefaultClientTimeOut 客户端超时时间
	DefaultClientTimeOut time.Duration = 30 * time.Second
	// LatestComic 最新的数量
	LatestComic ComicNumber = 0
)

// XKCDClient 客户端结构体
type XKCDClient struct {
	client  *http.Client
	baseURL string
}

// NewXKCDClient get XKCDClient obj
func NewXKCDClient() *XKCDClient {
	return &XKCDClient{
		client: &http.Client{
			Timeout: DefaultClientTimeOut,
		},
		baseURL: BaseURL,
	}
}

// SetTimeout 设置过期时间
func (hc *XKCDClient) SetTimeout(d time.Duration) {
	hc.client.Timeout = d
}

// Fetch 检索漫画
func (hc *XKCDClient) Fetch(n ComicNumber, save bool) (model.Comic, error) {
	resp, err := hc.client.Get(hc.buildURL(n))
	if err != nil {
		return model.Comic{}, err
	}
	defer resp.Body.Close()

	var comicResp model.ComicResponse
	if err := json.NewDecoder(resp.Body).Decode(&comicResp); err != nil {
		return model.Comic{}, err
	}
	if save {
		if err := hc.SaveToDisk(comicResp.Img, "."); err != nil {
			fmt.Println("Failed to save image !")
		}
	}
	return comicResp.Comic(), err
}

// buildURL 构建api链接
func (hc *XKCDClient) buildURL(n ComicNumber) string {
	var finalURL string
	if n == LatestComic {
		finalURL = fmt.Sprintf("%s/info.0.json", hc.baseURL)
	} else {
		finalURL = fmt.Sprintf("%s/%d/info.0.json", hc.baseURL, n)
	}
	return finalURL
}

// SaveToDisk 将漫画保存到磁盘
func (hc *XKCDClient) SaveToDisk(url, savePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	absSaePath, _ := filepath.Abs(savePath)
	filePath := fmt.Sprintf("%s/%s", absSaePath, path.Base(url))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
