package telnet

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func New(user string, pass string, ip string) *Telnet {
	return &Telnet{
		user: user,
		pass: pass,
		ip:   ip,
	}
}

func (t *Telnet) PermTelnet() error {
	conn, err := net.Dial("tcp", t.ip+":telnet")
	if err != nil {
		return err
	}
	defer conn.Close()
	t.conn = conn

	if err := t.loginTelnet(); err != nil {
		return err
	}

	if err := t.modifyDB(); err != nil {
		return err
	}

	// reboot device
	fmt.Println("wait reboot..")
	time.Sleep(time.Second)
	if err := t.Reboot(); err != nil {
		return err
	}

	return nil
}

func (t *Telnet) loginTelnet() error {
	return t.sendCmd(t.user, t.pass)
}

func (t *Telnet) modifyDB() error {
	// set DB data
	prefix := "sendcmd 1 DB set TelnetCfg 0 "
	lanEnable := prefix + "Lan_Enable 1"
	tsLanUser := prefix + "TSLan_UName root"
	tsLanPwd := prefix + "TSLan_UPwd Zte521"
	maxConn := prefix + "Max_Con_Num 3"
	initSecLvl := prefix + "InitSecLvl 3"

	// save DB
	save := "sendcmd 1 DB save"

	if err := t.sendCmd(lanEnable, tsLanUser, tsLanPwd, maxConn, initSecLvl, save); err != nil {
		return err
	}
	fmt.Println("Permanent Telnet succeed\r\nuser: root, pass: Zte521")

	return nil
}

func (t *Telnet) sendCmd(commands ...string) error {
	cmd := []byte(strings.Join(commands, ctrl) + ctrl)
	n, err := t.conn.Write(cmd)
	if err != nil {
		return err
	}

	if expected, actual := len(cmd), n; expected != actual {
		err := fmt.Errorf("transmission problem: tried sending %d bytes, but actually only sent %d bytes", expected, actual)
		return err
	}

	return nil
}

func (t *Telnet) Reboot() error {
	return t.sendCmd("reboot")
}