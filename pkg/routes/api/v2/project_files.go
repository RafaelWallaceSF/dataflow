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

package apiv2

import (
	"context"
	"net/http"

	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/models"
	webfiles "code.vikunja.io/api/pkg/web/files"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
)

type projectFileListBody struct {
	Body Paginated[*models.ProjectFile]
}

type projectFileUploadInput struct {
	ProjectID int64 `path:"project" doc:"The id of the project to attach the files to."`
	RawBody   huma.MultipartFormFiles[struct {
		Files []huma.FormFile `form:"files" required:"true" doc:"One or more files to upload as project files. Send multiple parts under the same \"files\" field to upload several at once."`
	}]
}

type projectFileUploadBody struct {
	Body *webfiles.ProjectFileUploadResult
}

// RegisterProjectFileRoutes wires project-level file list/upload/download/delete onto the Huma API.
func RegisterProjectFileRoutes(api huma.API) {
	if !config.ServiceEnableTaskAttachments.GetBool() {
		return
	}

	tags := []string{"projects"}

	Register(api, huma.Operation{
		OperationID: "project-files-list",
		Summary:     "List a project's files",
		Description: "Returns project-level file metadata, paginated. Requires read access to the project. The file bytes are not included; fetch them from the download endpoint.",
		Method:      http.MethodGet,
		Path:        "/projects/{project}/files",
		Tags:        tags,
	}, projectFilesList)

	Register(api, huma.Operation{
		OperationID:  "project-files-upload",
		Summary:      "Upload project files",
		Description:  "Uploads one or more files directly to a project via multipart/form-data under the \"files\" field. Requires write access to the project. Each file is processed independently.",
		Method:       http.MethodPost,
		Path:         "/projects/{project}/files",
		Tags:         tags,
		MaxBodyBytes: (int64(config.GetMaxFileSizeInMBytes()) + 2) * 1024 * 1024,
	}, projectFilesUpload)

	Register(api, huma.Operation{
		OperationID: "project-files-download",
		Summary:     "Download a project file",
		Description: "Returns the raw bytes of one project-level file. Requires read access to the project.",
		Method:      http.MethodGet,
		Path:        "/projects/{project}/files/{file}",
		Tags:        tags,
		Responses: map[string]*huma.Response{
			"200": {
				Description: "The project file bytes. The Content-Type header carries the file's mime type.",
				Content: map[string]*huma.MediaType{
					"application/octet-stream": {Schema: &huma.Schema{Type: huma.TypeString, Format: "binary"}},
				},
			},
		},
	}, projectFilesDownload)

	Register(api, huma.Operation{
		OperationID: "project-files-delete",
		Summary:     "Delete a project file",
		Description: "Deletes one project-level file and its underlying stored file. Requires write access to the project.",
		Method:      http.MethodDelete,
		Path:        "/projects/{project}/files/{file}",
		Tags:        tags,
	}, projectFilesDelete)
}

func init() { AddRouteRegistrar(RegisterProjectFileRoutes) }

func projectFilesList(ctx context.Context, in *struct {
	ProjectID int64 `path:"project" doc:"The id of the project whose files to list."`
	ListParams
}) (*projectFileListBody, error) {
	a, err := authFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	s := db.NewSession()
	defer s.Close()
	items, total, err := models.ReadProjectFiles(s, a, in.ProjectID, in.Page, in.PerPage)
	if err != nil {
		_ = s.Rollback()
		return nil, translateDomainError(err)
	}
	if err := s.Commit(); err != nil {
		_ = s.Rollback()
		return nil, translateDomainError(err)
	}
	return &projectFileListBody{Body: NewPaginated(items, total, in.Page, in.PerPage)}, nil
}

func projectFilesUpload(ctx context.Context, in *projectFileUploadInput) (*projectFileUploadBody, error) {
	a, err := authFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	s := db.NewSession()
	defer s.Close()

	formFiles := in.RawBody.Data().Files
	uploads := make([]*models.AttachmentToUpload, 0, len(formFiles))
	for _, file := range formFiles {
		uploads = append(uploads, &models.AttachmentToUpload{Reader: file, Filename: file.Filename, Size: uint64(file.Size)})
	}

	success, failures, err := models.UploadProjectFiles(s, a, in.ProjectID, uploads)
	if err != nil {
		_ = s.Rollback()
		return nil, translateDomainError(err)
	}
	if err := s.Commit(); err != nil {
		_ = s.Rollback()
		return nil, translateDomainError(err)
	}

	return &projectFileUploadBody{Body: webfiles.BuildProjectFileUploadResult(success, failures)}, nil
}

func projectFilesDownload(ctx context.Context, in *struct {
	ProjectID int64 `path:"project" doc:"The id of the project the file belongs to."`
	FileID    int64 `path:"file" doc:"The id of the project file to download."`
}) (*huma.StreamResponse, error) {
	a, err := authFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	pf, err := models.GetProjectFileForDownload(a, in.ProjectID, in.FileID)
	if err != nil {
		return nil, translateDomainError(err)
	}

	return &huma.StreamResponse{Body: func(hctx huma.Context) {
		c := humaecho.Unwrap(hctx)
		defer func() { _ = pf.File.File.Close() }()
		webfiles.WriteFileDownload((*c).Response(), (*c).Request(), pf.File)
	}}, nil
}

func projectFilesDelete(ctx context.Context, in *struct {
	ProjectID int64 `path:"project" doc:"The id of the project the file belongs to."`
	FileID    int64 `path:"file" doc:"The id of the project file to delete."`
}) (*emptyBody, error) {
	a, err := authFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	s := db.NewSession()
	defer s.Close()
	if err := models.DeleteProjectFile(s, a, in.ProjectID, in.FileID); err != nil {
		_ = s.Rollback()
		return nil, translateDomainError(err)
	}
	if err := s.Commit(); err != nil {
		_ = s.Rollback()
		return nil, translateDomainError(err)
	}
	return &emptyBody{}, nil
}
