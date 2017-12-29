package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/leaguej/fabric-cli/btcd/btcec"
	"github.com/leaguej/fabric-cli/btcutil"
	"github.com/leaguej/fabric-cli/btcutil/base58"
)

func StandardizePublicKey(pubKey []byte) []byte {
	curve := btcec.S256()

	pubKeyObj, err := btcec.ParsePubKey(pubKey, curve)
	if err != nil {
		return nil
	}

	data := pubKeyObj.SerializeUncompressed()

	return data
}

/*
 * PublicKeyToAddress 公钥数据到地址的转换
 */
func PublicKeyToAddress(pubKey []byte) (string, error) {
	data := StandardizePublicKey(pubKey)

	hashed := btcutil.Hash160(data)

	addr := "hc" + base58.CheckEncode(hashed, 0x00)

	return addr, nil
}

/*
**数字签名验证的步骤**

1. 对原始数据进行Hash
    hashed := sha256.Sum256([]byte(msg))
2. 从签名中得到公钥，将它和存放的公钥进行比对
    pk, wasCompressed, err := btcec.RecoverCompact(curve, compactSign, hashed)
3. 从签名构造签名对象
    sig = &btcec.Signature{
        R: new(big.Int).SetBytes(compactSign[1 : 32+1]),
        S: new(big.Int).SetBytes(compactSign[32+1:]),
    }
4. 调用 Verify 来验签
    bVerified = sig.Verify(hashed, pk)
*/
func VerifySignature(pubKey, msg, signData []byte) (bool, error) {
	curve := btcec.S256()

	hashed := sha256.Sum256(msg)

	pk, _, err := btcec.RecoverCompact(curve, signData, hashed[:])
	if err != nil {
		return false, errors.New("Failed to call RecoverCompact " + err.Error())
	}

	pubKeyOrg, err := btcec.ParsePubKey(pubKey, curve)
	if err != nil {
		return false, errors.New("Public key format is not right, " + err.Error() + ", pubKey=" + hex.EncodeToString(pubKey))
	}

	if pk.X.Cmp(pubKeyOrg.X) != 0 || pk.Y.Cmp(pubKeyOrg.Y) != 0 {
		return false, errors.New("Public keys are not equal: 04" +
			hex.EncodeToString(pk.X.Bytes()) + hex.EncodeToString(pk.Y.Bytes()) + " != " + hex.EncodeToString(pubKey))
	}

	sig := &btcec.Signature{
		R: new(big.Int).SetBytes(signData[1 : 32+1]),
		S: new(big.Int).SetBytes(signData[32+1:]),
	}

	var bVerified = sig.Verify(hashed[:], pk)

	return bVerified, nil
}
