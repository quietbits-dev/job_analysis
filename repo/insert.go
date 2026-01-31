package repo

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

func GcsUpload(bucketName, objectName string, data io.Reader) error {
	ctx := context.Background()

	// Recommended: Remove the hardcoded path.
	// The client library will automatically find credentials from the GOOGLE_APPLICATION_CREDENTIALS environment variable.
	client, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// 取得 GCS 物件的控制代碼 (handle)。
	obj := client.Bucket(bucketName).Object(objectName)

	// 為物件建立一個 writer。
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, data); err != nil {
		// 在回傳錯誤前關閉 writer 以進行清理是很重要的。
		wc.Close()
		return fmt.Errorf("io.Copy: %w", err)
	}

	// 關閉 writer 來完成上傳。
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Printf("Blob %v uploaded to bucket %v.\n", objectName, bucketName)
	return nil
}
