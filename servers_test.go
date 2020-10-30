package binarylane

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestServers_ListServers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"servers": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			],
			"meta": {
				"total": 2
			}
		}`)
	})

	servers, resp, err := client.Servers.List(ctx, nil)
	if err != nil {
		t.Errorf("Servers.List returned error: %v", err)
	}

	expectedServers := []Server{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(servers, expectedServers) {
		t.Errorf("Servers.List\nServers: got=%#v\nwant=%#v", servers, expectedServers)
	}
	expectedMeta := &Meta{Total: 2}
	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Servers.List\nMeta: got=%#v\nwant=%#v", resp.Meta, expectedMeta)
	}
}

func TestServers_ListServersByTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("Servers.ListByTag did not request with a tag parameter")
		}

		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"servers": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			],
			"meta": {
				"total": 2
			}
		}`)
	})

	servers, resp, err := client.Servers.ListByTag(ctx, "testing-1", nil)
	if err != nil {
		t.Errorf("Servers.ListByTag returned error: %v", err)
	}

	expectedServers := []Server{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(servers, expectedServers) {
		t.Errorf("Servers.ListByTag returned servers %+v, expected %+v", servers, expectedServers)
	}
	expectedMeta := &Meta{Total: 2}
	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Servers.ListByTag returned meta %+v, expected %+v", resp.Meta, expectedMeta)
	}
}

func TestServers_ListServersMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dr := serversRoot{
			Servers: []Server{
				{ID: 1},
				{ID: 2},
			},
			Links: &Links{
				Pages: &Pages{Next: "http://example.com/v2/servers/?page=2"},
			},
		}

		b, err := json.Marshal(dr)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprint(w, string(b))
	})

	_, resp, err := client.Servers.List(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 1)
}

func TestServers_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"servers": [{"id":1},{"id":2}],
		"links":{
			"pages":{
				"next":"http://example.com/v2/servers/?page=3",
				"prev":"http://example.com/v2/servers/?page=1",
				"last":"http://example.com/v2/servers/?page=3",
				"first":"http://example.com/v2/servers/?page=1"
			}
		}
	}`

	mux.HandleFunc("/v2/servers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{Page: 2}
	_, resp, err := client.Servers.List(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestServers_GetServer(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"server":{"id":12345}}`)
	})

	servers, _, err := client.Servers.Get(ctx, 12345)
	if err != nil {
		t.Errorf("Server.Get returned error: %v", err)
	}

	expected := &Server{ID: 12345}
	if !reflect.DeepEqual(servers, expected) {
		t.Errorf("Servers.Get\n got=%#v\nwant=%#v", servers, expected)
	}
}

func TestServers_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &ServerCreateRequest{
		Name:   "name",
		Region: "region",
		Size:   "size",
		Image: ServerCreateImage{
			ID: 1,
		},
		Volumes: []ServerCreateVolume{
			{Name: "hello-im-a-volume"},
			{ID: "hello-im-another-volume"},
			{Name: "hello-im-still-a-volume", ID: "should be ignored due to Name"},
		},
		Tags:  []string{"one", "two"},
		VPCID: 2,
	}

	mux.HandleFunc("/v2/servers", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":               "name",
			"region":             "region",
			"size":               "size",
			"image":              float64(1),
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"volumes": []interface{}{
				map[string]interface{}{"name": "hello-im-a-volume"},
				map[string]interface{}{"id": "hello-im-another-volume"},
				map[string]interface{}{"name": "hello-im-still-a-volume"},
			},
			"tags":   []interface{}{"one", "two"},
			"vpc_id": float64(2),
		}
		jsonBlob := `
{
  "server": {
    "id": 1,
    "vpc_id": 2
  },
  "links": {
    "actions": [
      {
        "id": 1,
        "href": "http://example.com",
        "rel": "create"
      }
    ]
  }
}
`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, jsonBlob)
	})

	server, resp, err := client.Servers.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("Servers.Create returned error: %v", err)
	}

	if id := server.ID; id != 1 {
		t.Errorf("expected id '%d', received '%d'", 1, id)
	}

	vpcid := 2
	if id := server.VPCID; id != vpcid {
		t.Errorf("expected VPC id '%d', received '%d'", vpcid, id)
	}

	if a := resp.Links.Actions[0]; a.ID != 1 {
		t.Errorf("expected action id '%d', received '%d'", 1, a.ID)
	}
}

