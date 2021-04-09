package binarylane

import (
	"context"
	"fmt"
	"net/http"
)

const loadBalancersBasePath = "/v2/load_balancers"
const forwardingRulesPath = "forwarding_rules"

const serversPath = "servers"

// LoadBalancersService is an interface for managing load balancers with the BinaryLane API.
// See: https://api.binarylane.com.au/reference#load-balancers
type LoadBalancersService interface {
	Get(context.Context, int) (*LoadBalancer, *Response, error)
	List(context.Context, *ListOptions) ([]LoadBalancer, *Response, error)
	Create(context.Context, *LoadBalancerRequest) (*LoadBalancer, *Response, error)
	Update(ctx context.Context, lbID int, lbr *LoadBalancerRequest) (*LoadBalancer, *Response, error)
	Delete(ctx context.Context, lbID int) (*Response, error)
	AddServers(ctx context.Context, lbID int, serverIDs ...int) (*Response, error)
	RemoveServers(ctx context.Context, lbID int, serverIDs ...int) (*Response, error)
	AddForwardingRules(ctx context.Context, lbID int, rules ...ForwardingRule) (*Response, error)
	RemoveForwardingRules(ctx context.Context, lbID int, rules ...ForwardingRule) (*Response, error)
}

// LoadBalancer represents a BinaryLane load balancer configuration.
// Tags can only be provided upon the creation of a Load Balancer.
type LoadBalancer struct {
	ID                     int              `json:"id,float64,omitempty"`
	Name                   string           `json:"name,omitempty"`
	IP                     string           `json:"ip,omitempty"`
	SizeSlug               string           `json:"size,omitempty"`
	Algorithm              string           `json:"algorithm,omitempty"`
	Status                 string           `json:"status,omitempty"`
	Created                string           `json:"created_at,omitempty"`
	ForwardingRules        []ForwardingRule `json:"forwarding_rules,omitempty"`
	HealthCheck            *HealthCheck     `json:"health_check,omitempty"`
	StickySessions         *StickySessions  `json:"sticky_sessions,omitempty"`
	Region                 *Region          `json:"region,omitempty"`
	ServerIDs              []int            `json:"server_ids,omitempty"`
	Tag                    string           `json:"tag,omitempty"`
	Tags                   []string         `json:"tags,omitempty"`
	RedirectHttpToHttps    bool             `json:"redirect_http_to_https,omitempty"`
	EnableProxyProtocol    bool             `json:"enable_proxy_protocol,omitempty"`
	EnableBackendKeepalive bool             `json:"enable_backend_keepalive,omitempty"`
	VPCID                  int              `json:"vpc_id,float64,omitempty"`
}

// String creates a human-readable description of a LoadBalancer.
func (l LoadBalancer) String() string {
	return Stringify(l)
}

// URN returns the load balancer ID in a valid BL API URN form.
func (l LoadBalancer) URN() string {
	return ToURN("LoadBalancer", l.ID)
}

// AsRequest creates a LoadBalancerRequest that can be submitted to Update with the current values of the LoadBalancer.
// Modifying the returned LoadBalancerRequest will not modify the original LoadBalancer.
func (l LoadBalancer) AsRequest() *LoadBalancerRequest {
	r := LoadBalancerRequest{
		Name:                   l.Name,
		Algorithm:              l.Algorithm,
		SizeSlug:               l.SizeSlug,
		ForwardingRules:        append([]ForwardingRule(nil), l.ForwardingRules...),
		ServerIDs:              append([]int(nil), l.ServerIDs...),
		Tag:                    l.Tag,
		RedirectHttpToHttps:    l.RedirectHttpToHttps,
		EnableProxyProtocol:    l.EnableProxyProtocol,
		EnableBackendKeepalive: l.EnableBackendKeepalive,
		HealthCheck:            l.HealthCheck,
		VPCID:                  l.VPCID,
	}

	if l.HealthCheck != nil {
		r.HealthCheck = &HealthCheck{}
		*r.HealthCheck = *l.HealthCheck
	}
	if l.StickySessions != nil {
		r.StickySessions = &StickySessions{}
		*r.StickySessions = *l.StickySessions
	}
	if l.Region != nil {
		r.Region = l.Region.Slug
	}
	return &r
}

// ForwardingRule represents load balancer forwarding rules.
type ForwardingRule struct {
	EntryProtocol  string `json:"entry_protocol,omitempty"`
	EntryPort      int    `json:"entry_port,omitempty"`
	TargetProtocol string `json:"target_protocol,omitempty"`
	TargetPort     int    `json:"target_port,omitempty"`
	CertificateID  string `json:"certificate_id,omitempty"`
	TlsPassthrough bool   `json:"tls_passthrough,omitempty"`
}

// String creates a human-readable description of a ForwardingRule.
func (f ForwardingRule) String() string {
	return Stringify(f)
}

// HealthCheck represents optional load balancer health check rules.
type HealthCheck struct {
	Protocol               string `json:"protocol,omitempty"`
	Port                   int    `json:"port,omitempty"`
	Path                   string `json:"path,omitempty"`
	CheckIntervalSeconds   int    `json:"check_interval_seconds,omitempty"`
	ResponseTimeoutSeconds int    `json:"response_timeout_seconds,omitempty"`
	HealthyThreshold       int    `json:"healthy_threshold,omitempty"`
	UnhealthyThreshold     int    `json:"unhealthy_threshold,omitempty"`
}

// String creates a human-readable description of a HealthCheck.
func (h HealthCheck) String() string {
	return Stringify(h)
}

