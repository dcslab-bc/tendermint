package mempool

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func BenchmarkTxMempool_CheckTx(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	txmp := setup(ctx, b, 10000)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		prefix := make([]byte, 20)
		_, err := rng.Read(prefix)
		require.NoError(b, err)

		priority := int64(rng.Intn(9999-1000) + 1000)
		tx := []byte(fmt.Sprintf("%X=%d", prefix, priority))
		b.StartTimer()

		require.NoError(b, txmp.CheckTx(ctx, tx, nil, TxInfo{}))
	}
}
