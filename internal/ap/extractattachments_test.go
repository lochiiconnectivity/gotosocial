// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ap_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/superseriousbusiness/activity/streams"
	"github.com/superseriousbusiness/gotosocial/internal/ap"
)

type ExtractAttachmentsTestSuite struct {
	APTestSuite
}

func (suite *ExtractAttachmentsTestSuite) TestExtractAttachmentMissingURL() {
	d1 := suite.document1
	d1.SetActivityStreamsUrl(streams.NewActivityStreamsUrlProperty())

	attachment, err := ap.ExtractAttachment(d1)
	suite.EqualError(err, "ExtractAttachment: error extracting attachment URL: ExtractURL: no valid URL property found")
	suite.Nil(attachment)
}

func TestExtractAttachmentsTestSuite(t *testing.T) {
	suite.Run(t, &ExtractAttachmentsTestSuite{})
}
