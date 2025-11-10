package ipfs

import (
	"context"
	"os"
	"time"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/root-N-root/webipfs/types"
)

type Client struct {
	node      *core.IpfsNode
	api       coreiface.CoreAPI
	ctx       context.Context
	connector *types.Connector
}

var REPO_NAME = "rootNroot-webipfs-repo-*"

func NewClient(ctx context.Context) (*Client, error) {
	repoPath, err := os.MkdirTemp("", REPO_NAME)
	if err != nil {
		return nil, err
	}
	cfg, err := config.Init(os.Stdout, 2048)
	if err != nil {
		return nil, err
	}

	cfg.Addresses.Swarm = []string{
		"/ip4/0.0.0.0/tcp/0",
		"/ip6/::/tcp/0",
	}
	cfg.Swarm.RelayClient.Enabled = config.True
	cfg.Swarm.RelayService.Enabled = config.True
	cfg.Swarm.EnableHolePunching = config.True
	cfg.Routing.Type = config.NewOptionalString("dht")
	cfg.Swarm.Transports.Network.Relay = config.True

	if err := fsrepo.Init(repoPath, cfg); err != nil {
		return nil, err
	}

	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}
	nodeOptions := &core.BuildCfg{
		Online: true,
		Repo:   r,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": false,
		},
	}
	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		node.Close()
		return nil, err
	}
	bootstrapPeers := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMj3LVhPSu",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTbnCfv8dhxnb1bfRePG",
	}

	for _, addrStr := range bootstrapPeers {
		addr, err := multiaddr.NewMultiaddr(addrStr)
		if err != nil {
			continue
		}
		peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			continue
		}
		if err := node.PeerHost.Connect(ctx, *peerInfo); err != nil {
			// Не критично для bootstrap, продолжаем
		}
	}

	// Запускаем фоновые процессы для поддержания соединений
	go func() {
		// Периодически проверяем соединения с bootstrap-пирами
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				for _, addrStr := range bootstrapPeers {
					addr, err := multiaddr.NewMultiaddr(addrStr)
					if err != nil {
						continue
					}
					peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
					if err != nil {
						continue
					}
					_ = node.PeerHost.Connect(ctx, *peerInfo)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return &Client{
		node: node,
		api:  api,
		ctx:  ctx,
	}, nil
}

func (c *Client) GetSwarmPeers(ctx context.Context) (int, error) {
	peers, err := c.api.Swarm().Peers(ctx)
	if err != nil {
		return 0, err
	}
	return len(peers), nil
}

func (c *Client) GetPeerID() string {
	if c.node.PeerHost != nil {
		return c.node.PeerHost.ID().String()
	}
	return ""
}

func (c *Client) PinFile(ctx context.Context, cidString string) error {
	cidObj, err := cid.Decode(cidString)
	if err != nil {
		return err
	}

	p := path.FromCid(cidObj)

	return c.api.Pin().Add(ctx, p)
}

func (c *Client) GetFile(ctx context.Context, cidString string) (types.FileUpdate, error) {
	cidObj, err := cid.Decode(cidString)
	if err != nil {
		return types.FileUpdate{}, err
	}

	p := path.FromCid(cidObj)

	_, err = c.api.Block().Get(ctx, p)
	if err != nil {
		return types.FileUpdate{}, err
	}

	peers, err := c.GetSwarmPeers(ctx)
	if err != nil {
		peers = 0 // Default if error
	}

	return types.NewFileUpdate(
		types.FuwCid(cidObj.String()),
		types.FuwPeers(peers),
		types.FuwStatus(types.StatusDownloading),
	), nil
}

var globalClient *Client

func CreateFile(filepath string, filename string) types.FileUpdate {
	if globalClient != nil {
		result, err := globalClient.AddFile(filepath, filename)
		if err != nil {
			fu := types.NewFileUpdate(
				types.FuwName(filename),
				types.FuwPath(filepath),
				types.FuwStatus(types.StatusError),
			)

			if globalClient.connector != nil {
				globalClient.connector.SendFileUp(fu)
			}

			return fu
		}

		if globalClient.connector != nil {
			globalClient.connector.SendFileUp(result)
		}

		return result
	}
	return types.NewFileUpdate(types.FuwName(filename), types.FuwPath(filepath), types.FuwCid("test"))
}

func SetGlobalClient(client *Client) {
	globalClient = client
}

func Initialize(ctx context.Context, con *types.Connector) (*Client, error) {
	client, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}

	SetGlobalClient(client)

	client.SetConnector(con)

	return client, nil
}

func (c *Client) SetConnector(con *types.Connector) {
	c.connector = con
}

func (c *Client) AddFile(filepath string, filename string) (types.FileUpdate, error) {
	ctx := context.Background()

	file, err := os.Open(filepath)
	if err != nil {
		return types.FileUpdate{}, err
	}
	defer file.Close()

	fileNode := files.NewReaderFile(file)

	res, err := c.api.Unixfs().Add(ctx, fileNode)
	if err != nil {
		return types.FileUpdate{}, err
	}

	cid := res.RootCid()

	return types.NewFileUpdate(
		types.FuwName(filename),
		types.FuwPath(filepath),
		types.FuwCid(cid.String()),
		types.FuwStatus(types.StatusComplete),
	), nil
}
