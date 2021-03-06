package binarylane

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRegions_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/regions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"regions": [
				{
					"slug": "1"
				},
				{
					"slug": "2"
				}
			],
			"meta": {
				"total": 2
			}
		}`)
	})

	regions, resp, err := client.Regions.List(ctx, nil)
	if err != nil {
		t.Errorf("Regions.List returned error: %v", err)
	}

	expectedRegions := []Region{{Slug: "1"}, {Slug: "2"}}
	if !reflect.DeepEqual(regions, expectedRegions) {
		t.Errorf("Regions.List returned regions %+v, expected %+v", regions, expectedRegions)
	}

	expectedMeta := &Meta{Total: 2}
	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Regions.List returned meta %+v, expected %+v", resp.Meta, expectedMeta)
	}
}

func TestRegions_ListRegionsMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/regions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"regions": [{"id":1},{"id":2}], "links":{"pages":{"next":"http://example.com/v2/regions/?page=2"}}}`)
	})

	_, resp, err := client.Regions.List(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 1)
}

func TestRegions_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"regions": [{"id":1},{"id":2}],
		"links":{
			"pages":{
				"next":"http://example.com/v2/regions/?page=3",
				"prev":"http://example.com/v2/regions/?page=1",
				"last":"http://example.com/v2/regions/?page=3",
				"first":"http://example.com/v2/regions/?page=1"
			}
		}
	}`

	mux.HandleFunc("/v2/regions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{Page: 2}
	_, resp, err := client.Regions.List(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestRegion_String(t *testing.T) {
	region := &Region{
		Slug:      "region",
		Name:      "Region",
		Sizes:     []string{"1", "2"},
		Available: true,
	}

	stringified := region.String()
	expected := `binarylane.Region{Slug:"region", Name:"Region", Sizes:["1" "2"], Available:true}`
	if expected != stringified {
		t.Errorf("Region.String returned %+v, expected %+v", stringified, expected)
	}
}
