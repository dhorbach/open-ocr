package ocrworker

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"testing"
)

func init() {
	logg.LogKeys["TEST"] = true
	logg.LogKeys["OCR_CLIENT"] = true
	logg.LogKeys["OCR_WORKER"] = true
	logg.LogKeys["OCR_HTTP"] = true
	logg.LogKeys["OCR_TESSERACT"] = true
}

func rabbitConfigForTests() RabbitConfig {

	rabbitConfig := RabbitConfig{
		AmqpURI:            "amqp://guest:guest@localhost:5672/",
		Exchange:           "test-exchange",
		ExchangeType:       "direct",
		RoutingKey:         "test-key",
		CallbackRoutingKey: "callback-key",
		Reliable:           true,
		QueueName:          "test-queue",
		CallbackQueueName:  "callback-queue",
	}
	return rabbitConfig
}

// This test assumes that rabbit mq is running
func TestOcrRpcClientIntegration(t *testing.T) {

	// TODO: serve this up through a fake webserver
	// that reads from the filesystem
	testImageUrl := "http://localhost:8080/img"

	rabbitConfig := rabbitConfigForTests()

	// kick off a worker
	// this would normally happen on a different machine ..
	ocrWorker, err := NewOcrRpcWorker(rabbitConfig)
	if err != nil {
		logg.LogTo("TEST", "err: %v", err)
	}
	ocrWorker.Run()

	ocrClient, err := NewOcrRpcClient(rabbitConfig)
	if err != nil {
		logg.LogTo("TEST", "err: %v", err)
	}
	assert.True(t, err == nil)
	decodeResult, err := ocrClient.DecodeImageUrl(testImageUrl, ENGINE_MOCK)
	if err != nil {
		logg.LogTo("TEST", "err: %v", err)
	}
	assert.True(t, err == nil)
	logg.LogTo("TEST", "decodeResult: %v", decodeResult)
	assert.Equals(t, decodeResult.Text, MOCK_ENGINE_RESPONSE)

}
