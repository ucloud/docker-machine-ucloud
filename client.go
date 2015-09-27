package ucloud

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/ucloud/ucloud-sdk-go/service/uhost"
	"github.com/ucloud/ucloud-sdk-go/service/unet"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

var (
	hostsvc *uhost.UHost
	unetsvc *unet.UNet
)

func init() {
	config := &ucloud.Config{
		Credentials: &auth.KeyPair{
			PublicKey:  "ucloudsomeone@example.com1296235120854146120",
			PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
		},
		Region:    "cn-north-03",
		ProjectID: "",
	}

	hostsvc = uhost.New(config)
	unetsvc = unet.New(config)
}

func (d *Driver) createUHost() error {
	//TODO: password should remove the "=" after base64 encoding, I don't why?
	password := strings.Replace(base64.StdEncoding.EncodeToString([]byte(d.Password)), "=", "", -1)

	createUhostParams := uhost.CreateUHostInstanceParams{

		Region:    d.Region,
		ImageId:   d.ImageId,
		LoginMode: "Password",
		Password:  password,
		CPU:       defaultCPU,
		Memory:    defaultMemory,
		Quantity:  1,
		Count:     1,
	}

	resp, err := hostsvc.CreateUHostInstance(&createUhostParams)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("response is empty")
	}

	if len(resp.UHostIds) == 0 {
		return fmt.Errorf("UHostIds is empty")
	}
	d.UhostID = resp.UHostIds[0]

	return nil
}

