package utilsx

import (
	"encoding/json"

	"github.com/mattbaird/jsonpatch"
)

// CheckDifference the parameters should be ptr
func CheckDifference(oldObj, newObj interface{}) ([]jsonpatch.JsonPatchOperation, error) {
	oldByte, _ := json.Marshal(oldObj)
	newByte, _ := json.Marshal(newObj)

	return jsonpatch.CreatePatch(oldByte, newByte)
}
