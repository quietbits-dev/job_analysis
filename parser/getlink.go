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
func Getlink(ctx context.Context, url string) (JobResult, error) {
	var htmlContent string

	var jobResult JobResult

	fmt.Println("正在開啟職稱網頁:", url)

	// 這裡直接使用傳入的 ctx，不需要重新 NewContext
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`h3`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return jobResult, fmt.Errorf("Chromedp 錯誤: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return jobResult, fmt.Errorf("Goquery 解析錯誤: %w", err)
	}

	// 修正回傳邏輯
	doc.Find("p.job-description__content").EachWithBreak(func(i int, s *goquery.Selection) bool {
		out := strings.TrimSpace(s.Text())
		if out != "" {
			jobResult.Description = out
			fmt.Println("職缺描述:", jobResult.Description)
			fmt.Println("--------------------------------------------------")
			return false // 找到第一個後就停止迴圈 (Break)
		}
		return true // 繼續尋找
	})
	doc.Find(".company__name").EachWithBreak(func(i int, s *goquery.Selection) bool {
		out2 := strings.TrimSpace(s.Text())
		if out2 != "" {
			jobResult.Company_Name = out2
			fmt.Println("公司名稱:", jobResult.Company_Name)
			fmt.Println("--------------------------------------------------")
			return false
		}
		return true
	})

	doc.Find("div[data-gtm-content='公司名稱'] a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		out3, exists := s.Attr("href")
		if exists {
			jobResult.Company_Link = out3
			fmt.Println("公司Link:", jobResult.Company_Link)
			fmt.Println("================================================")
			return false
		}
		return true
	})

	return jobResult, nil
}
