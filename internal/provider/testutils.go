package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/rs/xid"
)

// randResourceID generates unique id string
// id length must be longer than (prefix + uuid length)
func randResourceID(length int) string {
	// The first character must be alphabet for algolia resources
	uuid := "test_" + xid.New().String()

	if length < len(uuid) {
		panic(fmt.Sprintf("length must be longer than %d", len(uuid)))
	}

	return uuid + acctest.RandStringFromCharSet(length-len(uuid), acctest.CharSetAlphaNum)
}
