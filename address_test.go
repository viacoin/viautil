// Copyright (c) 2013, 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcutil_test

import (
	"bytes"
	"code.google.com/p/go.crypto/ripemd160"
	"encoding/hex"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
	"reflect"
	"testing"
)

// invalidNet is an invalid bitcoin network.
const invalidNet = btcwire.BitcoinNet(0xffffffff)

func TestAddresses(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		valid     bool
		canDecode bool
		result    btcutil.Address
		f         func() (btcutil.Address, error)
		net       btcwire.BitcoinNet
	}{
		// Positive P2PKH tests.
		{
			name:      "mainnet p2pkh",
			addr:      "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX",
			valid:     true,
			canDecode: true,
			result: btcutil.TstAddressPubKeyHash(
				[ripemd160.Size]byte{
					0xe3, 0x4c, 0xce, 0x70, 0xc8, 0x63, 0x73, 0x27, 0x3e, 0xfc,
					0xc5, 0x4c, 0xe7, 0xd2, 0xa4, 0x91, 0xbb, 0x4a, 0x0e, 0x84},
				btcutil.MainNetAddr),
			f: func() (btcutil.Address, error) {
				pkHash := []byte{
					0xe3, 0x4c, 0xce, 0x70, 0xc8, 0x63, 0x73, 0x27, 0x3e, 0xfc,
					0xc5, 0x4c, 0xe7, 0xd2, 0xa4, 0x91, 0xbb, 0x4a, 0x0e, 0x84}
				return btcutil.NewAddressPubKeyHash(pkHash, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			name:      "mainnet p2pkh 2",
			addr:      "12MzCDwodF9G1e7jfwLXfR164RNtx4BRVG",
			valid:     true,
			canDecode: true,
			result: btcutil.TstAddressPubKeyHash(
				[ripemd160.Size]byte{
					0x0e, 0xf0, 0x30, 0x10, 0x7f, 0xd2, 0x6e, 0x0b, 0x6b, 0xf4,
					0x05, 0x12, 0xbc, 0xa2, 0xce, 0xb1, 0xdd, 0x80, 0xad, 0xaa},
				btcutil.MainNetAddr),
			f: func() (btcutil.Address, error) {
				pkHash := []byte{
					0x0e, 0xf0, 0x30, 0x10, 0x7f, 0xd2, 0x6e, 0x0b, 0x6b, 0xf4,
					0x05, 0x12, 0xbc, 0xa2, 0xce, 0xb1, 0xdd, 0x80, 0xad, 0xaa}
				return btcutil.NewAddressPubKeyHash(pkHash, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			name:      "testnet p2pkh",
			addr:      "mrX9vMRYLfVy1BnZbc5gZjuyaqH3ZW2ZHz",
			valid:     true,
			canDecode: true,
			result: btcutil.TstAddressPubKeyHash(
				[ripemd160.Size]byte{
					0x78, 0xb3, 0x16, 0xa0, 0x86, 0x47, 0xd5, 0xb7, 0x72, 0x83,
					0xe5, 0x12, 0xd3, 0x60, 0x3f, 0x1f, 0x1c, 0x8d, 0xe6, 0x8f},
				btcutil.TestNetAddr),
			f: func() (btcutil.Address, error) {
				pkHash := []byte{
					0x78, 0xb3, 0x16, 0xa0, 0x86, 0x47, 0xd5, 0xb7, 0x72, 0x83,
					0xe5, 0x12, 0xd3, 0x60, 0x3f, 0x1f, 0x1c, 0x8d, 0xe6, 0x8f}
				return btcutil.NewAddressPubKeyHash(pkHash, btcwire.TestNet3)
			},
			net: btcwire.TestNet3,
		},

		// Negative P2PKH tests.
		{
			name:      "p2pkh wrong byte identifier/net",
			addr:      "MrX9vMRYLfVy1BnZbc5gZjuyaqH3ZW2ZHz",
			valid:     false,
			canDecode: true,
			f: func() (btcutil.Address, error) {
				pkHash := []byte{
					0x78, 0xb3, 0x16, 0xa0, 0x86, 0x47, 0xd5, 0xb7, 0x72, 0x83,
					0xe5, 0x12, 0xd3, 0x60, 0x3f, 0x1f, 0x1c, 0x8d, 0xe6, 0x8f}
				return btcutil.NewAddressPubKeyHash(pkHash, invalidNet)
			},
		},
		{
			name:      "p2pkh wrong hash length",
			addr:      "",
			valid:     false,
			canDecode: true,
			f: func() (btcutil.Address, error) {
				pkHash := []byte{
					0x00, 0x0e, 0xf0, 0x30, 0x10, 0x7f, 0xd2, 0x6e, 0x0b, 0x6b,
					0xf4, 0x05, 0x12, 0xbc, 0xa2, 0xce, 0xb1, 0xdd, 0x80, 0xad,
					0xaa}
				return btcutil.NewAddressPubKeyHash(pkHash, btcwire.MainNet)
			},
		},
		{
			name:      "p2pkh bad checksum",
			addr:      "1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gY",
			valid:     false,
			canDecode: true,
		},

		// Positive P2SH tests.
		{
			// Taken from transactions:
			// output: 3c9018e8d5615c306d72397f8f5eef44308c98fb576a88e030c25456b4f3a7ac
			// input:  837dea37ddc8b1e3ce646f1a656e79bbd8cc7f558ac56a169626d649ebe2a3ba.
			name:      "mainnet p2sh",
			addr:      "3QJmV3qfvL9SuYo34YihAf3sRCW3qSinyC",
			valid:     true,
			canDecode: true,
			result: btcutil.TstAddressScriptHash(
				[ripemd160.Size]byte{
					0xf8, 0x15, 0xb0, 0x36, 0xd9, 0xbb, 0xbc, 0xe5, 0xe9, 0xf2,
					0xa0, 0x0a, 0xbd, 0x1b, 0xf3, 0xdc, 0x91, 0xe9, 0x55, 0x10},
				btcutil.MainNetScriptHash),
			f: func() (btcutil.Address, error) {
				script := []byte{
					0x52, 0x41, 0x04, 0x91, 0xbb, 0xa2, 0x51, 0x09, 0x12, 0xa5,
					0xbd, 0x37, 0xda, 0x1f, 0xb5, 0xb1, 0x67, 0x30, 0x10, 0xe4,
					0x3d, 0x2c, 0x6d, 0x81, 0x2c, 0x51, 0x4e, 0x91, 0xbf, 0xa9,
					0xf2, 0xeb, 0x12, 0x9e, 0x1c, 0x18, 0x33, 0x29, 0xdb, 0x55,
					0xbd, 0x86, 0x8e, 0x20, 0x9a, 0xac, 0x2f, 0xbc, 0x02, 0xcb,
					0x33, 0xd9, 0x8f, 0xe7, 0x4b, 0xf2, 0x3f, 0x0c, 0x23, 0x5d,
					0x61, 0x26, 0xb1, 0xd8, 0x33, 0x4f, 0x86, 0x41, 0x04, 0x86,
					0x5c, 0x40, 0x29, 0x3a, 0x68, 0x0c, 0xb9, 0xc0, 0x20, 0xe7,
					0xb1, 0xe1, 0x06, 0xd8, 0xc1, 0x91, 0x6d, 0x3c, 0xef, 0x99,
					0xaa, 0x43, 0x1a, 0x56, 0xd2, 0x53, 0xe6, 0x92, 0x56, 0xda,
					0xc0, 0x9e, 0xf1, 0x22, 0xb1, 0xa9, 0x86, 0x81, 0x8a, 0x7c,
					0xb6, 0x24, 0x53, 0x2f, 0x06, 0x2c, 0x1d, 0x1f, 0x87, 0x22,
					0x08, 0x48, 0x61, 0xc5, 0xc3, 0x29, 0x1c, 0xcf, 0xfe, 0xf4,
					0xec, 0x68, 0x74, 0x41, 0x04, 0x8d, 0x24, 0x55, 0xd2, 0x40,
					0x3e, 0x08, 0x70, 0x8f, 0xc1, 0xf5, 0x56, 0x00, 0x2f, 0x1b,
					0x6c, 0xd8, 0x3f, 0x99, 0x2d, 0x08, 0x50, 0x97, 0xf9, 0x97,
					0x4a, 0xb0, 0x8a, 0x28, 0x83, 0x8f, 0x07, 0x89, 0x6f, 0xba,
					0xb0, 0x8f, 0x39, 0x49, 0x5e, 0x15, 0xfa, 0x6f, 0xad, 0x6e,
					0xdb, 0xfb, 0x1e, 0x75, 0x4e, 0x35, 0xfa, 0x1c, 0x78, 0x44,
					0xc4, 0x1f, 0x32, 0x2a, 0x18, 0x63, 0xd4, 0x62, 0x13, 0x53,
					0xae}
				return btcutil.NewAddressScriptHash(script, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			// Taken from transactions:
			// output: b0539a45de13b3e0403909b8bd1a555b8cbe45fd4e3f3fda76f3a5f52835c29d
			// input: (not yet redeemed at time test was written)
			name:      "mainnet p2sh 2",
			addr:      "3NukJ6fYZJ5Kk8bPjycAnruZkE5Q7UW7i8",
			valid:     true,
			canDecode: true,
			result: btcutil.TstAddressScriptHash(
				[ripemd160.Size]byte{
					0xe8, 0xc3, 0x00, 0xc8, 0x79, 0x86, 0xef, 0xa8, 0x4c, 0x37,
					0xc0, 0x51, 0x99, 0x29, 0x01, 0x9e, 0xf8, 0x6e, 0xb5, 0xb4},
				btcutil.MainNetScriptHash),
			f: func() (btcutil.Address, error) {
				hash := []byte{
					0xe8, 0xc3, 0x00, 0xc8, 0x79, 0x86, 0xef, 0xa8, 0x4c, 0x37,
					0xc0, 0x51, 0x99, 0x29, 0x01, 0x9e, 0xf8, 0x6e, 0xb5, 0xb4}
				return btcutil.NewAddressScriptHashFromHash(hash, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			// Taken from bitcoind base58_keys_valid.
			name:      "testnet p2sh",
			addr:      "2NBFNJTktNa7GZusGbDbGKRZTxdK9VVez3n",
			valid:     true,
			canDecode: true,
			result: btcutil.TstAddressScriptHash(
				[ripemd160.Size]byte{
					0xc5, 0x79, 0x34, 0x2c, 0x2c, 0x4c, 0x92, 0x20, 0x20, 0x5e,
					0x2c, 0xdc, 0x28, 0x56, 0x17, 0x04, 0x0c, 0x92, 0x4a, 0x0a},
				btcutil.TestNetScriptHash),
			f: func() (btcutil.Address, error) {
				hash := []byte{
					0xc5, 0x79, 0x34, 0x2c, 0x2c, 0x4c, 0x92, 0x20, 0x20, 0x5e,
					0x2c, 0xdc, 0x28, 0x56, 0x17, 0x04, 0x0c, 0x92, 0x4a, 0x0a}
				return btcutil.NewAddressScriptHashFromHash(hash, btcwire.TestNet3)
			},
			net: btcwire.TestNet3,
		},

		// Negative P2SH tests.
		{
			name:      "p2sh wrong hash length",
			addr:      "",
			valid:     false,
			canDecode: true,
			f: func() (btcutil.Address, error) {
				hash := []byte{
					0x00, 0xf8, 0x15, 0xb0, 0x36, 0xd9, 0xbb, 0xbc, 0xe5, 0xe9,
					0xf2, 0xa0, 0x0a, 0xbd, 0x1b, 0xf3, 0xdc, 0x91, 0xe9, 0x55,
					0x10}
				return btcutil.NewAddressScriptHashFromHash(hash, btcwire.MainNet)
			},
		},
		{
			name:      "p2sh wrong byte identifier/net",
			addr:      "0NBFNJTktNa7GZusGbDbGKRZTxdK9VVez3n",
			valid:     false,
			canDecode: true,
			f: func() (btcutil.Address, error) {
				hash := []byte{
					0xc5, 0x79, 0x34, 0x2c, 0x2c, 0x4c, 0x92, 0x20, 0x20, 0x5e,
					0x2c, 0xdc, 0x28, 0x56, 0x17, 0x04, 0x0c, 0x92, 0x4a, 0x0a}
				return btcutil.NewAddressScriptHashFromHash(hash, invalidNet)
			},
		},

		// Positive P2PK tests.
		{
			name:      "mainnet p2pk compressed (0x02)",
			addr:      "13CG6SJ3yHUXo4Cr2RY4THLLJrNFuG3gUg",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x02, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4},
				btcutil.PKFCompressed, btcutil.MainNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x02, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			name:      "mainnet p2pk compressed (0x03)",
			addr:      "15sHANNUBSh6nDp8XkDPmQcW6n3EFwmvE6",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x03, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65},
				btcutil.PKFCompressed, btcutil.MainNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x03, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			name:      "mainnet p2pk uncompressed (0x04)",
			addr:      "12cbQLTFMXRnSzktFkuoG3eHoMeFtpTu3S",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a, 0x01, 0x6b,
					0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e, 0xb6, 0x8a, 0x38,
					0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca, 0xd7, 0xb1, 0x48, 0xa6,
					0x90, 0x9a, 0x5c, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
					0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b,
					0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43,
					0xf6, 0x56, 0xb4, 0x12, 0xa3},
				btcutil.PKFUncompressed, btcutil.MainNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a, 0x01, 0x6b,
					0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e, 0xb6, 0x8a, 0x38,
					0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca, 0xd7, 0xb1, 0x48, 0xa6,
					0x90, 0x9a, 0x5c, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
					0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b,
					0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43,
					0xf6, 0x56, 0xb4, 0x12, 0xa3}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			name:      "mainnet p2pk hybrid (0x06)",
			addr:      "1Ja5rs7XBZnK88EuLVcFqYGMEbBitzchmX",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x06, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4, 0x0d, 0x45, 0x26, 0x48, 0x38, 0xc0, 0xbd,
					0x96, 0x85, 0x26, 0x62, 0xce, 0x6a, 0x84, 0x7b, 0x19, 0x73,
					0x76, 0x83, 0x01, 0x60, 0xc6, 0xd2, 0xeb, 0x5e, 0x6a, 0x4c,
					0x44, 0xd3, 0x3f, 0x45, 0x3e},
				btcutil.PKFHybrid, btcutil.MainNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x06, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4, 0x0d, 0x45, 0x26, 0x48, 0x38, 0xc0, 0xbd,
					0x96, 0x85, 0x26, 0x62, 0xce, 0x6a, 0x84, 0x7b, 0x19, 0x73,
					0x76, 0x83, 0x01, 0x60, 0xc6, 0xd2, 0xeb, 0x5e, 0x6a, 0x4c,
					0x44, 0xd3, 0x3f, 0x45, 0x3e}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			name:      "mainnet p2pk hybrid (0x07)",
			addr:      "1ExqMmf6yMxcBMzHjbj41wbqYuqoX6uBLG",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x07, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65, 0x37, 0xa5, 0x76, 0x78, 0x2e, 0xba, 0x66,
					0x8a, 0x7e, 0xf8, 0xbd, 0x3b, 0x3c, 0xfb, 0x1e, 0xdb, 0x71,
					0x17, 0xab, 0x65, 0x12, 0x9b, 0x8a, 0x2e, 0x68, 0x1f, 0x3c,
					0x1e, 0x09, 0x08, 0xef, 0x7b},
				btcutil.PKFHybrid, btcutil.MainNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x07, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65, 0x37, 0xa5, 0x76, 0x78, 0x2e, 0xba, 0x66,
					0x8a, 0x7e, 0xf8, 0xbd, 0x3b, 0x3c, 0xfb, 0x1e, 0xdb, 0x71,
					0x17, 0xab, 0x65, 0x12, 0x9b, 0x8a, 0x2e, 0x68, 0x1f, 0x3c,
					0x1e, 0x09, 0x08, 0xef, 0x7b}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.MainNet)
			},
			net: btcwire.MainNet,
		},
		{
			name:      "testnet p2pk compressed (0x02)",
			addr:      "mhiDPVP2nJunaAgTjzWSHCYfAqxxrxzjmo",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x02, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4},
				btcutil.PKFCompressed, btcutil.TestNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x02, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.TestNet3)
			},
			net: btcwire.TestNet3,
		},
		{
			name:      "testnet p2pk compressed (0x03)",
			addr:      "mkPETRTSzU8MZLHkFKBmbKppxmdw9qT42t",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x03, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65},
				btcutil.PKFCompressed, btcutil.TestNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x03, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.TestNet3)
			},
			net: btcwire.TestNet3,
		},
		{
			name:      "testnet p2pk uncompressed (0x04)",
			addr:      "mh8YhPYEAYs3E7EVyKtB5xrcfMExkkdEMF",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a, 0x01, 0x6b,
					0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e, 0xb6, 0x8a, 0x38,
					0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca, 0xd7, 0xb1, 0x48, 0xa6,
					0x90, 0x9a, 0x5c, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
					0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b,
					0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43,
					0xf6, 0x56, 0xb4, 0x12, 0xa3},
				btcutil.PKFUncompressed, btcutil.TestNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a, 0x01, 0x6b,
					0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e, 0xb6, 0x8a, 0x38,
					0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca, 0xd7, 0xb1, 0x48, 0xa6,
					0x90, 0x9a, 0x5c, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
					0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b,
					0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43,
					0xf6, 0x56, 0xb4, 0x12, 0xa3}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.TestNet3)
			},
			net: btcwire.TestNet3,
		},
		{
			name:      "testnet p2pk hybrid (0x06)",
			addr:      "my639vCVzbDZuEiX44adfTUg6anRomZLEP",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x06, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4, 0x0d, 0x45, 0x26, 0x48, 0x38, 0xc0, 0xbd,
					0x96, 0x85, 0x26, 0x62, 0xce, 0x6a, 0x84, 0x7b, 0x19, 0x73,
					0x76, 0x83, 0x01, 0x60, 0xc6, 0xd2, 0xeb, 0x5e, 0x6a, 0x4c,
					0x44, 0xd3, 0x3f, 0x45, 0x3e},
				btcutil.PKFHybrid, btcutil.TestNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x06, 0x19, 0x2d, 0x74, 0xd0, 0xcb, 0x94, 0x34, 0x4c, 0x95,
					0x69, 0xc2, 0xe7, 0x79, 0x01, 0x57, 0x3d, 0x8d, 0x79, 0x03,
					0xc3, 0xeb, 0xec, 0x3a, 0x95, 0x77, 0x24, 0x89, 0x5d, 0xca,
					0x52, 0xc6, 0xb4, 0x0d, 0x45, 0x26, 0x48, 0x38, 0xc0, 0xbd,
					0x96, 0x85, 0x26, 0x62, 0xce, 0x6a, 0x84, 0x7b, 0x19, 0x73,
					0x76, 0x83, 0x01, 0x60, 0xc6, 0xd2, 0xeb, 0x5e, 0x6a, 0x4c,
					0x44, 0xd3, 0x3f, 0x45, 0x3e}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.TestNet3)
			},
			net: btcwire.TestNet3,
		},
		{
			name:      "testnet p2pk hybrid (0x07)",
			addr:      "muUnepk5nPPrxUTuTAhRqrpAQuSWS5fVii",
			valid:     true,
			canDecode: false,
			result: btcutil.TstAddressPubKey(
				[]byte{
					0x07, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65, 0x37, 0xa5, 0x76, 0x78, 0x2e, 0xba, 0x66,
					0x8a, 0x7e, 0xf8, 0xbd, 0x3b, 0x3c, 0xfb, 0x1e, 0xdb, 0x71,
					0x17, 0xab, 0x65, 0x12, 0x9b, 0x8a, 0x2e, 0x68, 0x1f, 0x3c,
					0x1e, 0x09, 0x08, 0xef, 0x7b},
				btcutil.PKFHybrid, btcutil.TestNetAddr),
			f: func() (btcutil.Address, error) {
				serializedPubKey := []byte{
					0x07, 0xb0, 0xbd, 0x63, 0x42, 0x34, 0xab, 0xbb, 0x1b, 0xa1,
					0xe9, 0x86, 0xe8, 0x84, 0x18, 0x5c, 0x61, 0xcf, 0x43, 0xe0,
					0x01, 0xf9, 0x13, 0x7f, 0x23, 0xc2, 0xc4, 0x09, 0x27, 0x3e,
					0xb1, 0x6e, 0x65, 0x37, 0xa5, 0x76, 0x78, 0x2e, 0xba, 0x66,
					0x8a, 0x7e, 0xf8, 0xbd, 0x3b, 0x3c, 0xfb, 0x1e, 0xdb, 0x71,
					0x17, 0xab, 0x65, 0x12, 0x9b, 0x8a, 0x2e, 0x68, 0x1f, 0x3c,
					0x1e, 0x09, 0x08, 0xef, 0x7b}
				return btcutil.NewAddressPubKey(serializedPubKey, btcwire.TestNet3)
			},
			net: btcwire.TestNet3,
		},
	}

	for _, test := range tests {
		var decoded btcutil.Address
		var err error
		if test.canDecode {
			// Decode addr and compare error against valid.
			decoded, err = btcutil.DecodeAddr(test.addr)
			if (err == nil) != test.valid {
				t.Errorf("%v: decoding test failed", test.name)
				return
			}
		} else {
			// The address can't be decoded directly, so instead
			// call the creation function.
			decoded, err = test.f()
			if (err == nil) != test.valid {
				t.Errorf("%v: creation test failed", test.name)
				return
			}
		}

		// If decoding succeeded, encode again and compare against the original.
		if err == nil {
			encoded := decoded.EncodeAddress()

			// Compare encoded addr against the original encoding.
			if test.addr != encoded {
				t.Errorf("%v: decoding and encoding produced different addressess: %v != %v",
					test.name, test.addr, encoded)
				return
			}

			// Perform type-specific calculations.
			var saddr []byte
			switch d := decoded.(type) {
			case *btcutil.AddressPubKeyHash:
				saddr = btcutil.TstAddressSAddr(encoded)

			case *btcutil.AddressScriptHash:
				saddr = btcutil.TstAddressSAddr(encoded)

			case *btcutil.AddressPubKey:
				// Ignore the error here since the script
				// address is checked below.
				saddr, _ = hex.DecodeString(d.String())
			}

			// Check script address.
			if !bytes.Equal(saddr, decoded.ScriptAddress()) {
				t.Errorf("%v: script addresses do not match:\n%x != \n%x",
					test.name, saddr, decoded.ScriptAddress())
				return
			}

			// Check networks.  This check always succeeds for non-P2PKH and
			// non-P2SH addresses as both nets will be Go's default zero value.
			if !decoded.IsForNet(test.net) {
				t.Errorf("%v: calculated network does not match expected",
					test.name)
				return
			}
		}

		if !test.valid {
			// If address is invalid, but a creation function exists,
			// verify that it returns a nil addr and non-nil error.
			if test.f != nil {
				_, err := test.f()
				if err == nil {
					t.Errorf("%v: address is invalid but creating new address succeeded",
						test.name)
					return
				}
			}
			continue
		}

		// Valid test, compare address created with f against expected result.
		addr, err := test.f()
		if err != nil {
			t.Errorf("%v: address is valid but creating new address failed with error %v",
				test.name, err)
			return
		}

		if !reflect.DeepEqual(addr, test.result) {
			t.Errorf("%v: created address does not match expected result",
				test.name)
			return
		}
	}
}

