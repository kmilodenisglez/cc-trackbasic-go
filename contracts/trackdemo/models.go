package trackdemo

const (
	DocTypeAsset = "org.asset"
)

const (
	ContractNameLeaseTerm = "org.trackdemo"
)

// Status enum
type Status uint

// status in leasing terms
const (
	AVAILABLE  Status = iota + 1 // disponible
	INSPECTING                   // inspeccionando
	REPAIRING                    // reparando
)

func (term Status) String() string {
	names := []string{"available", "inspecting", "repairing"}
	if term < AVAILABLE || term > REPAIRING {
		return "unknown"
	}

	return names[term-1]
}

// Condition enum
type Condition uint

// Enumerate asset condition values
const (
	NEW          Condition = iota + 1 // nuevo
	REFURBISHED                       // renovado
	NEEDS_REPAIR                      // necesita reparación
)

func (term Condition) String() string {
	names := []string{"news", "refurbished", "needs_repair"}
	if term < NEW || term > NEEDS_REPAIR {
		return "unknown"
	}

	return names[term-1]
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DocType           string `json:"docType"` // DocType is used to distinguish different object types in the same chaincode namespace
	ID                string `json:"ID"`      // ID is a predefined unique UUID
	AssetType         string `json:"assetType"`
	Owner             string `json:"owner"`
	CurrentState      string `json:"currentState"`
	Location          string `json:"location"` // Location is the current location of the product, included in any requests sent to the blockchain
	Manufacturer      string `json:"manufacturer"`
	PublicDescription string `json:"publicDescription"`
	Data              []byte `json:"data,omitempty"` // Data save a custom json in []byte format
	// Data              string `json:"data,omitempty"` // Data save a custom json in string format
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Asset
}

type ReadAssetRequest struct {
	ID string `json:"id"`
}

type ManufactureAssetRequest struct {
	ID                string `json:"id"`
	AssetType         string `json:"assetType"`
	Location          string `json:"location"`
	Manufacturer      string `json:"manufacturer"`
	PublicDescription string `json:"publicDescription" metadata:",optional"`
}

type TransferAssetRequest struct {
	ID           string `json:"id"`
	AssetType    string `json:"assetType"`
	Location     string `json:"location"`
	Manufacturer string `json:"manufacturer"`
	NewOwner     string `json:"newOwner"`
}

type UpdateDataAssetRequest struct {
	ID                string `json:"id"`
	Location          string `json:"location" metadata:",optional"`
	PublicDescription string `json:"publicDescription" metadata:",optional"`
	Data              string `json:"data" metadata:",optional"` // Data save a custom json
}
