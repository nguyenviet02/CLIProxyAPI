package responses

import (
	"context"
	"testing"

	"github.com/tidwall/gjson"
)

func TestConvertOpenAIChatCompletionsResponseToOpenAIResponsesNonStream_PassesThroughDirectResponseObject(t *testing.T) {
	input := []byte(`{"id":"resp_passthrough","object":"response","status":"completed","output":[{"type":"message","id":"msg_1","content":[{"type":"output_text","text":"hello"}]}],"usage":{"input_tokens":11,"output_tokens":7,"total_tokens":18}}`)

	output := ConvertOpenAIChatCompletionsResponseToOpenAIResponsesNonStream(context.Background(), "", nil, nil, input, nil)
	if gjson.GetBytes(output, "id").String() != "resp_passthrough" {
		t.Fatalf("unexpected id: %s", string(output))
	}
	if gjson.GetBytes(output, "output.0.content.0.text").String() != "hello" {
		t.Fatalf("unexpected output: %s", string(output))
	}
	if gjson.GetBytes(output, "usage.output_tokens").Int() != 7 {
		t.Fatalf("unexpected usage: %s", string(output))
	}
}
