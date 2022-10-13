/**
 * Copyright 2018, Xiaomi.
 * All rights reserved.
 * Author: wangfan8@xiaomi.com
 */

package main

import (
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/MapleLove2014/talos-sdk-golang/client"
	"github.com/MapleLove2014/talos-sdk-golang/producer"
	"github.com/MapleLove2014/talos-sdk-golang/thrift/message"
	"github.com/MapleLove2014/talos-sdk-golang/utils"
	"github.com/sirupsen/logrus"
)

type MyMessageCallback struct {
	log *logrus.Logger
}

func NewMyMessageCallback(logger *logrus.Logger) *MyMessageCallback {
	return &MyMessageCallback{log: logger}
}

func (c *MyMessageCallback) OnSuccess(userMessageResult *producer.UserMessageResult) {
	atomic.StoreInt64(successPutNumber, atomic.LoadInt64(successPutNumber)+int64(len(userMessageResult.GetMessageList())))
	count := atomic.LoadInt64(successPutNumber)

	for _, msg := range userMessageResult.GetMessageList() {
		c.log.Infof("success to put message: %s", string(msg.GetMessage()))
	}
	c.log.Infof("success to put message: %d for partition: %d so far.", count, userMessageResult.GetPartitionId())
}

func (c *MyMessageCallback) OnError(userMessageResult *producer.UserMessageResult) {
	for _, msg := range userMessageResult.GetMessageList() {
		c.log.Infof("failed to put message: %d , will retry to put it.", msg)
	}
	err := talosProducer.AddUserMessage(userMessageResult.GetMessageList())
	if err != nil {
		c.log.Errorf("put message retry failed: %s", err.Error())
	}
}

var successPutNumber *int64
var talosProducer *producer.TalosProducer

func main() {
	successPutNumber = new(int64)
	atomic.StoreInt64(successPutNumber, 0)
	// init client config by put $your_propertyFile in your classpath
	// with the content of:
	/*
	   galaxy.talos.service.endpoint=$talosServiceURI
	*/
	var propertyFilename string
	flag.StringVar(&propertyFilename, "conf", "talosProducer.conf", "conf: talosProducer.conf'")
	flag.Parse()

	var err error
	log := utils.InitLogger()
	talosProducer, err = producer.NewTalosProducerWithLogger(propertyFilename,
		client.NewSimpleTopicAbnormalCallback(), NewMyMessageCallback(log), log)
	if err != nil {
		log.Errorf("Init talosProducer failed: %s", err.Error())
		return
	}
	toPutMsgNumber := 20
	messageList := make([]*message.Message, 0)
	for i := 0; i < toPutMsgNumber; i++ {
		messageStr := fmt.Sprintf("This message is a text string. messageId: %d", i)
		msg := &message.Message{Message: []byte(messageStr)}
		messageList = append(messageList, msg)
	}

	for {
		if err := talosProducer.AddUserMessage(messageList); err != nil {
			log.Error(err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
