package gamq

import (
	"bufio"
	"bytes"
	"github.com/FireEater64/gamq/message"
	"testing"

	"github.com/onsi/gomega"
)

func TestMessageShipper_SuccessfullyForwardsMessages(t *testing.T) {
	gomega.RegisterTestingT(t)

	inputChannel := make(chan *message.Message, 0)

	writerBuffer := new(bytes.Buffer)
	dummyWriter := bufio.NewWriter(writerBuffer)
	closedChannel := make(chan bool)
	dummyClient := Client{Name: "Test", Writer: dummyWriter, Closed: &closedChannel}
	dummyMetricsChannel := make(chan *Metric)

	underTest := newMessageShipper(inputChannel, &dummyClient, dummyMetricsChannel, "test")

	testMessagePayload := []byte("This is a test!")
	expectedMessagePayload := []byte("This is a test!\r\n.\r\n")
	testMessage := message.NewHeaderlessMessage(&testMessagePayload)
	underTest.messageChannel <- testMessage

	gomega.Eventually(func() []byte {
		return writerBuffer.Bytes()
	}).Should(gomega.Equal(expectedMessagePayload))
}
