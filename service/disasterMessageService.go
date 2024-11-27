package service

import (
	"encoding/json"
	"fmt"
	"github.com/GDG-on-Campus-KHU/SC2_BE/models"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

const (
	BaseURL     = "https://www.safetydata.go.kr/V2/api/DSSP-IF-00247"
	ServiceKey  = "72ZBE332C1399B51" // 발급받은 서비스키
	PollingTime = 30 * time.Second
)

var lastSN string // 마지막으로 처리한 재난 문자의 SN(일련번호)

// fetchLatestMessage: 최신 데이터를 가져오는 함수
func FetchLatestDisasterMessage() (*models.DisasterMessage, error) {
	client := resty.New()

	// API 호출
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"serviceKey": ServiceKey,
			"numOfRows":  "1", // 하나의 데이터만 요청
			"pageNo":     "1",
			"returnType": "json",
			"crtDt":      time.Now().UTC().Format(time.RFC3339),
		}).
		Get(BaseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %w", err)
	}

	// 응답이 비어있는 경우 처리
	if resp == nil || resp.Body() == nil {
		return nil, fmt.Errorf("empty response from API")
	}

	// JSON 응답 파싱
	var disasterResponse models.DisasterResponse
	if err := json.Unmarshal(resp.Body(), &disasterResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Items가 비어있는 경우 처리
	if len(disasterResponse.Items) == 0 {
		log.Println("No new disaster messages available from the API.")
		return nil, nil // 데이터 없음은 에러로 처리하지 않음
	}

	// 최신 데이터 반환
	return &disasterResponse.Items[0], nil
}

// 주기적으로 API를 호출하여 업데이트 확인
func PollForUpdates() {
	ticker := time.NewTicker(PollingTime) // 주기 설정
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Checking for updates from Disaster Message API...")
		message, err := FetchLatestDisasterMessage()
		if err != nil {
			log.Printf("Error fetching latest disaster message: %v", err)
			continue
		}

		// 메시지가 없는 경우 처리
		if message == nil {
			log.Println("No new disaster message found.")
			continue
		}

		if message.SN == lastSN {
			log.Println("No new updates.")
			continue
		}

		// 새 메시지가 있는 경우 처리
		log.Printf("New disaster message detected: SN=%s\n", message.SN)
		lastSN = message.SN // 마지막 처리된 메시지 SN 업데이트

		// 새로운 데이터를 처리하는 로직
		processNewMessage(message)
	}
}

func processNewMessage(message *models.DisasterMessage) {
	// 메시지가 nil인지 확인
	if message == nil {
		log.Println("Received nil message. Skipping processing.")
		return
	}

	// 새 메시지 처리 로직
	log.Printf("Processing new disaster message: %v\n", message)
}
