package schema

import (
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/pkg/trace"
	"github.com/tendermint/tendermint/types"
)

// MempoolTables returns the list of tables for mempool tracing.
func MempoolTables() []string {
	return []string{
		MempoolTxTable,
		MempoolPeerStateTable,
	}
}

// Schema constants for the mempool_tx table
const (
	// MempoolTxTable is the tracing "measurement" (aka table) for the mempool
	// that stores tracing data related to gossiping transactions.
	//
	// The schema for this table is:
	// | time | peerID | tx size | tx hash | transfer type | mempool version |
	MempoolTxTable = "mempool_tx"

	// TxFieldKey is the tracing field key for receiving for sending a
	// tx. This should take the form of a tx hash as the value.
	TxFieldKey = "tx"

	// SizeFieldKey is the tracing field key for the size of a tx. This
	// should take the form of the size of the tx as the value.
	SizeFieldKey = "size"

	// VersionFieldKey is the tracing field key for the version of the mempool.
	// This is used to distinguish between versions of the mempool.
	VersionFieldKey = "version"

	// V1VersionFieldValue is a tracing field value for the version of
	// the mempool. This value is used by the "version" field key.
	V1VersionFieldValue = "v1"

	// CatVersionFieldValue is a tracing field value for the version of
	// the mempool. This value is used by the "version" field key.
	CatVersionFieldValue = "cat"
)

// WriteMempoolTx writes a tracing point for a tx using the predetermined
// schema for mempool tracing. This is used to create a table in the following
// schema:
//
// | time | peerID | tx size | tx hash | transfer type | mempool version |
func WriteMempoolTx(client *trace.Client, peer p2p.ID, tx []byte, transferType, version string) {
	// this check is redundant to what is checked during WritePoint, although it
	// is an optimization to avoid allocations from the map of fields.
	if !client.IsCollecting(MempoolTxTable) {
		return
	}
	client.WritePoint(MempoolTxTable, map[string]interface{}{
		TxFieldKey:           bytes.HexBytes(types.Tx(tx).Hash()).String(),
		PeerFieldKey:         peer,
		SizeFieldKey:         len(tx),
		TransferTypeFieldKey: transferType,
		VersionFieldKey:      version,
	})
}

const (
	// MempoolPeerState is the tracing "measurement" (aka table) for the mempool
	// that stores tracing data related to mempool state, specifically
	// the gossipping of "SeenTx" and "WantTx".
	//
	// The schema for this table is:
	// | time | peerID | update type | mempool version |
	MempoolPeerStateTable = "mempool_peer_state"

	// StateUpdateFieldKey is the tracing field key for state updates of the mempool.
	StateUpdateFieldKey = "update"

	// SeenTxStateUpdateFieldValue is a tracing field value for the state
	// update of the mempool. This value is used by the "update" field key.
	SeenTxStateUpdateFieldValue = "seen_tx"

	// WantTxStateUpdateFieldValue is a tracing field value for the state
	// update of the mempool. This value is used by the "update" field key.
	WantTxStateUpdateFieldValue = "want_tx"

	// RemovedTxStateUpdateFieldValue is a tracing field value for the local
	// state update of the mempool. This value is used by the "update" field
	// key.
	RemovedTxStateUpdateFieldValue = "removed_tx"

	// AddedTxStateUpdateFieldValue is a tracing field value for the local state
	// update of the mempool. This value is used by the "update" field key.
	AddedTxStateUpdateFieldValue = "added_tx"
)

// WriteMempoolPeerState writes a tracing point for the mempool state using
// the predetermined schema for mempool tracing. This is used to create a table
// in the following schema:
//
// | time | peerID | transfer type | state update | mempool version |
func WriteMempoolPeerState(client *trace.Client, peer p2p.ID, stateUpdate, transferType, version string) {
	// this check is redundant to what is checked during WritePoint, although it
	// is an optimization to avoid allocations from creating the map of fields.
	if !client.IsCollecting(RoundStateTable) {
		return
	}
	client.WritePoint(RoundStateTable, map[string]interface{}{
		PeerFieldKey:         peer,
		TransferTypeFieldKey: transferType,
		StateUpdateFieldKey:  stateUpdate,
		VersionFieldKey:      version,
	})
}

const (
// LocalTable is the tracing "measurement" (aka table) for the local mempool
// updates, such as when a tx is added or removed.
// TODO: actually implement the local mempool tracing
// LocalTable = "mempool_local"
)
