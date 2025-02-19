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

package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	apiutil "github.com/superseriousbusiness/gotosocial/internal/api/util"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/oauth"
)

const (
	domainBlockListPath = aboutPath + "/suspended"
)

func (m *Module) domainBlockListGETHandler(c *gin.Context) {
	authed, err := oauth.Authed(c, false, false, false, false)
	if err != nil {
		apiutil.WebErrorHandler(c, gtserror.NewErrorUnauthorized(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	if !config.GetInstanceExposeSuspendedWeb() && (authed.Account == nil || authed.User == nil) {
		err := fmt.Errorf("this instance does not expose the list of suspended domains publicly")
		apiutil.WebErrorHandler(c, gtserror.NewErrorUnauthorized(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	instance, err := m.processor.InstanceGetV1(c.Request.Context())
	if err != nil {
		apiutil.WebErrorHandler(c, gtserror.NewErrorInternalError(err), m.processor.InstanceGetV1)
		return
	}

	domainBlocks, errWithCode := m.processor.InstancePeersGet(c.Request.Context(), true, false, false)
	if errWithCode != nil {
		apiutil.WebErrorHandler(c, errWithCode, m.processor.InstanceGetV1)
		return
	}

	c.HTML(http.StatusOK, "domain-blocklist.tmpl", gin.H{
		"instance":  instance,
		"ogMeta":    ogBase(instance),
		"blocklist": domainBlocks,
		"stylesheets": []string{
			assetsPathPrefix + "/Fork-Awesome/css/fork-awesome.min.css",
		},
		"javascript": []string{distPathPrefix + "/frontend.js"},
	})
}
