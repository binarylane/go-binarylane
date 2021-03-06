package binarylane

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lbListJSONResponse = `
{
	"load_balancers":[
        {
            "id":3,
            "name":"example-lb-01",
            "ip":"46.214.185.203",
            "algorithm":"round_robin",
            "status":"active",
            "created_at":"2016-12-15T14:16:36Z",
            "forwarding_rules":[
                {
                    "entry_protocol":"https",
                    "entry_port":443,
                    "target_protocol":"http",
                    "target_port":80,
                    "certificate_id":"a-b-c"
                }
            ],
            "health_check":{
                "protocol":"http",
                "port":80,
                "path":"/index.html",
                "check_interval_seconds":10,
                "response_timeout_seconds":5,
                "healthy_threshold":5,
                "unhealthy_threshold":3
            },
            "sticky_sessions":{
                "type":"cookies",
                "cookie_name":"DO-LB",
                "cookie_ttl_seconds":5
            },
            "region":{
            	"name":"New York 1",
                "slug":"nyc1",
                "sizes":[
                    "512mb",
                    "1gb",
                    "2gb",
                    "4gb",
                    "8gb",
                    "16gb"
                ],
                "features":[
                    "private_networking",
                    "backups",
                    "ipv6",
                    "metadata",
                    "storage"
                ],
                "available":true
            },
            "server_ids":[
                2,
                21
            ]
        }
    ],
    "links":{
        "pages":{
            "last":"http://localhost:3001/v2/load_balancers?page=3&per_page=1",
            "next":"http://localhost:3001/v2/load_balancers?page=2&per_page=1"
        }
    },
    "meta":{
        "total":3
    }
}
`

var lbCreateJSONResponse = `
{
    "load_balancer":{
        "id":3,
        "name":"example-lb-01",
        "ip":"",
        "algorithm":"round_robin",
        "status":"new",
        "created_at":"2016-12-15T14:19:09Z",
        "forwarding_rules":[
            {
                "entry_protocol":"https",
                "entry_port":443,
                "target_protocol":"http",
                "target_port":80,
                "certificate_id":"a-b-c"
            },
            {
                "entry_protocol":"https",
                "entry_port":444,
                "target_protocol":"https",
                "target_port":443,
                "tls_passthrough":true
            }
        ],
        "health_check":{
            "protocol":"http",
            "port":80,
            "path":"/index.html",
            "check_interval_seconds":10,
            "response_timeout_seconds":5,
            "healthy_threshold":5,
            "unhealthy_threshold":3
        },
        "sticky_sessions":{
            "type":"cookies",
            "cookie_name":"DO-LB",
            "cookie_ttl_seconds":5
        },
        "region":{
            "name":"New York 1",
            "slug":"nyc1",
            "sizes":[
                "512mb",
                "1gb",
                "2gb",
                "4gb",
                "8gb",
                "16gb"
            ],
            "features":[
                "private_networking",
                "backups",
                "ipv6",
                "metadata",
                "storage"
            ],
            "available":true
		},
		"tags": ["my-tag"],
        "server_ids":[
            2,
            21
        ],
        "redirect_http_to_https":true,
        "vpc_id":2
    }
}
`

var lbGetJSONResponse = `
{
    "load_balancer":{
        "id":3,
        "name":"example-lb-01",
        "ip":"46.214.185.203",
        "algorithm":"round_robin",
        "status":"active",
        "created_at":"2016-12-15T14:16:36Z",
        "forwarding_rules":[
            {
                "entry_protocol":"https",
                "entry_port":443,
                "target_protocol":"http",
                "target_port":80,
                "certificate_id":"a-b-c"
            }
        ],
        "health_check":{
            "protocol":"http",
            "port":80,
            "path":"/index.html",
            "check_interval_seconds":10,
            "response_timeout_seconds":5,
            "healthy_threshold":5,
            "unhealthy_threshold":3
        },
        "sticky_sessions":{
            "type":"cookies",
            "cookie_name":"DO-LB",
            "cookie_ttl_seconds":5
        },
        "region":{
            "name":"New York 1",
            "slug":"nyc1",
            "sizes":[
                "512mb",
                "1gb",
                "2gb",
                "4gb",
                "8gb",
                "16gb"
            ],
            "features":[
                "private_networking",
                "backups",
                "ipv6",
                "metadata",
                "storage"
            ],
            "available":true
        },
        "server_ids":[
            2,
            21
        ]
    }
}
`

