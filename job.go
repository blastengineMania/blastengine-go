package blastengine

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Job struct {
	Client *Client

	JobId        int    `json:"job_id"`
	Percentage   int    `json:"percentage"`
	Status       string `json:"status"`
	SuccessCount int    `json:"success_count"`
	FailedCount  int    `json:"failed_count"`
	TotalCount   int    `json:"total_count"`
	ErrorFileUrl string `json:"error_file_url"`
}

func (j *Job) Finished() (bool, error) {
	if j.JobId == 0 {
		return false, fmt.Errorf("job ID is required")
	}
	url := fmt.Sprintf("https://app.engn.jp/api/v1/deliveries/-/emails/import/%d", j.JobId)
	bodyBytes, err := j.Client.sendRequest("GET", url, nil, nil, false, nil)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bodyBytes, &j)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return j.Percentage == 100, nil
}

func (j *Job) Download() (string, error) {
	if j.ErrorFileUrl == "" {
		return "", nil
	}
	// ZIPをメモリ上で展開
	zipData, err := j.Client.sendRequest("GET", j.ErrorFileUrl, nil, nil, false, nil)
	if err != nil {
		return "", err
	}
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return "", err
	}

	// ZIPにファイルが含まれているか確認
	if len(reader.File) == 0 {
		return "", fmt.Errorf("zip archive is empty")
	}

	// 最初のファイルを取得
	file := reader.File[0]

	// ZIP内のファイルを開く
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	// ファイルの内容を読み込む
	content, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
