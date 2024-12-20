package service

import (
	"encoding/json"
	"fmt"
	"github.com/GDG-on-Campus-KHU/SC2_BE/config"
	"github.com/GDG-on-Campus-KHU/SC2_BE/models"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"time"
)

const PollingTime = 10 * time.Second

var lastSN string // 마지막으로 처리한 재난 문자의 SN(일련번호)

// API를 호출해서 재난 안전 문자 반환
func FetchLatestDisasterMessage() (*models.DisasterMessage, error) {
	client := resty.New()

	// API 호출
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"serviceKey": config.ServiceKey,
			"numOfRows":  "1", // 하나의 데이터만 요청
			"pageNo":     "1",
			"returnType": "json",
			"crtDt":      time.Now().Format("20241127"),
		}).
		Get(config.BaseURL)

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

	log.Println("Background polling for disaster messages started...")

	for range ticker.C {
		log.Println("Checking for updates from Disaster Message API...")
		message, err := FetchLatestDisasterMessage()
		if err != nil {
			log.Printf("Error fetching latest disaster message: %v", err)
			continue
		}

		if message == nil {
			log.Println("[INFO] No new disaster messages available from the API.")
			continue // message가 nil일 경우 이후 로직 실행하지 않음
		}

		// 새 메시지가 없는 경우
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

// 새로운 재난 안전 문자 데이터 처리
func processNewMessage(message *models.DisasterMessage) {
	// 1. 메시지가 nil인지 확인
	if message == nil {
		log.Println("No message to process.")
		return
	}

	// 2. AI 모델에 메시지 전송하고 response를 받음
	response, err := SendDisasterMessage(*message)
	if err != nil {
		log.Printf("Error sending disaster message: %v", err)
		return
	}

	// 3. AI 응답 데이터를 푸시 알림으로 처리
	err = SendNotification(response.Results.HotspotResults.PushAlarming)
	if err != nil {
		log.Printf("[ERROR] Failed to send push notification: %v", err)
		return
	}

	log.Println("[INFO] Push notification sent successfully.")

}

// AI 모델에 재난 문자 request로 전송
func SendDisasterMessage(data models.DisasterMessage) (*models.DisasterGuideResponse, error) {
	// Resty 클라이언트 생성
	client := resty.New()

	// JSON 요청 전송
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		Post(os.Getenv("AI_MODEL_URL"))

	if err != nil {
		return nil, fmt.Errorf("failed to send disaster message: %w", err)
	}

	// HTTP 상태 코드 확인
	if resp.StatusCode() != 200 {
		log.Printf("[ERROR] Unexpected status code: %d, response body: %s", resp.StatusCode(), resp.String())
		return nil, fmt.Errorf("failed to send disaster message, status code: %d", resp.StatusCode())
	}

	// 응답 데이터를 구조체로 디코딩
	var result models.DisasterGuideResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON into DisasterGuideResponse: %w", err)
	}

	// 성공적으로 구조체 반환
	return &result, nil
}

func GetActRmksList(response models.DisasterGuideResponse) []string {
	var actRmksList []string

	// action_plan 순회하며 actRmks 값을 수집
	for _, action := range response.Results.HotspotResults.ActionPlan {
		actRmksList = append(actRmksList, action.ActRmks)
	}

	return actRmksList
}
