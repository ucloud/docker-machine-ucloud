package ucloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
)

type Driver struct {
	*drivers.BaseDriver

	PublicKey  string
	PrivateKey string
	Region     string
	ImageId    string
	Password   string
	UhostID    string

	CPU       int
	Memory    int
	DiskSpace int

	PrivateIPOnly 	bool

	//	SecurityGroupId string
	//	Address         string
	//	UHostType       string
}

const (
	defaultTimeout   = 1 * time.Second
	defaultCPU       = 1
	defaultMemory    = 1024
	defaultDiskSpace = 20000
	defaultRegion    = "cn-north-03"
	defaultRetries   = 10
)

func NewDriver(hostName, artifactPath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName:  hostName,
			ArtifactPath: artifactPath,
		},
		Region:    defaultRegion,
		Memory:    defaultMemory,
		CPU:       defaultCPU,
		DiskSpace: defaultDiskSpace,
	}
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		{
			Name:   "ucloud-public-key",
			Usage:  "UCloud Public Key",
			Value:  "",
			EnvVar: "UCLOUD_PUBLIC_KEY",
		},
		{
			Name:   "ucloud-private-key",
			Usage:  "UCloud Private Key",
			Value:  "",
			EnvVar: "UCLOUD_PRIVATE_KEY",
		},
		{
			Name:  "ucloud-imageid",
			Usage: "UHost image id",
			Value: "",
		},
		{
			Name:  "ucloud-user-password",
			Usage: "Password of ucloud user",
			Value: "",
		},
		{
			Name:   "ucloud-region",
			Usage:  "Region of ucloud idc",
			Value:  "cn-north-03",
			EnvVar: "UCLOUD_REGION",
		},
		{
			Name:  "ucloud-ssh-user",
			Usage: "SSH user",
			Value: "root",
		},
		{
			Name:  "ucloud-ssh-port",
			Usage: "SSH port",
			Value: 22,
		},
		{
			Name: "ucloud-private-address-only",
			Usage: "Only use a private IP address",
		},
	}
}

func (d *Driver) DriverName() string {
	return "ucloud"
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetSSHUsername() string {
	if d.SSHUser == "" {
		d.SSHUser = "root"
	}

	return d.SSHUser
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	region, err := validateUCloudRegion(flags.String("ucloud-region"))
	if err != nil {
		return err
	}
	d.Region = region

	d.PublicKey = flags.String("ucloud-public-key")
	if d.PublicKey == "" {
		return fmt.Errorf("ucloud driver requires the --ucloud-public-key option")
	}
	log.Debugf("ucloud public key: %s", d.PublicKey)

	d.PrivateKey = flags.String("ucloud-private-key")
	if d.PrivateKey == "" {
		return fmt.Errorf("ucloud driver requires the --ucloud-private-key option")
	}
	log.Debugf("ucloud private key: %s", d.PrivateKey)

	image := flags.String("ucloud-imageid")
	if len(image) == 0 {
		image = "uimage-j4fbrn"
	}
	d.ImageId = image

	d.PrivateIPOnly = flags.Bool("ucloud-private-address-only")

	d.SSHUser = strings.ToLower(flags.String("ucloud-ssh-user"))
	if d.SSHUser == "" {
		d.SSHUser = "root"
	}
	d.Password = flags.String("ucloud-user-password")
	d.SSHPort = 22


	return nil
}

func (d *Driver) PreCreateCheck() error {
	return nil
}

func (d *Driver) Create() error {
	log.Infof("Create UHost instance...")

	if d.Password == "" {
		return fmt.Errorf("ucloud driver requires --ucloud-user-password options.")
	}
	//TODO: create an keypairs when uhost api is supported, we have to use expect to
	// create the ssh keypairs

	// create uhost instance
	hostId, err := createUHost(d.Region, d.ImageId, d.Password)
	if err != nil {
		return fmt.Errorf("create UHost failed:%s", err)
	}
	d.UhostID = hostId

	if !d.PrivateIPOnly {
		//TODO: user the exist eip and security group to configure network
		d.createUNet()
	}

	for i := 0; i < defaultRetries; i++ {
		// get details of host
		details, err := getHostDescription(d.Region, hostId)
		if err != nil {
			return fmt.Errorf("get UHost details failed:%s", err)
		}

		if details != nil && details.ipAddress != "" && details.cpu != 0 {
			d.IPAddress = details.ipAddress
			d.CPU = details.cpu
			d.Memory = details.memory

			break
		}

		time.Sleep(1 * time.Second)
	}

	// TODO: install docker

	return nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}

	return d.IPAddress, nil
}

func (d *Driver) GetState() (state.State, error) {
	log.Debugf("Driver:%+v", d)
	log.Info("Get Uhost details")
	if d.UhostID == "" || d.Region == "" {
		return state.None, fmt.Errorf("region or uhost is empty")
	}

	details, err := getHostDescription(d.Region, d.UhostID)
	if err != nil {
		return state.None, fmt.Errorf("get UHost details failed:%s", err)
	}

	var st state.State
	if details != nil && details.state != "" {
		switch details.state {
		case "Initializing", "Starting", "Rebooting":
			st = state.Starting
		case "Running":
			st = state.Running
		case "Stopped":
			st = state.Stopped
		case "Stopping":
			st = state.Stopping
		case "Install Fail":
			st = state.Error
		default:
			st = state.None
		}
	}

	return st, nil
}

func (d *Driver) Start() error {
	log.Info("Start UHost...")
	err := startUHost(d.Region, d.UhostID)
	if err != nil {
		return fmt.Errorf("Cannot start Machine:%s, with UHost: %s.", d.MachineName, d.UhostID)
	}

	return nil
}

func (d *Driver) Stop() error {
	log.Info("Stop UHost...")
	if len(d.UhostID) == 0 {
		return fmt.Errorf("UHost is not exist for Machine:%s", d.MachineName)
	}

	err := stopUHost(d.Region, d.UhostID)
	if err != nil {
		return fmt.Errorf("Cannot start Machine:%s, with UHost: %s.", d.MachineName, d.UhostID)
	}

	return nil
}

func (d *Driver) Remove() error {
	log.Debug("Removing...")
	if err := terminateUHost(d.Region, d.UhostID); err != nil {
		return fmt.Errorf("Unable to terminate the UHost instance:%s", err)
	}

	//TODO: any cleanup ?
	return nil
}

func (d *Driver) Restart() error {
	log.Debug("Restarting...")
	if err := rebootUHost(d.Region, d.UhostID); err != nil {
		return fmt.Errorf("Unable to restart the UHost instance:%s", err)
	}

	return nil
}

func (d *Driver) Kill() error {
	log.Debug("Killing...")
	if err := killUHost(d.Region, d.UhostID); err != nil {
		return fmt.Errorf("Unable to kill the UHost instance:%s", err)
	}

	return nil
}
