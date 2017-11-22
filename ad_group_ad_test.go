package gads

import (
	"fmt"
	"testing"
)

func testAdGroupAdService(t *testing.T) (service *AdGroupAdService) {
	return &AdGroupAdService{Auth: testAuthSetup(t)}
}

func TestAdGroupAd(t *testing.T) {
	adGroup, cleanupAdGroup := testAdGroup(t)
	defer cleanupAdGroup()

	agas := testAdGroupAdService(t)
	adGroupAds, err := agas.Mutate(
		AdGroupAdOperations{
			"ADD": {
				NewAdGroupExpandedTextAd(adGroup.Id, "https://classdo.com/en", "test headline "+rand_word(10), "test line one", "test line two", "PAUSED"),
				NewAdGroupExpandedTextAd(adGroup.Id, "https://classdo.com/en", "test   teStTo "+rand_word(10), "test line one", "test line two", "PAUSED"),
				NewAdGroupExpandedTextAd(adGroup.Id, "https://classdo.com/en", "test headline "+rand_word(10), "test line one", "test line two", "PAUSED"),
				NewAdGroupExpandedTextAd(adGroup.Id, "https://classdo.com/en", "test headline "+rand_word(10), "test line one", "test line two", "PAUSED"),
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err = agas.Mutate(AdGroupAdOperations{"REMOVE": adGroupAds})
		if err != nil {
			t.Error(err)
		}
	}()

	adGroupIdStr := fmt.Sprintf("%d", adGroup.Id)
	_, _, err = agas.Get(
		Selector{
			Fields: []string{
				"AdGroupId",
				"Id",
				"Status",
			},
			Predicates: []Predicate{
				{"AdGroupId", "EQUALS", []string{adGroupIdStr}},
			},
			Ordering: []OrderBy{
				{"AdGroupId", "ASCENDING"},
				{"Id", "ASCENDING"},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

}
