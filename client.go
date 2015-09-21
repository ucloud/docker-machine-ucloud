package ucloud

import (
	"encoding/base64"
	"fmt"

	"github.com/ucloud/ucloud-sdk-go/service/uhost"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

var (
	hostsvc *uhost.UHost
)

func init() {
	hostsvc = uhost.New(&ucloud.Config{
		Credentials: &auth.KeyPair{
			PublicKey:  "ucloudsomeone@example.com1296235120854146120",
			PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
		},
		Region:    "cn-north-03",
		ProjectID: "",
	})
}

func createUHost(region, imageId, password string) (string, error) {
	password = base64.StdEncoding.EncodeToString([]byte(password))

	createUhostParams := uhost.CreateUHostInstanceParams{

		Region:    region,
		ImageId:   imageId,
		LoginMode: "Password",
		Password:  password,
		CPU:       1,
		Memory:    2048,
		Quantity:  1,
		Count:     1,
	}

	resp, err := hostsvc.CreateUHostInstance(&createUhostParams)
	if err != nil {
		return "", err
	}

	if resp.RetCode != 0 {
		return "", fmt.Errorf("Create UHost error, RetCode:%d, err:%s", resp.RetCode, err)
	}

	return resp.HostIds[0], nil
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
	//
	//	killUHostParams := uhost.{
	//		Region:  region,
	//		UHostId: hostId,
	//	}
	//
	//	resp, err := hostsvc.PoweroffUHostInstance(&killUHostParams)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if resp.RetCode != 0 {
	//		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	//	}
	//
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

	state     string
	ipAddress string
	cpu       int
	memory    int
}

func getHostDescription(region, hostId string) (*UHostDetail, error) {

	describeParams := uhost.DescribeUHostInstanceParams{
		Region: region,
		Offset: 0,
		Limit:  10,
	}

	resp, err := hostsvc.DescribeUHostInstance(&describeParams)
	if err != nil {
		return nil, err
	}

	if resp.RetCode != 0 {
		return nil, fmt.Errorf("Describe UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	hostState := resp.UHostSet[0].State
	ip := resp.UHostSet[0].IPSet[0].IP

	details := &UHostDetail{
		region:    region,
		hostID:    hostId,
		state:     hostState,
		ipAddress: ip,
		cpu:       resp.UHostSet[0].CPU,
		memory:    resp.UHostSet[0].Memory,
	}

	return details, nil
}
