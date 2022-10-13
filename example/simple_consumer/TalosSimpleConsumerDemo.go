/**
 * Copyright 2018, Xiaomi.
 * All rights reserved.
 * Author: wangfan8@xiaomi.com
 */

package main

import (
	"flag"
	"time"

	"github.com/MapleLove2014/talos-sdk-golang/consumer"
	"github.com/MapleLove2014/talos-sdk-golang/thrift/common"
	"github.com/MapleLove2014/talos-sdk-golang/utils"
)

func main() {
	log := utils.InitLogger()
	// init client config by put $your_propertyFile in current directory
	// with the content of:
	/*
		    galaxy.talos.service.endpoint=$talosServiceURI
				set your conf path, AK:SK, topicName, and partitionId
	*/
	var propertyFilename string
	flag.StringVar(&propertyFilename, "conf", "simpleConsumer.conf", "conf: simpleConsumer.conf'")
	flag.Parse()

	consumerConfig := consumer.NewTalosConsumerConfigByFilename(propertyFilename)

	//finishedOffset=-2 -> actualStartOffset=-1 and maxFetchNum=1000 set as default

	finishedOffset := int64(-2)
	maxFetchNum := consumerConfig.GetMaxFetchRecords()

	simpleConsumer, err := consumer.NewSimpleConsumerWithLogger(propertyFilename, log)
	if err != nil {
		log.Infof("Init simpleConsumer failed: %s", err.Error())
	}

	stopChan := make(chan utils.StopSign)
	ticker := time.NewTicker(time.Duration(1 * time.Second))
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				messageList, err := simpleConsumer.FetchMessage(finishedOffset+1, maxFetchNum)
				if err != nil {
					if te, ok := err.(*common.GalaxyTalosException); ok {
						log.Errorf("FetchMessage error: %s ", te.GetDetails())
					}
				}
				if len(messageList) > 0 {
					finishedOffset = messageList[len(messageList)-1].GetMessageOffset()
					for i := 0; i < len(messageList); i++ {
						log.Infof("get message: %s success", messageList[i].GetMessage().GetMessage())
					}
					log.Infof("total process %d message", len(messageList))
				}
			case <-stopChan:
				stopChan <- utils.Shutdown
			}
		}
	}()

	<-stopChan
}
