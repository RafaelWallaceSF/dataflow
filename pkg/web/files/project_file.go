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

package files

import "code.vikunja.io/api/pkg/models"

// ProjectFileUploadResult is the outcome of a project-file upload.
type ProjectFileUploadResult struct {
	Errors  []AttachmentUploadError `json:"errors" doc:"Per-file failures. A file that fails here does not fail the whole request; the others still upload."`
	Success []*models.ProjectFile   `json:"success" doc:"The project files that were created successfully."`
}

// BuildProjectFileUploadResult turns the domain function's plain return values into the wire DTO.
func BuildProjectFileUploadResult(success []*models.ProjectFile, failures []error) *ProjectFileUploadResult {
	r := &ProjectFileUploadResult{Success: success}
	for _, err := range failures {
		r.Errors = append(r.Errors, toAttachmentUploadError(err))
	}
	return r
}
