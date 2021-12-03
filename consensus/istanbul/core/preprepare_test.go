// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/celo-org/celo-blockchain/consensus/istanbul"
	"github.com/stretchr/testify/require"
)

func newTestPreprepare(v *istanbul.View) *istanbul.Preprepare {
	return &istanbul.Preprepare{
		View:     v,
		Proposal: newTestProposal(),
	}
}

func TestHandlePreprepare(t *testing.T) {
	N := uint64(4) // replica 0 is the proposer, it will send messages to others
	F := uint64(1) // F does not affect tests

	getRoundState := func(c *core) *roundStateImpl {
		return c.current.(*rsSaveDecorator).rs.(*roundStateImpl)
	}

	testCases := []struct {
		name            string
		system          func() *testSystem
		getCert         func(*testSystem) istanbul.RoundChangeCertificate
		expectedRequest istanbul.Proposal
		expectedErr     error
	}{
		{
			"normal case",
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for _, backend := range sys.backends {
					backend.engine.(*core).Start()
				}
				return sys
			},
			func(_ *testSystem) istanbul.RoundChangeCertificate {
				return istanbul.RoundChangeCertificate{}
			},
			newTestProposal(),
			nil,
		},
		{
			"ROUND CHANGE certificate missing",
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for _, backend := range sys.backends {
					backend.engine.(*core).Start()
					c := backend.engine.(*core)
					getRoundState(c).state = StatePreprepared
					getRoundState(c).round = big.NewInt(int64(N))
					getRoundState(c).desiredRound = getRoundState(c).round
				}
				return sys
			},
			func(_ *testSystem) istanbul.RoundChangeCertificate {
				return istanbul.RoundChangeCertificate{}
			},
			makeBlock(1),
			errMissingRoundChangeCertificate,
		},
		{
			"ROUND CHANGE certificate invalid, duplicate messages.",
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for _, backend := range sys.backends {
					backend.engine.(*core).Start()
					c := backend.engine.(*core)
					getRoundState(c).state = StatePreprepared
					getRoundState(c).round = big.NewInt(int64(N))
					getRoundState(c).desiredRound = getRoundState(c).round
				}
				return sys
			},
			func(sys *testSystem) istanbul.RoundChangeCertificate {
				// Duplicate messages
				roundChangeCertificate := sys.getRoundChangeCertificate(t, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, istanbul.EmptyPreparedCertificate())
				roundChangeCertificate.RoundChangeMessages[1] = roundChangeCertificate.RoundChangeMessages[0]
				return roundChangeCertificate
			},
			makeBlock(1),
			errInvalidRoundChangeCertificateDuplicate,
		},
		{
			"ROUND CHANGE certificate contains PREPARED certificate with inconsistent views among the cert's messages",
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for _, backend := range sys.backends {
					backend.engine.(*core).Start()
					c := backend.engine.(*core)
					getRoundState(c).state = StatePreprepared
					getRoundState(c).round = big.NewInt(int64(N))
					getRoundState(c).desiredRound = getRoundState(c).round
					getRoundState(c).sequence = big.NewInt(1)
					c.current.TransitionToPreprepared(&istanbul.Preprepare{
						View: &istanbul.View{
							Round:    big.NewInt(int64(N)),
							Sequence: big.NewInt(1),
						},
						Proposal:               makeBlock(1),
						RoundChangeCertificate: istanbul.RoundChangeCertificate{},
					})
				}
				return sys
			},
			func(sys *testSystem) istanbul.RoundChangeCertificate {
				view1 := *(sys.backends[0].engine.(*core).current.View())

				var view2 istanbul.View
				view2.Sequence = big.NewInt(view1.Sequence.Int64())
				view2.Round = big.NewInt(view1.Round.Int64() + 1)

				preparedCertificate := sys.getPreparedCertificate(t, []istanbul.View{view1, view2}, makeBlock(1))
				roundChangeCertificate := sys.getRoundChangeCertificate(t, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, preparedCertificate)
				return roundChangeCertificate
			},
			makeBlock(1),
			errInvalidPreparedCertificateInconsistentViews,
		},
		{
			"ROUND CHANGE certificate contains PREPARED certificate for a different block.",
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for _, backend := range sys.backends {
					backend.engine.(*core).Start()
					c := backend.engine.(*core)
					getRoundState(c).state = StatePreprepared
					getRoundState(c).round = big.NewInt(int64(N))
					getRoundState(c).desiredRound = getRoundState(c).round
					c.current.TransitionToPreprepared(&istanbul.Preprepare{
						View: &istanbul.View{
							Round:    big.NewInt(int64(N)),
							Sequence: big.NewInt(0),
						},
						Proposal:               makeBlock(2),
						RoundChangeCertificate: istanbul.RoundChangeCertificate{},
					})
				}
				return sys
			},
			func(sys *testSystem) istanbul.RoundChangeCertificate {
				preparedCertificate := sys.getPreparedCertificate(t, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, makeBlock(2))
				roundChangeCertificate := sys.getRoundChangeCertificate(t, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, preparedCertificate)
				return roundChangeCertificate
			},
			makeBlock(1),
			errInvalidPreparedCertificateDigestMismatch,
		},
		{
			"ROUND CHANGE certificate for N+1 round with valid PREPARED certificates",
			// Round is N+1 to match the correct proposer.
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					backend.engine.(*core).Start()
					c := backend.engine.(*core)
					getRoundState(c).round = big.NewInt(int64(N))
					getRoundState(c).desiredRound = getRoundState(c).round
					if i != 0 {
						getRoundState(c).state = StateAcceptRequest
					}
				}
				return sys
			},
			func(sys *testSystem) istanbul.RoundChangeCertificate {
				preparedCertificate := sys.getPreparedCertificate(t, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, makeBlock(1))
				roundChangeCertificate := sys.getRoundChangeCertificate(t, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, preparedCertificate)
				return roundChangeCertificate
			},
			makeBlock(1),
			nil,
		},
		{
			"ROUND CHANGE certificate for N+1 round with empty PREPARED certificates",
			// Round is N+1 to match the correct proposer.
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					backend.engine.(*core).Start()
					c := backend.engine.(*core)
					getRoundState(c).round = big.NewInt(int64(N))
					getRoundState(c).desiredRound = getRoundState(c).round
					if i != 0 {
						getRoundState(c).state = StateAcceptRequest
					}
				}
				return sys
			},
			func(sys *testSystem) istanbul.RoundChangeCertificate {
				roundChangeCertificate := sys.getRoundChangeCertificate(t, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, istanbul.EmptyPreparedCertificate())
				return roundChangeCertificate
			},
			makeBlock(1),
			nil,
		},
		{
			"ROUND CHANGE certificate for N+1 or later rounds with empty PREPARED certificates",
			// Round is N+1 to match the correct proposer.
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					backend.engine.(*core).Start()
					c := backend.engine.(*core)
					getRoundState(c).round = big.NewInt(int64(N))
					getRoundState(c).desiredRound = getRoundState(c).round
					if i != 0 {
						getRoundState(c).state = StateAcceptRequest
					}
				}
				return sys
			},
			func(sys *testSystem) istanbul.RoundChangeCertificate {
				v1 := *(sys.backends[0].engine.(*core).current.View())
				v2 := istanbul.View{Sequence: v1.Sequence, Round: big.NewInt(v1.Round.Int64() + 1)}
				v3 := istanbul.View{Sequence: v1.Sequence, Round: big.NewInt(v1.Round.Int64() + 2)}
				roundChangeCertificate := sys.getRoundChangeCertificate(t, []istanbul.View{v1, v2, v3}, istanbul.EmptyPreparedCertificate())
				return roundChangeCertificate
			},
			makeBlock(1),
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			t.Log("Running", "test", test.name)
			sys := test.system()

			// Even those we pass core=false here, cores do get started and need stopping.
			sys.Run(false)
			defer sys.Stop(true)

			v0 := sys.backends[0]
			r0 := v0.engine.(*core)

			curView := r0.current.View()

			preprepareView := curView

			msg := istanbul.NewPreprepareMessage(
				&istanbul.Preprepare{
					View:                   preprepareView,
					Proposal:               test.expectedRequest,
					RoundChangeCertificate: test.getCert(sys),
				},
				v0.Address(),
			)
			err := msg.Sign(v0.Sign)
			require.NoError(t, err)
			payload, err := msg.Payload()
			require.NoError(t, err)

			for i, v := range sys.backends {
				// i == 0 is primary backend, it is responsible for send PRE-PREPARE messages to others.
				if i == 0 {
					continue
				}

				c := v.engine.(*core)

				// run each backends and verify handlePreprepare function.
				if err := c.handleMsg(payload); err != nil {
					if !errors.Is(err, test.expectedErr) {
						t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
					}
					return
				}

				if c.current.State() != StatePreprepared {
					t.Errorf("state mismatch: have %v, want %v", c.current.State(), StatePreprepared)
				}

				// verify prepare messages
				decodedMsg := new(istanbul.Message)
				err = decodedMsg.FromPayload(v.sentMsgs[0], nil)
				if err != nil {
					t.Errorf("error mismatch: have %v, want nil", err)
				}

				if decodedMsg.Code != istanbul.MsgPrepare {
					t.Errorf("message code mismatch: have %v, want %v", decodedMsg.Code, istanbul.MsgPrepare)
				}

				subject := decodedMsg.Prepare()
				if decodedMsg.Code == istanbul.MsgCommit {
					subject = decodedMsg.Commit().Subject
				}

				expectedSubject := c.current.Subject()
				if !reflect.DeepEqual(subject, expectedSubject) {
					t.Errorf("subject mismatch: have %v, want %v", subject, expectedSubject)
				}

			}
		})
	}
}