// StickySessions represents optional load balancer session affinity rules.
type StickySessions struct {
	Type             string `json:"type,omitempty"`
	CookieName       string `json:"cookie_name,omitempty"`
	CookieTtlSeconds int    `json:"cookie_ttl_seconds,omitempty"`
}

// String creates a human-readable description of a StickySessions instance.
func (s StickySessions) String() string {
	return Stringify(s)
}

// LoadBalancerRequest represents the configuration to be applied to an existing or a new load balancer.
type LoadBalancerRequest struct {
	Name                   string           `json:"name,omitempty"`
	Algorithm              string           `json:"algorithm,omitempty"`
	Region                 string           `json:"region,omitempty"`
	SizeSlug               string           `json:"size,omitempty"`
	ForwardingRules        []ForwardingRule `json:"forwarding_rules,omitempty"`
	HealthCheck            *HealthCheck     `json:"health_check,omitempty"`
	StickySessions         *StickySessions  `json:"sticky_sessions,omitempty"`
	ServerIDs              []int            `json:"server_ids,omitempty"`
	Tag                    string           `json:"tag,omitempty"`
	Tags                   []string         `json:"tags,omitempty"`
	RedirectHttpToHttps    bool             `json:"redirect_http_to_https,omitempty"`
	EnableProxyProtocol    bool             `json:"enable_proxy_protocol,omitempty"`
	EnableBackendKeepalive bool             `json:"enable_backend_keepalive,omitempty"`
	VPCID                  int              `json:"vpc_id,omitempty"`
}

// String creates a human-readable description of a LoadBalancerRequest.
func (l LoadBalancerRequest) String() string {
	return Stringify(l)
}

type forwardingRulesRequest struct {
	Rules []ForwardingRule `json:"forwarding_rules,omitempty"`
}

func (l forwardingRulesRequest) String() string {
	return Stringify(l)
}

type serverIDsRequest struct {
	IDs []int `json:"server_ids,omitempty"`
}

func (l serverIDsRequest) String() string {
	return Stringify(l)
}

type loadBalancersRoot struct {
	LoadBalancers []LoadBalancer `json:"load_balancers"`
	Links         *Links         `json:"links"`
	Meta          *Meta          `json:"meta"`
}

type loadBalancerRoot struct {
	LoadBalancer *LoadBalancer `json:"load_balancer"`
}

// LoadBalancersServiceOp handles communication with load balancer-related methods of the BinaryLane API.
type LoadBalancersServiceOp struct {
	client *Client
}

var _ LoadBalancersService = &LoadBalancersServiceOp{}

// Get an existing load balancer by its identifier.
func (l *LoadBalancersServiceOp) Get(ctx context.Context, lbID int) (*LoadBalancer, *Response, error) {
	path := fmt.Sprintf("%s/%d", loadBalancersBasePath, lbID)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancerRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LoadBalancer, resp, err
}

// List load balancers, with optional pagination.
func (l *LoadBalancersServiceOp) List(ctx context.Context, opt *ListOptions) ([]LoadBalancer, *Response, error) {
	path, err := addOptions(loadBalancersBasePath, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancersRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.LoadBalancers, resp, err
}

// Create a new load balancer with a given configuration.
func (l *LoadBalancersServiceOp) Create(ctx context.Context, lbr *LoadBalancerRequest) (*LoadBalancer, *Response, error) {
	req, err := l.client.NewRequest(ctx, http.MethodPost, loadBalancersBasePath, lbr)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancerRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LoadBalancer, resp, err
}

// Update an existing load balancer with new configuration.
func (l *LoadBalancersServiceOp) Update(ctx context.Context, lbID int, lbr *LoadBalancerRequest) (*LoadBalancer, *Response, error) {
	path := fmt.Sprintf("%s/%d", loadBalancersBasePath, lbID)

	req, err := l.client.NewRequest(ctx, "PUT", path, lbr)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancerRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LoadBalancer, resp, err
}

// Delete a load balancer by its identifier.
func (l *LoadBalancersServiceOp) Delete(ctx context.Context, ldID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", loadBalancersBasePath, ldID)

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// AddServers adds servers to a load balancer.
func (l *LoadBalancersServiceOp) AddServers(ctx context.Context, lbID int, serverIDs ...int) (*Response, error) {
	path := fmt.Sprintf("%s/%d/%s", loadBalancersBasePath, lbID, serversPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, &serverIDsRequest{IDs: serverIDs})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// RemoveServers removes servers from a load balancer.
func (l *LoadBalancersServiceOp) RemoveServers(ctx context.Context, lbID int, serverIDs ...int) (*Response, error) {
	path := fmt.Sprintf("%s/%d/%s", loadBalancersBasePath, lbID, serversPath)

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, &serverIDsRequest{IDs: serverIDs})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// AddForwardingRules adds forwarding rules to a load balancer.
func (l *LoadBalancersServiceOp) AddForwardingRules(ctx context.Context, lbID int, rules ...ForwardingRule) (*Response, error) {
	path := fmt.Sprintf("%s/%d/%s", loadBalancersBasePath, lbID, forwardingRulesPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, &forwardingRulesRequest{Rules: rules})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// RemoveForwardingRules removes forwarding rules from a load balancer.
func (l *LoadBalancersServiceOp) RemoveForwardingRules(ctx context.Context, lbID int, rules ...ForwardingRule) (*Response, error) {
	path := fmt.Sprintf("%s/%d/%s", loadBalancersBasePath, lbID, forwardingRulesPath)

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, &forwardingRulesRequest{Rules: rules})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}
