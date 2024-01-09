package ossutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"hash"
	"io"
	"time"

	"github.com/pkg/errors"
)

/* const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
} */

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicySign struct {
	AccessKeyId string `json:"OSSAccessKeyId"` //
	Host        string `json:"host"`           //
	Expire      int64  `json:"expire"`         //
	Signature   string `json:"signature"`      //
	Policy      string `json:"policy"`         //
	Directory   string `json:"dir"`            //
	Callback    string `json:"callback"`       //
	Key         string `json:"key"`            // 上传的objectName
}

type PolicySignParam struct {
	Host        string `json:"host"`        // host的格式为 bucketname.endpoint, 如: "https://bucket-name.oss-cn-hangzhou.aliyuncs.com", 留空则自动拼接
	ExpiresIn   int64  `json:"expiresIn"`   // 过期时间(秒), 如: 30
	UploadDir   string `json:"uploadDir"`   // 用户上传文件时指定的前缀, 如: "user-dir-prefix/"
	CallbackUrl string `json:"callbackUrl"` // 上传回调服务器的URL, 请将IP和Port配置为您自己的真实信息
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

// GetPolicySign web直传服务端签名
func (p *OssUtil) GetPolicySign(param PolicySignParam) (policySign PolicySign, err error) {

	now := time.Now()
	var accessKeyId = p.accessKeyId
	var accessKeySecret = p.accessKeySecret
	var host = param.Host
	if host == "" {
		host = "https://" + p.bucket + "." + p.endpoint
	}
	var upload_dir = param.UploadDir
	var callbackUrl = param.CallbackUrl
	var expire_end = now.Unix() + param.ExpiresIn
	var tokenExpire = get_gmt_iso8601(expire_end)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, upload_dir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	if err != nil {
		err = errors.Wrap(err, "jsonMarshal 1")
		return
	}
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = callbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		err = errors.Wrap(err, "callback json err:")
		return
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	policySign.AccessKeyId = accessKeyId
	policySign.Host = host
	policySign.Expire = expire_end
	policySign.Signature = string(signedStr)
	policySign.Directory = upload_dir
	policySign.Policy = string(debyte)
	policySign.Callback = string(callbackBase64)

	return policySign, nil
}
