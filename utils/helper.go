package utils

import (
	"archive/zip"
	"backend_gui/dto"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*
fuction to get and check username, check app name from decode redis
by irma 30/12/2021
*/
func DecodeSessionRedis(data string) (map[string]interface{}, error) {
	prices := make(map[string]interface{})
	s := strings.Split(data, ";")
	isAnyDataUsername := false
	isAnyDataNameApp := false
	for i := 0; i < len(s); i++ {
		tData := strings.Split(s[i], "|")
		if len(tData) > 1 {
			if strings.Contains(tData[0], "username") {
				tValue := strings.Split(fmt.Sprintf("%v", tData[1]), ":")
				valData := tValue[len(tValue)-1]
				prices[tData[0]] = strings.Trim(valData, "\"")
				isAnyDataUsername = true
			} else if strings.Contains(tData[0], os.Getenv("APP_NAME")) {
				tValue := strings.Split(fmt.Sprintf("%v", tData[1]), ":")
				valData := tValue[len(tValue)-1]
				if valData == "1" {
					prices[tData[0]] = valData
					isAnyDataNameApp = true
				}

			}
		}

	}
	if isAnyDataUsername && isAnyDataNameApp {
		return prices, nil
	}
	return nil, errors.New("not found")

}

/*
fuction to get cookie session  from redis,
decode redis, and validate result decode cookie session
by Irma 30/12/2021
*/
func GetSession(c *gin.Context, keyRedis string) (interface{}, error) {
	cookie, err := c.Request.Cookie(os.Getenv("COOKIE_NAME"))
	if err != nil {
		log.Error("Error redis : ", err.Error())
		return nil, err
	}
	redisData, errRedis := RedisConnection.Get(c, cookie.Value).Result()
	if errRedis != nil {
		log.Error("Error redis : ", errRedis.Error())
		return nil, errRedis
	}
	sessionData, errDecode := DecodeSessionRedis(redisData)

	if errDecode != nil {
		log.Error("Error redis : ", errRedis.Error())
		return nil, errDecode
	}

	if sessionData[keyRedis] == nil {
		log.Error("Error redis : null pointer key username")
		return nil, errors.New("null pointer")
	}

	if sessionData[os.Getenv("APP_NAME")] == nil {
		log.Error("Error redis : null pointer key app name")
		return nil, errors.New("null pointer")
	}
	return sessionData[keyRedis], nil
}

func CreatSessionToken(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie("session_token")
	if err != nil {
		return "", err
	}
	var expirestToken = jwt.NewNumericDate(time.Now().Add(time.Minute * 60))
	claims := JwtToken{
		cookie.Value,
		jwt.RegisteredClaims{
			ExpiresAt: expirestToken,
			Issuer:    "fakturpajak",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errTknString := token.SignedString([]byte("fakturP4J$1<"))
	if errTknString != nil {
		return "", errTknString
	}
	return tokenString, nil
}

type JwtToken struct {
	Token string `json:"token"`
	jwt.RegisteredClaims
}

func CheckAnyDuplicate(tempFiles []string, file string) bool {
	duplicate := false
	for x := 0; x < len(tempFiles); x++ {
		if tempFiles[x] == file {
			duplicate = true
		}
	}

	return duplicate
}
func CheckHeaderFile(row []string) bool {
	result := false
	for idx, header := range row {
		if idx == 0 {
			//fmt.Println(strings.ToLower(strings.Trim(header, " ")))
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "doc.number" {
				result = true
			}
		} else if idx == 1 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "e-mailaddress" {
				result = true
			}
		} else if idx == 2 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "createdon" {
				result = true
			}
		} else if idx == 3 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "time" {
				result = true
			}
		} else if idx == 4 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "createdby" {
				result = true
			}

		} else if idx == 5 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "emailto/cc" {
				result = true
			}
		} else if idx == 6 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "payer" {
				result = true
			}
		} else if idx == 7 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "payername" {
				result = true
			}
		} else if idx == 8 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "ship" {
				result = true
			}
		} else if idx == 9 {
			if strings.ToLower(strings.ReplaceAll(header, " ", "")) == "shiptoname" {
				result = true
			}
		}
	}

	return result

}

