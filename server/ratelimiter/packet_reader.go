package ratelimiter

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

// PacketReader 是基于 bufio.Reader 的自定义 Reader
type PacketReader struct {
	reader *bufio.Reader
}

// PacketHeader 表示数据包的包头信息
type PacketHeader struct {
	TotalSize    uint32
	PacketSize   uint16
	TotalPackets uint16
	Seq          uint16
}

// NewPacketReader 构造函数，基于 bufio.Reader 创建自定义的 PacketReader
func NewPacketReader(r io.Reader) *PacketReader {
	return &PacketReader{
		reader: bufio.NewReader(r),
	}
}

// ReadPacket 读取并返回一个完整的数据包（包含包头和包体）
func (cr *PacketReader) ReadPacket() ([]byte, error) {
	// 读取包头（10 字节）
	headerBytes := make([]byte, 10)
	if _, err := io.ReadFull(cr.reader, headerBytes); err != nil {
		return make([]byte, 0), err
	}

	// 解析包头
	header, err := parseHeader(headerBytes)
	if err != nil {
		return make([]byte, 0), err
	}

	// 根据包头中的 packetSize 读取包体
	body := make([]byte, header.PacketSize)
	if _, err := io.ReadFull(cr.reader, body); err != nil {
		return make([]byte, 0), err
	}

	// 返回完整的数据包
	return body, nil
}

// 解析包头
func parseHeader(header []byte) (*PacketHeader, error) {
	reader := bytes.NewReader(header)
	packetHeader := &PacketHeader{}

	// 按大端序读取包头的各个字段
	if err := binary.Read(reader, binary.BigEndian, &packetHeader.TotalSize); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &packetHeader.PacketSize); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &packetHeader.TotalPackets); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &packetHeader.Seq); err != nil {
		return nil, err
	}

	return packetHeader, nil
}