func startUHost(region, hostId string) error {

	startUhostParams := uhost.StartUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.StartUHostInstance(&startUhostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func killUHost(region, hostId string) error {
	killUHostParams := uhost.PoweroffUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.PoweroffUHostInstance(&killUHostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func rebootUHost(region, hostId string) error {

	killUHostParams := uhost.PoweroffUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.PoweroffUHostInstance(&killUHostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func terminateUHost(region, hostId string) error {

	terminateUHostParams := uhost.TerminateUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.TerminateUHostInstance(&terminateUHostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func stopUHost(region, hostId string) error {
	stopUhostParams := uhost.StopUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.StopUHostInstance(&stopUhostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Stop UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

type UHostDetail struct {
	region string
	hostID string

	state           string
	publicIPAddress string
	privateIPAdress string
	cpu             int
	memory          int
}

func (d *Driver) getHostDescription() (*UHostDetail, error) {

	describeParams := uhost.DescribeUHostInstanceParams{
		Region:   d.Region,
		UHostIds: []string{d.UhostID},
		Offset:   0,
		Limit:    10,
	}

	resp, err := hostsvc.DescribeUHostInstance(&describeParams)
	if err != nil {
		return nil, err
	}

	if len(resp.UHostSet) == 0 {
		return nil, fmt.Errorf("UHostSet is empty")
	}

	if len(resp.UHostSet[0].IPSet) == 0 {
		return nil, fmt.Errorf("IPSet is empty")
	}

	var publicIpAddress string
	var privateIPAddress string
	for _, ip := range resp.UHostSet[0].IPSet {
		switch ip.Type {
		case "Private":
			privateIPAddress = ip.IP
		case "Bgp":
			publicIpAddress = ip.IP
		}
	}

	details := &UHostDetail{
		region:          d.Region,
		hostID:          resp.UHostSet[0].UHostId,
		state:           resp.UHostSet[0].State,
		publicIPAddress: publicIpAddress,
		privateIPAdress: privateIPAddress,
		cpu:             resp.UHostSet[0].CPU,
		memory:          resp.UHostSet[0].Memory,
	}

	return details, nil
}

// createUNet create network for uhost
func (d *Driver) createUNet() error {
	if err := d.configureIPAddress(); err != nil {
		return fmt.Errorf("configure IPAddress error:%s", err)
	}

	if err := d.configureSecurityGroup(); err != nil {
		return fmt.Errorf("configure security group error:%s", err)
	}

	return nil
}

// createKeyPair create keypair for ssh to docker-machine
func (d *Driver) createKeyPair() error {
	log.Debugf("SSH key path:%s", d.GetSSHKeyPath())

	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	return nil
}

// uploadKeyPair upload the public key to docker-machine
func (d *Driver) uploadKeyPair() error {

	ipAddr := d.IPAddress
	port, _ := d.GetSSHPort()

	tcpAddr := fmt.Sprintf("%s:%d", ipAddr, port)
	ssh.WaitForTCP(tcpAddr)

	auth := ssh.Auth{
		Passwords: []string{d.Password},
	}

	ssh.SetDefaultClient(ssh.Native)
	sshClient, err := ssh.NewClient(d.GetSSHUsername(), ipAddr, port, &auth)
	if err != nil {
		return err
	}

	publicKey, err := ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return err
	}

	command := fmt.Sprintf("mkdir -p ~/.ssh; echo '%s' > ~/.ssh/authorized_keys", string(publicKey))
	log.Debugf("Upload the public key with command: %s", command)

	output, err := sshClient.Output(command)
	log.Debugf("Upload command err, output: %v: %s", err, output)
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) configureIPAddress() error {

	// create an EIP and bind it to host
	if !d.PrivateIPOnly {
		createEIPParams := unet.AllocateEIPParams{
			Region:       d.Region,
			OperatorName: "Bgp",
			Bandwidth:    2,
			ChargeType:   "Dynamic",
			Quantity:     1,
		}

		resp, err := unetsvc.AllocateEIP(&createEIPParams)
		if err != nil {
			return fmt.Errorf("Allocate EIP failed:%s", err)
		}
		log.Debug(resp)

		if len(*resp.EIPSet) == 0 {
			return fmt.Errorf("EIP is empty")
		}
		eipId := (*resp.EIPSet)[0].EIPId
		if len(*(*resp.EIPSet)[0].EIPAddr) == 0 {
			return fmt.Errorf("IP Address is empty")
		}
		d.IPAddress = (*(*resp.EIPSet)[0].EIPAddr)[0].IP

		bindHostParams := unet.BindEIPParams{
			Region:       d.Region,
			EIPId:        eipId,
			ResourceType: "uhost",
			ResourceId:   d.UhostID,
		}

		bindEIPResp, err := unetsvc.BindEIP(&bindHostParams)
		if err != nil {
			return fmt.Errorf("Bind EIP failed:%s", err)
		}
		log.Debug(bindEIPResp)
	} else {
		hostDetails, err := d.getHostDescription()
		if err != nil {
			return fmt.Errorf("get host detail failed: %s", err)
		}
		d.IPAddress = hostDetails.publicIPAddress
		d.PrivateIPAddress = hostDetails.privateIPAdress
	}

	return nil
}

func (d *Driver) getSecurityGroup(name string) (int, error) {
	log.Debugf("get security group for group:%s", name)
	describeSecurityGroupsParams := unet.DescribeSecurityGroupParams{
		Region: d.Region,
	}
	describeSecurityGroupsResp, err := unetsvc.DescribeSecurityGroup(&describeSecurityGroupsParams)
	if err != nil {
		return 0, fmt.Errorf("get security groups failed:%s", err)
	}

	if len(describeSecurityGroupsResp.DataSet) == 0 {
		return 0, fmt.Errorf("security groups is empty")
	}

	for _, groups := range describeSecurityGroupsResp.DataSet {
		log.Debugf("name:%s, group id:%d", name, groups.GroupId)
		if groups.GroupName == name {
			log.Debugf("groups:%+v", groups)
			return groups.GroupId, nil
		}
	}

	return 0, fmt.Errorf("group:%s is not exist", name)
}

func (d *Driver) securityGroupAvailableFunc(name string) func() bool {
	return func() bool {
		_, err := d.getSecurityGroup(name)
		if err == nil {
			return true
		}
		return false
	}
}

func (d *Driver) configureSecurityGroup() error {
	var groupId int
	groupId, err := d.getSecurityGroup(d.SecurityGroupName)
	if err != nil {
		log.Debugf("get security group error:%s", err)
	}
	log.Debugf("groupId:%s", groupId)
	if groupId == 0 {
		log.Infof("security group is not found, create a new one")

		securityGroupParams := unet.CreateSecurityGroupParams{
			Region:      d.Region,
			GroupName:   "docker-machine",
			Description: "docker machine to open 2379 and 22 port of tcp",
			Rule: []string{"TCP|22|0.0.0.0/0|ACCEPT|50",
				"TCP|3389|0.0.0.0/0|ACCEPT|50",
				"TCP|2376|0.0.0.0/0|ACCEPT|50"},
		}
		_, err := unetsvc.CreateSecurityGroup(&securityGroupParams)
		if err != nil {
			return fmt.Errorf("create security group failed:%s", err)
		}

		log.Debug("waiting for security group to become avaliable")
		if err := mcnutils.WaitFor(d.securityGroupAvailableFunc(d.SecurityGroupName)); err != nil {
			return err
		}
		groupId, err = d.getSecurityGroup(d.SecurityGroupName)
	}
	d.SecurityGroupId = groupId

	grantSecurityGroupParams := unet.GrantSecurityGroupParams{
		Region:       d.Region,
		GroupId:      groupId,
		ResourceType: "uhost",
		ResourceId:   d.UhostID,
	}
	log.Debugf("grant security group(%d) to uhost(%s)", groupId, d.UhostID)
	_, err = unetsvc.GrantSecurityGroup(&grantSecurityGroupParams)
	if err != nil {
		return fmt.Errorf("grant security group failed:%s", err)
	}

	return nil
}
