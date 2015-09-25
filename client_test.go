package ucloud


import (
	"testing"
)

func TestCreateUNet(t *testing.T) {
	d := Driver{
		Region: "cn-north-03",
		UhostID: "uhost-ygtk1p",
	}

	err := d.createUNet()
	if  err != nil {
		t.Errorf("create UNet failed:%s", err)
	}
}
