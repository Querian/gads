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
