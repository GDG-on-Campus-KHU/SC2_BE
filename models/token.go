package models

import "time"

// Token 모델 정의
type Token struct {
	ID        uint      `gorm:"primaryKey"`      // 기본 키
	UserID    string    `gorm:"index;not null"`  // 사용자 ID (인덱스 생성)
	Token     string    `gorm:"unique;not null"` // FCM 토큰 (고유값)
	CreatedAt time.Time `gorm:"autoCreateTime"`  // 생성 시간
	UpdatedAt time.Time `gorm:"autoUpdateTime"`  // 갱신 시간
}
