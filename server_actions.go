package binarylane

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ActionRequest represents BinaryLane Action Request
type ActionRequest map[string]interface{}

// ServerActionsService is an interface for interfacing with the Server actions
// endpoints of the BinaryLane API
// See: https://api.binarylane.com.au/reference#server-actions
type ServerActionsService interface {
	Shutdown(context.Context, int) (*Action, *Response, error)
	ShutdownByTag(context.Context, string) ([]Action, *Response, error)
	PowerOff(context.Context, int) (*Action, *Response, error)
	PowerOffByTag(context.Context, string) ([]Action, *Response, error)
	PowerOn(context.Context, int) (*Action, *Response, error)
	PowerOnByTag(context.Context, string) ([]Action, *Response, error)
	PowerCycle(context.Context, int) (*Action, *Response, error)
	PowerCycleByTag(context.Context, string) ([]Action, *Response, error)
	Reboot(context.Context, int) (*Action, *Response, error)
	Restore(context.Context, int, int) (*Action, *Response, error)
	Resize(context.Context, int, string, bool) (*Action, *Response, error)
	Rename(context.Context, int, string) (*Action, *Response, error)
	Snapshot(context.Context, int, string) (*Action, *Response, error)
	SnapshotByTag(context.Context, string, string) ([]Action, *Response, error)
	EnableBackups(context.Context, int) (*Action, *Response, error)
	EnableBackupsByTag(context.Context, string) ([]Action, *Response, error)
	DisableBackups(context.Context, int) (*Action, *Response, error)
	DisableBackupsByTag(context.Context, string) ([]Action, *Response, error)
	PasswordReset(context.Context, int) (*Action, *Response, error)
	RebuildByImageID(context.Context, int, int) (*Action, *Response, error)
	RebuildByImageSlug(context.Context, int, string) (*Action, *Response, error)
	ChangeKernel(context.Context, int, int) (*Action, *Response, error)
	EnableIPv6(context.Context, int) (*Action, *Response, error)
	EnableIPv6ByTag(context.Context, string) ([]Action, *Response, error)
	EnablePrivateNetworking(context.Context, int) (*Action, *Response, error)
	EnablePrivateNetworkingByTag(context.Context, string) ([]Action, *Response, error)
	Get(context.Context, int, int) (*Action, *Response, error)
	GetByURI(context.Context, string) (*Action, *Response, error)
}

// ServerActionsServiceOp handles communication with the Server action related
// methods of the BinaryLane API.
type ServerActionsServiceOp struct {
	client *Client
}

var _ ServerActionsService = &ServerActionsServiceOp{}

// Shutdown a Server
func (s *ServerActionsServiceOp) Shutdown(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "shutdown"}
	return s.doAction(ctx, id, request)
}

// ShutdownByTag shuts down Servers matched by a Tag.
func (s *ServerActionsServiceOp) ShutdownByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "shutdown"}
	return s.doActionByTag(ctx, tag, request)
}

// PowerOff a Server
func (s *ServerActionsServiceOp) PowerOff(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "power_off"}
	return s.doAction(ctx, id, request)
}

// PowerOffByTag powers off Servers matched by a Tag.
func (s *ServerActionsServiceOp) PowerOffByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "power_off"}
	return s.doActionByTag(ctx, tag, request)
}

// PowerOn a Server
func (s *ServerActionsServiceOp) PowerOn(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "power_on"}
	return s.doAction(ctx, id, request)
}

// PowerOnByTag powers on Servers matched by a Tag.
func (s *ServerActionsServiceOp) PowerOnByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "power_on"}
	return s.doActionByTag(ctx, tag, request)
}

// PowerCycle a Server
func (s *ServerActionsServiceOp) PowerCycle(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "power_cycle"}
	return s.doAction(ctx, id, request)
}

// PowerCycleByTag power cycles Servers matched by a Tag.
func (s *ServerActionsServiceOp) PowerCycleByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "power_cycle"}
	return s.doActionByTag(ctx, tag, request)
}

// Reboot a Server
func (s *ServerActionsServiceOp) Reboot(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "reboot"}
	return s.doAction(ctx, id, request)
}

// Restore an image to a Server
func (s *ServerActionsServiceOp) Restore(ctx context.Context, id, imageID int) (*Action, *Response, error) {
	requestType := "restore"
	request := &ActionRequest{
		"type":  requestType,
		"image": float64(imageID),
	}
	return s.doAction(ctx, id, request)
}

// Resize a Server
func (s *ServerActionsServiceOp) Resize(ctx context.Context, id int, sizeSlug string, resizeDisk bool) (*Action, *Response, error) {
	requestType := "resize"
	request := &ActionRequest{
		"type": requestType,
		"size": sizeSlug,
		"disk": resizeDisk,
	}
	return s.doAction(ctx, id, request)
}

// Rename a Server
func (s *ServerActionsServiceOp) Rename(ctx context.Context, id int, name string) (*Action, *Response, error) {
	requestType := "rename"
	request := &ActionRequest{
		"type": requestType,
		"name": name,
	}
	return s.doAction(ctx, id, request)
}