func TestServers_CreateMultiple(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &ServerMultiCreateRequest{
		Names:  []string{"name1", "name2"},
		Region: "region",
		Size:   "size",
		Image: ServerCreateImage{
			ID: 1,
		},
		Tags:  []string{"one", "two"},
		VPCID: 2,
	}

	mux.HandleFunc("/v2/servers", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"names":              []interface{}{"name1", "name2"},
			"region":             "region",
			"size":               "size",
			"image":              float64(1),
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"one", "two"},
			"vpc_id":             float64(2),
		}
		jsonBlob := `
{
  "servers": [
    {
      "id": 1,
	  "vpc_id": 2
    },
    {
      "id": 2,
	  "vpc_id": 2
    }
  ],
  "links": {
    "actions": [
      {
        "id": 1,
        "href": "http://example.com",
        "rel": "multiple_create"
      }
    ]
  }
}
`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body = %#v, expected %#v", v, expected)
		}

		fmt.Fprintf(w, jsonBlob)
	})

	servers, resp, err := client.Servers.CreateMultiple(ctx, createRequest)
	if err != nil {
		t.Errorf("Servers.CreateMultiple returned error: %v", err)
	}

	if id := servers[0].ID; id != 1 {
		t.Errorf("expected id '%d', received '%d'", 1, id)
	}
	if id := servers[1].ID; id != 2 {
		t.Errorf("expected id '%d', received '%d'", 2, id)
	}

	vpcid := 2
	if id := servers[0].VPCID; id != vpcid {
		t.Errorf("expected VPC id '%d', received '%d'", vpcid, id)
	}
	if id := servers[1].VPCID; id != vpcid {
		t.Errorf("expected VPC id '%d', received '%d'", vpcid, id)
	}

	if a := resp.Links.Actions[0]; a.ID != 1 {
		t.Errorf("expected action id '%d', received '%d'", 1, a.ID)
	}
}

func TestServers_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Servers.Delete(ctx, 12345)
	if err != nil {
		t.Errorf("Server.Delete returned error: %v", err)
	}
}

func TestServers_DestroyByTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("Servers.DeleteByTag did not request with a tag parameter")
		}

		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Servers.DeleteByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("Server.Delete returned error: %v", err)
	}
}

func TestServers_Kernels(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/12345/kernels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"kernels": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			],
			"meta": {
				"total": 2
			}
		}`)
	})

	opt := &ListOptions{Page: 2}
	kernels, resp, err := client.Servers.Kernels(ctx, 12345, opt)
	if err != nil {
		t.Errorf("Servers.Kernels returned error: %v", err)
	}

	expectedKernels := []Kernel{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(kernels, expectedKernels) {
		t.Errorf("Servers.Kernels\nKernels got=%#v\nwant=%#v", kernels, expectedKernels)
	}
	expectedMeta := &Meta{Total: 2}
	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Servers.Kernels\nMeta: got=%#v\nwant=%#v", resp.Meta, expectedMeta)
	}
}

func TestServers_Snapshots(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/12345/snapshots", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"snapshots": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			],
			"meta": {
				"total": 2
			}
		}`)
	})

	opt := &ListOptions{Page: 2}
	snapshots, resp, err := client.Servers.Snapshots(ctx, 12345, opt)
	if err != nil {
		t.Errorf("Servers.Snapshots returned error: %v", err)
	}

	expectedSnapshots := []Image{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(snapshots, expectedSnapshots) {
		t.Errorf("Servers.Snapshots\nSnapshots got=%#v\nwant=%#v", snapshots, expectedSnapshots)
	}
	expectedMeta := &Meta{Total: 2}
	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Servers.Snapshots\nMeta: got=%#v\nwant=%#v", resp.Meta, expectedMeta)
	}
}

