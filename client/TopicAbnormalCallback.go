/**
 * Copyright 2018, Xiaomi.
 * All rights reserved.
 * Author: wangfan8@xiaomi.com
 */

package client

import (
	"github.com/MapleLove2014/talos-sdk-golang/thrift/topic"
)

type TopicAbnormalCallback interface {

	/**
	 * User implement this method to process topic abnormal status such as 'TopicNotExist'
	 */
	AbnormalHandler(topicTalosResourceName *topic.TopicTalosResourceName, err error)
}
