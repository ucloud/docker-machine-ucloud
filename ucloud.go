package ucloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
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

	CPU        int
	Memory     int
	DiskSpace  int
	ChargeType string

	PrivateIPOnly     bool
	PrivateIPAddress  string
	SecurityGroupId   int
	SecurityGroupName string
}

const (
	defaultTimeout    = 1 * time.Second
	defaultCPU        = 1
	defaultMemory     = 1024
	defaultDiskSpace  = 20
	defaultRegion     = "cn-north-03"
	defaultChargeType = "Month"
	defaultRetries    = 10
	defaultImageId    = "uimage-5yt2b0" // we use CentOS 7.0 default
)

func NewDriver(hostName, artifactPath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			//			ArtifactPath: artifactPath,
		},
		Region:    defaultRegion,
		Memory:    defaultMemory,
		CPU:       defaultCPU,
		DiskSpace: defaultDiskSpace,
	}
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{

		mcnflag.StringFlag{
			Name:   "ucloud-public-key",
			Usage:  "UCloud Public Key",
			Value:  "",
			EnvVar: "UCLOUD_PUBLIC_KEY",
		},
		mcnflag.StringFlag{
			Name:   "ucloud-private-key",
			Usage:  "UCloud Private Key",
			Value:  "",
			EnvVar: "UCLOUD_PRIVATE_KEY",
		},
		mcnflag.StringFlag{
			Name:  "ucloud-imageid",
			Usage: "UHost image id",
			Value: "",
		},
		mcnflag.StringFlag{
			Name:   "ucloud-region",
			Usage:  "Region of ucloud idc",
			Value:  "cn-north-03",
			EnvVar: "UCLOUD_REGION",
		},
		mcnflag.StringFlag{
			Name:  "ucloud-ssh-user",
			Usage: "SSH user",
			Value: "root",
		},
		mcnflag.IntFlag{
			Name:  "ucloud-cpu-core",
			Usage: "Number of CPU cores,default is 1",
			Value: defaultCPU,
		},
		mcnflag.IntFlag{
			Name:  "ucloud-memory-size",
			Usage: "Size of memory, unit(MB), default 1024M",
			Value: defaultMemory,
		},
		mcnflag.IntFlag{
			Name:  "ucloud-disk-space",
			Usage: "Disk size, unit(GB),default is 20G",
			Value: defaultDiskSpace,
		},
		mcnflag.StringFlag{
			Name:  "ucloud-charge-type",
			Usage: "How to pay for, you can chose from (Year,Month,Dynamic,Trial),default is Month",
			Value: defaultChargeType,
		},
		mcnflag.IntFlag{
			Name:  "ucloud-ssh-port",
			Usage: "SSH port",
			Value: 22,
		},
		mcnflag.StringFlag{
			Name:  "ucloud-user-password",
			Usage: "Password of ucloud user, random password will be used if not set",
			Value: "",
		},
		mcnflag.BoolFlag{
			Name:  "ucloud-private-address-only",
			Usage: "Only use a private IP address",
		},
		mcnflag.StringFlag{
			Name:  "ucloud-security-group",
			Usage: "UCloud security group",
			Value: "docker-machine",
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

func (d *Driver) setDefaultConfig() {
	d.Memory = defaultMemory
	d.CPU = defaultCPU
	d.ChargeType = defaultChargeType
	d.DiskSpace = defaultDiskSpace
	d.Region = defaultRegion
	d.ImageId = defaultImageId
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.setDefaultConfig()
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
		image = defaultImageId
	}
	d.ImageId = image
	d.CPU = flags.Int("ucloud-cpu-core")
	d.Memory = flags.Int("ucloud-memory-size")
	d.DiskSpace = flags.Int("ucloud-disk-space")
	d.ChargeType = flags.String("ucloud-charge-type")

	d.PrivateIPOnly = flags.Bool("ucloud-private-address-only")
	d.SecurityGroupName = flags.String("ucloud-security-group")

	d.SSHUser = strings.ToLower(flags.String("ucloud-ssh-user"))
	if d.SSHUser == "" {
		d.SSHUser = "root"
	}
	d.Password = flags.String("ucloud-user-password")
	d.SSHPort = flags.Int("ucloud-ssh-port")

	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")

	return nil
}

func (d *Driver) PreCreateCheck() error {
	if d.CPU < 1 || d.CPU > 16 {
		return fmt.Errorf("CPU cores must be in set of (1,2,4,8,16)")
	}
	if d.Memory < 1024 || d.Memory > 65536 {
		return fmt.Errorf("Memory must be in range of [2048, 65536) with step of 2048MB, you can set 1024 in beijing-BGP-C")
	}
	if d.DiskSpace > 1000 {
		return fmt.Errorf("Disk space must in range of [0, 1000) with step of 10GB")
	}
	return nil
}

func (d *Driver) Create() error {
	log.Infof("Create UHost instance...")

	if d.Password == "" {
		d.Password = generateRandomPassword(16)
		log.Infof("password is not set, we use the random password instead, password:%s", d.Password)
	}

	// create keypair
	log.Infof("Creating key pair for instances...")
	if err := d.createKeyPair(); err != nil {
		return fmt.Errorf("unable to create key pair: %s", err)
	}

	// create uhost instance
	log.Infof("Creating uhost instance...")
	if err := d.createUHost(); err != nil {
		return fmt.Errorf("create UHost failed:%s", err)
	}

	// waiting for creating successful
	if err := mcnutils.WaitForSpecific(drivers.MachineInState(d, state.Running), 120, 3*time.Second); err != nil {
		return fmt.Errorf("wait for machine running failed: %s", err)
	}

	// create networks, like private ip, eip, and security group
	log.Infof("Creating networks...")
	//TODO: user the exist eip and security group to configure network
	if err := d.createUNet(); err != nil {
		return fmt.Errorf("create networks failed:%s", err)
	}

	// upload keypair
	if err := d.uploadKeyPair(); err != nil {
		return fmt.Errorf("upload keypair failed:%s", err)
	}

	return nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}

	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetIP() (string, error) {
	if !d.PrivateIPOnly && d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	if d.PrivateIPOnly && d.PrivateIPAddress == "" {
		return "", fmt.Errorf("Private address is not set")
	}

	s, err := d.GetState()
	if err != nil {
		return "", err
	}
	if s != state.Running {
		return "", drivers.ErrHostIsNotRunning
	}

	if d.PrivateIPOnly {
		return d.PrivateIPAddress, nil
	}

	return d.IPAddress, nil
}

func (d *Driver) GetState() (state.State, error) {
	log.Debugf("Get Machine State")
	if d.UhostID == "" || d.Region == "" {
		return state.None, fmt.Errorf("region or uhost is empty")
	}

	details, err := d.getHostDescription()
	if err != nil {
		return state.None, err
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
	if err := d.startUHost(); err != nil {
		return fmt.Errorf("Cannot start Machine:%s, with UHost: %s.", d.MachineName, d.UhostID)
	}

	return nil
}

func (d *Driver) Stop() error {
	log.Info("Stop UHost...")
	if len(d.UhostID) == 0 {
		return fmt.Errorf("UHost is not exist for Machine: %s", d.MachineName)
	}

	if err := d.stopUHost(); err != nil {
		return fmt.Errorf("Cannot start Machine:%s, with UHost: %s.", d.MachineName, d.UhostID)
	}

	return nil
}

func (d *Driver) Remove() error {
	log.Debug("Removing...")
	if err := d.terminateUHost(); err != nil {
		return fmt.Errorf("Unable to terminate the UHost instance: %s", err)
	}

	//TODO: any cleanup ?
	return nil
}

func (d *Driver) Restart() error {
	log.Debug("Restarting...")
	if err := d.rebootUHost(); err != nil {
		return fmt.Errorf("Unable to restart the UHost instance: %s", err)
	}

	return nil
}

func (d *Driver) Kill() error {
	log.Debug("Killing...")
	if err := d.killUHost(); err != nil {
		return fmt.Errorf("Unable to kill the UHost instance: %s", err)
	}

	return nil
}
