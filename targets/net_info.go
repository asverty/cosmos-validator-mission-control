package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func GetNetInfo(ops HTTPOptions, cfg *config.Config) {
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error getting node_info: %v", err)
		return
	}
	var ni NetInfo
	err = json.Unmarshal(resp.Body, &ni)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	numPeers, err := strconv.Atoi(ni.Result.NumPeers)
	if err != nil {
		log.Printf("Error converting num_peers to int: %v", err)
	} else if int64(numPeers) < cfg.NumPeersThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Number of peers has fallen below %d", cfg.NumPeersThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Number of peers has fallen below %d", cfg.NumPeersThreshold), cfg)
	}

	peerAddrs := make([]string, len(ni.Result.Peers))
	for i, peer := range ni.Result.Peers {
		peerAddrs[i] = peer.RemoteIP
	}

	log.Printf("No. of peers: %s \n Peer Addresses: %v", numPeers, peerAddrs)
}
