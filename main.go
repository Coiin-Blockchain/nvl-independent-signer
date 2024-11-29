// Copyright 2023 Coiin
// Licensed under the Apache License, Version 2.0 (the "Apache License")
// with the following modification; you may not use this file except in
// compliance with the Apache License and the following modification to it:
// Section 6. Trademarks. is deleted and replaced with:
//      6. Trademarks. This License does not grant permission to use the trade
//         names, trademarks, service marks, or product names of the Licensor
//         and its affiliates, except as required to comply with Section 4(c) of
//         the License and to reproduce the content of the NOTICE file.
// You may obtain a copy of the Apache License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the Apache License with the above modification is
// distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the Apache License for the specific
// language governing permissions and limitations under the Apache License.

package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	signingKeyFilename     = "signing-key"
	priorBlockHashFilename = "prior-block-hash"
)

var (
	Version = "v0.0.0"

	dataDir string

	signingKeyFilePath     string
	priorBlockHashFilePath string

	nvlBaseURL string
)

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.Getwd()
		if err != nil {
			log.Fatal("could not find working directory")
		}
	}
	dataDir = filepath.Join(configDir, "coiin", "nvl", "independent-signer")

	signingKeyFilePath = filepath.Join(dataDir, signingKeyFilename)
	priorBlockHashFilePath = filepath.Join(dataDir, priorBlockHashFilename)

	flag.StringVar(&nvlBaseURL, "nvlBaseURL", "https://nvl.api.coiin.ai", "Host that would be called to sign blocks to")
	flag.Parse()
}

type NVLBlockHeader struct {
	Type        string `json:"type"`
	PriorBlock  string `json:"priorBlock"`
	Timestamp   string `json:"timestamp"`
	PublicKey   string `json:"publicKey"`
	CoiinSupply string `json:"coiinSupply" datastore:"coiinSupply"`
}

type NVLBlockSeal struct {
	Proofs    string `json:"proofs"`
	Signature string `json:"signature,omitempty"`
}

type NVLBlock struct {
	Version string          `json:"version"`
	Header  *NVLBlockHeader `json:"header"`
	Blocks  []string        `json:"blocks"`
	Seal    *NVLBlockSeal   `json:"signature"`

	raw string
}

func (b *NVLBlock) MarshalForSigning() ([]byte, error) {
	blocks := b.Blocks
	if blocks == nil {
		blocks = make([]string, 0)
	}

	data := map[string]interface{}{
		"header": map[string]string{
			"type":       b.Header.Type,
			"priorBlock": b.Header.PriorBlock,
			"timestamp":  b.Header.Timestamp,
			"publicKey":  b.Header.PublicKey,
		},
		"blocks":  blocks,
		"version": b.Version,
	}

	if b.Header.CoiinSupply != "" {
		data["header"].(map[string]string)["coiinSupply"] = b.Header.CoiinSupply
	}

	return json.Marshal(data)
}

func main() {
	log.Printf("Starting NVL independent signer %s\n", Version)

	signingKey, err := loadSigningKey()
	if err != nil {
		log.Fatalf("failed to load signing key: %s", err)
	}

	verifyingKey, err := loadVerifyingKey()
	if err != nil {
		log.Fatalf("failed to load verifying key: %s", err)
	}

	nvlBlock, err := fetchLatestNVLBlock()
	if err != nil {
		log.Fatalf("failed to fetch NVL block: %s", err)
	} else if nvlBlock == nil {
		return
	}

	if valid, err := verifyNVLBlock(verifyingKey, nvlBlock); err != nil {
		log.Fatalf("error verifying NVL block: %s", err)
	} else if !valid {
		log.Fatalf("NVL Proxy block failed validation")
	} else {
		log.Println("NVL Proxy block passed verification")
	}

	priorBlockHash, err := loadPriorBlockHash()
	if err != nil {
		log.Fatalf("failed to load prior block hash %s", err)
	}

	indNVLBlock := createIndependentNVLBlock(signingKey, nvlBlock, priorBlockHash)

	hash, sig, err := signIndependentNVLBlock(signingKey, indNVLBlock)
	if err != nil {
		log.Fatalf("failed to sign independent block: %s", err)
	}
	indNVLBlock.Seal.Proofs = hash
	indNVLBlock.Seal.Signature = sig

	if err := postIndependentNVLBlock(indNVLBlock); err != nil {
		log.Fatalf("failed to post independent block to NVL proxy: %s", err)
	}

	if err := savePriorBlockHash(indNVLBlock.Seal.Proofs); err != nil {
		log.Fatalf("failed to save prior block hash: %s", err)
	}

	log.Println("Complete!")
}

func loadSigningKey() (*ecdsa.PrivateKey, error) {
	log.Println("Loading signing key")
	if _, err := os.Stat(signingKeyFilePath); err != nil {
		log.Println("Signing key not found")
		if err := generateSigningKey(); err != nil {
			return nil, err
		}
	}

	fileData, err := os.ReadFile(signingKeyFilePath)
	if err != nil {
		return nil, err
	}

	signingKey, err := crypto.HexToECDSA(strings.TrimSpace(string(fileData)))

	publicKey := signingKey.Public().(*ecdsa.PublicKey)
	log.Printf("Public Key: %s\n", strings.ToLower(hex.EncodeToString(crypto.FromECDSAPub(publicKey))))

	return signingKey, err
}

