package database

import (
	"encoding/csv"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
)

func LoadDataFromFile(db *gorm.DB, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines[1:] {

		firstName := line[0]
		lastName := line[1]
		email := line[2]
		password := line[3]

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		id := uuid.New().String()

		user := Account{
			Email: email,
		}

		log.Print(user)

		if err := db.First(&user, "email=?", email).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user.ID = id
				user.FirstName = firstName
				user.LastName = lastName
				user.Password = string(hashedPassword)
				db.Create(&user)
			} else {
				return err
			}
		}
	}
	return nil
}
