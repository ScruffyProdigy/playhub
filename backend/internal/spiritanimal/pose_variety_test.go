package spiritanimal

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestEffectivePoseForImageKeepsSpecificLLMPose(t *testing.T) {
	readingID := uuid.MustParse("4c6e80f2-b911-4c3e-9bb3-b599dcd922a6")
	got := effectivePoseForImage(readingID, "fox", "playful pounce crouch with ears perked")
	if got != "playful pounce crouch with ears perked" {
		t.Fatalf("expected LLM pose to be kept, got %q", got)
	}
}

func TestEffectivePoseForImageReplacesGenericPose(t *testing.T) {
	readingID := uuid.MustParse("4c6e80f2-b911-4c3e-9bb3-b599dcd922a6")
	got := effectivePoseForImage(readingID, "fox", "standing upright facing forward")
	if isGenericMascotPose(normalizeAnswerLabel(got)) {
		t.Fatalf("expected pool pose, got generic %q", got)
	}
}

func TestPickPoseVariantSpreadsAcrossReadings(t *testing.T) {
	seen := map[string]struct{}{}
	for i := range 20 {
		readingID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("reading-%d", i)))
		seen[pickPoseVariant(readingID, "fox")] = struct{}{}
	}
	if len(seen) < 3 {
		t.Fatalf("expected at least 3 distinct fox poses across 20 readings, got %d", len(seen))
	}
}

func TestPickPoseVariantStableForSameReading(t *testing.T) {
	readingID := uuid.MustParse("4c6e80f2-b911-4c3e-9bb3-b599dcd922a6")
	a := pickPoseVariant(readingID, "fox")
	b := pickPoseVariant(readingID, "fox")
	if a != b {
		t.Fatalf("expected stable pose, got %q and %q", a, b)
	}
}