var lbUpdateJSONResponse = `
{
    "load_balancer":{
        "id":3,
        "name":"example-lb-01",
        "ip":"12.34.56.78",
        "algorithm":"least_connections",
        "status":"active",
        "created_at":"2016-12-15T14:19:09Z",
        "forwarding_rules":[
            {
                "entry_protocol":"http",
                "entry_port":80,
                "target_protocol":"http",
                "target_port":80
            },
            {
                "entry_protocol":"https",
                "entry_port":443,
                "target_protocol":"http",
                "target_port":80,
                "certificate_id":"a-b-c"
            }
        ],
        "health_check":{
            "protocol":"tcp",
            "port":80,
            "path":"",
            "check_interval_seconds":10,
            "response_timeout_seconds":5,
            "healthy_threshold":5,
            "unhealthy_threshold":3
        },
        "sticky_sessions":{
            "type":"none"
        },
        "region":{
            "name":"New York 1",
            "slug":"nyc1",
            "sizes":[
                "512mb",
                "1gb",
                "2gb",
                "4gb",
                "8gb",
                "16gb"
            ],
            "features":[
                "private_networking",
                "backups",
                "ipv6",
                "metadata",
                "storage"
            ],
            "available":true
        },
        "server_ids":[
            2,
            21
        ]
    }
}
`

func TestLoadBalancers_Get(t *testing.T) {
	setup()
	defer teardown()

	path := "/v2/load_balancers"
	loadBalancerID := 3
	path = fmt.Sprintf("%s/%d", path, loadBalancerID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, lbGetJSONResponse)
	})

	loadBalancer, _, err := client.LoadBalancers.Get(ctx, loadBalancerID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error: %v", err)
	}

	expected := &LoadBalancer{
		ID:        3,
		Name:      "example-lb-01",
		IP:        "46.214.185.203",
		Algorithm: "round_robin",
		Status:    "active",
		Created:   "2016-12-15T14:16:36Z",
		ForwardingRules: []ForwardingRule{
			{
				EntryProtocol:  "https",
				EntryPort:      443,
				TargetProtocol: "http",
				TargetPort:     80,
				CertificateID:  "a-b-c",
				TlsPassthrough: false,
			},
		},
		HealthCheck: &HealthCheck{
			Protocol:               "http",
			Port:                   80,
			Path:                   "/index.html",
			CheckIntervalSeconds:   10,
			ResponseTimeoutSeconds: 5,
			HealthyThreshold:       5,
			UnhealthyThreshold:     3,
		},
		StickySessions: &StickySessions{
			Type:             "cookies",
			CookieName:       "DO-LB",
			CookieTtlSeconds: 5,
		},
		Region: &Region{
			Slug:      "nyc1",
			Name:      "New York 1",
			Sizes:     []string{"512mb", "1gb", "2gb", "4gb", "8gb", "16gb"},
			Available: true,
			Features:  []string{"private_networking", "backups", "ipv6", "metadata", "storage"},
		},
		ServerIDs: []int{2, 21},
	}

	assert.Equal(t, expected, loadBalancer)
}

