package networking

import (
	"encoding/json"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func NetworkingReadAttributesTags(d *schema.ResourceData, tags []string) {
	util.ExpandObjectReadTags(d, tags)
}

func NetworkingV2UpdateAttributesTags(d *schema.ResourceData) []string {
	return util.ExpandObjectUpdateTags(d)
}

func NetworkingAttributesTags(d *schema.ResourceData) []string {
	return util.ExpandObjectTags(d)
}

type neutronErrorWrap struct {
	NeutronError neutronError
}

type neutronError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Detail  string `json:"detail"`
}

func retryOn409(err error) bool {
	switch err := err.(type) {
	case gophercloud.ErrDefault409:
		neutronError, e := decodeNeutronError(err.ErrUnexpectedResponseCode.Body)
		if e != nil {
			// retry, when error type cannot be detected
			log.Printf("[DEBUG] failed to decode a neutron error: %s", e)
			return true
		}
		if neutronError.Type == "IpAddressGenerationFailure" {
			return true
		}

		// don't retry on quota or other errors
		return false
	case gophercloud.ErrDefault400:
		neutronError, e := decodeNeutronError(err.ErrUnexpectedResponseCode.Body)
		if e != nil {
			// retry, when error type cannot be detected
			log.Printf("[DEBUG] failed to decode a neutron error: %s", e)
			return true
		}
		if neutronError.Type == "ExternalIpAddressExhausted" {
			return true
		}

		// don't retry on quota or other errors
		return false
	case gophercloud.ErrDefault404: // this case is handled mostly for functional tests
		return true
	}

	return false
}

func decodeNeutronError(body []byte) (*neutronError, error) {
	e := &neutronErrorWrap{}
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}

	return &e.NeutronError, nil
}

func GetSDN(d *schema.ResourceData) string {
	if v, ok := d.GetOk("sdn"); ok {
		return v.(string)
	}

	return networking.SearchInAllSDNs
}

func ValidateSDN() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"neutron", "sprut"}, true))
}
