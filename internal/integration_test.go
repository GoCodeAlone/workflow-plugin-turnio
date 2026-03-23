package internal_test

import (
	"testing"

	"github.com/GoCodeAlone/workflow/wftest"
)

// TestIntegration_SendText verifies that a pipeline using step.turnio_send_text
// can be mocked and that downstream step.set output is visible in the result.
func TestIntegration_SendText(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  notify:
    steps:
      - name: send
        type: step.turnio_send_text
        config:
          to: "+15555555555"
          body: "hello"
      - name: confirm
        type: step.set
        config:
          values:
            sent: true
`),
		wftest.MockStep("step.turnio_send_text", wftest.Returns(map[string]any{
			"response": `{"messages":[{"id":"msg1"}]}`,
			"to":       "+15555555555",
		})),
	)

	result := h.ExecutePipeline("notify", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["sent"] != true {
		t.Errorf("expected sent=true, got %v", result.Output["sent"])
	}
}

// TestIntegration_CheckContact verifies that step.turnio_check_contact can be
// mocked and that its output key is available to subsequent steps.
func TestIntegration_CheckContact(t *testing.T) {
	rec := wftest.RecordStep("step.turnio_check_contact")
	rec.WithOutput(map[string]any{
		"contacts": `{"contacts":[{"input":"+15555555555","status":"valid"}]}`,
	})

	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  check:
    steps:
      - name: lookup
        type: step.turnio_check_contact
        config:
          phone: "+15555555555"
      - name: mark
        type: step.set
        config:
          values:
            checked: true
`),
		rec,
	)

	result := h.ExecutePipeline("check", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["checked"] != true {
		t.Errorf("expected checked=true, got %v", result.Output["checked"])
	}
	if rec.CallCount() != 1 {
		t.Errorf("expected 1 call to turnio_check_contact, got %d", rec.CallCount())
	}
}

// TestIntegration_ListMessages verifies that step.turnio_list_messages can be
// mocked and that a pipeline runs end-to-end without error.
func TestIntegration_ListMessages(t *testing.T) {
	h := wftest.New(t, wftest.WithYAML(`
pipelines:
  fetch:
    steps:
      - name: list
        type: step.turnio_list_messages
        config:
          contact_id: "contact-abc"
      - name: done
        type: step.set
        config:
          values:
            fetched: true
`),
		wftest.MockStep("step.turnio_list_messages", wftest.Returns(map[string]any{
			"messages": `{"messages":[{"id":"m1","text":{"body":"hi"}}]}`,
		})),
	)

	result := h.ExecutePipeline("fetch", nil)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.Output["fetched"] != true {
		t.Errorf("expected fetched=true, got %v", result.Output["fetched"])
	}
	stepOut := result.StepOutput("list")
	if stepOut == nil {
		t.Fatal("expected step output for 'list'")
	}
	if stepOut["messages"] == nil {
		t.Errorf("expected messages in step output, got %v", stepOut)
	}
}
