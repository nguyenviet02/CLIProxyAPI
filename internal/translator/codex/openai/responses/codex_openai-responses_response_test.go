package responses

import (
	"context"
	"testing"

	"github.com/tidwall/gjson"
)

func TestConvertCodexResponseToOpenAIResponsesNonStream_UnwrapsCompletedEvent(t *testing.T) {
	input := []byte(`{"type":"response.completed","response":{"id":"resp_1","object":"response","status":"completed","output":[{"type":"message","id":"msg_1"}]}}`)

	output := ConvertCodexResponseToOpenAIResponsesNonStream(context.Background(), "", nil, nil, input, nil)
	if gjson.GetBytes(output, "id").String() != "resp_1" {
		t.Fatalf("unexpected id: %s", string(output))
	}
	if gjson.GetBytes(output, "output.0.id").String() != "msg_1" {
		t.Fatalf("unexpected output: %s", string(output))
	}
}

func TestConvertCodexResponseToOpenAIResponsesNonStream_PassesThroughDirectResponseObject(t *testing.T) {
	input := []byte(`{"id":"resp_2","object":"response","status":"completed","output":[{"type":"message","id":"msg_2","content":[{"type":"output_text","text":"hello"}]}]}`)

	output := ConvertCodexResponseToOpenAIResponsesNonStream(context.Background(), "", nil, nil, input, nil)
	if gjson.GetBytes(output, "id").String() != "resp_2" {
		t.Fatalf("unexpected id: %s", string(output))
	}
	if gjson.GetBytes(output, "output.0.content.0.text").String() != "hello" {
		t.Fatalf("unexpected output text: %s", string(output))
	}
}