func generateSigningKey() error {
	log.Println("Generating signing key")
	// Ensure the data directory exists
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return err
	}

	signingKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(signingKeyFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer mustClose(file)

	_, err = file.WriteString(strings.ToLower(hex.EncodeToString(crypto.FromECDSA(signingKey))))
	if err != nil {
		return err
	}

	log.Println("New signing key generated")
	return nil
}

func loadVerifyingKey() ([]byte, error) {
	log.Println("Loading NVL Proxy verifying key")

	resp, err := http.Get(nvlBaseURL + "/api/v1/status")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("NVL returned non 200 status code: Status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer mustClose(resp.Body)

	status := &struct {
		PublicKey string `json:"publicKey"`
	}{}
	if err := json.Unmarshal(body, status); err != nil {
		return nil, err
	}

	return hexutil.Decode("0x" + status.PublicKey)
}

func fetchLatestNVLBlock() (*NVLBlock, error) {
	log.Println("Fetching latest NVL Proxy block")

	blockHash, err := fetchLatestNVLBlockHash()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(nvlBaseURL + "/api/v1/blocks/" + blockHash + "?raw=true")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("NVL returned non 200 status code: Status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer mustClose(resp.Body)

	block := new(NVLBlock)
	if err := json.Unmarshal(body, block); err != nil {
		return nil, err
	}

	block.raw = string(body)

	log.Printf("Latest NVL Proxy Block hash: %s\n", block.Seal.Proofs)

	return block, nil
}

func fetchLatestNVLBlockHash() (string, error) {
	resp, err := http.Get(nvlBaseURL + "/api/v1/blocks?size=1")
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("NVL returned non 200 status code: Status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer mustClose(resp.Body)

	blocks := &struct {
		Blocks []struct {
			Hash string `json:"hash"`
		}
	}{}
	if err := json.Unmarshal(body, blocks); err != nil {
		return "", err
	}

	if len(blocks.Blocks) != 1 {
		log.Println("NVL did not return any blocks to sign")
		return "", nil
	}

	return blocks.Blocks[0].Hash, nil
}

func verifyNVLBlock(publicKey []byte, block *NVLBlock) (bool, error) {
	sig, err := hexutil.Decode("0x" + block.Seal.Signature)
	if err != nil {
		return false, err
	}
	data, err := block.MarshalForSigning()
	if err != nil {
		return false, err
	}
	hash := crypto.Keccak256Hash(data)
	return crypto.VerifySignature(publicKey, hash.Bytes(), sig[:len(sig)-1]), nil
}

func loadPriorBlockHash() (string, error) {
	log.Println("Loading prior block hash")
	if _, err := os.Stat(priorBlockHashFilePath); err != nil {
		log.Println("No prior block hash, must be first time executed")
		return "", nil
	}

	fileData, err := os.ReadFile(priorBlockHashFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(fileData)), nil
}

func createIndependentNVLBlock(signingKey *ecdsa.PrivateKey, block *NVLBlock, priorHash string) *NVLBlock {
	log.Println("Creating independent NVL block")

	publicKey := strings.ToLower(hex.EncodeToString(crypto.FromECDSAPub(signingKey.Public().(*ecdsa.PublicKey))))
	return &NVLBlock{
		Version: "1",
		Header: &NVLBlockHeader{
			Type:        "INDEPENDENT",
			PriorBlock:  priorHash,
			Timestamp:   fmt.Sprintf("%d", time.Now().Unix()),
			PublicKey:   publicKey,
			CoiinSupply: block.Header.CoiinSupply,
		},
		Blocks: []string{block.Seal.Proofs},
		Seal:   &NVLBlockSeal{},
	}
}

func signIndependentNVLBlock(signingKey *ecdsa.PrivateKey, block *NVLBlock) (string, string, error) {
	data, err := block.MarshalForSigning()
	if err != nil {
		return "", "", err
	}
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), signingKey)
	if err != nil {
		return "", "", err
	}

	hashStr := fmt.Sprintf("%064x", hash.Bytes())
	sigStr := fmt.Sprintf("%0130x", signature)
	log.Println("New independent NVL block signed!")
	log.Printf("Hash: %s\n", hashStr)
	log.Printf("Signature: %s\n", sigStr)

	return hashStr, sigStr, nil
}

func postIndependentNVLBlock(block *NVLBlock) error {
	log.Println("Posting new block to NVL Proxy")

	body := struct {
		Version                  string    `json:"version"`
		Block                    *NVLBlock `json:"block"`
		IndependentSignerVersion string    `json:"independentSignerVersion"`
	}{
		Version:                  "1",
		Block:                    block,
		IndependentSignerVersion: Version,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		nvlBaseURL+"/api/v1/independent/enqueue",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	log.Printf("NVL Proxy resp code: %d\n", resp.StatusCode)

	if resp.StatusCode > 299 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer mustClose(resp.Body)
		log.Println(string(respBody))
	}

	return nil
}

func savePriorBlockHash(hash string) error {
	log.Printf("Saving prior block hash: %s\n", hash)
	return os.WriteFile(priorBlockHashFilePath, []byte(hash), 0600)
}

type Closer interface {
	Close() error
}

func mustClose(f Closer) {
	if err := f.Close(); err != nil {
		log.Fatalf("Close() returned error: %s", err)
	}
}
