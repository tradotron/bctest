/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	cb "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// GetMockEndorserClient return a endorser client return specified ProposalResponse and err(nil or error)
func GetMockEndorserClient(repponse *pb.ProposalResponse, err error) pb.EndorserClient {
	return &mockEndorserClient{
		repponse: repponse,
		err:      err,
	}
}

type mockEndorserClient struct {
	repponse *pb.ProposalResponse
	err      error
}

func (m *mockEndorserClient) ProcessProposal(ctx context.Context, in *pb.SignedProposal, opts ...grpc.CallOption) (*pb.ProposalResponse, error) {
	return m.repponse, m.err
}

func GetMockBroadcastClient(err error) BroadcastClient {
	return &mockBroadcastClient{err: err}
}

// mockBroadcastClient return success immediately
type mockBroadcastClient struct {
	err error
}

func (m *mockBroadcastClient) Send(env *cb.Envelope) error {
	return m.err
}

func (m *mockBroadcastClient) Close() error {
	return nil
}
