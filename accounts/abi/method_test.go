// Copyright 2018 The go-odf Authors
// This file is part of the go-odf library.
//
// The go-odf library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-odf library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-odf library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"strings"
	"testing"
)

const modfoddata = `
[
	{"type": "function", "name": "balance", "stateMutability": "view"},
	{"type": "function", "name": "send", "inputs": [{ "name": "amount", "type": "uint256" }]},
	{"type": "function", "name": "transfer", "inputs": [{"name": "from", "type": "address"}, {"name": "to", "type": "address"}, {"name": "value", "type": "uint256"}], "outputs": [{"name": "success", "type": "bool"}]},
	{"constant":false,"inputs":[{"components":[{"name":"x","type":"uint256"},{"name":"y","type":"uint256"}],"name":"a","type":"tuple"}],"name":"tuple","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"components":[{"name":"x","type":"uint256"},{"name":"y","type":"uint256"}],"name":"a","type":"tuple[]"}],"name":"tupleSlice","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"components":[{"name":"x","type":"uint256"},{"name":"y","type":"uint256"}],"name":"a","type":"tuple[5]"}],"name":"tupleArray","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"constant":false,"inputs":[{"components":[{"name":"x","type":"uint256"},{"name":"y","type":"uint256"}],"name":"a","type":"tuple[5][]"}],"name":"complexTuple","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},
	{"stateMutability":"nonpayable","type":"fallback"},
	{"stateMutability":"payable","type":"receive"}
]`

func TestModfodString(t *testing.T) {
	var table = []struct {
		modfod      string
		expectation string
	}{
		{
			modfod:      "balance",
			expectation: "function balance() view returns()",
		},
		{
			modfod:      "send",
			expectation: "function send(uint256 amount) returns()",
		},
		{
			modfod:      "transfer",
			expectation: "function transfer(address from, address to, uint256 value) returns(bool success)",
		},
		{
			modfod:      "tuple",
			expectation: "function tuple((uint256,uint256) a) returns()",
		},
		{
			modfod:      "tupleArray",
			expectation: "function tupleArray((uint256,uint256)[5] a) returns()",
		},
		{
			modfod:      "tupleSlice",
			expectation: "function tupleSlice((uint256,uint256)[] a) returns()",
		},
		{
			modfod:      "complexTuple",
			expectation: "function complexTuple((uint256,uint256)[5][] a) returns()",
		},
		{
			modfod:      "fallback",
			expectation: "fallback() returns()",
		},
		{
			modfod:      "receive",
			expectation: "receive() payable returns()",
		},
	}

	abi, err := JSON(strings.NewReader(modfoddata))
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range table {
		var got string
		if test.modfod == "fallback" {
			got = abi.Fallback.String()
		} else if test.modfod == "receive" {
			got = abi.Receive.String()
		} else {
			got = abi.Modfods[test.modfod].String()
		}
		if got != test.expectation {
			t.Errorf("expected string to be %s, got %s", test.expectation, got)
		}
	}
}

func TestModfodSig(t *testing.T) {
	var cases = []struct {
		modfod string
		expect string
	}{
		{
			modfod: "balance",
			expect: "balance()",
		},
		{
			modfod: "send",
			expect: "send(uint256)",
		},
		{
			modfod: "transfer",
			expect: "transfer(address,address,uint256)",
		},
		{
			modfod: "tuple",
			expect: "tuple((uint256,uint256))",
		},
		{
			modfod: "tupleArray",
			expect: "tupleArray((uint256,uint256)[5])",
		},
		{
			modfod: "tupleSlice",
			expect: "tupleSlice((uint256,uint256)[])",
		},
		{
			modfod: "complexTuple",
			expect: "complexTuple((uint256,uint256)[5][])",
		},
	}
	abi, err := JSON(strings.NewReader(modfoddata))
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range cases {
		got := abi.Modfods[test.modfod].Sig
		if got != test.expect {
			t.Errorf("expected string to be %s, got %s", test.expect, got)
		}
	}
}
