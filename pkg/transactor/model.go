package transactor

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
)

// The platform model this transactor serves, read from the data directory.
//
// It is NOT shipped in this repository and NOT compiled in. The model is a
// third-party artifact under its own licence, and embedding it made this
// repository a redistributor of it — every clone and every binary carried a
// copy. Reading it at run time means a deployment that has the artifact can
// supply it and one that does not simply serves no model, which is the honest
// answer for a repository that does not own it.
//
// Absent is not an error. buildHierarchy answers an empty hierarchy for input it
// cannot parse, and loadModel below returns the same nil for a file that is not
// there, so a server without one starts and serves an empty model rather than
// refusing to boot over a file it cannot legally carry.
const modelFile = "model.json"

// loadModel reads the model from dir, or returns nil when there is none.
func loadModel(dir string) []byte {
	b, err := os.ReadFile(filepath.Join(dir, modelFile))
	if err != nil {
		return nil
	}
	return b
}

// hashModel labels the served model in hello/loadModel. The client uses it to
// skip re-pulling an unchanged model; any stable value is correct since the full
// set is always returned (a mismatch triggers a full pull, never a stale model).
// A content hash is stable across restarts, and an absent model hashes its own
// emptiness rather than being a special case.
func hashModel(model []byte) string {
	sum := sha1.Sum(model)
	return hex.EncodeToString(sum[:])
}
