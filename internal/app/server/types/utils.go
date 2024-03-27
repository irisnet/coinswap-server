package types

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func MarshalJsonIgnoreErr(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func Get(url string) (bz []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode != 200")
	}
	defer resp.Body.Close()

	bz, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr) // 输出加密结果
}
