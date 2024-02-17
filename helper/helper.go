package helper

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"io"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// 伟大始于渺小！
const key = "From small beginnings comes great things"

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"name"`
	Identity string `json:"identity"`
	Is_admin int    `json:"is_admin"`
}

func GetMd5(data string) string {
	sumData := md5.Sum([]byte(data))

	return fmt.Sprintf("%x", sumData)
}

func GetToken(username, identity string, isAdmin int) (tokenString string, e error) {
	data := &UserClaims{
		StandardClaims: jwt.StandardClaims{},
		Username:       username,
		Identity:       identity,
		Is_admin:       isAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	tokenString, e = token.SignedString([]byte(key))

	if e != nil {
		return "", e
	}

	return tokenString, e
}

func AnalyseToken(tokenString string) (data *UserClaims, e error) {
	data = new(UserClaims)
	token, e := jwt.ParseWithClaims(tokenString, data, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})

	if e != nil {
		return nil, e
	}

	if !token.Valid {
		return nil, errors.New("analyse token failure")
	}

	return data, nil
}

func SendEmailCode(emailString, code string) (err error) {
	e := email.NewEmail()
	e.From = "test_xtu_oj <test_xtu_oj@yeah.net>"
	e.To = []string{emailString}

	e.Subject = "以下是给您的消息"
	e.Text = []byte("您的消息是")
	e.HTML = []byte("<h1>" + code + "</h1>")
	if err := e.Send("smtp.yeah.net:25", smtp.PlainAuth("", "test_xtu_oj@yeah.net", "GETTHQETCMCVALET", "smtp.yeah.net")); err != nil {
		return err
	}
	return err
}

func GetUuid() (uuidstr string, e error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uid.String(), err
}

func GetRandCode() (code string) {
	for i := 0; i < 5; i++ {
		code = code + strconv.Itoa(rand.Intn(10))
	}
	return code
}

func GetDataFromConfigFile(path string) (result map[string]string) {
	result = map[string]string{}

	// 打开配置文件
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		log.Print("open MysqlConfig.txt file failure:", err)
		return result
	}
	defer f.Close()

	inputReader := bufio.NewReader(f)
	for {
		// 逐行读取配置文件信息
		inputstring, err := inputReader.ReadString('\n')
		// 如果读到文件末尾
		if err == io.EOF {
			break
		}

		sc := strings.Split(inputstring, ":")
		// 如果缺少参数或者是格式错误
		if len(sc) < 2 {
			log.Panicln("File format error")
			return result
		}
		sc[1] = strings.TrimSpace(sc[1])

		// 存入结果
		result[sc[0]] = sc[1]
	}
	return result
}

// 保存代码
func SaveCode(code []byte) (path string, e error) {
	// 新建路径
	uuidPath, _ := GetUuid()
	dirname := "./code/" + uuidPath
	err := os.Mkdir(dirname, 0777)
	if err != nil {
		return "", err
	}

	// 新建文件
	path = dirname + "/main.go"
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}

	f.Write(code)
	defer f.Close()
	return path, nil
}