func TestLoadBalancers_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &LoadBalancerRequest{
		Name:      "example-lb-01",
		Algorithm: "round_robin",
		Region:    "nyc1",
		ForwardingRules: []ForwardingRule{
			{
				EntryProtocol:  "https",
				EntryPort:      443,
				TargetProtocol: "http",
				TargetPort:     80,
				CertificateID:  "a-b-c",
			},
		},
		HealthCheck: &HealthCheck{
			Protocol:               "http",
			Port:                   80,
			Path:                   "/index.html",
			CheckIntervalSeconds:   10,
			ResponseTimeoutSeconds: 5,
			UnhealthyThreshold:     3,
			HealthyThreshold:       5,
		},
		StickySessions: &StickySessions{
			Type:             "cookies",
			CookieName:       "DO-LB",
			CookieTtlSeconds: 5,
		},
		Tag:                 "my-tag",
		Tags:                []string{"my-tag"},
		ServerIDs:           []int{2, 21},
		RedirectHttpToHttps: true,
		VPCID:               2,
	}

	path := "/v2/load_balancers"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := new(LoadBalancerRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		assert.Equal(t, createRequest, v)

		fmt.Fprint(w, lbCreateJSONResponse)
	})

	loadBalancer, _, err := client.LoadBalancers.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("LoadBalancers.Create returned error: %v", err)
	}

	expected := &LoadBalancer{
		ID:        3,
		Name:      "example-lb-01",
		Algorithm: "round_robin",
		Status:    "new",
		Created:   "2016-12-15T14:19:09Z",
		ForwardingRules: []ForwardingRule{
			{
				EntryProtocol:  "https",
				EntryPort:      443,
				TargetProtocol: "http",
				TargetPort:     80,
				CertificateID:  "a-b-c",
				TlsPassthrough: false,
			},
			{
				EntryProtocol:  "https",
				EntryPort:      444,
				TargetProtocol: "https",
				TargetPort:     443,
				CertificateID:  "",
				TlsPassthrough: true,
			},
		},
		HealthCheck: &HealthCheck{
			Protocol:               "http",
			Port:                   80,
			Path:                   "/index.html",
			CheckIntervalSeconds:   10,
			ResponseTimeoutSeconds: 5,
			HealthyThreshold:       5,
			UnhealthyThreshold:     3,
		},
		StickySessions: &StickySessions{
			Type:             "cookies",
			CookieName:       "DO-LB",
			CookieTtlSeconds: 5,
		},
		Region: &Region{
			Slug:      "nyc1",
			Name:      "New York 1",
			Sizes:     []string{"512mb", "1gb", "2gb", "4gb", "8gb", "16gb"},
			Available: true,
			Features:  []string{"private_networking", "backups", "ipv6", "metadata", "storage"},
		},
		Tags:                []string{"my-tag"},
		ServerIDs:           []int{2, 21},
		RedirectHttpToHttps: true,
		VPCID:               2,
	}

	assert.Equal(t, expected, loadBalancer)
}

