package gads

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// AWQLClient struct for AWQL caller
type AWQLClient struct {
	Auth
}

// AWQLFormat type for AWQL format
type AWQLFormat string

// AWQL formats
const (
	AWQLFormatCSVForExcel AWQLFormat = "CSVFOREXCEL"
	AWQLFormatCSV         AWQLFormat = "CSV"
	AWQLFormatTSV         AWQLFormat = "TSV"
	AWQLFormatXML         AWQLFormat = "XML"
	AWQLFormatGzippedCSV  AWQLFormat = "GZIPPED_CSV"
	AWQLFormatGzippedXML  AWQLFormat = "GZIPPED_XML"
)

// AWQLRequest struct for awql request
type AWQLRequest struct {
	Query                  string
	Format                 AWQLFormat
	SkipReportHeader       bool
	SkipColumnHeader       bool
	SkipReportSummary      bool
	IncludeZeroImpressions bool
	UseRawEnumValues       bool
}

type AwqlError struct {
	XMLName xml.Name     `xml:"reportDownloadError"`
	Error   AwqlApiError `xml:"ApiError"`
}

type AwqlApiError struct {
	Type      string `xml:"type"`
	Trigger   string `xml:"trigger"`
	Fieldpath string `xml:"fieldPath"`
}

// NewAWQLClient is a constructor for AWQLClient
func NewAWQLClient(auth *Auth) *AWQLClient {
	return &AWQLClient{Auth: *auth}
}

const (
	maxSizeForNotValidBody = 30
)

// Download downloads a report by awql request
func (a *AWQLClient) Download(awqlReq AWQLRequest) (io.ReadCloser, error) {
	req, err := http.NewRequest(
		"POST",
		reportAPIURL,
		strings.NewReader(url.Values{"__rdquery": {awqlReq.Query}, "__fmt": {string(awqlReq.Format)}}.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("developerToken", a.DeveloperToken)
	req.Header.Add("clientCustomerId", a.CustomerId)
	req.Header.Add("skipReportHeader", strconv.FormatBool(awqlReq.SkipReportHeader))
	req.Header.Add("skipColumnHeader", strconv.FormatBool(awqlReq.SkipColumnHeader))
	req.Header.Add("skipReportSummary", strconv.FormatBool(awqlReq.SkipReportSummary))
	req.Header.Add("includeZeroImpressions", strconv.FormatBool(awqlReq.IncludeZeroImpressions))
	req.Header.Add("useRawEnumValues", strconv.FormatBool(awqlReq.UseRawEnumValues))
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		respBody, errRead := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if errRead != nil {
			return nil, fmt.Errorf("unexpected StatusCode [%d]", resp.StatusCode)
		}

		v := AwqlError{}
		errRead = xml.Unmarshal(respBody, &v)
		if errRead != nil {
			bl := len(respBody)
			if bl > maxSizeForNotValidBody {
				bl = maxSizeForNotValidBody
			}
			return nil, fmt.Errorf("unexpected StatusCode [%d] with content [%s]", resp.StatusCode, string(respBody[:bl]))
		}

		return nil, fmt.Errorf("[%s @ ; trigger:%s]", v.Error.Type, v.Error.Trigger)
	}

	return resp.Body, nil
}