func MoveFile(sourcePath, destPath string, changePermission bool) error {
	log.Info("move file ", "sourcePath  "+sourcePath+" destPath "+destPath)
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		log.Error("Couldn't open source file: %s", err)
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)

	if err != nil {
		inputFile.Close()
		log.Error("Couldn't open source file: %s", err)
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		log.Error("Writing to output file failed: %s", err)
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		log.Error("Failed removing original file: %s", err)
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	//off sementara
	if changePermission {
		log.Info("proses change permissions ", destPath)
		cmd := exec.Command("chmod", "666", destPath)
		_, errChangePermission := cmd.Output()
		if errChangePermission != nil {
			log.Info("Error update permission file, "+destPath+", with error :", errChangePermission.Error())
		}
	}

	return nil
}

func CopyFile(source string, destination string, changePermission bool) bool {
	successCopy := true
	// Open original file
	original, errOpenOriginalFile := os.Open(source)
	if errOpenOriginalFile != nil {
		log.Error("Error copy file :", errOpenOriginalFile.Error())
		successCopy = false
	}
	defer original.Close()

	// Create new file
	newFile, errNewFile := os.Create(destination)
	if errNewFile != nil {
		log.Error("Error copy file :", errNewFile.Error())
		successCopy = false
	}
	defer newFile.Close()

	//This will copy
	_, errCopyFile := io.Copy(newFile, original)
	if errCopyFile != nil {
		log.Error("Error copy file :", errCopyFile.Error())
		successCopy = false
	}
	if changePermission {
		log.Info("proses change permissions ", destination)
		cmd := exec.Command("chmod", "666", destination)
		_, errChangePermission := cmd.Output()
		if errChangePermission != nil {
			log.Info("Error update permission file, "+destination+", with error :", errChangePermission.Error())
		}
	}
	return successCopy
}

func CheckDayLog() {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
	)
	fileLog := "logs/" + "fpm-backend" + (strings.Replace(formatted, "-", "", -1)) + ".log"
	if _, err := os.Stat(fileLog); os.IsNotExist(err) {
		WriteLogFpmBackend()
	}
}

func CheckChannel(dc string) bool {
	result := false
	channel := []string{"20", "35", "40", "45"}
	for k := 0; k < len(channel); k++ {
		if dc == channel[k] {
			result = true
			break
		}
	}
	return result
}

func ConvertDateToString(currentTime time.Time, format string) string {
	dateConvert := currentTime.String()
	dateConvert = currentTime.Format(format)
	return dateConvert
}