func TestLoadBalancers_Update(t *testing.T) {
	setup()
	defer teardown()

	updateRequest := &LoadBalancerRequest{
		Name:      "example-lb-01",
		Algorithm: "least_connections",
		Region:    "nyc1",
		ForwardingRules: []ForwardingRule{
			{
				EntryProtocol:  "http",
				EntryPort:      80,
				TargetProtocol: "http",
				TargetPort:     80,
			},
			{
				EntryProtocol:  "https",
				EntryPort:      443,
				TargetProtocol: "http",
				TargetPort:     80,
				CertificateID:  "a-b-c",
			},
		},
		HealthCheck: &HealthCheck{
			Protocol:               "tcp",
			Port:                   80,
			Path:                   "",
			CheckIntervalSeconds:   10,
			ResponseTimeoutSeconds: 5,
			UnhealthyThreshold:     3,
			HealthyThreshold:       5,
		},
		StickySessions: &StickySessions{
			Type: "none",
		},
		ServerIDs: []int{2, 21},
	}

	path := "/v2/load_balancers"
	loadBalancerID := 3
	path = fmt.Sprintf("%s/%d", path, loadBalancerID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := new(LoadBalancerRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, "PUT")
		assert.Equal(t, updateRequest, v)

		fmt.Fprint(w, lbUpdateJSONResponse)
	})

	loadBalancer, _, err := client.LoadBalancers.Update(ctx, loadBalancerID, updateRequest)
	if err != nil {
		t.Errorf("LoadBalancers.Update returned error: %v", err)
	}

	expected := &LoadBalancer{
		ID:        3,
		Name:      "example-lb-01",
		IP:        "12.34.56.78",
		Algorithm: "least_connections",
		Status:    "active",
		Created:   "2016-12-15T14:19:09Z",
		ForwardingRules: []ForwardingRule{
			{
				EntryProtocol:  "http",
				EntryPort:      80,
				TargetProtocol: "http",
				TargetPort:     80,
			},
			{
				EntryProtocol:  "https",
				EntryPort:      443,
				TargetProtocol: "http",
				TargetPort:     80,
				CertificateID:  "a-b-c",
			},
		},
		HealthCheck: &HealthCheck{
			Protocol:               "tcp",
			Port:                   80,
			Path:                   "",
			CheckIntervalSeconds:   10,
			ResponseTimeoutSeconds: 5,
			UnhealthyThreshold:     3,
			HealthyThreshold:       5,
		},
		StickySessions: &StickySessions{
			Type: "none",
		},
		Region: &Region{
			Slug:      "nyc1",
			Name:      "New York 1",
			Sizes:     []string{"512mb", "1gb", "2gb", "4gb", "8gb", "16gb"},
			Available: true,
			Features:  []string{"private_networking", "backups", "ipv6", "metadata", "storage"},
		},
		ServerIDs: []int{2, 21},
	}

	assert.Equal(t, expected, loadBalancer)
}

func TestLoadBalancers_List(t *testing.T) {
	setup()
	defer teardown()

	path := "/v2/load_balancers"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, lbListJSONResponse)
	})

	loadBalancers, resp, err := client.LoadBalancers.List(ctx, nil)

	if err != nil {
		t.Errorf("LoadBalancers.List returned error: %v", err)
	}

	expectedLBs := []LoadBalancer{
		{
			ID:        3,
			Name:      "example-lb-01",
			IP:        "46.214.185.203",
			Algorithm: "round_robin",
			Status:    "active",
			Created:   "2016-12-15T14:16:36Z",
			ForwardingRules: []ForwardingRule{
				{
					EntryProtocol:  "https",
					EntryPort:      443,
					TargetProtocol: "http",
					TargetPort:     80,
					CertificateID:  "a-b-c",
				},
			},
			HealthCheck: &HealthCheck{
				Protocol:               "http",
				Port:                   80,
				Path:                   "/index.html",
				CheckIntervalSeconds:   10,
				ResponseTimeoutSeconds: 5,
				HealthyThreshold:       5,
				UnhealthyThreshold:     3,
			},
			StickySessions: &StickySessions{
				Type:             "cookies",
				CookieName:       "DO-LB",
				CookieTtlSeconds: 5,
			},
			Region: &Region{
				Slug:      "nyc1",
				Name:      "New York 1",
				Sizes:     []string{"512mb", "1gb", "2gb", "4gb", "8gb", "16gb"},
				Available: true,
				Features:  []string{"private_networking", "backups", "ipv6", "metadata", "storage"},
			},
			ServerIDs: []int{2, 21},
		},
	}

	assert.Equal(t, expectedLBs, loadBalancers)

	expectedMeta := &Meta{Total: 3}
	assert.Equal(t, expectedMeta, resp.Meta)
}

func TestLoadBalancers_List_Pagination(t *testing.T) {
	setup()
	defer teardown()

	path := "/v2/load_balancers"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, map[string]string{"page": "2"})
		fmt.Fprint(w, lbListJSONResponse)
	})

	opts := &ListOptions{Page: 2}
	_, resp, err := client.LoadBalancers.List(ctx, opts)

	if err != nil {
		t.Errorf("LoadBalancers.List returned error: %v", err)
	}

	assert.Equal(t, "http://localhost:3001/v2/load_balancers?page=2&per_page=1", resp.Links.Pages.Next)
	assert.Equal(t, "http://localhost:3001/v2/load_balancers?page=3&per_page=1", resp.Links.Pages.Last)
}

