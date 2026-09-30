// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package webtests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uploadProjectFileRequest(t *testing.T, e http.Handler, projectID string, files map[string][]byte, token string) *httptest.ResponseRecorder {
	t.Helper()
	body, contentType := multipartFilesBody(t, files)
	req := httptest.NewRequest(http.MethodPost, "/api/v2/projects/"+projectID+"/files", body)
	req.Header.Set("Content-Type", contentType)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func uploadOneProjectFile(t *testing.T, e http.Handler, token, filename string, content []byte) int64 {
	t.Helper()
	rec := uploadProjectFileRequest(t, e, "1", map[string][]byte{filename: content}, token)
	require.Equal(t, http.StatusCreated, rec.Code, "body: %s", rec.Body.String())

	var resp struct {
		Success []struct {
			ID   int64 `json:"id"`
			File struct {
				Name string `json:"name"`
			} `json:"file"`
		} `json:"success"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Empty(t, resp.Errors, "upload reported per-file errors: %+v", resp.Errors)
	require.Len(t, resp.Success, 1)
	require.NotZero(t, resp.Success[0].ID)
	assert.Equal(t, filename, resp.Success[0].File.Name)
	return resp.Success[0].ID
}

func TestProjectFilesV2(t *testing.T) {
	t.Run("Upload list download and delete project file", func(t *testing.T) {
		e, err := setupTestEnv()
		require.NoError(t, err)
		token := humaTokenFor(t, &testuser1)

		content := []byte("project-level document")
		id := uploadOneProjectFile(t, e, token, "brief.txt", content)

		list := humaRequest(t, e, http.MethodGet, "/api/v2/projects/1/files", "", token, "")
		require.Equal(t, http.StatusOK, list.Code, "body: %s", list.Body.String())
		assert.Contains(t, list.Body.String(), "brief.txt")

		download := humaRequest(t, e, http.MethodGet, "/api/v2/projects/1/files/"+strconv.FormatInt(id, 10), "", token, "")
		require.Equal(t, http.StatusOK, download.Code, "body: %s", download.Body.String())
		assert.Equal(t, content, download.Body.Bytes())
		assert.Contains(t, download.Header().Get("Content-Disposition"), "brief.txt")

		deleted := humaRequest(t, e, http.MethodDelete, "/api/v2/projects/1/files/"+strconv.FormatInt(id, 10), "", token, "")
		require.Equal(t, http.StatusNoContent, deleted.Code, "body: %s", deleted.Body.String())

		missing := humaRequest(t, e, http.MethodGet, "/api/v2/projects/1/files/"+strconv.FormatInt(id, 10), "", token, "")
		assert.Equal(t, http.StatusNotFound, missing.Code, "body: %s", missing.Body.String())
	})

	t.Run("Project file routes enforce project permissions", func(t *testing.T) {
		e, err := setupTestEnv()
		require.NoError(t, err)
		token := humaTokenFor(t, &testuser1)

		rec := uploadProjectFileRequest(t, e, "2", map[string][]byte{"secret.txt": []byte("nope")}, token)
		assert.Equal(t, http.StatusForbidden, rec.Code, "body: %s", rec.Body.String())

		list := humaRequest(t, e, http.MethodGet, "/api/v2/projects/2/files", "", token, "")
		assert.Equal(t, http.StatusForbidden, list.Code, "body: %s", list.Body.String())
	})
}
