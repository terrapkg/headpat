package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type FedoraAuth struct {
}

func (auth *FedoraAuth) Mechanism() string {
	return "EXTERNAL"
}

func (auth *FedoraAuth) Response() string {
	return "guest\nguest"
}

func anitnya_conn() {
	l := log.New(os.Stdout, "[anitnya] ", 0)

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
	// "org.release-monitoring.prod.anitya.project.version.update.v2",
	for _, topic := range []string{"org.fedoraproject.prod.copr"} {
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
