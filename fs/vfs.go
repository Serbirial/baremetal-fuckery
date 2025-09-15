// fs/vfs.go
package fs

import (
	"encoding/binary"
	"errors"
	"strings"
)

// FileEntry represents a FAT32 directory entry
type FileEntry struct {
	Name    string
	Cluster uint32
	Size    uint32
}

// VFS represents a very simple filesystem
type VFS struct {
	FAT  *FAT32
	Root []FileEntry
}

// MountVFS reads the root directory and populates file entries
func MountVFS(fat *FAT32) (*VFS, error) {
	vfs := &VFS{
		FAT: fat,
	}

	// Root directory is at FAT32.RootCluster
	cluster := fat.RootCluster
	var buf = make([]byte, int(fat.BytesPerSector)*int(fat.SectorsPerCluster))

	for {
		// Read cluster data
		for i := uint32(0); i < uint32(fat.SectorsPerCluster); i++ {
			if err := ReadBlock(fat.ClusterToLBA(cluster)+i, buf[int(i)*int(fat.BytesPerSector):int(i+1)*int(fat.BytesPerSector)]); err != nil {
				return nil, err
			}
		}

		// Each directory entry is 32 bytes
		for offset := 0; offset < len(buf); offset += 32 {
			entry := buf[offset : offset+32]
			if entry[0] == 0x00 {
				// End of directory
				return vfs, nil
			}
			if entry[0] == 0xE5 || entry[11]&0x0F == 0x0F {
				// Deleted entry or LFN skip
				continue
			}

			name := strings.TrimSpace(string(entry[0:8]))
			ext := strings.TrimSpace(string(entry[8:11]))
			if ext != "" {
				name += "." + ext
			}

			highCluster := uint32(entry[20]) | uint32(entry[21])<<8 | uint32(entry[22])<<16 | uint32(entry[23])<<24
			lowCluster := uint32(entry[26]) | uint32(entry[27])<<8
			clusterNum := (highCluster << 16) | lowCluster
			size := binary.LittleEndian.Uint32(entry[28:32])

			vfs.Root = append(vfs.Root, FileEntry{
				Name:    name,
				Cluster: clusterNum,
				Size:    size,
			})
		}

		// TODO: follow cluster chain for larger directories
		break // currently only single-cluster root
	}

	return vfs, nil
}

// OpenFile returns the FileEntry for a given filename
func (v *VFS) OpenFile(filename string) (*FileEntry, error) {
	for _, f := range v.Root {
		if strings.EqualFold(f.Name, filename) {
			return &f, nil
		}
	}
	return nil, errors.New("file not found")
}

// ReadFile reads the entire file content into a buffer
func (v *VFS) ReadFile(f *FileEntry, buf []byte) error {
	cluster := f.Cluster
	bytesLeft := f.Size
	offset := 0
	clusterBuf := make([]byte, int(v.FAT.BytesPerSector)*int(v.FAT.SectorsPerCluster))

	for bytesLeft > 0 {
		// Read entire cluster
		for i := uint32(0); i < uint32(v.FAT.SectorsPerCluster); i++ {
			start := int(i) * int(v.FAT.BytesPerSector)
			end := int(i+1) * int(v.FAT.BytesPerSector)
			if err := ReadBlock(v.FAT.ClusterToLBA(cluster)+i, clusterBuf[start:end]); err != nil {
				return err
			}
		}

		// Copy to output buffer
		toCopy := len(clusterBuf)
		if uint32(toCopy) > bytesLeft {
			toCopy = int(bytesLeft)
		}
		copy(buf[offset:], clusterBuf[:toCopy])
		offset += toCopy
		bytesLeft -= uint32(toCopy)

		// TODO: follow FAT table for next cluster
		break // currently only single-cluster files
	}

	return nil
}