func TestLoadBalancers_Delete(t *testing.T) {
	setup()
	defer teardown()

	lbID := 3
	path := "/v2/load_balancers"
	path = fmt.Sprintf("%s/%d", path, lbID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.LoadBalancers.Delete(ctx, lbID)

	if err != nil {
		t.Errorf("LoadBalancers.Delete returned error: %v", err)
	}
}

func TestLoadBalancers_AddServers(t *testing.T) {
	setup()
	defer teardown()

	serverIdsRequest := &serverIDsRequest{
		IDs: []int{42, 44},
	}

	lbID := 3
	path := fmt.Sprintf("/v2/load_balancers/%d/servers", lbID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := new(serverIDsRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		assert.Equal(t, serverIdsRequest, v)

		fmt.Fprint(w, nil)
	})

	_, err := client.LoadBalancers.AddServers(ctx, lbID, serverIdsRequest.IDs...)

	if err != nil {
		t.Errorf("LoadBalancers.AddServers returned error: %v", err)
	}
}

func TestLoadBalancers_RemoveServers(t *testing.T) {
	setup()
	defer teardown()

	serverIdsRequest := &serverIDsRequest{
		IDs: []int{2, 21},
	}

	lbID := 3
	path := fmt.Sprintf("/v2/load_balancers/%d/servers", lbID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := new(serverIDsRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodDelete)
		assert.Equal(t, serverIdsRequest, v)

		fmt.Fprint(w, nil)
	})

	_, err := client.LoadBalancers.RemoveServers(ctx, lbID, serverIdsRequest.IDs...)

	if err != nil {
		t.Errorf("LoadBalancers.RemoveServers returned error: %v", err)
	}
}

func TestLoadBalancers_AddForwardingRules(t *testing.T) {
	setup()
	defer teardown()

	frr := &forwardingRulesRequest{
		Rules: []ForwardingRule{
			{
				EntryProtocol:  "https",
				EntryPort:      444,
				TargetProtocol: "http",
				TargetPort:     81,
				CertificateID:  "b2abc00f-d3c4-426c-9f0b-b2f7a3ff7527",
			},
			{
				EntryProtocol:  "tcp",
				EntryPort:      8080,
				TargetProtocol: "tcp",
				TargetPort:     8081,
			},
		},
	}

	lbID := 3
	path := fmt.Sprintf("/v2/load_balancers/%d/forwarding_rules", lbID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := new(forwardingRulesRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		assert.Equal(t, frr, v)

		fmt.Fprint(w, nil)
	})

	_, err := client.LoadBalancers.AddForwardingRules(ctx, lbID, frr.Rules...)

	if err != nil {
		t.Errorf("LoadBalancers.AddForwardingRules returned error: %v", err)
	}
}

func TestLoadBalancers_RemoveForwardingRules(t *testing.T) {
	setup()
	defer teardown()

	frr := &forwardingRulesRequest{
		Rules: []ForwardingRule{
			{
				EntryProtocol:  "https",
				EntryPort:      444,
				TargetProtocol: "http",
				TargetPort:     81,
			},
			{
				EntryProtocol:  "tcp",
				EntryPort:      8080,
				TargetProtocol: "tcp",
				TargetPort:     8081,
			},
		},
	}

	lbID := 3
	path := fmt.Sprintf("/v2/load_balancers/%d/forwarding_rules", lbID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		v := new(forwardingRulesRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodDelete)
		assert.Equal(t, frr, v)

		fmt.Fprint(w, nil)
	})

	_, err := client.LoadBalancers.RemoveForwardingRules(ctx, lbID, frr.Rules...)

	if err != nil {
		t.Errorf("LoadBalancers.RemoveForwardingRules returned error: %v", err)
	}
}

