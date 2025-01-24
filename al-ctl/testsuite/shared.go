package testsuite

import (
	"encoding/json"
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model/conversion"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var (
	BaseSchemaPath              string
	SchemaPath                  string
	assetPath                   string
	TargetClass                 string
	discoveryFile               string
	semanticIdentifierInputType bool
)

func transformSemanticIdentifierToAsset() error {
	var discoveryResponse generated.DiscoverResponse
	// Read the asset JSON file
	file, err := os.Open(shared.AssetJsonPath)
	if err != nil {
		log.Err(err).Msg("failed to read asset json file")
		return err
	}
	// unmarshal the asset JSON file
	fileInfo, err := file.Stat()
	if err != nil {
		log.Err(err).Msg("failed to read asset json file structure")
		return err
	}
	byteBuffer := make([]byte, fileInfo.Size())
	_, _ = file.Read(byteBuffer)
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	if err := unmarshalOptions.Unmarshal(byteBuffer, &discoveryResponse); err != nil {
		log.Err(err).Msg("failed to unmarshal asset")
		return err
	}
	// Transform the semantic-identifiers to asset
	testDevice := conversion.TransformDevice(discoveryResponse.Devices[0], "URI", "")
	file, err = os.Create("testsuite.json")
	if err != nil {
		log.Err(err).Msg("failed to create testsuite json file")
		return err
	}
	// change the path to the asset json file
	arr := strings.Split(shared.AssetJsonPath, "/")
	newArr := arr[:len(arr)-1]
	shared.AssetJsonPath = strings.Join(newArr, "/")
	shared.AssetJsonPath = "/testsuite.json"
	// Write the transformed asset to a file
	jsonWriter := json.NewEncoder(file)
	if err := jsonWriter.Encode(testDevice); err != nil {
		log.Err(err).Msg("failed to write transformed asset to file")
		return err
	}
	return nil
}

func GetBaseSchemaVersionFromExtendedSchema() (string, error) {
	file, err := os.Open(SchemaPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var data map[string]interface{}

	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
		return "", err
	}
	imports, ok := data["imports"].([]interface{})
	if !ok {
		return "", fmt.Errorf("imports not found in extended schema")
	}
	if len(imports) < 2 {
		return "", fmt.Errorf("reference to base schema not found in extended schema")
	}
	baseSchemaVersion := imports[1].(string)
	return baseSchemaVersion, nil
}