func TestServers_Backups(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/12345/backups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"backups": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			],
			"meta": {
				"total": 2
			}
		}`)
	})

	opt := &ListOptions{Page: 2}
	backups, resp, err := client.Servers.Backups(ctx, 12345, opt)
	if err != nil {
		t.Errorf("Servers.Backups returned error: %v", err)
	}

	expectedBackups := []Image{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(backups, expectedBackups) {
		t.Errorf("Servers.Backups\nBackups got=%#v\nwant=%#v", backups, expectedBackups)
	}
	expectedMeta := &Meta{Total: 2}
	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Servers.Backups\nMeta: got=%#v\nwant=%#v", resp.Meta, expectedMeta)
	}
}

func TestServers_Actions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/12345/actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"actions": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			],
			"meta": {
				"total": 2
			}
		}`)
	})

	opt := &ListOptions{Page: 2}
	actions, resp, err := client.Servers.Actions(ctx, 12345, opt)
	if err != nil {
		t.Errorf("Servers.Actions returned error: %v", err)
	}

	expectedActions := []Action{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(actions, expectedActions) {
		t.Errorf("Servers.Actions\nActions got=%#v\nwant=%#v", actions, expectedActions)
	}
	expectedMeta := &Meta{Total: 2}
	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Servers.Actions\nMeta: got=%#v\nwant=%#v", resp.Meta, expectedMeta)
	}
}

func TestServers_Neighbors(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/12345/neighbors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"servers": [{"id":1},{"id":2}]}`)
	})

	neighbors, _, err := client.Servers.Neighbors(ctx, 12345)
	if err != nil {
		t.Errorf("Servers.Neighbors returned error: %v", err)
	}

	expected := []Server{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(neighbors, expected) {
		t.Errorf("Servers.Neighbors\n got=%#v\nwant=%#v", neighbors, expected)
	}
}

func TestNetworkV4_String(t *testing.T) {
	network := &NetworkV4{
		IPAddress: "192.168.1.2",
		Netmask:   "255.255.255.0",
		Gateway:   "192.168.1.1",
	}

	stringified := network.String()
	expected := `binarylane.NetworkV4{IPAddress:"192.168.1.2", Netmask:"255.255.255.0", Gateway:"192.168.1.1", Type:""}`
	if expected != stringified {
		t.Errorf("NetworkV4.String\n got=%#v\nwant=%#v", stringified, expected)
	}

}

func TestNetworkV6_String(t *testing.T) {
	network := &NetworkV6{
		IPAddress: "2604:A880:0800:0010:0000:0000:02DD:4001",
		Netmask:   64,
		Gateway:   "2604:A880:0800:0010:0000:0000:0000:0001",
	}
	stringified := network.String()
	expected := `binarylane.NetworkV6{IPAddress:"2604:A880:0800:0010:0000:0000:02DD:4001", Netmask:64, Gateway:"2604:A880:0800:0010:0000:0000:0000:0001", Type:""}`
	if expected != stringified {
		t.Errorf("NetworkV6.String\n got=%#v\nwant=%#v", stringified, expected)
	}
}

func TestServers_IPMethods(t *testing.T) {
	var d Server

	ipv6 := "1000:1000:1000:1000:0000:0000:004D:B001"

	d.Networks = &Networks{
		V4: []NetworkV4{
			{IPAddress: "192.168.0.1", Type: "public"},
			{IPAddress: "10.0.0.1", Type: "private"},
		},
		V6: []NetworkV6{
			{IPAddress: ipv6, Type: "public"},
		},
	}

	ip, err := d.PublicIPv4()
	if err != nil {
		t.Errorf("unknown error")
	}

	if got, expected := ip, "192.168.0.1"; got != expected {
		t.Errorf("Server.PublicIPv4 returned %s; expected %s", got, expected)
	}

	ip, err = d.PrivateIPv4()
	if err != nil {
		t.Errorf("unknown error")
	}

	if got, expected := ip, "10.0.0.1"; got != expected {
		t.Errorf("Server.PrivateIPv4 returned %s; expected %s", got, expected)
	}

	ip, err = d.PublicIPv6()
	if err != nil {
		t.Errorf("unknown error")
	}

	if got, expected := ip, ipv6; got != expected {
		t.Errorf("Server.PublicIPv6 returned %s; expected %s", got, expected)
	}
}
