package antnet

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
)

const (
	cryptA uint32 = 214003
	cryptB uint32 = 2531012
)

func DefaultNetEncrypt(seed uint32, buf []byte, offset uint32, len uint32) []byte {
	if len <= offset {
		return buf
	}
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.LittleEndian, seed)
	key := b_buf.Bytes()
	k := int32(0)
	c := byte(0)
	for i := offset; i < len; i++ {
		k %= 4
		x := (buf[i] ^ key[k]) + c
		k++
		c = x
		buf[i] = c
	}
	return buf
}

func DefaultNetDecrypt(seed uint32, buf []byte, offset uint32, len uint32) []byte {
	if len <= offset {
		return buf
	}
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.LittleEndian, seed)
	key := b_buf.Bytes()
	k := int32(0)
	c := byte(0)
	for i := offset; i < len; i++ {
		k %= 4
		x := (buf[i] - c) ^ key[k]
		k++
		c = buf[i]
		buf[i] = x
	}
	return buf
}

func CountBCC(buf []byte, offset uint32, len uint32) byte {
	value := byte(0x00)
	for i := offset; i < offset+len; i++ {
		value ^= buf[i]
	}
	return value
}

const HttpISeed uint32 = 3820485974
const HttpOSeed uint32 = 1776134994

var HttpNetDecryptLenNotEnough = NewErr("http net decrypt len not enough")
var HttpNetDecryptBccCheckFailed = NewErr("http net decrypt bcc error")

func HttpNetDecrypt(seed uint32, data string) ([]byte, error) {
	if len(data) < 2 {
		return nil, HttpNetDecryptLenNotEnough
	}
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	dLen := len(dataBytes) - 1
	bcc := dataBytes[0]
	content := DefaultNetDecrypt(seed, dataBytes[1:], 0, uint32(dLen))
	dataBcc := CountBCC(content, 0, uint32(dLen))
	if dataBcc != bcc {
		LogInfo("HttpNetDecrypt bcc error: %v, %v, %v", bcc, dataBcc, content)
		return nil, HttpNetDecryptBccCheckFailed
	}
	return content, nil
}

func HttpNetEncrypt(seed uint32, data []byte) []byte {
	dLen := uint32(len(data))
	cipherData := make([]byte, dLen+1)
	cipherData[0] = CountBCC(data, 0, dLen)
	copy(cipherData[1:], data)
	DefaultNetEncrypt(seed, cipherData[1:], 0, dLen)
	content := []byte(base64.StdEncoding.EncodeToString(cipherData))
	//LogInfo("HttpNetEncrypt: %v, %v", cipherData[0], content)
	return content
}

func HttpNetEncryptJsonPack(seed uint32, msg interface{}) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	newData := make([]byte, len(data)+1)
	copy(newData[1:], data)
	cipherData := HttpNetEncrypt(seed, newData)
	return cipherData, nil
}
