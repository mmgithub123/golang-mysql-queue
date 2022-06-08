package makeData

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type VsTask struct {
	ID       uint
	Sha256   string
	FileType uint8
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString returns a random string with a fixed length
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func InsertTableRand(db *gorm.DB) {
	fmt.Println("start insert into")

	for {
		sha256Str := RandomString(64)
		fileTypeNum := rand.Intn(len(defaultLetters))
		sleepTime := rand.Intn(3)

		user := VsTask{Sha256: sha256Str, FileType: uint8(fileTypeNum)}

		result := db.Create(&user)

		fmt.Println(result.RowsAffected)
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}
