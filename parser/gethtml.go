package parser

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// 定義資料結構，方便之後轉 JSON
type JobResult struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Description  string `json:"description"`
	Company_Name string `json:"company"`
	Company_Link string `json:"company_link"`
}

func Gethtml(url string, index int) []JobResult {
	// 1. 處理 URL 分頁
	finalURL := strings.Replace(url, "page=1", fmt.Sprintf("page=%d", index), 1)

	// 2. 建立 Context (建議把 Allocator 放在這層)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var htmlContent string
	fmt.Printf("正在爬取第 %d 頁...\n", index)

	err := chromedp.Run(ctx,
		chromedp.Navigate(finalURL),
		chromedp.WaitVisible(`.info-name`, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &htmlContent),
	)
	if err != nil {
		log.Printf("Chromedp 錯誤: %v", err)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal("Goquery 解析錯誤:", err)
	}

	var results []JobResult

	// 3. 迭代每個職缺項目
	doc.Find("a.info-job.jb-link").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".info-name").Text())
		href, exists := s.Attr("href")

		if exists && title != "" {
			// 修正點：呼叫 Getlink 取得詳細內容
			// 注意：這裡建議傳入 ctx 以節省開啟瀏覽器的開銷 (如前次討論)
			singlejob, err := Getlink(ctx, href)
			if err == nil {
				// 2. 使用 append 將這一筆資料「存入」你的結果清單中
				results = append(results, singlejob)
			}

		}

	})

	return results
}
