package adapter

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/sirupsen/logrus"
)

type AdapterImagem struct {
	Cfg *config.Config
}

func NewImage(cfg *config.Config) *AdapterImagem {
	return &AdapterImagem{
		Cfg: cfg,
	}
}

func (ai *AdapterImagem) Upload(file []byte, id, name string, wg *sync.WaitGroup, ch chan error) error {
	defer wg.Done()
	logrus.Info("[Upload] image started")
	defer logrus.Info("[Upload] image finished")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(ai.Cfg.S3Region),
		Credentials: credentials.NewStaticCredentials(ai.Cfg.S3AccessKeyID, ai.Cfg.S3SecretAccessKey, ""),
	})
	logrus.WithFields(logrus.Fields{
		"region": ai.Cfg.S3Region,
		"bucket": ai.Cfg.S3BucketName,
		"key":    ai.Cfg.S3AccessKeyID,
	}).Info("CREDENTIALS")

	if err != nil {
		logrus.Infof("[Upload] 1 %s", err.Error())
		return err
	}
	extensao := strings.Split(name, ".")
	nameFile := fmt.Sprintf("%s.%s", id, extensao[len(extensao)-1])
	if strings.Split(http.DetectContentType(file), "/")[0] != "image" {
		return fmt.Errorf("invalid file type")
	}
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(ai.Cfg.S3BucketName),
		Key:           aws.String(nameFile),
		Body:          bytes.NewReader(file),
		ContentLength: aws.Int64(int64(len(file))),
		ACL:           aws.String("private"),
		ContentType:   aws.String(http.DetectContentType(file)),
	})
	if err != nil {
		logrus.Infof("[Upload] 2 %s", err.Error())
		return err
	}
	return nil
}

func (a *AdapterImagem) Read() (url string, err error) {
	return "", nil
}
