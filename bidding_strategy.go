package gads

import (
	"encoding/xml"
)

// BiddingStrategyService servie for bidding strategies
type BiddingStrategyService struct {
	Auth
}

// NewBiddingStrategyService returns new instance of BiddingStrategyService
func NewBiddingStrategyService(auth *Auth) *BiddingStrategyService {
	return &BiddingStrategyService{Auth: *auth}
}

// SharedBiddingStrategy struct of entity for BiddingStrategyService
type SharedBiddingStrategy struct {
	BiddingScheme BiddingSchemeInterface `xml:"biddingScheme,omitempty"`
	ID            int64                  `xml:"id,omitempty"` // A unique identifier
	Name          string                 `xml:"name,omitempty"`
	Status        string                 `xml:"status,omitempty"`
	Type          string                 `xml:"type,omitempty"`
}

// UnmarshalXML special unmarshal for the different bidding schemes
func (b *SharedBiddingStrategy) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "id":
				if err := dec.DecodeElement(&b.ID, &start); err != nil {
					return err
				}
			case "name":
				if err := dec.DecodeElement(&b.Name, &start); err != nil {
					return err
				}
			case "status":
				if err := dec.DecodeElement(&b.Status, &start); err != nil {
					return err
				}
			case "type":
				if err := dec.DecodeElement(&b.Type, &start); err != nil {
					return err
				}
			case "biddingScheme":
				bs, err := biddingSchemeUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				b.BiddingScheme = bs
			}
		}
	}
	return nil
}

// A BiddingStrategyOperations maps operations to the bidding strategy they will be performed
// on.  BiddingStrategy operations can be 'ADD', 'REMOVE' or 'SET'
type BiddingStrategyOperations map[string][]SharedBiddingStrategy

// Get returns budgets matching a given selector and the total count of matching budgets.
func (s *BiddingStrategyService) Get(selector Selector) ([]SharedBiddingStrategy, int64, error) {
	selector.XMLName = xml.Name{Space: "", Local: "selector"}
	respBody, err := s.Auth.request(
		biddingStrategyServiceUrl,
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
		return nil, 0, err
	}
	getResp := struct {
		Size          int64                   `xml:"rval>totalNumEntries"`
		BiddingStrats []SharedBiddingStrategy `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return nil, 0, err
	}
	return getResp.BiddingStrats, getResp.Size, err
}

// Mutate takes a budgetOperations and creates, modifies or destroys the associated budgets.
func (s *BiddingStrategyService) Mutate(bidOperations BiddingStrategyOperations) ([]SharedBiddingStrategy, error) {
	type bidStratOperation struct {
		Action   string                `xml:"operator"`
		BidStrat SharedBiddingStrategy `xml:"operand"`
	}
	operations := []bidStratOperation{}
	for action, bidstrats := range bidOperations {
		for _, bidstrat := range bidstrats {
			operations = append(operations,
				bidStratOperation{Action: action, BidStrat: bidstrat},
			)
		}
	}
	respBody, err := s.Auth.request(
		biddingStrategyServiceUrl,
		"mutate",
		struct {
			XMLName xml.Name
			Ops     []bidStratOperation `xml:"operations"`
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "mutate",
			},
			Ops: operations,
		},
	)
	if err != nil {
		return nil, err
	}
	mutateResp := struct {
		BaseResponse
		BidStrats []SharedBiddingStrategy `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return nil, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.BidStrats, err
}
