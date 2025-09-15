// fs/fat32.go
package fs

import (
	"encoding/binary"
	"errors"
)

// FAT32 boot sector structure (first 512 bytes of partition)
type BootSector struct {
	BytesPerSector    uint16
	SectorsPerCluster byte
	ReservedSectors   uint16
	NumFATs           byte
	RootEntryCount    uint16
	TotalSectors16    uint16
	Media             byte
	FATSize16         uint16
	SectorsPerTrack   uint16
	NumHeads          uint16
	HiddenSectors     uint32
	TotalSectors32    uint32
	FATSize32         uint32
	ExtFlags          uint16
	FSVersion         uint16
	RootCluster       uint32
	FSInfo            uint16
	BackupBoot        uint16
}

// Minimal FAT32 driver
type FAT32 struct {
	Boot              BootSector
	FATStart          uint32
	DataStart         uint32
	SectorsPerCluster uint8
	BytesPerSector    uint16
	RootCluster       uint32
}

// --- Block device read function ---
// You need to implement this for your Pi Zero 2W SD card driver
var ReadBlock func(blockNum uint32, buf []byte) error

// Initialize FAT32 from boot sector
func Mount() (*FAT32, error) {
	buf := make([]byte, 512)
	if ReadBlock(0, buf) != nil {
		return nil, errors.New("cannot read boot sector")
	}

	bs := BootSector{
		BytesPerSector:    binary.LittleEndian.Uint16(buf[11:13]),
		SectorsPerCluster: buf[13],
		ReservedSectors:   binary.LittleEndian.Uint16(buf[14:16]),
		NumFATs:           buf[16],
		RootEntryCount:    binary.LittleEndian.Uint16(buf[17:19]),
		TotalSectors16:    binary.LittleEndian.Uint16(buf[19:21]),
		Media:             buf[21],
		FATSize16:         binary.LittleEndian.Uint16(buf[22:24]),
		SectorsPerTrack:   binary.LittleEndian.Uint16(buf[24:26]),
		NumHeads:          binary.LittleEndian.Uint16(buf[26:28]),
		HiddenSectors:     binary.LittleEndian.Uint32(buf[28:32]),
		TotalSectors32:    binary.LittleEndian.Uint32(buf[32:36]),
		FATSize32:         binary.LittleEndian.Uint32(buf[36:40]),
		ExtFlags:          binary.LittleEndian.Uint16(buf[40:42]),
		FSVersion:         binary.LittleEndian.Uint16(buf[42:44]),
		RootCluster:       binary.LittleEndian.Uint32(buf[44:48]),
		FSInfo:            binary.LittleEndian.Uint16(buf[48:50]),
		BackupBoot:        binary.LittleEndian.Uint16(buf[50:52]),
	}

	fatStart := uint32(bs.ReservedSectors)
	dataStart := fatStart + uint32(bs.NumFATs)*bs.FATSize32

	return &FAT32{
		Boot:              bs,
		FATStart:          fatStart,
		DataStart:         dataStart,
		SectorsPerCluster: bs.SectorsPerCluster,
		BytesPerSector:    bs.BytesPerSector,
		RootCluster:       bs.RootCluster,
	}, nil
}

// Convert cluster number to absolute LBA block
func (f *FAT32) ClusterToLBA(cluster uint32) uint32 {
	return f.DataStart + (cluster-2)*uint32(f.SectorsPerCluster)
}
