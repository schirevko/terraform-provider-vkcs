package blockstorage_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

func TestAccBlockStorageSnapshot_basic(t *testing.T) {
	var snapshot snapshots.Snapshot

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageSnapshotBasic, map[string]string{"TestAccBlockStorageVolumeBasic": acctest.AccTestRenderConfig(testAccBlockStorageVolumeBasic)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageSnapshotExists("vkcs_blockstorage_snapshot.snapshot_1", &snapshot),
					testAccCheckBlockStorageSnapshotMetadata(&snapshot, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_snapshot.snapshot_1", "name", "snapshot_1"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_snapshot.snapshot_1", "description", "first test snapshot"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageSnapshotUpdate, map[string]string{"TestAccBlockStorageVolumeBasic": acctest.AccTestRenderConfig(testAccBlockStorageVolumeBasic)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageSnapshotExists("vkcs_blockstorage_snapshot.snapshot_1", &snapshot),
					testAccCheckBlockStorageSnapshotMetadata(&snapshot, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_snapshot.snapshot_1", "name", "snapshot_1-updated"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_snapshot.snapshot_1", "description", "first test snapshot-updated"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageSnapshotExists(n string, snapshot *snapshots.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("volume snapshot not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		blockStorageClient, err := config.BlockStorageV3Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating block storage client: %s", err)
		}

		found, err := snapshots.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("volume not found")
		}

		*snapshot = *found
		return nil
	}
}

func testAccCheckBlockStorageSnapshotDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_blockstorage_snapshot" {
			continue
		}

		_, err := snapshots.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("snapshot still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageSnapshotMetadata(
	snapshot *snapshots.Snapshot, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if snapshot.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range snapshot.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

const testAccBlockStorageSnapshotBasic = `
{{.TestAccBlockStorageVolumeBasic}}

resource "vkcs_blockstorage_snapshot" "snapshot_1" {
  volume_id = vkcs_blockstorage_volume.volume_1.id
  name = "snapshot_1"
  description = "first test snapshot"
  metadata = {
    foo = "bar"
  }
}
`

const testAccBlockStorageSnapshotUpdate = `
{{.TestAccBlockStorageVolumeBasic}}

resource "vkcs_blockstorage_snapshot" "snapshot_1" {
  volume_id = vkcs_blockstorage_volume.volume_1.id
  name = "snapshot_1-updated"
  description = "first test snapshot-updated"
  metadata = {
    foo = "bar"
  }
}
`
