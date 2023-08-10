package avro

import (
	"fmt"
	"github.com/linkedin/goavro"
	"io"
	"log"
	"net/http"
)

func GetSchema(group, artifactId string) (*goavro.Codec, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/apis/registry/v2/groups/%s/artifacts/%s", group, artifactId))
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Got schema: %s", string(bytes))

	return goavro.NewCodec(string(bytes))
}

func DecodeMessage(data []byte, codec *goavro.Codec) (map[string]interface{}, error) {
	decoded, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Avro data: %w", err)
	}

	decodedMap, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("decoded Avro data is not a map[string]interface{}")
	}

	return decodedMap, nil
}

func ValidateMessage(user interface{}, codec *goavro.Codec) ([]byte, error) {
	return codec.BinaryFromNative(nil, user)
}
