package grpc

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/NethermindEth/juno/db"
	"github.com/NethermindEth/juno/grpc/gen"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient(t *testing.T) {
	t.Skip("manual testing")

	conn, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := gen.NewKVClient(conn)
	stream, err := client.Tx(context.Background())
	require.NoError(t, err)

	err = stream.Send(&gen.Cursor{
		Op: gen.Op_SEEK,
		K:  db.ChainHeight.Key(),
	})
	require.NoError(t, err)

	err = stream.Send(&gen.Cursor{
		Op: gen.Op_CURRENT,
	})
	require.NoError(t, err)

	pair, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			fmt.Println("disconnected from server")
		} else {
			spew.Dump("error", err)
		}
	}

	spew.Dump(pair)
}