// benchMarkHandleRoundChange benchmarks handling a round change messages with n prepare or commit messages in the prepared certificate
func benchMarkHandleRoundChange(n int, b *testing.B) {
	// Setup
	N := uint64(n)
	F := uint64(1) // F does not affect tests
	sys := NewMutedTestSystemWithBackend(N, F)
	c := sys.backends[0].engine.(*core)
	c.Start()
	sys.backends[1].engine.(*core).Start()
	// getPreparedCertificate defaults to 50% commits, 50% prepares. Modify the function to change the ratio.
	preparedCertificate := sys.getPreparedCertificate(b, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, makeBlock(1))
	msg, err := sys.backends[1].getRoundChangeMessage(istanbul.View{Round: big.NewInt(1), Sequence: big.NewInt(1)}, preparedCertificate)
	if err != nil {
		b.Errorf("Error creating a round change message. err: %v", err)
	}

	// benchmarked portion
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = c.handleRoundChange(&msg)
		if err != nil {
			b.Errorf("Error handling the round change message. err: %v", err)
		}
	}
}

func BenchmarkHandleRoundChange_10(b *testing.B) { benchMarkHandleRoundChange(10, b) }
func BenchmarkHandleRoundChange_50(b *testing.B) { benchMarkHandleRoundChange(50, b) }
func BenchmarkHandleRoundChange_90(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandleRoundChange(90, b)
}
func BenchmarkHandleRoundChange_100(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandleRoundChange(100, b)
}
func BenchmarkHandleRoundChange_120(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandleRoundChange(120, b)
}
func BenchmarkHandleRoundChange_150(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandleRoundChange(150, b)
}
func BenchmarkHandleRoundChange_200(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandleRoundChange(200, b)
}

