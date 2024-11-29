package service

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/messaging"
	"github.com/GDG-on-Campus-KHU/SC2_BE/config"
	"github.com/GDG-on-Campus-KHU/SC2_BE/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

var tokenCollection *mongo.Collection

func InitTokenCollection(client *mongo.Client) {
	tokenCollection = client.Database("SC2_DB").Collection("tokens")
	log.Println("[INFO] Token collection initialized")
}

// 토큰 저장
func SaveToken(token models.TokenRequest) error {
	if tokenCollection == nil {
		log.Println("[ERROR] Token collection is not initialized")
		return errors.New("token collection is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 중복 확인
	filter := bson.M{"token": token.Token}
	count, err := tokenCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("[ERROR] Failed to check token existence: %v", err)
		return err
	}
	if count > 0 {
		log.Println("[WARN] Token already exists, skipping save")
		return nil // 중복된 경우 저장하지 않음
	}

	// 토큰 저장
	_, err = tokenCollection.InsertOne(ctx, bson.M{"token": token.Token})
	if err != nil {
		log.Printf("[ERROR] Failed to save token: %v", err)
		return err
	}

	log.Println("[INFO] Token saved successfully")
	return nil
}

// FCM 메시지 전송
func SendNotification(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 모든 토큰 조회
	cursor, err := tokenCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("[ERROR] Failed to fetch tokens: %v", err)
		return err
	}
	defer cursor.Close(ctx)

	// FCM 클라이언트 초기화
	client, err := config.FirebaseApp.Messaging(context.Background())
	if err != nil {
		log.Printf("[ERROR] FCM 클라이언트 초기화 실패: %v", err)
		return err
	}

	// 토큰별 메시지 전송
	for cursor.Next(ctx) {
		var tokenDoc struct {
			Token string `bson:"token"`
		}
		if err := cursor.Decode(&tokenDoc); err != nil {
			log.Printf("[ERROR] Failed to decode token: %v", err)
			continue
		}

		// 메시지 구성
		fcmMessage := &messaging.Message{
			Token: tokenDoc.Token,
			Notification: &messaging.Notification{
				Title: "재난 안전",
				Body:  message,
			},
		}

		// 메시지 전송
		_, err = client.Send(context.Background(), fcmMessage)
		if err != nil {
			log.Printf("[ERROR] Failed to send notification to token %s: %v", tokenDoc.Token, err)
			continue
		}

		log.Printf("[INFO] Notification sent to token: %s", tokenDoc.Token)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("[ERROR] Cursor iteration error: %v", err)
		return err
	}

	log.Println("[INFO] All notifications sent successfully")
	return nil
}
