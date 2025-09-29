package types

import "cosmossdk.io/collections"

// StoredFileKey is the prefix to retrieve all StoredFile
var StoredFileKey = collections.NewPrefix("storedFile/value/")
