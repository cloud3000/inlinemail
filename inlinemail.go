package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/blackjack/syslog"
	mail "github.com/xhit/go-simple-mail"
)

// Smtpconf defines the smtpconf.json config file.
type Smtpconf struct {
	Smtpuser string `json:"smtpUser"`
	Smtppass string `json:"smtpPass"`
	Smtpserv string `json:"smtpServ"`
	Smtpport int    `json:"smtpPort"`
}

func getconf() (Smtpconf, error) {
	var obj Smtpconf
	data, err := ioutil.ReadFile("./smtpconf.json")
	if err != nil {
		return obj, err
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, err
}

func iEmail(mailfrom *string, mailto *string, mailsub *string, mailmsg *string) int {
	conf, conferr := getconf()
	if conferr != nil {
		fmt.Printf("getconf error: %s\n", conferr.Error())
		os.Exit(2)
	}

	// Set up authentication information.
	server := mail.NewSMTPClient()

	//SMTP Server
	server.Host = conf.Smtpserv
	server.Port = conf.Smtpport
	server.Username = conf.Smtpuser
	server.Password = conf.Smtppass
	server.Encryption = mail.EncryptionTLS

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	//Variable to keep alive connection
	server.KeepAlive = false

	//Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	//Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	//SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		//log.Fatal(err)
	}

	//New email simple html with inline and CC
	email := mail.NewMSG()

	email.SetFrom(*mailfrom).
		AddTo(*mailto).
		SetSubject(*mailsub)

	email.SetBody(mail.TextHTML, *mailmsg)

	email.AddInline("/home/michael/Pictures/j3k_logo.jpg", "j3k_logo.jpg")

	//Call Send and pass the client
	err = email.Send(smtpClient)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Email Sent")
	return 0
}
func main() {
	var (
		inlineFrom    = flag.String("f", "", "FROM Email address")
		inlineTo      = flag.String("t", "", "TO Email address")
		inlineSubject = flag.String("s", "", "SUBJECT of Email")
	)
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [flags] command [command args…]\n", os.Args[0])
		flag.PrintDefaults()
	}

	if len(*inlineSubject) == 0 || len(*inlineFrom) == 0 || len(*inlineTo) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	syslog.Openlog("inlinemail", syslog.LOG_PID, syslog.LOG_USER)
	syslog.Syslog(syslog.LOG_INFO, "inlinemail begin")

	htmlBody, err := ioutil.ReadFile("./inlinemail.html")
	if err != nil {
		fmt.Printf("Unable to open HTML templete ./inlinemail.html %s\n", err.Error())
		os.Exit(1)
	}

	inlineMsg := string(htmlBody)
	iEmail(inlineFrom, inlineTo, inlineSubject, &inlineMsg)
	syslog.Syslog(syslog.LOG_INFO, "inlinemail end")
	syslog.Closelog()

}
