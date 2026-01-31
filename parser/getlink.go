package parser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// 建議修改簽名，傳入 ctx 以重用瀏覽器，並回傳 error
func Getlink(ctx context.Context, url string) (string, error) {
	var htmlContent string
	var description string

	fmt.Println("正在開啟職稱網頁:", url)

	// 這裡直接使用傳入的 ctx，不需要重新 NewContext
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`h3`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return "", fmt.Errorf("Chromedp 錯誤: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("Goquery 解析錯誤: %w", err)
	}

	// 修正回傳邏輯
	doc.Find("p.job-description__content").EachWithBreak(func(i int, s *goquery.Selection) bool {
		out := strings.TrimSpace(s.Text())
		if out != "" {
			description = out
			//fmt.Println("職缺描述:", description)
			//fmt.Println("--------------------------------------------------")
			return false // 找到第一個後就停止迴圈 (Break)
		}
		return true // 繼續尋找
	})

	return description, nil
}
