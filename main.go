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
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	signingKeyFilename     = "signing-key"
	registrationIdFilename = "registration-id"

	nvlBaseUrl = "https://nvl.api.coiin.io"
)

var (
	dataDir string

	signingKeyFilePath     string
	registrationIdFilePath string
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
	registrationIdFilePath = filepath.Join(dataDir, registrationIdFilename)
}

type NVLBlockHeader struct {
	Type       string `json:"type"`
	PriorBlock string `json:"priorBlock"`
	Timestamp  string `json:"timestamp"`
	PublicKey  string `json:"publicKey"`
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

	return json.Marshal(data)
}

func main() {
	log.Println("NVL independent signer starting")

	signingKey, err := loadSigningKey()
	if err != nil {
		log.Fatalf("failed to load signing key: %s", err)
	}

	regId, err := loadRegistrationId()
	if err != nil {
		log.Fatalf("failed to load registration id: %s", err)
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
		log.Fatalf("NVL block failed validation")
	} else {
		log.Println("NVL passed verification")
	}

	identNVLBlock := createIndependentNVLBlock(signingKey, nvlBlock)

	hash, sig, err := signIndependentNVLBlock(signingKey, identNVLBlock)
	if err != nil {
		log.Fatalf("failed to sign independent block: %s", err)
	}
	identNVLBlock.Seal.Proofs = hash
	identNVLBlock.Seal.Signature = sig

	if err := postIndependentNVLBlock(regId, identNVLBlock); err != nil {
		log.Fatalf("failed to post idependent block to NVL proxy: %s", err)
	}
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

	privateKey, err := crypto.HexToECDSA(strings.TrimSpace(string(fileData)))
	return privateKey, err
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

	publicKey := signingKey.Public().(*ecdsa.PublicKey)
	log.Printf("New public Key: %s\n", strings.ToLower(hex.EncodeToString(crypto.FromECDSAPub(publicKey))))

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

func loadRegistrationId() (string, error) {
	log.Println("Loading registration ID")
	if _, err := os.Stat(registrationIdFilePath); err != nil {
		log.Println("Registration ID not found")
		if err := promptForRegistrationId(); err != nil {
			return "", err
		}
	}

	fileData, err := os.ReadFile(registrationIdFilePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(fileData)), nil
}

func promptForRegistrationId() error {
	regId := ""

	fmt.Print("Please enter the registration ID: ")
	if _, err := fmt.Scan(&regId); err != nil {
		return err
	}

	file, err := os.OpenFile(registrationIdFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer mustClose(file)

	_, err = file.WriteString(regId)
	if err != nil {
		return err
	}

	fmt.Print("Registration ID saved!")
	return nil
}

func loadVerifyingKey() ([]byte, error) {
	resp, err := http.Get(nvlBaseUrl + "/api/v1/status")
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
	blockHash, err := fetchLatestNVLBlockHash()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(nvlBaseUrl + "/api/v1/blocks/" + blockHash)
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

	return block, nil
}

func fetchLatestNVLBlockHash() (string, error) {
	resp, err := http.Get(nvlBaseUrl + "/api/v1/blocks?size=1")
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

func createIndependentNVLBlock(signingKey *ecdsa.PrivateKey, block *NVLBlock) *NVLBlock {
	publicKey := strings.ToLower(hex.EncodeToString(crypto.FromECDSAPub(signingKey.Public().(*ecdsa.PublicKey))))

	return &NVLBlock{
		Version: "1",
		Header: &NVLBlockHeader{
			Type:       "INDEPENDENT",
			PriorBlock: block.Header.PriorBlock,
			Timestamp:  block.Header.Timestamp,
			PublicKey:  publicKey,
		},
		Blocks: []string{block.raw},
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
	return fmt.Sprintf("%064x", hash.Bytes()), fmt.Sprintf("%0130x", signature), err
}

func postIndependentNVLBlock(regId string, block *NVLBlock) error {
	return nil
}

type Closer interface {
	Close() error
}

func mustClose(f Closer) {
	if err := f.Close(); err != nil {
		log.Fatalf("Close() returned error: %s", err)
	}
}
