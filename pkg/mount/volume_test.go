package mount

import (
	"fmt"
	"testing"

	"code.storageos.net/scm/storageos/control/testutil"
)

func TestParseFSType(t *testing.T) {

	type data struct {
		out      string
		expected string
		err      error
	}
	// NOTE: The string keys in these test fixtures have no special meaning,
	// they look the way they do as most were taken from a file dump of an
	// active /var/lib/storageos/volumes directory.
	var input = map[string]data{
		"/var/lib/storageos/volumes/raw":              {"/var/lib/storageos/volumes/raw: data", "raw", nil},
		"/var/lib/storageos/volumes/bench1":           {"/var/lib/storageos/volumes/bench1: Linux rev 1.0 ext2 filesystem data, UUID=85e2943e-d7d2-466f-b2ae-03f3dea5f37b (large files)", "ext2", nil},
		"/var/lib/storageos/volumes/bench2":           {"/var/lib/storageos/volumes/bench2: Linux rev 1.0 ext3 filesystem data, UUID=5ab21cda-efd7-4772-9453-62be762ea824 (large files)", "ext3", nil},
		"/var/lib/storageos/volumes/bench3":           {"/var/lib/storageos/volumes/bench3: Linux rev 1.0 ext4 filesystem data, UUID=510227b8-e2c0-49ac-a17a-86fa55431328 (extents) (large files) (huge files)", "ext4", nil},
		"/var/lib/storageos/volumes/bench4":           {"/var/lib/storageos/volumes/bench4: SGI XFS filesystem data (blksz 4096, inosz 512, v2 dirs)", "xfs", nil},
		"/var/lib/storageos/volumes/bench5":           {"/var/lib/storageos/volumes/bench5: BTRFS Filesystem sectorsize 4096, nodesize 16384, leafsize 16384, UUID=ea1d9357-ac2c-4d62-923f-e77e8028442a, 114688/10737418240 bytes used, 1 devices", "btrfs", nil},
		"/var/lib/storageos/volumes/bench6":           {"/var/lib/storageos/volumes/bench6: DOS/MBR boot sector, code offset 0x58+2, OEM-ID \"mkfs.fat\", sectors/cluster 16, Media descriptor 0xf8, sectors/track 63, heads 255, sectors 20971520 (volumes > 32 MB) , FAT (32 bit), sectors/FAT 10231, serial number 0x3bb990a5, unlabeled", "fat", nil},
		"/var/lib/storageos/volumes/bench7":           {"/var/lib/storageos/volumes/bench7: DOS/MBR boot sector, code offset 0x58+2, OEM-ID \"mkfs.fat\", sectors/cluster 16, Media descriptor 0xf8, sectors/track 63, heads 255, sectors 20971520 (volumes > 32 MB) , FAT (32 bit), sectors/FAT 10231, serial number 0x3d140f7c, unlabeled", "fat", nil},
		"/var/lib/storageos/volumes/bench8":           {"/var/lib/storageos/volumes/bench8: DOS/MBR boot sector, code offset 0x52+2, OEM-ID \"NTFS    \", sectors/cluster 8, Media descriptor 0xf8, sectors/track 0, dos < 4.0 BootSector (0x80), FAT (1Y bit by descriptor); NTFS, sectors 20971519, $MFT start cluster 4, $MFTMirror start cluster 1310719, bytes/RecordSegment 2^(-1*246), clusters/index block 1, serial number 02abadec550ea2443; contains Microsoft Windows XP/VISTA bootloader BOOTMGR", "ntfs", nil},
		"/var/lib/storageos/volumes/bench-r1":         {"/var/lib/storageos/volumes/bench-r1: Linux rev 1.0 ext2 filesystem data, UUID=85e2943e-d7d2-466f-b2ae-03f3dea5f37b (large files)", "ext2", nil},
		"/var/lib/storageos/volumes/bench-r2":         {"/var/lib/storageos/volumes/bench-r2: Linux rev 1.0 ext3 filesystem data, UUID=5ab21cda-efd7-4772-9453-62be762ea824 (large files)", "ext3", nil},
		"/var/lib/storageos/volumes/bench-r3":         {"/var/lib/storageos/volumes/bench-r3: Linux rev 1.0 ext4 filesystem data, UUID=510227b8-e2c0-49ac-a17a-86fa55431328 (extents) (large files) (huge files)", "ext4", nil},
		"/var/lib/storageos/volumes/bench-r4":         {"/var/lib/storageos/volumes/bench-r4: SGI XFS filesystem data (blksz 4096, inosz 512, v2 dirs)", "xfs", nil},
		"/var/lib/storageos/volumes/bench-r5":         {"/var/lib/storageos/volumes/bench-r5: BTRFS Filesystem sectorsize 4096, nodesize 16384, leafsize 16384, UUID=ea1d9357-ac2c-4d62-923f-e77e8028442a, 114688/10737418240 bytes used, 1 devices", "btrfs", nil},
		"/var/lib/storageos/volumes/bench-r6":         {"/var/lib/storageos/volumes/bench-r6: DOS/MBR boot sector, code offset 0x58+2, OEM-ID \"mkfs.fat\", sectors/cluster 16, Media descriptor 0xf8, sectors/track 63, heads 255, sectors 20971520 (volumes > 32 MB) , FAT (32 bit), sectors/FAT 10231, serial number 0x3bb990a5, unlabeled", "fat", nil},
		"/var/lib/storageos/volumes/bench-r7":         {"/var/lib/storageos/volumes/bench-r7: DOS/MBR boot sector, code offset 0x58+2, OEM-ID \"mkfs.fat\", sectors/cluster 16, Media descriptor 0xf8, sectors/track 63, heads 255, sectors 20971520 (volumes > 32 MB) , FAT (32 bit), sectors/FAT 10231, serial number 0x3d140f7c, unlabeled", "fat", nil},
		"/var/lib/storageos/volumes/bench-r8":         {"/var/lib/storageos/volumes/bench-r8: DOS/MBR boot sector, code offset 0x52+2, OEM-ID \"NTFS    \", sectors/cluster 8, Media descriptor 0xf8, sectors/track 0, dos < 4.0 BootSector (0x80), FAT (1Y bit by descriptor); NTFS, sectors 20971519, $MFT start cluster 4, $MFTMirror start cluster 1310719, bytes/RecordSegment 2^(-1*246), clusters/index block 1, serial number 02abadec550ea2443; contains Microsoft Windows XP/VISTA bootloader BOOTMGR", "ntfs", nil},
		"/var/lib/storageos/volumes/rdb":              {"/var/lib/storageos/volumes/rdb: block special", "", fmt.Errorf("not nil")},
		"/var/lib/storageos/volumes/some_empty_thing": {"/var/lib/storageos/volumes/some_empty_thing: empty", "raw", nil},
	}

	for path, d := range input {
		fstype, err := parseFileOutput(path, d.out)
		if d.err == nil {
			testutil.Expect(t, err, nil)
		} else {
			testutil.Refute(t, err, nil)
		}
		testutil.Expect(t, fstype, d.expected)
	}

}
