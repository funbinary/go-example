package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/bin-work/go-example/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	rep, err := Signature("406EB3F3F93D1D445A8AF8DFEFBF50D1FCDFA322F9F3E718F4CD3258F868FCC390C73F9C0C6F58A2F150D8D2C52A3062F73F634949C057D01C87A5F6DCF622141D513F4491A9EE6AB8D445C2385B24A4E410218DF92CADF1E26DDEBFB4F692346DB1A56DEBDE3D7F593CEF372694D4587CB0CF95DC0ED32B7064C52FDA796B38EBA132D41623296AAC41D10CB18C06B7D869A7A92F25AED11AFD6187A77ACDBCAD564D633CBE65D1517D21FC2F26623D05D8DE5AC00A7302CD1B9D156A7AD139C98EF979AA162AC22ADB394A0972BB042913E5A3E3DAEDC811AA2BA1D58A3B2F12BC5537FF9FEB6D719651C258E6EADB22893BCA98C1CA5AA0F23068101FE4D7")
	fmt.Println(rep.IsSignature())
	fmt.Println(err)

}

type SignatureRep struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		XAccessToken  string `json:"X_Access_Token"`
		RoomId        string `json:"room_id"`
		LiveUserToken string `json:"live_user_token"`
	} `json:"result"`
}

func (self SignatureRep) IsSignature() bool {
	return self.Code == 200
}

func (self SignatureRep) Error() string {
	return self.Message
}

func Signature(baseinfo string) (SignatureRep, error) {
	var rep SignatureRep

	url := "https://erp1.neoscholar.com/neo-education/eduRoomDetection/noToken/getSignatureCallback"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("baseinfo", baseinfo)

	err := writer.Close()
	if err != nil {
		rep.Message = err.Error()
		logrus.Infof("write close err %v ", err)
		return rep, err
	}

	logrus.Infof("url: %v baseinfo ", url, baseinfo)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		rep.Message = err.Error()
		logrus.Errorf("NewRequest %v", err)
		return rep, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		rep.Message = err.Error()
		return rep,
			errors.Wrapf(err, "http.Post %s req param %s", url, baseinfo)
	}
	defer res.Body.Close()
	respjson, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(respjson, &rep); err != nil {
		rep.Message = errors.Wrapf(err, "Unmarshal %v", string(respjson)).Error()
		return rep,
			errors.Wrapf(err, "Unmarshal %v", string(respjson))
	}
	return rep, nil
}
