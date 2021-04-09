package binarylane

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestServerActions_Shutdown(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "shutdown",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.Shutdown(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.Shutdown returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Shutdown returned %+v, expected %+v", action, expected)
	}
}

func TestServerActions_ShutdownByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "shutdown",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.ShutdownByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.ShutdownByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.ShutdownByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.ShutdownByTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_PowerOff(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "power_off",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.PowerOff(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.PowerOff returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Poweroff returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_PowerOffByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "power_off",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.PowerOffByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.PowerOffByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.PowerOffByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.PoweroffByTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_PowerOn(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "power_on",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.PowerOn(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.PowerOn returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.PowerOn returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_PowerOnByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "power_on",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.PowerOnByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.PowerOnByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.PowerOnByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.PowerOnByTag returned %+v, expected %+v", action, expected)
	}
}
func TestServerAction_Reboot(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "reboot",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)

	})

	action, _, err := client.ServerActions.Reboot(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.Reboot returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Reboot returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_Restore(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type":  "restore",
		"image": float64(1),
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)

	})

	action, _, err := client.ServerActions.Restore(ctx, 1, 1)
	if err != nil {
		t.Errorf("ServerActions.Restore returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Restore returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_Resize(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "resize",
		"size": "1024mb",
		"disk": true,
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)

	})

	action, _, err := client.ServerActions.Resize(ctx, 1, "1024mb", true)
	if err != nil {
		t.Errorf("ServerActions.Resize returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Resize returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_Rename(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "rename",
		"name": "Server-Name",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.Rename(ctx, 1, "Server-Name")
	if err != nil {
		t.Errorf("ServerActions.Rename returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Rename returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_PowerCycle(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "power_cycle",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)

	})

	action, _, err := client.ServerActions.PowerCycle(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.PowerCycle returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.PowerCycle returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_PowerCycleByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "power_cycle",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.PowerCycleByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.PowerCycleByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.PowerCycleByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.PowerCycleByTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_Snapshot(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "snapshot",
		"name": "Image-Name",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.Snapshot(ctx, 1, "Image-Name")
	if err != nil {
		t.Errorf("ServerActions.Snapshot returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Snapshot returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_SnapshotByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "snapshot",
		"name": "Image-Name",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.SnapshotByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.SnapshotByTag(ctx, "testing-1", "Image-Name")
	if err != nil {
		t.Errorf("ServerActions.SnapshotByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.SnapshotByTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_EnableBackups(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "enable_backups",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.EnableBackups(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.EnableBackups returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.EnableBackups returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_EnableBackupsByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "enable_backups",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.EnableBackupByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.EnableBackupsByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.EnableBackupsByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.EnableBackupsByTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_DisableBackups(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "disable_backups",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.DisableBackups(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.DisableBackups returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.DisableBackups returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_DisableBackupsByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "disable_backups",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.DisableBackupsByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.DisableBackupsByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.DisableBackupsByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.DisableBackupsByTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_PasswordReset(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "password_reset",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.PasswordReset(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.PasswordReset returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.PasswordReset returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_RebuildByImageID(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type":  "rebuild",
		"image": float64(2),
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = \n%#v, expected \n%#v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.RebuildByImageID(ctx, 1, 2)
	if err != nil {
		t.Errorf("ServerActions.RebuildByImageID returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.RebuildByImageID returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_RebuildByImageSlug(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type":  "rebuild",
		"image": "Image-Name",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.RebuildByImageSlug(ctx, 1, "Image-Name")
	if err != nil {
		t.Errorf("ServerActions.RebuildByImageSlug returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.RebuildByImageSlug returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_ChangeKernel(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type":   "change_kernel",
		"kernel": float64(2),
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.ChangeKernel(ctx, 1, 2)
	if err != nil {
		t.Errorf("ServerActions.ChangeKernel returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.ChangeKernel returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_EnableIPv6(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "enable_ipv6",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.EnableIPv6(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.EnableIPv6 returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.EnableIPv6 returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_EnableIPv6ByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "enable_ipv6",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.EnableIPv6ByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.EnableIPv6ByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.EnableIPv6ByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.EnableIPv6byTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_EnablePrivateNetworking(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "enable_private_networking",
	}

	mux.HandleFunc("/v2/servers/1/actions", func(w http.ResponseWriter, r *http.Request) {
		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.EnablePrivateNetworking(ctx, 1)
	if err != nil {
		t.Errorf("ServerActions.EnablePrivateNetworking returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.EnablePrivateNetworking returned %+v, expected %+v", action, expected)
	}
}

func TestServerAction_EnablePrivateNetworkingByTag(t *testing.T) {
	setup()
	defer teardown()

	request := &ActionRequest{
		"type": "enable_private_networking",
	}

	mux.HandleFunc("/v2/servers/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tag_name") != "testing-1" {
			t.Errorf("ServerActions.EnablePrivateNetworkingByTag did not request with a tag parameter")
		}

		v := new(ActionRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, request) {
			t.Errorf("Request body = %+v, expected %+v", v, request)
		}

		fmt.Fprint(w, `{"actions": [{"status":"in-progress"},{"status":"in-progress"}]}`)
	})

	action, _, err := client.ServerActions.EnablePrivateNetworkingByTag(ctx, "testing-1")
	if err != nil {
		t.Errorf("ServerActions.EnablePrivateNetworkingByTag returned error: %v", err)
	}

	expected := []Action{{Status: "in-progress"}, {Status: "in-progress"}}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.EnablePrivateNetworkingByTag returned %+v, expected %+v", action, expected)
	}
}

func TestServerActions_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/servers/123/actions/456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"action":{"status":"in-progress"}}`)
	})

	action, _, err := client.ServerActions.Get(ctx, 123, 456)
	if err != nil {
		t.Errorf("ServerActions.Get returned error: %v", err)
	}

	expected := &Action{Status: "in-progress"}
	if !reflect.DeepEqual(action, expected) {
		t.Errorf("ServerActions.Get returned %+v, expected %+v", action, expected)
	}
}
