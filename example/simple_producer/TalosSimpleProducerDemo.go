/**
 * Copyright 2018, Xiaomi.
 * All rights reserved.
 * Author: wangfan8@xiaomi.com
 */

package main

import (
	"flag"
	"sync/atomic"
	"time"

	"github.com/MapleLove2014/talos-sdk-golang/producer"
	"github.com/MapleLove2014/talos-sdk-golang/thrift/message"
	"github.com/MapleLove2014/talos-sdk-golang/thrift/thrift"
	"github.com/MapleLove2014/talos-sdk-golang/utils"
)

func main() {
	var propertyFilename string
	flag.StringVar(&propertyFilename, "conf", "simpleProducer.conf", "conf: simpleProducer.conf'")
	flag.Parse()

	log := utils.InitLogger()
	simpleProducer, err := producer.NewSimpleProducerWithLogger(propertyFilename, log)
	if err != nil {
		log.Errorf("Init simpleProducer failed: %s", err.Error())
		return
	}

	msg := &message.Message{
		Message: []byte("test message: this message is a text string."),
	}
	msgList := make([]*message.Message, 0)
	msgList = append(msgList, msg)

	successPutNumber := thrift.Int64Ptr(0)
	// a toy demo for putting messages to Talos server continuously
	// by using a infinite loop
	for true {
		err := simpleProducer.PutMessageList(msgList)
		if err != nil {
			log.Errorf("putMessageList error: %s, try again", err.Error())
			return
		}
		time.Sleep(2000 * time.Millisecond)
		log.Infof("success put message count: %d", atomic.LoadInt64(successPutNumber))
		atomic.StoreInt64(successPutNumber, (*successPutNumber)+1)
	}
}
