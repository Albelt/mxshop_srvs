package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:create_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"uniqueIndex:idx_mobile;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(256);not null"` //哈希加密后的密码
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   uint8      `gorm:"default:1;type:tinyint comment '1男 2女'"`
	Role     int        `gorm:"default:1;type:tinyint comment '1表示普通用户, 2表示管理员'"`
}

// 实际DB中存储有关密码的json字符串
type Password struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
	Alg  string `json:"alg"`
}

func (u *User) GenPassword(rawPwd string) {
	salt, hash := encodePassword(rawPwd)

	pwd := Password{
		Hash: hash,
		Salt: salt,
		Alg:  "sha256",
	}
	tmp, _ := json.Marshal(pwd)

	u.Password = string(tmp)
}

func (u *User) VerifyPassword(inputPwd string) bool {
	var pwd Password

	err := json.Unmarshal([]byte(u.Password), &pwd)
	if err != nil {
		return false
	}

	return verifyPassword(inputPwd, pwd.Salt, pwd.Hash)
}

// 对原始密码进行编码 (sha256WithSalt)
//  @return string salt
//  @return string hash
func encodePassword(rawPwd string) (string, string) {
	salt := generateSalt()

	h := sha256.New()
	h.Write([]byte(salt + rawPwd))

	return salt, hex.EncodeToString(h.Sum(nil))
}

// 校验密码
//  @param rawPwd 用户输入的密码
//  @param salt
//  @param encodedPwd 编码后的密码
//  @return bool
func verifyPassword(inputPwd, salt, encodedPwd string) bool {
	h := sha256.New()
	h.Write([]byte(salt + inputPwd))

	return hex.EncodeToString(h.Sum(nil)) == encodedPwd
}

func generateSalt() string {
	const (
		letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
		lettersLen = len(letters)
		saltLen    = 32
	)

	salt := make([]byte, saltLen)
	for i := range salt {
		salt[i] = letters[rand.Intn(lettersLen)]
	}

	return string(salt)
}