// benchMarkHandlePreprepare benchmarks handling a preprepare with a round change certificate that has
// filled round change messages (i.e. the round change messages have prepared certificates that are not empty)
func benchMarkHandlePreprepare(n int, b *testing.B) {
	// Setup
	getRoundState := func(c *core) *roundStateImpl {
		return c.current.(*rsSaveDecorator).rs.(*roundStateImpl)
	}
	N := uint64(n)
	F := uint64(1) // F does not affect tests
	sys := NewMutedTestSystemWithBackend(N, F)

	for i, backend := range sys.backends {
		backend.engine.(*core).Start()
		c := backend.engine.(*core)
		getRoundState(c).round = big.NewInt(int64(N))
		getRoundState(c).desiredRound = getRoundState(c).round
		if i != 0 {
			getRoundState(c).state = StateAcceptRequest
		}
	}
	c := sys.backends[0].engine.(*core)

	// Create pre-prepare
	block := makeBlock(1)
	nextView := istanbul.View{Round: big.NewInt(int64(N)), Sequence: big.NewInt(1)}
	// getPreparedCertificate defaults to 50% commits, 50% prepares. Modify the function to change the ratio.
	preparedCertificate := sys.getPreparedCertificate(b, []istanbul.View{*(sys.backends[0].engine.(*core).current.View())}, block)
	roundChangeCertificate := sys.getRoundChangeCertificate(b, []istanbul.View{nextView}, preparedCertificate)
	msg, err := sys.backends[0].getPreprepareMessage(nextView, roundChangeCertificate, block)
	if err != nil {
		b.Errorf("Error creating a pre-prepare message. err: %v", err)
	}

	err = msg.Sign(sys.backends[0].Sign)
	require.NoError(b, err)
	payload, err := msg.Payload()
	require.NoError(b, err)

	// benchmarked portion
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = c.handleMsg(payload)
		if err != nil {
			b.Errorf("Error handling the pre-prepare message. err: %v", err)
		}
	}
}

func BenchmarkHandlePreprepare_10(b *testing.B) { benchMarkHandlePreprepare(10, b) }
func BenchmarkHandlePreprepare_50(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandlePreprepare(50, b)
}
func BenchmarkHandlePreprepare_90(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandlePreprepare(90, b)
}
func BenchmarkHandlePreprepare_100(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandlePreprepare(100, b)
}
func BenchmarkHandlePreprepare_120(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandlePreprepare(120, b)
}
func BenchmarkHandlePreprepare_150(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandlePreprepare(150, b)
}
func BenchmarkHandlePreprepare_200(b *testing.B) {
	b.Skip("Skipping slow benchmark")
	benchMarkHandlePreprepare(200, b)
}