// Snapshot a Server.
func (s *ServerActionsServiceOp) Snapshot(ctx context.Context, id int, name string) (*Action, *Response, error) {
	requestType := "snapshot"
	request := &ActionRequest{
		"type": requestType,
		"name": name,
	}
	return s.doAction(ctx, id, request)
}

// SnapshotByTag snapshots Servers matched by a Tag.
func (s *ServerActionsServiceOp) SnapshotByTag(ctx context.Context, tag string, name string) ([]Action, *Response, error) {
	requestType := "snapshot"
	request := &ActionRequest{
		"type": requestType,
		"name": name,
	}
	return s.doActionByTag(ctx, tag, request)
}

// EnableBackups enables backups for a Server.
func (s *ServerActionsServiceOp) EnableBackups(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "enable_backups"}
	return s.doAction(ctx, id, request)
}

// EnableBackupsByTag enables backups for Servers matched by a Tag.
func (s *ServerActionsServiceOp) EnableBackupsByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "enable_backups"}
	return s.doActionByTag(ctx, tag, request)
}

// DisableBackups disables backups for a Server.
func (s *ServerActionsServiceOp) DisableBackups(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "disable_backups"}
	return s.doAction(ctx, id, request)
}

// DisableBackupsByTag disables backups for Server matched by a Tag.
func (s *ServerActionsServiceOp) DisableBackupsByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "disable_backups"}
	return s.doActionByTag(ctx, tag, request)
}

// PasswordReset resets the password for a Server.
func (s *ServerActionsServiceOp) PasswordReset(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "password_reset"}
	return s.doAction(ctx, id, request)
}

// RebuildByImageID rebuilds a Server from an image with a given id.
func (s *ServerActionsServiceOp) RebuildByImageID(ctx context.Context, id, imageID int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "rebuild", "image": imageID}
	return s.doAction(ctx, id, request)
}

// RebuildByImageSlug rebuilds a Server from an Image matched by a given Slug.
func (s *ServerActionsServiceOp) RebuildByImageSlug(ctx context.Context, id int, slug string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "rebuild", "image": slug}
	return s.doAction(ctx, id, request)
}

// ChangeKernel changes the kernel for a Server.
func (s *ServerActionsServiceOp) ChangeKernel(ctx context.Context, id, kernelID int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "change_kernel", "kernel": kernelID}
	return s.doAction(ctx, id, request)
}

// EnableIPv6 enables IPv6 for a Server.
func (s *ServerActionsServiceOp) EnableIPv6(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "enable_ipv6"}
	return s.doAction(ctx, id, request)
}

// EnableIPv6ByTag enables IPv6 for Servers matched by a Tag.
func (s *ServerActionsServiceOp) EnableIPv6ByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "enable_ipv6"}
	return s.doActionByTag(ctx, tag, request)
}

// EnablePrivateNetworking enables private networking for a Server.
func (s *ServerActionsServiceOp) EnablePrivateNetworking(ctx context.Context, id int) (*Action, *Response, error) {
	request := &ActionRequest{"type": "enable_private_networking"}
	return s.doAction(ctx, id, request)
}

// EnablePrivateNetworkingByTag enables private networking for Servers matched by a Tag.
func (s *ServerActionsServiceOp) EnablePrivateNetworkingByTag(ctx context.Context, tag string) ([]Action, *Response, error) {
	request := &ActionRequest{"type": "enable_private_networking"}
	return s.doActionByTag(ctx, tag, request)
}

func (s *ServerActionsServiceOp) doAction(ctx context.Context, id int, request *ActionRequest) (*Action, *Response, error) {
	if id < 1 {
		return nil, nil, NewArgError("id", "cannot be less than 1")
	}

	if request == nil {
		return nil, nil, NewArgError("request", "request can't be nil")
	}

	path := serverActionPath(id)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

func (s *ServerActionsServiceOp) doActionByTag(ctx context.Context, tag string, request *ActionRequest) ([]Action, *Response, error) {
	if tag == "" {
		return nil, nil, NewArgError("tag", "cannot be empty")
	}

	if request == nil {
		return nil, nil, NewArgError("request", "request can't be nil")
	}

	path := serverActionPathByTag(tag)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Actions, resp, err
}

// Get an action for a particular Server by id.
func (s *ServerActionsServiceOp) Get(ctx context.Context, serverID, actionID int) (*Action, *Response, error) {
	if serverID < 1 {
		return nil, nil, NewArgError("serverID", "cannot be less than 1")
	}

	if actionID < 1 {
		return nil, nil, NewArgError("actionID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", serverActionPath(serverID), actionID)
	return s.get(ctx, path)
}

// GetByURI gets an action for a particular Server by id.
func (s *ServerActionsServiceOp) GetByURI(ctx context.Context, rawurl string) (*Action, *Response, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, nil, err
	}

	return s.get(ctx, u.Path)

}

func (s *ServerActionsServiceOp) get(ctx context.Context, path string) (*Action, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err

}

func serverActionPath(serverID int) string {
	return fmt.Sprintf("v2/servers/%d/actions", serverID)
}

func serverActionPathByTag(tag string) string {
	return fmt.Sprintf("v2/servers/actions?tag_name=%s", tag)
}
