// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.
package params

// EticaBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the Etica network.
var CrucibleBootnodes = []string{
"enode://b20956879ebba0b1df9fab12c059cc905e61626d1d36e9c9109bc477ccc937d3c36d4daccb790fc41fbc86b94b80c723c56875e7d2cec4f32b0d04fc8a95f79d@173.212.202.226:30303",
}

var dnsPrefixCRUCIBLE = "enrtree://AJE62Q4DUX4QMMXEHCSSCSC65TDHZYSMONSD64P3WULVLSF6MRQ3K@" //not used set random value

var CrucibleDNSNetwork1 = dnsPrefixCRUCIBLE + "all.crucible.blockd.info" //not used set random value
