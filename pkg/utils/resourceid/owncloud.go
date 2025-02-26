// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package resourceid

import (
	"encoding/base64"
	"errors"
	"strings"
	"unicode/utf8"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

const (
	idDelimiter string = ":"
)

// OwnCloudResourceIDUnwrap returns the wrapped resource id
// by OwnCloudResourceIDWrap and returns nil if not possible
func OwnCloudResourceIDUnwrap(rid string) *provider.ResourceId {
	id, err := unwrap(rid)
	if err != nil {
		return nil
	}
	return id
}

func unwrap(rid string) (*provider.ResourceId, error) {
	decodedID, err := base64.URLEncoding.DecodeString(rid)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(string(decodedID), idDelimiter, 2)
	if len(parts) != 2 {
		return nil, errors.New("could not find two parts with given delimiter")
	}

	if !utf8.ValidString(parts[0]) || !utf8.ValidString(parts[1]) {
		return nil, errors.New("invalid utf8 string found")
	}

	return &provider.ResourceId{
		StorageId: parts[0],
		OpaqueId:  parts[1],
	}, nil
}

// OwnCloudResourceIDWrap wraps a resource id into a xml safe string
// which can then be passed to the outside world
func OwnCloudResourceIDWrap(r *provider.ResourceId) string {
	return wrap(r.StorageId, r.OpaqueId)
}

// The fileID must be encoded
// - XML safe, because it is going to be used in the propfind result
// - url safe, because the id might be used in a url, eg. the /dav/meta nodes
// which is why we base64 encode it
func wrap(sid string, oid string) string {
	return base64.URLEncoding.EncodeToString([]byte(sid + idDelimiter + oid))
}
