package main

import (
	"bytes"
	parser "crawler/parser"
	repo "crawler/repo"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var CurrentURL string
	var bucketName string

	file, _ := os.Open("./config/Getlink.json")
	defer file.Close()

	var config map[string]interface{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	} else {
		fmt.Println("", config["link"])
		CurrentURL = config["link"].(string)
		bucketName = config["bucket"].(string)
	}

	for i := 1; i <= 10; i++ { //並發爬取 5 頁
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			pageResults := parser.Gethtml(CurrentURL, index)

			// 將爬取到的每個職缺資訊轉換成 JSON 並上傳
			for _, job := range pageResults {
				//fmt.Printf("%s", pageResults[0])
				// 將整個 job 結構轉換為 JSON 格式，方便未來擴充與使用
				jsonData, err := json.MarshalIndent(job, "", "  ")
				if err != nil {
					fmt.Printf("無法將 job %s 轉換為 JSON: %v\n", job.Title, err)
					continue
				}

				// 使用 MD5 將職缺連結轉換為唯一的 ID 作為檔名
				hash := md5.Sum([]byte(job.Link)) // 計算 MD5 hash，回傳 [16]byte 陣列
				jobID := fmt.Sprintf("%x", hash)  // 將 byte 陣列格式化為十六進位字串

				objectName := fmt.Sprintf("%s.json", jobID)
				data := bytes.NewReader(jsonData)

				if err := repo.GcsUpload(bucketName, objectName, data); err != nil {
					fmt.Printf("上傳 %s 失敗: %v\n", objectName, err)
				}
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("所有頁面爬取完成")

}
