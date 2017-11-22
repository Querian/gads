package gads

import "encoding/xml"

type FeedItemService struct {
	Auth
}

type FeedItem struct {
	FeedID          int64                    `xml:"feedId"`
	FeedItemID      int64                    `xml:"feedItemId,omitempty"`
	Status          string                   `xml:"status,omitempty"`
	StartTime       string                   `xml:"startTime,omitempty"`
	EndTime         string                   `xml:"endTime,omitempty"`
	AttributeValues []FeedItemAttributeValue `xml:"attributeValues"`
}

type FeedItemAttributeValue struct {
	AttributeID   int64      `xml:"feedAttributeId"`
	IntegerValue  *int64     `xml:"integerValue,omitempty"`
	DoubleValue   *float64   `xml:"doubleValue,omitempty"`
	BooleanValue  *bool      `xml:"booleanValue,omitempty"`
	StringValue   *string    `xml:"stringValue,omitempty"`
	IntegerValues *[]int64   `xml:"integerValues,omitempty"`
	DoubleValues  *[]float64 `xml:"doubleValues,omitempty"`
	BooleanValues *[]bool    `xml:"booleanValues,omitempty"`
	StringValues  *[]string  `xml:"stringValues,omitempty"`
}

type FeedItemOperations map[string][]FeedItem

func NewSitelinkFeedItemAttributeText(value string) FeedItemAttributeValue {
	return FeedItemAttributeValue{AttributeID: 1, StringValue: &value}
}

func NewSitelinkFeedItemAttributeURL(value string) FeedItemAttributeValue {
	return FeedItemAttributeValue{AttributeID: 2, StringValue: &value}
}

func NewSitelinkFeedItemAttributeLine1(value string) FeedItemAttributeValue {
	return FeedItemAttributeValue{AttributeID: 3, StringValue: &value}
}

func NewSitelinkFeedItemAttributeLine2(value string) FeedItemAttributeValue {
	return FeedItemAttributeValue{AttributeID: 4, StringValue: &value}
}

func NewSitelinkFeedItemAttributeFinalURLs(value []string) FeedItemAttributeValue {
	return FeedItemAttributeValue{AttributeID: 5, StringValues: &value}
}

func NewSitelinkFeedItemAttributeFinalMobileURLs(value []string) FeedItemAttributeValue {
	return FeedItemAttributeValue{AttributeID: 6, StringValues: &value}
}

func NewSitelinkFeedItemAttributeTrackingURL(value string) FeedItemAttributeValue {
	return FeedItemAttributeValue{AttributeID: 7, StringValue: &value}
}

func NewFeedItemService(auth *Auth) *FeedItemService {
	return &FeedItemService{Auth: *auth}
}

func (s *FeedItemService) Get(selector Selector) (feedItems []FeedItem, totalCount int64, err error) {
	selector.XMLName = xml.Name{"", "selector"}
	respBody, err := s.Auth.request(
		feedItemServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)
	if err != nil {
		return
	}
	getResp := struct {
		Size      int64      `xml:"rval>totalNumEntries"`
		FeedItems []FeedItem `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return
	}
	return getResp.FeedItems, getResp.Size, err
}

func (s *FeedItemService) Mutate(feedItemOperations FeedItemOperations) (feedItems []FeedItem, err error) {
	type feedItemOperation struct {
		Action   string   `xml:"operator"`
		FeedItem FeedItem `xml:"operand"`
	}
	operations := []feedItemOperation{}
	for action, feeditems := range feedItemOperations {
		for _, feedItem := range feeditems {
			operations = append(operations, feedItemOperation{Action: action, FeedItem: feedItem})
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []feedItemOperation `xml:"operations"`
	}{
		XMLName: xml.Name{Space: baseUrl, Local: "mutate"},
		Ops:     operations,
	}
	respBody, err := s.Auth.request(feedItemServiceUrl, "mutate", mutation)
	if err != nil {
		return
	}
	mutateResp := struct {
		BaseResponse
		FeedItems []FeedItem `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.FeedItems, err
}
