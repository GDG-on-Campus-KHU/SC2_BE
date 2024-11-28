package repository

import (
	"gorm.io/gorm"
	"log"
)

type TokenRepository struct {
	DB *gorm.DB
}

// 클라이언트 토큰 저장
func (repo *TokenRepository) SaveToken(userID int, token string) error {
	err := repo.DB.Clauses()
}

func (repo *TokenRepository) GetToken(userID int) ([]string, error) {
	query := `
		SELECT token 
		FROM client_tokens 
		WHERE user_id = ?`

	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error fetching tokens for user %s: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			log.Printf("Error scanning token: %v", err)
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if len(tokens) == 0 {
		log.Printf("No tokens found for user %s", userID)
	}
	return tokens, nil
}
