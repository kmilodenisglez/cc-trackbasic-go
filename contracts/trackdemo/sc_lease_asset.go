package trackdemo

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	owner, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	assets := []Asset{
		{
			DocType:           DocTypeAsset,
			ID:                "asset1",
			AssetType:         "",
			Owner:             owner,
			CurrentState:      AVAILABLE.String(),
			Location:          "27.1,78.5",
			Manufacturer:      "user 2",
			PublicDescription: "",
			Data:              make([]byte, 0),
			//Data:              "",
		},
		{
			DocType:           DocTypeAsset,
			ID:                "asset2",
			AssetType:         "",
			Owner:             owner,
			CurrentState:      AVAILABLE.String(),
			Location:          "27.1,78.5",
			Manufacturer:      "user 1",
			PublicDescription: "",
			Data:              make([]byte, 0),
			//Data: "",
		},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// Manufacture asset
func (s *SmartContract) Manufacture(ctx contractapi.TransactionContextInterface, request ManufactureAssetRequest) error {
	fmt.Println("1.  ", request)
	exists, err := s.AssetExists(ctx, request.ID)
	fmt.Println("2.  ")
	if err != nil {
		return err
	} else if exists {
		return fmt.Errorf("the asset %s already exists", request.ID)
	}

	owner, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	asset := Asset{
		DocType:           DocTypeAsset,
		ID:                request.ID,
		AssetType:         request.AssetType,
		Owner:             owner,
		CurrentState:      AVAILABLE.String(),
		Location:          request.Location,
		Manufacturer:      request.Manufacturer,
		PublicDescription: request.PublicDescription,
		Data:              make([]byte, 0),
		//Data: "",
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(request.ID, assetJSON)
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, request TransferAssetRequest) (string, error) {
	asset, err := s.ReadAsset(ctx, ReadAssetRequest{request.ID})
	if err != nil {
		return "", err
	}

	ownerIdentity, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	if ownerIdentity != asset.Owner {
		return "", fmt.Errorf("you do not have permission to perform this operation")
	}
	asset.Owner = request.NewOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(request.ID, assetJSON)
	if err != nil {
		return "", err
	}

	return asset.Owner, nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, request ReadAssetRequest) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	} else if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", request.ID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, request UpdateDataAssetRequest) error {
	fmt.Println("1.1 *****    ", request)
	asset, err := s.ReadAsset(ctx, ReadAssetRequest{request.ID})
	fmt.Println("1.2 *****    ", err)
	if err != nil {
		return err
	}

	fmt.Println("1.3 *****    ")
	if request.Location == "" && request.Data == "" && request.PublicDescription == "" {
		return fmt.Errorf("invalid request")
	}
	fmt.Println("1 *****    ", request)

	// overwritting original asset with new asset
	if request.Location != "" {
		asset.Location = request.Location
	}

	if request.PublicDescription != "" {
		asset.PublicDescription = request.PublicDescription
	}

	fmt.Println("2 *****    ")
	if request.Data != "" {
		if valid := json.Valid([]byte(request.Data)); !valid {
			return errors.New("invalid JSON encoding: payload missing or invalid")
		}
		fmt.Println("3 *****    ")
		asset.Data = []byte(request.Data)
		fmt.Println("4 *****    ")
		//asset.Data = request.Data
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(request.ID, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	// range query with empty string for startKey and endKey does an open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &asset}
		results = append(results, queryResult)
	}

	return results, nil
}

func (s *SmartContract) GetEvaluateTransactions() []string {
	return []string{"ReadAsset", "AssetExists", "GetAllAssets"}
}

// GetSubmittingClientIdentity returns the name and issuer of the identity that
// invokes the smart contract. This function base64 decodes the identity string
// before returning the value to the client or smart contract.
func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}