func TemplateEmail(data []map[string]interface{}, typeEmail string) string {
	var bodyEmail string = ``
	var bodyTable string = ``

	switch {
	case typeEmail == TYPE_EMAIL_INV_CANCEL:
		bodyTable = BodyEmailInvCancel(data)
	case typeEmail == TYPE_EMAIL_NO_FP:
		bodyTable = BodyEmailNoFp(data)
	case typeEmail == TYPE_EMAIL_INV_NOT_SEND:
		bodyTable = BodyEmailInvNoSend(data)
	case typeEmail == TYPE_EMAIL_NO_FILE_FP:
		bodyTable = BodyEmailNoFileFp(data)
	case typeEmail == TYPE_EMAIL_FP_CANCEL:
		bodyTable = BodyEmailFpCancel(data)
	case typeEmail == TYPE_EMAIL_FP_NOT_SEND:
		bodyTable = BodyEmailFpNotSend(data)
	default:
		log.Error("error write body email : ", "invalid param type body email")
		return bodyEmail
	}

	var footer string = ``
	footer = `</tbody>
					</table>
					<div class="footer">
						<p style="margin-top: 42px; margin-bottom: 2px;">Thank you</p>
						<span style='font: 10px Arial Helvetica, sans-serif'>*)email ini dikirim automatis dari sistem</span>
					</div>
				</div>

			</div>`

	bodyEmail = bodyTable + footer
	return bodyEmail
}
func BodyEmailInvCancel(data []map[string]interface{}) string {
	var bodyTable string = ``
	var no int = 1
	var headerContent = `<div class="container">
						<div class="tittle">
							<h4>Dear All,</h4>
						</div>
						<div style="margin-top: -9px;" class="contentFpm">
							<p>` + os.Getenv("EMAIL_CONTENT_CANCEL") + `</p>
							<table border="1" class="table">
								<thead style='text-align: center; background-color:#5DADE2'>
								<tr>
									<th scope="col">No</th>
									<th scope="col">Billing Doc</th>
									<th scope="col">Reference</th>
									<th scope="col">Payer Name</th>
									<th scope="col">Channel</th>
                                    <th scope="col">No. FP</th>
									<th scope="col">Period</th>
								</tr>
								</thead>
								<tbody>`

	for _, v := range data {
		dateTemp := []rune(fmt.Sprintf("%v", v["billing_date_zsd001n"]))
		bilDate := string(dateTemp[6:len(dateTemp)]) + "." + string(dateTemp[4:6]) + "." + string(dateTemp[0:4])
		noConvert := strconv.Itoa(no)
		dataFp := fmt.Sprintf("%v", v["fp_number_zv60"])
		if v["fp_number_zv60"] == "" {
			dataFp = "-"
		}

		bodyTable += ` <tr>
						<td>` + noConvert + `</td>
						<td>` + fmt.Sprintf("%v", v["billing_doc_cancel"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["billing_document_zsd001n"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["payer_name_zsd001n"]) + `</td>
						<td style='text-align: center;'>` + fmt.Sprintf("%v", v["dc_zsd001n"]) + `</td>
						<td style='text-align: center;'>` + dataFp + `</td>
						<td>` + bilDate + `</td>
                      </tr>`

		no++
	}
	return headerContent + bodyTable
}
func BodyEmailNoFp(data []map[string]interface{}) string {
	var bodyTable string = ``
	var no int = 1
	var headerContent = `<div class="container">
						<div class="tittle">
							<h4>Dear All,</h4>
						</div>
						<div style="margin-top: -9px;" class="contentFpm">
							<p>` + os.Getenv("EMAIL_CONTENT_INV_NO_FP") + `</p>
							<table border="1" class="table">
								<thead style='text-align: center; background-color:#5DADE2'>
								<tr>
									<th scope="col">No</th>
									<th scope="col">Billing Doc</th>
									<th scope="col">Payer Name</th>
									<th scope="col">Channel</th>
									<th scope="col">Period</th>
								</tr>
								</thead>
								<tbody>`
	for _, v := range data {
		dateTemp := []rune(fmt.Sprintf("%v", v["billing_date_zsd001n"]))
		bilDate := string(dateTemp[6:len(dateTemp)]) + "." + string(dateTemp[4:6]) + "." + string(dateTemp[0:4])
		noConvert := strconv.Itoa(no)

		bodyTable += ` <tr>
						<td>` + noConvert + `</td>
						<td>` + fmt.Sprintf("%v", v["billing_document_zsd001n"]) + `</td>
                        <td>` + fmt.Sprintf("%v", v["payer_name_zsd001n"]) + `</td>
						<td style='text-align: center;'>` + fmt.Sprintf("%v", v["dc_zsd001n"]) + `</td>
						<td>` + bilDate + `</td>
                      </tr>`

		no++
	}
	return headerContent + bodyTable
}
func BodyEmailInvNoSend(data []map[string]interface{}) string {
	var bodyTable string = ``
	var no int = 1
	var headerContent = `<div class="container">
						<div class="tittle">
							<h4>Dear All,</h4>
						</div>
						<div style="margin-top: -9px;" class="contentFpm">
							<p>` + os.Getenv("EMAIL_CONTENT_INV_NOT_SEND") + `</p>
							<table border="1" class="table">
								<thead style='text-align: center; background-color:#5DADE2'>
								<tr>
									<th scope="col">No</th>
									<th scope="col">Billing Doc</th>
									<th scope="col">Payer Name</th>
									<th scope="col">Channel</th>
									<th scope="col">Period</th>
									<th scope="col">No. FP</th>
								</tr>
								</thead>
								<tbody>`
	for _, v := range data {
		dateTemp := []rune(fmt.Sprintf("%v", v["billing_date_zsd001n"]))
		bilDate := string(dateTemp[6:len(dateTemp)]) + "." + string(dateTemp[4:6]) + "." + string(dateTemp[0:4])
		noConvert := strconv.Itoa(no)

		bodyTable += ` <tr>
						<td>` + noConvert + `</td>
						<td>` + fmt.Sprintf("%v", v["billing_document_zsd001n"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["payer_name_zsd001n"]) + `</td>
						<td style='text-align: center;'>` + fmt.Sprintf("%v", v["dc_zsd001n"]) + `</td>
						<td>` + bilDate + `</td>
						<td>` + fmt.Sprintf("%v", v["faktur"]) + `</td>
                      </tr>`

		no++
	}
	return headerContent + bodyTable
}
func BodyEmailNoFileFp(data []map[string]interface{}) string {
	var bodyTable string = ``
	var no int = 1
	var headerContent = `<div class="container">
						<div class="tittle">
							<h4>Dear All,</h4>
						</div>
						<div style="margin-top: -9px;" class="contentFpm">
							<p>` + os.Getenv("EMAIL_CONTENT_NO_FILE_FP") + `</p>
							<table border="1" class="table">
								<thead style='text-align: center; background-color:#5DADE2'>
								<tr>
									<th scope="col">No</th>
									<th scope="col">Billing Doc</th>
									<th scope="col">Payer Name</th>
									<th scope="col">Channel</th>
									<th scope="col">No. FP</th>
									<th scope="col">Period</th>
								</tr>
								</thead>
								<tbody>`
	for _, v := range data {
		dateTemp := []rune(fmt.Sprintf("%v", v["billing_date_zsd001n"]))
		bilDate := string(dateTemp[6:len(dateTemp)]) + "." + string(dateTemp[4:6]) + "." + string(dateTemp[0:4])
		noConvert := strconv.Itoa(no)

		bodyTable += ` <tr>
						<td>` + noConvert + `</td>
						<td>` + fmt.Sprintf("%v", v["billing_document_zsd001n"]) + `</td>
                        <td>` + fmt.Sprintf("%v", v["payer_name_zsd001n"]) + `</td>
						<td style='text-align: center;'>` + fmt.Sprintf("%v", v["dc_zsd001n"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["faktur"]) + `</td>
						<td>` + bilDate + `</td>
                      </tr>`

		no++
	}
	return headerContent + bodyTable

}
func BodyEmailFpNotSend(data []map[string]interface{}) string {
	var bodyTable string = ``
	var no int = 1
	var headerContent = `<div class="container">
						<div class="tittle">
							<h4>Dear All,</h4>
						</div>
						<div style="margin-top: -9px;" class="contentFpm">
							<p>` + os.Getenv("EMAIL_CONTENT_FP_NOT_SEND") + `</p>
							<table border="1" class="table">
								<thead style='text-align: center; background-color:#5DADE2'>
								<tr>
									<th scope="col">No</th>
									<th scope="col">Billing Doc</th>
									<th scope="col">Payer Name</th>
									<th scope="col">Channel</th>
									<th scope="col">Period</th>
									<th scope="col">No. FP</th>
									<th scope="col">File FP</th>
									
								</tr>
								</thead>
								<tbody>`
	for _, v := range data {
		dateTemp := []rune(fmt.Sprintf("%v", v["billing_date_zsd001n"]))
		bilDate := string(dateTemp[6:len(dateTemp)]) + "." + string(dateTemp[4:6]) + "." + string(dateTemp[0:4])
		noConvert := strconv.Itoa(no)

		bodyTable += ` <tr>
						<td>` + noConvert + `</td>
						<td>` + fmt.Sprintf("%v", v["billing_document_zsd001n"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["payer_name_zsd001n"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["dc_zsd001n"]) + `</td>
						<td>` + bilDate + `</td>
						<td>` + fmt.Sprintf("%v", v["faktur"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["file_name_upload"]) + `</td>
                      </tr>`

		no++
	}
	return headerContent + bodyTable

}
func BodyEmailFpCancel(data []map[string]interface{}) string {
	var bodyTable string = ``
	var no int = 1
	var headerContent = `<div class="container">
						<div class="tittle">
							<h4>Dear All,</h4>
						</div>
						<div style="margin-top: -9px;" class="contentFpm">
							<p>` + os.Getenv("EMAIL_CONTENT_FP_CANCEL") + `</p>
							<table border="1" class="table">
								<thead style='text-align: center; background-color:#5DADE2'>
								<tr>
									<th scope="col">No</th>
									<th scope="col">Billing Doc</th>
									<th scope="col">Payer Name</th>
									<th scope="col">Channel</th>
									<th scope="col">No. FP</th>
									<th scope="col">Period</th>
								</tr>
								</thead>
								<tbody>`
	for _, v := range data {
		dateTemp := []rune(fmt.Sprintf("%v", v["billing_date_zsd001n"]))
		bilDate := string(dateTemp[6:len(dateTemp)]) + "." + string(dateTemp[4:6]) + "." + string(dateTemp[0:4])
		noConvert := strconv.Itoa(no)

		bodyTable += ` <tr>
						<td>` + noConvert + `</td>
						<td>` + fmt.Sprintf("%v", v["billing_document_zsd001n"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["payer_name_zsd001n"]) + `</td>
						<td style='text-align: center;'>` + fmt.Sprintf("%v", v["dc_zsd001n"]) + `</td>
						<td>` + fmt.Sprintf("%v", v["faktur"]) + `</td>
						<td>` + bilDate + `</td>
                      </tr>`

		no++
	}
	return headerContent + bodyTable

}

// email notification
func EmailNotificationService(dataEmailNotification dto.ResSaveEmailDto, idEmail int, attachFile []string) (dto.ResNotifDto, error) {
	nowTime := time.Now()
	var sendEmailNotif dto.ReqNotifDto
	var idDetailAcc string

	idDetailAcc = strconv.Itoa(int(idEmail))
	sendEmailNotif.Recipient = dataEmailNotification.Recipient
	sendEmailNotif.Cc = dataEmailNotification.Cc
	sendEmailNotif.Sender = os.Getenv("EMAIL_SENDER")
	sendEmailNotif.Subject = dataEmailNotification.Subject
	sendEmailNotif.Body = dataEmailNotification.Body
	sendEmailNotif.Timestamp = nowTime.Format(time.RFC3339)
	sendEmailNotif.RequestId, _ = strconv.Atoi(idDetailAcc)
	sendEmailNotif.Partner = os.Getenv("EMAIL_PATNER_CODE")
	if dataEmailNotification.PathAttachment != "" {
		sendEmailNotif.PathAttachment = attachFile
	} else {
		sendEmailNotif.PathAttachment = []string{}
	}
	sendEmailNotif.LogStatus = strconv.Itoa(dataEmailNotification.LogStatus)
	ConstPartnerKey := os.Getenv("EMAIL_PATNER_KEY")
	dataSignature := idDetailAcc + "" + sendEmailNotif.Partner + "" + sendEmailNotif.Timestamp + "" + ConstPartnerKey
	signature := sha256.Sum256([]byte(dataSignature))
	strSignature := hex.EncodeToString(signature[:])
	sendEmailNotif.Signature = (strSignature)
	djs, _ := json.Marshal(sendEmailNotif)

	postValue := bytes.NewBuffer(djs)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	request, errCurl := http.NewRequest("POST", os.Getenv("SEND_EMAIL_URL_FPM"), postValue)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var respEmail dto.ResNotifDto
	if errCurl != nil {
		log.Error("Failed send email notification with subject : "+dataEmailNotification.Subject+": ", errCurl.Error())
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Error("Failed send email notification with subject : "+dataEmailNotification.Subject+": ", err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	errMarshal := json.Unmarshal(body, &respEmail)
	if errMarshal != nil {
		return respEmail, err
	}

	return respEmail, nil
}

func MonthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

func GenerateToken(username string) (string, error) {
	var SECRETKEY = []byte(os.Getenv("JWT_SECRET_KEY"))
	claim := jwt.MapClaims{}
	claim["username"] = username

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRETKEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func SendUpload(url string, keyJwt string, fileName string,path string,  OtherFields map[string]interface{} ) (dto.RespUploadSourceData, error){
	log.Info("Prosess curl file upload source data with file : ",fileName)
	filePath := path+fileName
	file, errOpenFile := os.Open(filePath)
	var respUpload dto.RespUploadSourceData
	if errOpenFile != nil{
		log.Error("error upload source data : ", errOpenFile)
		return respUpload, errOpenFile
	}


	fileContents, errReadFile := ioutil.ReadAll(file)

	if errReadFile != nil{
		log.Error("error upload source data : ", errReadFile.Error())
		return respUpload, errReadFile
	}
	fi, err := file.Stat()
	if err != nil {
		log.Error("error upload source data : ",err.Error())
		return respUpload, err
	}
	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	partData, errWriter1 := writer.CreateFormFile("file", filepath.Base(fi.Name()))

	if errWriter1 != nil{
		log.Error("error upload source data : ", errWriter1.Error() )
		return  respUpload,errWriter1

	}
	
	partData.Write(fileContents)
	if len(OtherFields) > 0{
		for key,value := range OtherFields{
			errWriteFields := writer.WriteField(key, fmt.Sprintf("%s", value))
			if errWriteFields != nil{
				log.Error("error upload source data : ", errWriteFields.Error())
				return  respUpload,errWriteFields

			}
		}
	}

	_,errCopy:=io.Copy(partData, file)
	if errCopy != nil{
		log.Error("error upload source data : ", errCopy.Error())
		return respUpload, errCopy
	}
	writer.Close()
	fmt.Println(payload)
	req, errRequest := http.NewRequest("POST", url, payload)
	if errRequest != nil {
		log.Error("error upload source data : ", errRequest.Error())
		return respUpload,errRequest
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client {Transport:tr}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", keyJwt)


	// Submit the request
	res, err := client.Do(req)
	if err != nil{
		log.Error("error upload source data : ", err.Error())
		return respUpload,err
	}

	defer res.Body.Close()

	body, errRead := ioutil.ReadAll(res.Body)
	if errRead != nil {
		log.Error("error upload source data : ",errRead.Error())
		return respUpload,errRead
	}

	errMarshal := json.Unmarshal(body, &respUpload)

	if errMarshal != nil{
		log.Error("error upload source data : ",errMarshal.Error())
	}

	return respUpload,nil
}

func GetUrlUploadSourceData(typeFile string) (string, map[string]interface{}) {
	url :=""
	typeFileZsd081:=""
	var listTypeFile =[]string{"zv60","zsd001n","zsd081-fp","zsd081-inv"}
	var fields map[string]interface{}
	switch {
	case typeFile == listTypeFile[0]:
		url = os.Getenv("UPLOAD_SOURCE_DATA_ZV60")
		typeFileZsd081 =""
	case typeFile == listTypeFile[1]:
		url = os.Getenv("UPLOAD_SOURCE_DATA_ZSD001N")
		typeFileZsd081 =""
	case typeFile == listTypeFile[2]:
		url = os.Getenv("UPLOAD_SOURCE_DATA_ZSD081")
		typeFileZsd081 = "FPM"
	case typeFile == listTypeFile[3]:
		url = os.Getenv("UPLOAD_SOURCE_DATA_ZSD081")
		typeFileZsd081 = "INV"
	default:
		url = ""
		typeFileZsd081 =""
	}
	if typeFileZsd081 != ""{
		fields = map[string]interface{}{
			"type": typeFileZsd081,
		}
	}

	return url, fields

}


func AddFileToZip(zipWriter *zip.Writer, pathFilename string, fileName string) error {

	fileToZip, errOpen := os.Open(pathFilename)
	if errOpen != nil {
		log.Error("Add file to zip, filename : "+fileName+", error : ",errOpen.Error())
		return errOpen
	}
	defer fileToZip.Close()

	//// Get the file information
	info, errGetInfo := fileToZip.Stat()
	if errGetInfo != nil {
		log.Error("Add file to zip, filename : "+fileName+", error : ",errGetInfo.Error() )
		return errGetInfo
	}


	header, errInfoHeader := zip.FileInfoHeader(info)
	if errInfoHeader != nil {
		log.Error("Add file to zip, filename : "+fileName+", error : ",errInfoHeader.Error())
		return errInfoHeader
	}
	//
	//// Using FileInfoHeader() above only uses the basename of the file. If we want
	//// to preserve the folder structure we can overwrite this with the full path.
	header.Name = fileName
	//
	//// Change to deflate to gain better compression
	//// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, errZip := zipWriter.Create(fileName)
	if errZip != nil {
		log.Error("Add file to zip, filename : "+fileName+", error : ",errZip.Error())
		return errZip
	}
	_, errCopy := io.Copy(writer, fileToZip)
	if errCopy != nil{
		log.Error("Add file to zip, filename : "+fileName+", error : ",errCopy.Error())
	}

	return errCopy
}