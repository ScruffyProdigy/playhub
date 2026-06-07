package spiritanimal

import (
	"testing"
	"time"
)

func TestRecordQuestionsDurationUpdatesEstimate(t *testing.T) {
	questionsEstimateNanos.Store(defaultQuestionsEstimate.Nanoseconds())
	RecordQuestionsDuration(30 * time.Second)
	got := QuestionsEstimateSeconds()
	if got < 30 || got > 45 {
		t.Fatalf("expected estimate between 30 and 45 after 30s observation, got %d", got)
	}
}

func TestEstimateSecondsForStatus(t *testing.T) {
	if got := EstimateSecondsForStatus("generating_questions"); got <= 0 {
		t.Fatalf("expected positive questions estimate, got %d", got)
	}
	if got := EstimateSecondsForStatus("processing"); got <= 0 {
		t.Fatalf("expected positive mascots estimate, got %d", got)
	}
	if got := EstimateSecondsForStatus("ready"); got != 0 {
		t.Fatalf("expected 0 for ready status, got %d", got)
	}
}
