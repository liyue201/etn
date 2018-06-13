package main

import (
	"encoding/json"
	"github.com/liyue201/etn/client/config"
	"github.com/liyue201/etn/client/utils"
	"github.com/liyue201/go-logger"
	"io/ioutil"
	"net/http"
	"time"
	"path/filepath"
	"os"
)

const HOST = "http://127.0.0.1:7883"

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Url  string `json:"url"`
}

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func GetRemoteVersion() (int, error) {
	data, err := Get(HOST + "/api/v1/version")
	if err != nil {
		logger.Errorf("[GetRemoteVersion] %s", err)
		return 0, err
	}
	resp := struct {
		Code int `json:"code"`
		Data int `json:"data"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[GetRemoteVersion] %s", err)
		return 0, err
	}
	return resp.Data, nil
}

func GetFiles() ([]*File, error) {
	data, err := Get(HOST + "/api/v1/files")
	if err != nil {
		logger.Errorf("[GetFiles] %s", err)
		return nil, err
	}
	resp := struct {
		Code  int     `json:"code"`
		Files []*File `json:"data"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		logger.Errorf("[GetFiles] %s", err)
		return resp.Files, err
	}
	return resp.Files, nil
}

func DownloadFile(file *File) error {
	data, err := Get(file.Url)
	if err != nil{
		return err
	}
	dir := filepath.Dir(file.Path)
	if !utils.PathExist(dir) {
		err := os.MkdirAll(dir, 777)
		if err != nil{
			return err
		}
	}
	return  ioutil.WriteFile(file.Path, data, 777)
}

func main() {

	configPath := utils.GetExeDir() + "client.yml"
	err := config.InitConfig(configPath)
	if err != nil {
		logger.Errorf("[main] %s", err)
	}

	for {
		remoteVersion, err := GetRemoteVersion()
		if err != nil {
			logger.Errorf("[main] %s", err)

			time.Sleep(time.Second * 60)
			continue
		}

		if remoteVersion <= config.Cfg.Version {
			break
		}

		files, err := GetFiles()
		if err != nil {
			logger.Errorf("[main] %s", err)

			time.Sleep(time.Second * 60)
			continue
		}

		for _, file := range files {
			err = DownloadFile(file)
			if err != nil {
				logger.Errorf("[main] %s", err)

				time.Sleep(time.Second)
				continue
			}
		}
		break
	}
}