func TestEncodeDecodePrivateKey(t *testing.T) {
	tests := []struct {
		in         []byte
		net        btcwire.BitcoinNet
		compressed bool
		out        string
	}{
		{[]byte{
			0x0c, 0x28, 0xfc, 0xa3, 0x86, 0xc7, 0xa2, 0x27,
			0x60, 0x0b, 0x2f, 0xe5, 0x0b, 0x7c, 0xae, 0x11,
			0xec, 0x86, 0xd3, 0xbf, 0x1f, 0xbe, 0x47, 0x1b,
			0xe8, 0x98, 0x27, 0xe1, 0x9d, 0x72, 0xaa, 0x1d,
		}, btcwire.MainNet, false, "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"},
		{[]byte{
			0xdd, 0xa3, 0x5a, 0x14, 0x88, 0xfb, 0x97, 0xb6,
			0xeb, 0x3f, 0xe6, 0xe9, 0xef, 0x2a, 0x25, 0x81,
			0x4e, 0x39, 0x6f, 0xb5, 0xdc, 0x29, 0x5f, 0xe9,
			0x94, 0xb9, 0x67, 0x89, 0xb2, 0x1a, 0x03, 0x98,
		}, btcwire.TestNet3, true, "cV1Y7ARUr9Yx7BR55nTdnR7ZXNJphZtCCMBTEZBJe1hXt2kB684q"},
	}

	for x, test := range tests {
		wif, err := btcutil.EncodePrivateKey(test.in, test.net, test.compressed)
		if err != nil {
			t.Errorf("%x: %v", x, err)
			continue
		}
		if wif != test.out {
			t.Errorf("TestEncodeDecodePrivateKey failed: want '%s', got '%s'",
				test.out, wif)
			continue
		}

		key, _, compressed, err := btcutil.DecodePrivateKey(test.out)
		if err != nil {
			t.Error(err)
			continue
		}
		if !bytes.Equal(key, test.in) || compressed != test.compressed {
			t.Errorf("TestEncodeDecodePrivateKey failed: want '%x', got '%x'",
				test.out, key)
		}

	}
}