func TestLoadBalancers_AsRequest(t *testing.T) {
	lb := &LoadBalancer{
		ID:        3,
		Name:      "test-loadbalancer",
		IP:        "10.0.0.1",
		SizeSlug:  "lb-small",
		Algorithm: "least_connections",
		Status:    "active",
		Created:   "2011-06-24T12:00:00Z",
		HealthCheck: &HealthCheck{
			Protocol:               "http",
			Port:                   80,
			Path:                   "/ping",
			CheckIntervalSeconds:   30,
			ResponseTimeoutSeconds: 10,
			HealthyThreshold:       3,
			UnhealthyThreshold:     3,
		},
		StickySessions: &StickySessions{
			Type:             "cookies",
			CookieName:       "nomnom",
			CookieTtlSeconds: 32,
		},
		Region: &Region{
			Slug: "lon1",
		},
		RedirectHttpToHttps:    true,
		EnableProxyProtocol:    true,
		EnableBackendKeepalive: true,
		VPCID:                  2,
	}
	lb.ServerIDs = make([]int, 1, 2)
	lb.ServerIDs[0] = 12345
	lb.ForwardingRules = make([]ForwardingRule, 1, 2)
	lb.ForwardingRules[0] = ForwardingRule{
		EntryProtocol:  "http",
		EntryPort:      80,
		TargetProtocol: "http",
		TargetPort:     80,
	}

	want := &LoadBalancerRequest{
		Name:      "test-loadbalancer",
		Algorithm: "least_connections",
		Region:    "lon1",
		SizeSlug:  "lb-small",
		ForwardingRules: []ForwardingRule{{
			EntryProtocol:  "http",
			EntryPort:      80,
			TargetProtocol: "http",
			TargetPort:     80,
		}},
		HealthCheck: &HealthCheck{
			Protocol:               "http",
			Port:                   80,
			Path:                   "/ping",
			CheckIntervalSeconds:   30,
			ResponseTimeoutSeconds: 10,
			HealthyThreshold:       3,
			UnhealthyThreshold:     3,
		},
		StickySessions: &StickySessions{
			Type:             "cookies",
			CookieName:       "nomnom",
			CookieTtlSeconds: 32,
		},
		ServerIDs:              []int{12345},
		RedirectHttpToHttps:    true,
		EnableProxyProtocol:    true,
		EnableBackendKeepalive: true,
		VPCID:                  2,
	}

	r := lb.AsRequest()
	assert.Equal(t, want, r)
	assert.False(t, r.HealthCheck == lb.HealthCheck, "HealthCheck points to same struct")
	assert.False(t, r.StickySessions == lb.StickySessions, "StickySessions points to same struct")

	r.ServerIDs = append(r.ServerIDs, 54321)
	r.ForwardingRules = append(r.ForwardingRules, ForwardingRule{
		EntryProtocol:  "https",
		EntryPort:      443,
		TargetProtocol: "https",
		TargetPort:     443,
		TlsPassthrough: true,
	})

	// Check that original LoadBalancer hasn't changed
	lb.ServerIDs = append(lb.ServerIDs, 13579)
	lb.ForwardingRules = append(lb.ForwardingRules, ForwardingRule{
		EntryProtocol:  "tcp",
		EntryPort:      587,
		TargetProtocol: "tcp",
		TargetPort:     587,
	})
	assert.Equal(t, []int{12345, 54321}, r.ServerIDs)
	assert.Equal(t, []ForwardingRule{
		{
			EntryProtocol:  "http",
			EntryPort:      80,
			TargetProtocol: "http",
			TargetPort:     80,
		},
		{
			EntryProtocol:  "https",
			EntryPort:      443,
			TargetProtocol: "https",
			TargetPort:     443,
			TlsPassthrough: true,
		},
	}, r.ForwardingRules)
}
