package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/digitalocean/godo"
)

const CRED_FILE string = "/tmp/do_creds"
const DEVICE_PREFIX string = "/dev/disk/by-id/scsi-0DO_Volume_"

type Result struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Device  string `json:"device,omitempty"`
}

type AttachOptions struct {
	ApiKey   string `json:"kubernetes.io/secret/apiKey"`
	VolumeId string `json:"volumeId"`
}

func finish(result Result, code int) {
	res, err := json.Marshal(result)
	if err != nil {
		fmt.Println("{\"status\":\"Failure\",\"message\":\"JSON error\"}")
	} else {
		fmt.Println(string(res))
	}
	os.Exit(code)
}

func succeed(a ...interface{}) {
	finish(Result{
		Status:  "Success",
		Message: fmt.Sprint(a...),
	}, 0)
}

func fail(a ...interface{}) {
	finish(Result{
		Status:  "Failure",
		Message: fmt.Sprint(a...),
	}, 1)
}

func awaitAction(client *godo.Client, volId string, action *godo.Action) {
	start := time.Now()
	actionId := action.ID
	for {
		switch action.Status {
		case "errored":
			fail("attach failed: ", action.String())

		case "completed":
			return

		case "in-progress":

		default:
			fail("Unknown action status ", action.Status)
		}
		if 10*time.Minute < time.Since(start) {
			fail("Took too long to attach volume")
		}

		time.Sleep(250 * time.Millisecond)

		var err error
		action, _, err = client.StorageActions.Get(volId, actionId)
		if err != nil {
			fail(err.Error())
		}
	}
}

func attach(opt AttachOptions) {
	client := getClient(opt.ApiKey)
	if client == nil {
		fail("Could not create client")
	}

	vol, _, err := client.Storage.GetVolume(opt.VolumeId)
	if err != nil {
		fail("Could not get volume \"", opt.VolumeId, "\": ", err.Error(), " options were: ", os.Args[2])
	}

	droplet, err := detectDroplet(client)
	if err != nil {
		fail(err.Error())
	}

	action, _, err := client.StorageActions.Attach(vol.ID, droplet.ID)
	if err != nil {
		fail(err.Error())
	}

	awaitAction(client, vol.ID, action)

	if err = ioutil.WriteFile(CRED_FILE, []byte(opt.ApiKey), 0600); err != nil {
		fail("Could not save credentials: ", err.Error())
	}

	finish(Result{
		Status: "Success",
		Device: DEVICE_PREFIX + vol.Name,
	}, 0)
}

func detach(device string) {
	cred_bytes, err := ioutil.ReadFile(CRED_FILE)
	if err != nil {
		fail("Failed to read credentials: ", err.Error())
	}

	client := getClient(string(cred_bytes))
	if client == nil {
		fail("Could not create client")
	}

	if !strings.HasPrefix(device, DEVICE_PREFIX) {
		fail("Expected path starting with ", DEVICE_PREFIX, "; instead got ", device)
	}

	disk_name := device[len(DEVICE_PREFIX):]

	// TODO paginate?
	vols, _, err := client.Storage.ListVolumes(nil)
	if err != nil {
		fail("Failed to list volumes: ", err.Error())
	}

	for _, vol := range vols {
		if vol.Name == disk_name {
			action, _, err := client.StorageActions.Detach(vol.ID)
			if err != nil {
				fail("Failed to detach: ", err.Error())
			}

			awaitAction(client, vol.ID, action)
			succeed()
		}
	}
	fail("Could not find volume attached at ", device)
}

func main() {
	if len(os.Args) < 2 {
		fail("Expected at least one argument")
	}

	switch os.Args[1] {
	case "init":
		succeed()

	case "attach":
		if len(os.Args) != 3 {
			fail("attach expected exactly 3 arguments; got ", os.Args)
		}

		var opt AttachOptions
		if err := json.Unmarshal([]byte(os.Args[2]), &opt); err != nil {
			fail("Could not parse options for attach; got ", os.Args[2])
		}

		attach(opt)

	case "detach":
		if len(os.Args) != 3 {
			fail("detach expected exactly 3 arguments; got ", os.Args)
		}

		device := os.Args[2]
		detach(device)

	default:
		fmt.Println("Not sure what to do. Called with args: ", os.Args)
		os.Exit(1)
	}

}
