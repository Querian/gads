package gads

import "testing"

func testConversionTrackerService(t *testing.T) (service *ConversionTrackerService) {
	return &ConversionTrackerService{Auth: testAuthSetup(t)}
}

func TestConversionTracker(t *testing.T) {
	cs := testConversionTrackerService(t)
	conversionTrackers, err := cs.Mutate(
		ConversionTrackerOperations{
			"ADD": {
				AdWordsConversionTracker{
					CommonConversionTracker: &CommonConversionTracker{
						Name:                      "test adword conversion tracker" + rand_str(10),
						Status:                    "ENABLED",
						Category:                  "DEFAULT",
						ViewthroughLookbackWindow: 10,
						CtcLookbackWindow:         7,
						CountingType:              "ONE_PER_CLICK",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = cs.Get(
		Selector{
			Fields: []string{
				"Id",
				"Name",
			},
			Predicates: []Predicate{
				{"Id", "EQUALS", []string{conversionTrackers[0].GetID()}},
			},
			Ordering: []OrderBy{
				{"Id", "ASCENDING"},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	defer func(conversionTrackers ConversionTrackers) {
		conversionTrackers[0].SetStatus("REMOVED")

		_, err = cs.Mutate(ConversionTrackerOperations{"SET": conversionTrackers})
		if err != nil {
			t.Error(err)
		}
	}(conversionTrackers)
}
