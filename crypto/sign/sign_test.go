package sign

import (
	"fmt"
	"github.com/dyammarcano/utils/mocks"
	"testing"
)

var (
	key  = []byte("secret")
	data = []byte(`Hello, world!
This is a test file.`)
)

func TestInterfaces(t *testing.T) {
	encoders := []Encoder{Base64Encoder{}, Base58Encoder{}, HexEncoder{}, BynaryEncoder{}}

	for _, encoder := range encoders {
		encodedString := encoder.Encode(data)
		fmt.Printf("%T: %s\n", encoder, encodedString)
	}
}

func TestSignFile(t *testing.T) {
	testfile, err := mocks.NewLoremTestFile(t, 10)
	if err != nil {
		t.Error(err)
		return
	}

	defer testfile.CleanUp()

	signature := SignFile(key, testfile.File, Base64)
	fmt.Printf("%s\n", signature)
}

func TestSignString(t *testing.T) {
	signature := SignString(key, data, Base64)
	fmt.Printf("%s\n", signature)
}

func TestSignBytes(t *testing.T) {
	signature := SignBytes(key, data, Base64)
	fmt.Printf("%s\n", signature)
}

//func TestSignFileGPG(t *testing.T) {
//	f, err := os.Open("C:\\Users\\dyamm\\Downloads\\ojdbc11-full.tar.gz")
//	if err != nil {
//		t.Error(err)
//	}
//
//	signature := SignFileGPG(f)
//	fmt.Printf("%s\n", signature)
//}
