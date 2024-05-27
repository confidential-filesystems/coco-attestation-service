module cfs

go 1.21.7

require (
	github.com/confidential-filesystems/filesystem-toolchain v0.0.1
	github.com/ethereum/go-ethereum v1.10.4
	github.com/sigstore/sigstore v1.8.1
)

require (
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce // indirect
	github.com/go-resty/resty/v2 v2.12.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/letsencrypt/boulder v0.0.0-20231026200631-000cd05d5491 // indirect
	github.com/miguelmota/go-ethereum-hdwallet v0.1.1 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.8.0 // indirect
	github.com/spf13/cobra v1.8.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/titanous/rocacheck v0.0.0-20171023193734-afe73141d399 // indirect
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef // indirect
	github.com/zeebo/errs/v2 v2.0.5 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/go-jose/go-jose.v2 v2.6.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	confidentialfilesystems.com/cc/keyprovider => ../../../../../filesystem-toolchain/image/keyprovider
	github.com/confidential-filesystems/filesystem-toolchain => ../../../../../filesystem-toolchain
)
