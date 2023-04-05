package vkcs

import (
	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func publicDNSZoneStateRefreshFunc(client publicDNSClient, zoneID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		zone, err := zoneGet(client, zoneID).Extract()

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return zone, zoneStatusDeleted, nil
			}
			return nil, "", err
		}
		return zone, zone.Status, nil
	}
}
