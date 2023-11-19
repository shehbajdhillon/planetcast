// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type CheckoutSessionResponse struct {
	SessionID string `json:"sessionId"`
}

type LineItemInput struct {
	PriceData PriceDataInput `json:"priceData"`
	Quantity  int            `json:"quantity"`
}

type PriceDataInput struct {
	Currency    string `json:"currency"`
	UnitAmount  int    `json:"unitAmount"`
	ProductName string `json:"productName"`
}

type UploadOption string

const (
	UploadOptionFileUpload  UploadOption = "FILE_UPLOAD"
	UploadOptionYoutubeLink UploadOption = "YOUTUBE_LINK"
)

var AllUploadOption = []UploadOption{
	UploadOptionFileUpload,
	UploadOptionYoutubeLink,
}

func (e UploadOption) IsValid() bool {
	switch e {
	case UploadOptionFileUpload, UploadOptionYoutubeLink:
		return true
	}
	return false
}

func (e UploadOption) String() string {
	return string(e)
}

func (e *UploadOption) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UploadOption(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UploadOption", str)
	}
	return nil
}

func (e UploadOption) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
