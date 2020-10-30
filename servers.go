package binarylane

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const serverBasePath = "v2/servers"

var errNoNetworks = errors.New("no networks have been defined")

// ServersService is an interface for interfacing with the Server
// endpoints of the BinaryLane API
// See: https://api.binarylane.com.au/reference#servers
type ServersService interface {
	List(context.Context, *ListOptions) ([]Server, *Response, error)
	ListByTag(context.Context, string, *ListOptions) ([]Server, *Response, error)
	Get(context.Context, int) (*Server, *Response, error)
	Create(context.Context, *ServerCreateRequest) (*Server, *Response, error)
	CreateMultiple(context.Context, *ServerMultiCreateRequest) ([]Server, *Response, error)
	Delete(context.Context, int) (*Response, error)
	DeleteByTag(context.Context, string) (*Response, error)
	Kernels(context.Context, int, *ListOptions) ([]Kernel, *Response, error)
	Snapshots(context.Context, int, *ListOptions) ([]Image, *Response, error)
	Backups(context.Context, int, *ListOptions) ([]Image, *Response, error)
	Actions(context.Context, int, *ListOptions) ([]Action, *Response, error)
	Neighbors(context.Context, int) ([]Server, *Response, error)
}

// ServersServiceOp handles communication with the Server related methods of the
// BinaryLane API.
type ServersServiceOp struct {
	client *Client
}

var _ ServersService = &ServersServiceOp{}

// Server represents a BinaryLane Server
type Server struct {
	ID               int           `json:"id,float64,omitempty"`
	Name             string        `json:"name,omitempty"`
	Memory           int           `json:"memory,omitempty"`
	Vcpus            int           `json:"vcpus,omitempty"`
	Disk             int           `json:"disk,omitempty"`
	Region           *Region       `json:"region,omitempty"`
	Image            *Image        `json:"image,omitempty"`
	Size             *Size         `json:"size,omitempty"`
	SizeSlug         string        `json:"size_slug,omitempty"`
	BackupIDs        []int         `json:"backup_ids,omitempty"`
	NextBackupWindow *BackupWindow `json:"next_backup_window,omitempty"`
	SnapshotIDs      []int         `json:"snapshot_ids,omitempty"`
	Features         []string      `json:"features,omitempty"`
	Locked           bool          `json:"locked,bool,omitempty"`
	Status           string        `json:"status,omitempty"`
	Networks         *Networks     `json:"networks,omitempty"`
	Created          string        `json:"created_at,omitempty"`
	Kernel           *Kernel       `json:"kernel,omitempty"`
	Tags             []string      `json:"tags,omitempty"`
	VolumeIDs        []string      `json:"volume_ids"`
	VPCID            int           `json:"vpc_id,float64,omitempty"`
}

// PublicIPv4 returns the public IPv4 address for the Server.
func (s *Server) PublicIPv4() (string, error) {
	if s.Networks == nil {
		return "", errNoNetworks
	}

	for _, v4 := range s.Networks.V4 {
		if v4.Type == "public" {
			return v4.IPAddress, nil
		}
	}

	return "", nil
}

// PrivateIPv4 returns the private IPv4 address for the Server.
func (s *Server) PrivateIPv4() (string, error) {
	if s.Networks == nil {
		return "", errNoNetworks
	}

	for _, v4 := range s.Networks.V4 {
		if v4.Type == "private" {
			return v4.IPAddress, nil
		}
	}

	return "", nil
}

// PublicIPv6 returns the public IPv6 address for the Server.
func (s *Server) PublicIPv6() (string, error) {
	if s.Networks == nil {
		return "", errNoNetworks
	}

	for _, v6 := range s.Networks.V6 {
		if v6.Type == "public" {
			return v6.IPAddress, nil
		}
	}

	return "", nil
}

