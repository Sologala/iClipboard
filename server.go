package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.design/x/clipboard"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	apiVersion = "1"
)

func RunHTTPServer() {
	go func() {
		engin := gin.New()
		setupRoute(engin)
		if err := engin.Run(":" + Gconfig.Port); err != nil {
			log.Error().Msg("HTTP Server 启动失败 您的应用可能不能正常运行")
			systray.Quit()
			log.Error().Msg("failed to start http server")
			return
		}
	}()
}

func setupRoute(engin *gin.Engine) {
	engin.Use(clientName(), logger(), gin.Recovery(), apiVersionChecker(), auth())
	engin.GET("/", getHandler)
	engin.POST("/", setHandler)
	engin.NoRoute(notFoundHandler)
}

func clientName() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlEncodedClientName := c.GetHeader("X-Client-Name")
		clientName, err := url.PathUnescape(urlEncodedClientName)
		if err != nil || clientName == "" {
			clientName = "匿名设备"
		}
		c.Set("clientName", clientName)
		c.Next()
	}
}

func apiVersionChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.GetHeader("X-API-Version")
		if version == apiVersion {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "接口版本不匹配，请升级您的捷径",
		})
	}
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		if Gconfig.Authkey == "" {
			c.Next()
			return
		}

		reqAuth := c.GetHeader("X-Auth")

		timestamp := time.Now().Unix()
		timeKey := timestamp / Gconfig.Authkey_expired_timeout

		authCodeRaw := Gconfig.Authkey + "." + strconv.FormatInt(timeKey, 10)
		authCodeHash := md5.Sum([]byte(authCodeRaw))
		authCodeString := hex.EncodeToString(authCodeHash[:])

		if authCodeString == reqAuth {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "操作被拒绝：Authkey 验证失败",
		})
	}
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		clientName := c.GetString("clientName")
		log.Info().Msgf("method %s \n statusCode: %s\n clientIP %s\n path: %s\n duration %s\n clientName %s \n", c.Request.Method, statusCode, clientIP, path, duration, clientName)
	}
}

const (
	TypeText    = "text"
	TypeFile    = "file"
	TypeMedia   = "media"
	TypeBitmap  = "bitmap"
	TypeUnknown = "unknown"
)

type ResponseFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ResponseFiles []ResponseFile

func getHandler(c *gin.Context) {
	if ServiceFlag == false {
		c.Status(http.StatusBadRequest)
		return
	}
	var b []byte

	data_type := TypeText
	b = clipboard.Read(clipboard.FmtText)
	if b == nil {
		b = clipboard.Read(clipboard.FmtImage)
		data_type = TypeBitmap
	}

	log.Info().Msgf("The clipboard type is %s ", data_type)

	if b == nil {
		log.Error().Msg("failed to get content type of clipboard")
		c.Status(http.StatusBadRequest)
		return
	}

	if data_type == TypeText {
		log.Info().Msg("get clipboard text")
		c.JSON(http.StatusOK, gin.H{
			"type": "text",
			"data": string(b), // 将数据返回给客户端
		})
		return
	}

	if data_type == TypeBitmap {

		responseFiles := make([]ResponseFile, 0, 1)
		responseFiles = append(responseFiles, ResponseFile{
			"clipboard.png",
			base64.StdEncoding.EncodeToString(b),
		})

		c.JSON(http.StatusOK, gin.H{
			"type": "file",
			"data": responseFiles,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "无法识别剪切板内容"})
}

func readBase64FromFile(path string) (string, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

// Set clipboard handler

// TextBody is a struct of request body when iOS send files to windows
type TextBody struct {
	Text string `json:"data"`
}

func setHandler(c *gin.Context) {
	// if !Gconfig.ReserveHistory {
	// 	cleanTempFiles()
	// }
	if ServiceFlag == false {
		c.Status(http.StatusBadRequest)
		return
	}
	contentType := c.GetHeader("X-Content-Type")
	if contentType == TypeText {
		setTextHandler(c)
		return
	}

	setFileHandler(c)
}

func setTextHandler(c *gin.Context) {
	var body TextBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error().AnErr("failed to bind text body", err)
		c.Status(http.StatusBadRequest)
		return
	}

	clipboard.Write(clipboard.FmtText, []byte(body.Text))

	// if err := utils.Clipboard().SetText(body.Text); err != nil {
	// 	log.Error().AnErr("failed to set clipboard", err)
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }

	log.Error().Any("set clipboard text", body)
	c.Status(http.StatusOK)
}

// FileBody is a struct of request body when iOS send files to windows
type FileBody struct {
	Files []File `json:"data"`
}

// File is a struct represtents request file
type File struct {
	Name   string `json:"name"` // filename
	Base64 string `json:"base64"`
	_bytes []byte `json:"-"` // don't use this directly. use *File.Bytes() to get bytes
}

// Bytes returns byte slice of file
func (f *File) Bytes() ([]byte, error) {
	if len(f._bytes) > 0 {
		return f._bytes, nil
	}
	fileBytes, err := base64.StdEncoding.DecodeString(f.Base64)
	if err != nil {
		return []byte{}, nil
	}
	f._bytes = fileBytes
	return fileBytes, nil
}

func setFileHandler(c *gin.Context) {

	// contentType := c.GetHeader("X-Content-Type")

	// 将发送过来的json解析到结构体上。
	var body FileBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error().AnErr("failed to bind file body", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// paths := make([]string, 0, len(body.Files)) // 创建一个空切片（数组）
	// 枚举发送过来的base64编码的图片，或者文件
	for _, file := range body.Files {
		if file.Name == "-" && file.Base64 == "-" {
			continue
		}
		fmt.Println(file.Name) //Clipboard 2024年3月2日 10.52.png
		// 编译正则表达式
		// re_bmp := regexp.MustCompile(`^[a-zA-Z0-9_]+\.bmp$`)
		re_png := regexp.MustCompile(`.*\.png$`)

		// 判断字符串是否匹配正则表达式
		if re_png.MatchString(file.Name) {
			log.Info().Msg("Recived Image")

		} else {
			log.Error().Any("仅仅支持png|bmp数据", file.Name)
			c.Status(http.StatusBadRequest)

			fmt.Println("仅仅支持png|bmp数据") //Clipboard 2024年3月2日 10.52.png
			return
		}

		fileBytes, err := file.Bytes() // 在这里就完成了解码
		if err != nil {
			log.Error().AnErr("read buffer faild", err)
			c.Status(http.StatusBadRequest)
			return
		}
		clipboard.Write(clipboard.FmtImage, fileBytes)
		log.Info().Msgf("set Recived  %s  to buffer faild", file.Name)
	}
	c.Status(http.StatusOK)

	log.Info().Msg("success1")
}

func notFoundHandler(c *gin.Context) {
	log.Error().Msgf("user_ip : %s", c.Request.RemoteAddr)
	c.Status(http.StatusNotFound)
}
