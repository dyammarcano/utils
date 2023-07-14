package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
	"log"
	"os"

	"github.com/dyammarcano/utils/util"
)

func GzipCompressBytes(filepath string) (result []byte, err error) {
	metadata := util.OpenFileMetadata(filepath, false)

	if metadata.Error != nil {
		return nil, err
	}

	if metadata.Empty {
		return nil, err
	}

	defer func(File *os.File) {
		err := File.Close()
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}(metadata.File)

	buffer := new(bytes.Buffer)

	gzipWriter := gzip.NewWriter(buffer)

	inputData, err := io.ReadAll(metadata.File)
	if err != nil {
		return nil, err
	}

	_, err = gzipWriter.Write(inputData)
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func GzipUncompressBytes(byteArray []byte) (result []byte, err error) {
	buffer := new(bytes.Buffer)

	gzipReader, err := gzip.NewReader(bytes.NewReader(byteArray))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(buffer, gzipReader)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func GzipCompressString(str string) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write([]byte(str)); err != nil {
		return []byte{}, err
	}

	if err := gz.Flush(); err != nil {
		return []byte{}, err
	}

	if err := gz.Close(); err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

func GzipUncompressString(str []byte) (string, error) {
	out := new(bytes.Buffer)

	r, err := gzip.NewReader(bytes.NewBuffer(str))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(out, r); err != nil {
		return "", err
	}

	return out.String(), nil
}

func GzipCompressStruct(obj any) (result []byte, err error) {
	// Serialize the struct into a byte buffer
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(obj); err != nil {
		return nil, err
	}

	// Compress the byte buffer using gzip
	var compressed bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressed)

	if _, err = gzipWriter.Write(buf.Bytes()); err != nil {
		return nil, err
	}

	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}

func GzipUncompressStruct(byteArray []byte, obj any) (err error) {
	// Decompress the byte buffer using gzip
	var buf bytes.Buffer

	if _, err := buf.Write(byteArray); err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(&buf)
	if err != nil {
		return err
	}

	if err := gzipReader.Close(); err != nil {
		return err
	}

	// Deserialize the struct from the uncompressed byte buffer
	dec := gob.NewDecoder(gzipReader)

	if err = dec.Decode(obj); err != nil {
		return err
	}

	return nil
}

func CompressString(input string) (string, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return "", err
	}

	if _, err := gz.Write([]byte(input)); err != nil {
		return "", err
	}

	if err := gz.Flush(); err != nil {
		return "", err
	}

	if err := gz.Close(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func DecompressString(compressed string) (string, error) {
	out := new(bytes.Buffer)

	r, err := gzip.NewReader(bytes.NewBuffer([]byte(compressed)))

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(out, r); err != nil {
		return "", err
	}

	return out.String(), nil
}
