package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/streadway/amqp"
)

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
	cacert_ := []byte(<-dl3)

	l.Println("cfg cacert")
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(cacert_)

	l.Println("conn")
	conn, err := amqp.DialConfig("amqps://rabbitmq.fedoraproject.org/%2Fpublic_pubsub", amqp.Config{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
			RootCAs: certpool,
			// InsecureSkipVerify:       true,
			// PreferServerCipherSuites: true,
			// MinVersion:               tls.VersionTLS11,
			// MaxVersion:               tls.VersionTLS11,
			// CipherSuites: []uint16{
			// 	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			// 	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// },

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

	l.Println("ch consume")
	msgs, err := ch.Consume(
		"", //os.Getenv("QUEUE_UUID"),
		"",
		true, // autoAck
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		l.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			l.Printf("| %s\n", d.Body)
		}
	}()

	log.Println("conn ok")
	<-forever
}
