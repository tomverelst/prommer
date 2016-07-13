// Copyright (c) 2016 Tom Verelst <tomverelst@gmail.com>
// Copyright (c) 2016 Other contributors as noted in the AUTHORS file.
//
// This file is part of Prommer.
//
// Prommer is free software; you can redistribute it and/or modify it under the
// terms of the GNU Lesser General Public License as published by the Free
// Software Foundation; either version 3 of the License, or (at your option)
// any later version.
//
// Prommer is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more
// details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

/*
Package prommer contains Prommer core implementation details.

Prommer is a target/service discovery system for Prometheus using Docker.
It achieves this by listening to the Docker events stream,
fetching containers using the Docker Engine API
(https://github.com/docker/engine-api),
and writing configuration files to Prometheus.
*/
package prommer