// Kernel object
type Kernel struct {
	ID      int    `json:"id,float64,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// BackupWindow object
type BackupWindow struct {
	Start *Timestamp `json:"start,omitempty"`
	End   *Timestamp `json:"end,omitempty"`
}

// Convert Server to a string
func (s Server) String() string {
	return Stringify(s)
}

// URN returns the servers ID in a valid BL API URN form.
func (s Server) URN() string {
	return ToURN("Server", s.ID)
}

type serverRoot struct {
	Server *Server `json:"server"`
	Links  *Links  `json:"links,omitempty"`
}

type serversRoot struct {
	Servers []Server `json:"servers"`
	Links   *Links   `json:"links"`
	Meta    *Meta    `json:"meta"`
}

type kernelsRoot struct {
	Kernels []Kernel `json:"kernels,omitempty"`
	Links   *Links   `json:"links"`
	Meta    *Meta    `json:"meta"`
}

type serverSnapshotsRoot struct {
	Snapshots []Image `json:"snapshots,omitempty"`
	Links     *Links  `json:"links"`
	Meta      *Meta   `json:"meta"`
}

type backupsRoot struct {
	Backups []Image `json:"backups,omitempty"`
	Links   *Links  `json:"links"`
	Meta    *Meta   `json:"meta"`
}

// ServerCreateImage identifies an image for the create request. It prefers slug over ID.
type ServerCreateImage struct {
	ID   int
	Slug string
}

// MarshalJSON returns either the slug or id of the image. It returns the id
// if the slug is empty.
func (d ServerCreateImage) MarshalJSON() ([]byte, error) {
	if d.Slug != "" {
		return json.Marshal(d.Slug)
	}

	return json.Marshal(d.ID)
}

// ServerCreateVolume identifies a volume to attach for the create request. It
// prefers Name over ID,
type ServerCreateVolume struct {
	ID   string
	Name string
}

// MarshalJSON returns an object with either the name or id of the volume. It
// returns the id if the name is empty.
func (d ServerCreateVolume) MarshalJSON() ([]byte, error) {
	if d.Name != "" {
		return json.Marshal(struct {
			Name string `json:"name"`
		}{Name: d.Name})
	}

	return json.Marshal(struct {
		ID string `json:"id"`
	}{ID: d.ID})
}

// ServerCreateSSHKey identifies a SSH Key for the create request. It prefers fingerprint over ID.
type ServerCreateSSHKey struct {
	ID          int
	Fingerprint string
}

// MarshalJSON returns either the fingerprint or id of the ssh key. It returns
// the id if the fingerprint is empty.
func (d ServerCreateSSHKey) MarshalJSON() ([]byte, error) {
	if d.Fingerprint != "" {
		return json.Marshal(d.Fingerprint)
	}

	return json.Marshal(d.ID)
}

// ServerCreateRequest represents a request to create a Server.
type ServerCreateRequest struct {
	Name              string               `json:"name"`
	Region            string               `json:"region"`
	Size              string               `json:"size"`
	Image             ServerCreateImage    `json:"image"`
	SSHKeys           []ServerCreateSSHKey `json:"ssh_keys"`
	Backups           bool                 `json:"backups"`
	IPv6              bool                 `json:"ipv6"`
	PrivateNetworking bool                 `json:"private_networking"`
	Monitoring        bool                 `json:"monitoring"`
	UserData          string               `json:"user_data,omitempty"`
	Volumes           []ServerCreateVolume `json:"volumes,omitempty"`
	Tags              []string             `json:"tags"`
	VPCID             int                  `json:"vpc_id,omitempty"`
}

// ServerMultiCreateRequest is a request to create multiple Servers.
type ServerMultiCreateRequest struct {
	Names             []string             `json:"names"`
	Region            string               `json:"region"`
	Size              string               `json:"size"`
	Image             ServerCreateImage    `json:"image"`
	SSHKeys           []ServerCreateSSHKey `json:"ssh_keys"`
	Backups           bool                 `json:"backups"`
	IPv6              bool                 `json:"ipv6"`
	PrivateNetworking bool                 `json:"private_networking"`
	Monitoring        bool                 `json:"monitoring"`
	UserData          string               `json:"user_data,omitempty"`
	Tags              []string             `json:"tags"`
	VPCID             int                  `json:"vpc_id,omitempty"`
}

func (d ServerCreateRequest) String() string {
	return Stringify(d)
}

func (d ServerMultiCreateRequest) String() string {
	return Stringify(d)
}

// Networks represents the Server's Networks.
type Networks struct {
	V4 []NetworkV4 `json:"v4,omitempty"`
	V6 []NetworkV6 `json:"v6,omitempty"`
}

// NetworkV4 represents a BinaryLane IPv4 Network.
type NetworkV4 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (n NetworkV4) String() string {
	return Stringify(n)
}

// NetworkV6 represents a BinaryLane IPv6 network.
type NetworkV6 struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   int    `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (n NetworkV6) String() string {
	return Stringify(n)
}

// Performs a list request given a path.
func (s *ServersServiceOp) list(ctx context.Context, path string) ([]Server, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serversRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Servers, resp, err
}

// List all Servers.
func (s *ServersServiceOp) List(ctx context.Context, opt *ListOptions) ([]Server, *Response, error) {
	path := serverBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

// ListByTag lists all Servers matched by a Tag.
func (s *ServersServiceOp) ListByTag(ctx context.Context, tag string, opt *ListOptions) ([]Server, *Response, error) {
	path := fmt.Sprintf("%s?tag_name=%s", serverBasePath, tag)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

// Get individual Server.
func (s *ServersServiceOp) Get(ctx context.Context, serverID int) (*Server, *Response, error) {
	if serverID < 1 {
		return nil, nil, NewArgError("serverID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", serverBasePath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serverRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Server, resp, err
}

// Create Server
func (s *ServersServiceOp) Create(ctx context.Context, createRequest *ServerCreateRequest) (*Server, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := serverBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(serverRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Server, resp, err
}

// CreateMultiple creates multiple Servers.
func (s *ServersServiceOp) CreateMultiple(ctx context.Context, createRequest *ServerMultiCreateRequest) ([]Server, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := serverBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(serversRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Servers, resp, err
}

// Performs a delete request given a path
func (s *ServersServiceOp) delete(ctx context.Context, path string) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// Delete Server.
func (s *ServersServiceOp) Delete(ctx context.Context, serverID int) (*Response, error) {
	if serverID < 1 {
		return nil, NewArgError("serverID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", serverBasePath, serverID)

	return s.delete(ctx, path)
}

// DeleteByTag deletes Servers matched by a Tag.
func (s *ServersServiceOp) DeleteByTag(ctx context.Context, tag string) (*Response, error) {
	if tag == "" {
		return nil, NewArgError("tag", "cannot be empty")
	}

	path := fmt.Sprintf("%s?tag_name=%s", serverBasePath, tag)

	return s.delete(ctx, path)
}

// Kernels lists kernels available for a Server.
func (s *ServersServiceOp) Kernels(ctx context.Context, serverID int, opt *ListOptions) ([]Kernel, *Response, error) {
	if serverID < 1 {
		return nil, nil, NewArgError("serverID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/kernels", serverBasePath, serverID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(kernelsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Kernels, resp, err
}

// Actions lists the actions for a Server.
func (s *ServersServiceOp) Actions(ctx context.Context, serverID int, opt *ListOptions) ([]Action, *Response, error) {
	if serverID < 1 {
		return nil, nil, NewArgError("serverID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/actions", serverBasePath, serverID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Actions, resp, err
}

// Backups lists the backups for a Server.
func (s *ServersServiceOp) Backups(ctx context.Context, serverID int, opt *ListOptions) ([]Image, *Response, error) {
	if serverID < 1 {
		return nil, nil, NewArgError("serverID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/backups", serverBasePath, serverID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Backups, resp, err
}

// Snapshots lists the snapshots available for a Server.
func (s *ServersServiceOp) Snapshots(ctx context.Context, serverID int, opt *ListOptions) ([]Image, *Response, error) {
	if serverID < 1 {
		return nil, nil, NewArgError("serverID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/snapshots", serverBasePath, serverID)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serverSnapshotsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Snapshots, resp, err
}

// Neighbors lists the neighbors for a Server.
func (s *ServersServiceOp) Neighbors(ctx context.Context, serverID int) ([]Server, *Response, error) {
	if serverID < 1 {
		return nil, nil, NewArgError("serverID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d/neighbors", serverBasePath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serversRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Servers, resp, err
}

func (s *ServersServiceOp) serverActionStatus(ctx context.Context, uri string) (string, error) {
	action, _, err := s.client.ServerActions.GetByURI(ctx, uri)

	if err != nil {
		return "", err
	}

	return action.Status, nil
}
