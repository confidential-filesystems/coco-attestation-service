package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/confidential-filesystems/filesystem-toolchain/cert"
	"github.com/confidential-filesystems/filesystem-toolchain/cmd/common"
	"github.com/confidential-filesystems/filesystem-toolchain/resource"
	"github.com/confidential-filesystems/filesystem-toolchain/wallet"
	eCommon "github.com/ethereum/go-ethereum/common"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
)

const (
	defaultRepoDir = "/opt/confidential-containers/kbs/repository"
)

type File struct {
	Resource
}

func NewFile(config Config) (Resource, error) {
	file := &File{}
	return file, nil
}

func (f *File) SetResource(addr, typ, tag string, data []byte) error {
	fmt.Printf("confilesystem-go - File.SetResource(): addr = %v, typ = %v, tag = %v\n",
		addr, typ, tag)

	return nil
}

func (f *File) GetResource(addr, typ, tag string) ([]byte, error) {
	fmt.Printf("confilesystem-go - File.GetResource(): addr = %v, typ = %v, tag = %v\n",
		addr, typ, tag)

	return nil, nil
}

// utils api
func toGetResource(repoDir, addr, typ, tag string) ([]byte, error) {
	if eCommon.IsHexAddress(addr) {
		seeds, err := os.ReadFile(path.Join(repoDir, fmt.Sprintf(resource.ResSeeds, addr)))
		if err != nil {
			return nil, err
		}
		kl, err := resource.NewKeyLoad(string(seeds))
		if err != nil {
			return nil, err
		}
		return generateResource(kl, repoDir, addr, typ, tag)
	}
	data, err := os.ReadFile(path.Join(repoDir, addr, typ, tag))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func generateResource(kl *resource.KeyLoad, repoDir, addr, typ, tag string) ([]byte, error) {
	switch typ {
	case resource.ResTypeEcsk:
		ecki, err := resource.Str2Uint32(tag)
		if err != nil {
			return nil, err
		}
		seed, _, err := wallet.NewECKEY(kl.KeySeeds.ECSEED, ecki, true)
		if err != nil {
			return nil, err
		}
		return []byte(seed), nil
	case resource.ResTypeEcpk:
		ecki, err := resource.Str2Uint32(tag)
		if err != nil {
			return nil, err
		}
		_, priv, err := wallet.NewECKEY(kl.KeySeeds.ECSEED, ecki, true)
		if err != nil {
			return nil, err
		}
		pubPem, err := cryptoutils.MarshalPublicKeyToPEM(&priv.PublicKey)
		if err != nil {
			return nil, err
		}
		return pubPem, nil
	case resource.ResTypeIpk:
		iski, err := resource.Str2Uint32(tag)
		if err != nil {
			return nil, err
		}
		seed, _, err := wallet.NewIPK(kl.KeySeeds.ISEED, iski, true)
		if err != nil {
			return nil, err
		}
		return []byte(seed), nil
	case resource.ResTypeIvp:
		iski, imageRef, err := resource.ParseIskiAndImageReference(tag)
		if err != nil {
			return nil, err
		}
		_, pub, err := wallet.NewIPK(kl.KeySeeds.ISEED, iski, true)
		if err != nil {
			return nil, err
		}
		pubPem, err := cryptoutils.MarshalPublicKeyToPEM(pub)
		if err != nil {
			return nil, err
		}
		p := resource.NewDefaultPolicy(imageRef, string(pubPem))
		return json.MarshalIndent(p, "", "  ")
	case resource.ResTypeIkek:
		iski, err := resource.Str2Uint32(tag)
		if err != nil {
			return nil, err
		}
		seed, err := wallet.NewIKEK(kl.KeySeeds.ISEED, iski, true)
		if err != nil {
			return nil, err
		}
		return seed, nil
	case resource.ResTypeCerts:
		_, caPriv, err := wallet.NewCAKEY(kl.KeySeeds.CASEED, true)
		if err != nil {
			return nil, err
		}
		caPath := path.Join(repoDir, fmt.Sprintf(resource.ResCA, addr))
		var caPem string
		if !common.FileExists(caPath) {
			// write ca pem
			caPem, _, err = cert.CreateCaCertificate(caPriv, nil, nil)
			if err != nil {
				return nil, err
			}
			// create folder if need
			if err := os.MkdirAll(path.Join(repoDir, addr, "ca"), os.FileMode(0755)); err != nil {
				return nil, err
			}
			if err := os.WriteFile(caPath, []byte(caPem), os.FileMode(0644)); err != nil {
				return nil, err
			}
		} else {
			caPemBytes, err := os.ReadFile(caPath)
			if err != nil {
				return nil, err
			}
			caPem = string(caPemBytes)
		}
		switch tag {
		case resource.ResCertsClient:
			var certs resource.ClientCerts
			certs.Cert, certs.Key, err = cert.CreateClientCertificate(caPriv, caPem, nil, nil)
			if err != nil {
				return nil, err
			}
			certs.CA = caPem
			return json.Marshal(&certs)
		case resource.ResCertsServer:
			var certs resource.ServerCerts
			certs.Cert, certs.Key, err = cert.CreateServerCertificate(caPriv, caPem, nil, nil)
			if err != nil {
				return nil, err
			}
			certs.CA = caPem
			return json.Marshal(&certs)
		}
	case resource.ResTypeFsrk:
		tokenId, err := resource.Str2Uint32(tag)
		if err != nil {
			return nil, err
		}
		seed, _, err := wallet.NewFSRK(kl.KeySeeds.FSSEED, tokenId, true)
		if err != nil {
			return nil, err
		}
		return []byte(seed), nil
	case resource.ResTypeAssk:
		aski, err := resource.Str2Uint32(tag)
		if err != nil {
			return nil, err
		}
		_, assk, err := wallet.NewASSK(kl.KeySeeds.ASSEED, aski, true)
		if err != nil {
			return nil, err
		}
		return assk, nil
	}
	return nil, errors.New("invalid resource type")
}
