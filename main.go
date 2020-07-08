//-----------------------------------------------
// GoIPMailer
// 2020 by frickelblog.de
//-----------------------------------------------

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
	"gopkg.in/gomail.v2"
)

//-------------------------------------------------------------------------------------------------
// Service Configuration
var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run(s)
	return nil
}
func (p *program) run(s service.Service) {
	GoIPMailer_main()
	fmt.Println("Exit")
	s.Stop()
	os.Exit(1)
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "GoIPMailer",
		DisplayName: "GoIPMailer",
		Description: "GoIPMailer Service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}

//-------------------------------------------------------------------------------------------------

func GoIPMailer_main() {

	StartUpPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(StartUpPath)
	config := readConfig(path.Join(StartUpPath, "./config.json"))

	fmt.Println("Version: 1.0")
	fmt.Println("SMTPHost: " + config.SMTPHost)
	fmt.Println("SMTPFROM: " + config.SMTPFROM)
	fmt.Println("SMTPTO: " + config.SMTPTO)
	fmt.Println("IPFilter: " + config.IPFilter)
	fmt.Println("RDPUser: " + config.RDPUser)

	HostName, _ := os.Hostname()
	HostIP := resolveHostIP(config.IPFilter)
	fmt.Println(HostName)
	fmt.Println(HostIP)

	// RDP-Template auslesen und Platzhalter ersetzen
	bRDPTemplate, _ := Asset("assets/template.rdp")
	RDPTemplateString := bytes.NewBuffer(bRDPTemplate).String()
	RDPTemplateString = strings.Replace(RDPTemplateString, "<<HOSTIP>>", HostIP, -1)           // IP
	RDPTemplateString = strings.Replace(RDPTemplateString, "<<USERNAME>>", config.RDPUser, -1) // RDP-user
	// RDP-datei für spätere E-Mail Verarbeitung zwischenspeichern
	bRDPTemplate = []byte(RDPTemplateString)
	err := ioutil.WriteFile(path.Join(StartUpPath, HostIP+".rdp"), bRDPTemplate, 0644)
	check(err)

	// Text für die RDP-Zeile
	sMSTSCText := ""
	if config.RDPAttachBodyLine {
		sMSTSCText = "RDP: <code>mstsc.exe /v:" + HostIP + " /prompt</code></br>"
	}
	// Text für die SSH-Zeile
	sSSHText := ""
	if config.SSHAttachBodyLine {
		sSSHText = "SSH: <code>ssh " + config.SSHUser + "@" + HostIP + "</code></br>"
	}
	// Benutzerdefinierter Text
	sUserText := ""
	if config.UserAttachBodyLine {
		sUserText = config.UserBodyLine
	}

	// Neue GoMail Nachricht
	m := gomail.NewMessage()
	m.SetHeader("From", config.SMTPFROM)
	m.SetHeader("To", config.SMTPTO)
	m.SetHeader("Subject", "GoIPMailer - "+HostName+" - "+HostIP+"")
	m.SetBody("text/html", "<b>GoIPMailer</b> - "+HostName+" - "+HostIP+"</br></br>"+sMSTSCText+sSSHText+"</br></br>"+sUserText)
	if config.RDPAttachFile {
		m.Attach(HostIP + ".rdp")
	}
	if config.UserAttachFile {
		m.Attach(config.UserFile)
	}

	// Invalide Zertifikate ignorieren
	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// E-Mail senden
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func resolveHostIP(IPFilter string) string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if check(err) {
		for _, netInterfaceAddress := range netInterfaceAddresses {
			networkIP, ok := netInterfaceAddress.(*net.IPNet)
			if ok {
				ip := networkIP.IP.String()
				if strings.HasPrefix(ip, IPFilter) {
					return ip
				}
			}
		}
	}

	return ""
}

//--------------------------------------------------------------------------
// Typen
//--------------------------------------------------------------------------
type Configuration struct {
	BinPath            string
	SMTPHost           string
	SMTPPort           int
	SMTPFROM           string
	SMTPTO             string
	SMTPUser           string
	SMTPPass           string
	IPFilter           string
	RDPUser            string
	RDPAttachFile      bool
	RDPAttachBodyLine  bool
	SSHUser            string
	SSHAttachBodyLine  bool
	UserFile           string
	UserAttachFile     bool
	UserBodyLine       string
	UserAttachBodyLine bool
}

//--------------------------------------------------------------------------
// Hilfsfunktionen
//--------------------------------------------------------------------------

func check(e error) bool {
	if e != nil {
		fmt.Println(e)
		return false
	}
	return true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readConfig(filename string) *Configuration {
	// initialize conf with default values.
	conf := &Configuration{}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return conf
	}
	if err = json.Unmarshal(b, conf); err != nil {
		return conf
	}
	return conf
}
