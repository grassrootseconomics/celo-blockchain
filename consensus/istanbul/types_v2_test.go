package istanbul

import (
	"reflect"
	"testing"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/rlp"
	"gotest.tools/assert"
)

func dummyRoundChangeRequest() *RoundChangeRequest {
	req := RoundChangeRequest{
		View: View{
			Round:    common.Big1,
			Sequence: common.Big2,
		},
		PreparedCertificateV2: PreparedCertificateV2{
			PrepareOrCommitMessages: []Message{},
			ProposalHash:            dummyBlock(2).Hash(),
		},
	}
	// Set empty rather than nil signature since this is how rlp decodes non
	// existent slices.
	req.Signature = []byte{}
	return &req
}

func dummyRoundChangeCertificateV2() *RoundChangeCertificateV2 {
	return &RoundChangeCertificateV2{
		Requests: []RoundChangeRequest{*dummyRoundChangeRequest(), *dummyRoundChangeRequest(), *dummyRoundChangeRequest()},
	}
}

func TestRoundChangeCertificateV2RLPEncoding(t *testing.T) {
	var result, original *RoundChangeCertificateV2
	original = dummyRoundChangeCertificateV2()

	rawVal, err := rlp.EncodeToBytes(original)
	if err != nil {
		t.Fatalf("Error %v", err)
	}

	if err = rlp.DecodeBytes(rawVal, &result); err != nil {
		t.Fatalf("Error %v", err)
	}

	assert.Equal(t, len(original.Requests), len(original.Requests))
	o1 := original.Requests[0]
	r1 := result.Requests[0]

	if !reflect.DeepEqual(o1.Address, r1.Address) {
		t.Fatalf("RLP Encode/Decode mismatch at first Address")
	}

	if !reflect.DeepEqual(o1.Signature, r1.Signature) {
		t.Fatalf("RLP Encode/Decode mismatch at first Signature")
	}

	if !reflect.DeepEqual(o1.View, r1.View) {
		t.Fatalf("RLP Encode/Decode mismatch at first View")
	}

	if !reflect.DeepEqual(o1.PreparedCertificateV2, r1.PreparedCertificateV2) {
		t.Fatalf("RLP Encode/Decode mismatch at first PreparedCertificateV2")
	}

	if !reflect.DeepEqual(original, result) {
		t.Fatalf("RLP Encode/Decode mismatch. Got %v, expected %v", result, original)
	}
}
