// Copyright 2015 The go-odf Authors
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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/odf/go-odf/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("odf/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("odf/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("odf/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("odf/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("odf/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("odf/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("odf/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("odf/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("odf/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("odf/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("odf/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("odf/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("odf/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("odf/downloader/states/drop", nil)

	throttleCounter = metrics.NewRegisteredCounter("odf/downloader/throttle", nil)
)
