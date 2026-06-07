package spiritanimal

import (
	"hash/fnv"
	"strings"

	"github.com/google/uuid"
)

// mascotPoseVariants are thumbnail-readable poses shared across the player base.
// When many users pick the same animal, hash-based selection spreads poses out.
var mascotPoseVariants = []string{
	"seated with tail curled around paws",
	"one paw raised in greeting, weight shifted back",
	"leaning forward eagerly on front paws",
	"perched on haunches, alert and compact",
	"playful pounce crouch with hindquarters raised",
	"cozy curl, nose tucked toward tail",
	"mid-stride step with one paw lifted",
	"looking over shoulder, three-quarter body turn",
	"head tilted curious, ears perked",
	"reclining on side, relaxed and open",
	"rearing up playfully on hind legs",
	"wings or arms half-spread in a welcoming gesture",
}

func effectivePoseForImage(readingID uuid.UUID, animal, llmPose string) string {
	llmPose = strings.TrimSpace(llmPose)
	if llmPose != "" && !isGenericMascotPose(normalizeAnswerLabel(llmPose)) {
		return llmPose
	}
	return pickPoseVariant(readingID, animal)
}

func pickPoseVariant(readingID uuid.UUID, animal string) string {
	key := readingID.String() + "\x00" + strings.ToLower(strings.TrimSpace(animal))
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	idx := h.Sum64() % uint64(len(mascotPoseVariants))
	return mascotPoseVariants[idx]
}

func isGenericMascotPose(normalized string) bool {
	generic := []string{
		"standing facing forward",
		"standing upright facing forward",
		"standing facing viewer",
		"standing upright facing viewer",
		"standing facing camera",
		"standing upright",
		"facing forward",
		"facing viewer",
	}
	for _, phrase := range generic {
		if normalized == phrase || strings.Contains(normalized, phrase) {
			return true
		}
	}
	return false
}
