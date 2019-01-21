package gads

import (
	"encoding/xml"
	//	"fmt"
)

type AdGroupBidModifierService struct {
	Auth
}

func NewAdGroupBidModifierService(auth *Auth) *AdGroupBidModifierService {
	return &AdGroupBidModifierService{Auth: *auth}
}

type AdGroupBidModifierOperations map[string][]AdGroupBidModifier

type AdGroupBidModifier struct {
	CampaignId        int64     `xml:"campaignId"`
	AdGroupId         int64     `xml:"adGroupId"`
	Criterion         Criterion `xml:"criterion"`
	BidModifier       float64   `xml:"bidModifier"`
	BidModifierSource string    `xml:"bidModifierSource"`
}

func (cc *AdGroupBidModifier) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "campaignId":
				if err := dec.DecodeElement(&cc.CampaignId, &start); err != nil {
					return err
				}
			case "adGroupId":
				if err := dec.DecodeElement(&cc.AdGroupId, &start); err != nil {
					return err
				}
			case "criterion":
				criterion, err := criterionUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				cc.Criterion = criterion
			case "bidModifier":
				if err := dec.DecodeElement(&cc.BidModifier, &start); err != nil {
					return err
				}
			case "bidModifierSource":
				if err := dec.DecodeElement(&cc.BidModifierSource, &start); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get returns budgets matching a given selector and the total count of matching budgets.
func (s *AdGroupBidModifierService) Get(selector Selector) (bm []AdGroupBidModifier, totalCount int64, err error) {
	selector.XMLName = xml.Name{"", "selector"}
	respBody, err := s.Auth.request(
		adGroupBidModifierServiceUrl,
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
		return bm, totalCount, err
	}
	getResp := struct {
		Size                int64                `xml:"rval>totalNumEntries"`
		AdGroupBidModifiers []AdGroupBidModifier `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return bm, totalCount, err
	}
	return getResp.AdGroupBidModifiers, getResp.Size, err
}

// Mutate takes a budgetOperations and creates, modifies or destroys the associated budgets.
func (s *AdGroupBidModifierService) Mutate(bidmOperations AdGroupBidModifierOperations) (resp []AdGroupBidModifier, err error) {
	type bidmOperation struct {
		Action             string             `xml:"operator"`
		AdGroupBidModifier AdGroupBidModifier `xml:"operand"`
	}
	operations := []bidmOperation{}
	for action, bidms := range bidmOperations {
		for _, bidm := range bidms {
			operations = append(operations, bidmOperation{Action: action, AdGroupBidModifier: bidm})
		}
	}
	respBody, err := s.Auth.request(
		adGroupBidModifierServiceUrl,
		"mutate",
		struct {
			XMLName xml.Name
			Ops     []bidmOperation `xml:"operations"`
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "mutate",
			},
			Ops: operations,
		},
	)
	if err != nil {
		return resp, err
	}
	mutateResp := struct {
		BaseResponse
		AdGroupBidModifiers []AdGroupBidModifier `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return resp, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.AdGroupBidModifiers, err
}
