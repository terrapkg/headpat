package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var l = log.New(os.Stdout, "[anitnya] ", 0)

type FedoraAuth struct {
}

func (auth *FedoraAuth) Mechanism() string {
	return "EXTERNAL"
}

func (auth *FedoraAuth) Response() string {
	return ""
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func readf(path string) chan []byte {
	ch := make(chan []byte)
	go func() {
		f, err := os.Open(path)
		if err != nil {
			l.Fatal(err)
		}
		defer f.Close()
		s, err := io.ReadAll(f)
		if err != nil {
			l.Fatal(err)
		}
		ch <- s
	}()
	return ch
}

func init_tls() (tls.Certificate, *x509.CertPool) {
	x, err := exists("/etc/fedora-messaging/")
	if err != nil {
		l.Println("Fail to check /etc/fedora-messaging/: " + err.Error())
	}
	if x {
		l.Println("using keys and certs in /etc/fedora-messaging/")
		f1 := readf("/etc/fedora-messaging/fedora-cert.pem")
		f2 := readf("/etc/fedora-messaging/fedora-key.pem")
		f3 := readf("/etc/fedora-messaging/cacert.pem")
		cert, err := tls.X509KeyPair(<-f1, <-f2)
		if err != nil {
			l.Fatal(err)
		}
	
		l.Println("cfg cacert")
		certpool := x509.NewCertPool()
		certpool.AppendCertsFromPEM(<-f3)
		return cert, certpool
	}
	l.Println("dl cert, prikey, cacert")
	dl1 := dl("https://raw.githubusercontent.com/fedora-infra/fedora-messaging/stable/configs/fedora-cert.pem")
	dl2 := dl("https://raw.githubusercontent.com/fedora-infra/fedora-messaging/stable/configs/fedora-key.pem")
	dl3 := dl("https://raw.githubusercontent.com/fedora-infra/fedora-messaging/stable/configs/cacert.pem")
	cert, err := tls.X509KeyPair([]byte(<-dl1), []byte(<-dl2))
	if err != nil {
		l.Fatal(err)
	}

	l.Println("cfg cacert")
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM([]byte(<-dl3))
	return cert, certpool
}

func anitnya_conn() {
	cert, certpool := init_tls()

	l.Println("conn")
	conn, err := amqp.DialConfig("amqps://rabbitmq.fedoraproject.org/%2Fpublic_pubsub", amqp.Config{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      certpool,
		},
		SASL: []amqp.Authentication{
			&FedoraAuth{},
		},
	})
	if err != nil {
		l.Fatal(err)
	}
	l.Println("ch grab")
	ch, err := conn.Channel()
	if err != nil {
		l.Fatal(err)
	}
	defer conn.Close()
	// "org.release-monitoring.prod.anitya.project.version.update.v2", "org.fedoraproject.prod.copr"
	for _, topic := range []string{"org.fedoraproject.prod.ci.productmd-compose.test.complete"} {
		l.Println("xdcl " + topic)
		err = ch.ExchangeDeclare(topic, "topic", true, false, false, false, nil)
		if err != nil {
			l.Fatal(err)
		}
	}

	qname := os.Getenv("QUEUE_UUID")
	l.Println("qdcl " + qname)
	_, err = ch.QueueDeclare(qname, true, false, true, false, nil)
	if err != nil {
		l.Fatal(err)
	}

	l.Println("ch consume")
	msgs, err := ch.Consume(qname, "", true, false, false, false, nil)
	if err != nil {
		l.Fatal(err)
	}

	forever := make(chan bool, 1)
	go func() {
		for d := range msgs {
			l.Printf("| %s\n", d.Body)
		}
	}()

	l.Println("conn ok nya~")
	<-forever
}
