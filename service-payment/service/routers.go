package service

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	goavro "gopkg.in/linkedin/goavro.v1"
)

const (
	PRODUCER_URL string = "kafka:9092"
	KAFKA_TOPIC  string = "payment"
)

func PaymentRegister(router *gin.RouterGroup) {
	router.POST("/", PaymentCreate)
}

func ServiceCheckRegister(router *gin.RouterGroup) {
	router.GET("/__health__", HelthChceckRetrieve)
	router.GET("/__version__", VersionRetrieve)
}

func PaymentCreate(c *gin.Context) {
	var p CreatePaymentFeed

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordSchemaJSON := `
	{
	  "type": "record",
	  "name": "payments",
	  "doc:": "A basic schema for storing payments data",
	  "namespace": "banking",
	  "fields": [
		{
		  "doc": "UId of the paymane",
		  "type": "string",
		  "name": "paymentUId"
		},
		{
		  "doc": "Destination account UId",
		  "type": "string",
		  "name": "destinationAccountUId"
		},
		{
		  "doc": "Payee UId",
		  "type": "string",
		  "name": "payeeUId"
		},
		{
		  "doc": "Payment amount",
		  "type": "long",
		  "name": "amount"
		},
		{
		  "doc": "Event status",
		  "type": "string",
		  "name": "status"
		},
		{
		  "doc": "Event created time",
		  "type": "string",
		  "name": "created"
		},
		{
		  "doc": "Unix epoch time in milliseconds",
		  "type": "long",
		  "name": "timestamp"
		}
	  ]
	}
	`

	record, err := goavro.NewRecord(goavro.RecordSchema(recordSchemaJSON))
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":   err.Error(),
			"message": "Wrong schema is set",
		})
		return
	}
	event := Event(p)

        timeStamp := int64(time.Now().Unix())
	strTime := strconv.FormatInt(timeStamp, 10) // 10 is a magic number

	record.Set("paymentUId", event.PayementUId)
	record.Set("destinationAccountUId", event.DestinationAccountUid)
	record.Set("payeeUId", event.PayeeUId)
	record.Set("amount", event.Amount)
	record.Set("created", event.Created.Format("2006-01-0T15:04:05Z"))
	record.Set("status", event.Status)
	record.Set("timestamp", timeStamp)

	codec, err := goavro.NewCodec(recordSchemaJSON)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":   err.Error(),
			"message": "Not able to create codec",
		})
		return
	}

	bb := new(bytes.Buffer)
	if err = codec.Encode(bb, record); err != nil {
		fmt.Println(timeStamp)
		fmt.Println(record)
		c.IndentedJSON(http.StatusTooManyRequests, gin.H{
			"error":   err.Error(),
			"message": "Not able to encode to byte stream",
		})
		return
	}
	actual := bb.Bytes()

	dataString := string(actual)
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{PRODUCER_URL}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":   err.Error(),
			"message": "Not able to create async producer",
		})
		return
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// TODO: should write to logs
			panic(err)
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: KAFKA_TOPIC,
		Key:   sarama.StringEncoder(strTime),
		Value: sarama.StringEncoder(dataString),
	}
	producer.Input() <- msg

	// c.JSON(http.StatusOK, event)

	c.IndentedJSON(http.StatusOK, event)
}

func HelthChceckRetrieve(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func VersionRetrieve(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": Version})
}

// TODO: move to utils
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}


