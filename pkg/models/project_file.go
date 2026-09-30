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

package models

import (
	"io"
	"time"

	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/files"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/web"

	"xorm.io/xorm"
)

// ProjectFile is a file attached directly to a project, not to a task.
type ProjectFile struct {
	ID        int64 `xorm:"bigint autoincr not null unique pk" json:"id" param:"file" readOnly:"true" doc:"The unique, numeric id of this project file."`
	ProjectID int64 `xorm:"bigint not null INDEX" json:"project_id" param:"project" readOnly:"true" doc:"The id of the project this file belongs to. Taken from the URL, not the body."`
	FileID    int64 `xorm:"bigint not null" json:"-"`

	CreatedByID int64      `xorm:"bigint not null" json:"-"`
	CreatedBy   *user.User `xorm:"-" json:"created_by" readOnly:"true" doc:"The user who uploaded this file."`

	File *files.File `xorm:"-" json:"file" readOnly:"true" doc:"Metadata of the uploaded file. The bytes are fetched from the download endpoint."`

	Created time.Time `xorm:"created" json:"created" readOnly:"true" doc:"A timestamp when this file was uploaded."`

	web.CRUDable    `xorm:"-" json:"-"`
	web.Permissions `xorm:"-" json:"-"`
}

// TableName returns the table name for project files.
func (*ProjectFile) TableName() string {
	return "project_files"
}

// NewProjectFile stores a new project-level file and creates its project_files row.
func (pf *ProjectFile) NewProjectFile(s *xorm.Session, f io.ReadSeeker, realname string, realsize uint64, a web.Auth) error {
	file, err := files.CreateWithSession(s, f, realname, realsize, a)
	if err != nil {
		if files.IsErrFileIsTooLarge(err) {
			return ErrTaskAttachmentIsTooLarge{Size: realsize}
		}
		return err
	}
	pf.File = file
	pf.FileID = file.ID

	pf.CreatedBy, err = GetUserOrLinkShareUser(s, a)
	if err != nil {
		if err2 := file.Delete(s); err2 != nil {
			return err2
		}
		return err
	}
	pf.CreatedByID = pf.CreatedBy.ID

	_, err = s.Insert(pf)
	if err != nil {
		if err2 := file.Delete(s); err2 != nil {
			return err2
		}
		return err
	}
	return nil
}

// UploadProjectFiles checks write access to the project, then stores each file.
func UploadProjectFiles(s *xorm.Session, a web.Auth, projectID int64, uploads []*AttachmentToUpload) (success []*ProjectFile, failures []error, err error) {
	project := &Project{ID: projectID}
	can, err := project.CanWrite(s, a)
	if err != nil {
		return nil, nil, err
	}
	if !can {
		return nil, nil, ErrGenericForbidden{}
	}

	for _, upload := range uploads {
		projectFile := &ProjectFile{ProjectID: projectID}
		if err := projectFile.NewProjectFile(s, upload.Reader, upload.Filename, upload.Size, a); err != nil {
			failures = append(failures, err)
			continue
		}
		success = append(success, projectFile)
	}
	return success, failures, nil
}

// ReadOne loads one project file and verifies it belongs to the URL project.
func (pf *ProjectFile) ReadOne(s *xorm.Session, _ web.Auth) (err error) {
	query := s.Where("id = ?", pf.ID).NoAutoCondition()
	if pf.ProjectID != 0 {
		query = query.And("project_id = ?", pf.ProjectID)
	}

	exists, err := query.Get(pf)
	if err != nil {
		return err
	}
	if !exists {
		return ErrProjectFileDoesNotExist{ProjectID: pf.ProjectID, FileID: pf.ID}
	}

	pf.File = &files.File{ID: pf.FileID}
	if err := pf.File.LoadFileMetaByID(s); err != nil {
		return err
	}

	pf.CreatedBy, err = user.GetUserByID(s, pf.CreatedByID)
	if err != nil && !user.IsErrUserDoesNotExist(err) && !user.IsErrUserStatusError(err) {
		return err
	}
	return nil
}

// ReadProjectFiles returns project-level files, paginated.
func ReadProjectFiles(s *xorm.Session, a web.Auth, projectID int64, page int, perPage int) (projectFiles []*ProjectFile, total int64, err error) {
	project := &Project{ID: projectID}
	canRead, _, err := project.CanRead(s, a)
	if err != nil {
		return nil, 0, err
	}
	if !canRead {
		return nil, 0, ErrGenericForbidden{}
	}

	projectFiles = []*ProjectFile{}
	limit, start := getLimitFromPageIndex(page, perPage)
	query := s.Where("project_id = ?", projectID)
	if limit > 0 {
		query = query.Limit(limit, start)
	}
	if err := query.Find(&projectFiles); err != nil {
		return nil, 0, err
	}
	if len(projectFiles) == 0 {
		return projectFiles, 0, nil
	}

	fileIDs := make([]int64, 0, len(projectFiles))
	userIDs := make([]int64, 0, len(projectFiles))
	for _, pf := range projectFiles {
		fileIDs = append(fileIDs, pf.FileID)
		userIDs = append(userIDs, pf.CreatedByID)
	}

	fs := make(map[int64]*files.File)
	if err := s.In("id", fileIDs).Find(&fs); err != nil {
		return nil, 0, err
	}

	users, err := getUsersOrLinkSharesFromIDs(s, userIDs)
	if err != nil {
		return nil, 0, err
	}

	for _, pf := range projectFiles {
		if createdBy, has := users[pf.CreatedByID]; has {
			pf.CreatedBy = createdBy
		}
		pf.File = fs[pf.FileID]
	}

	total, err = s.Where("project_id = ?", projectID).Count(&ProjectFile{})
	return projectFiles, total, err
}

// GetProjectFileForDownload returns one project file with its storage reader opened.
func GetProjectFileForDownload(a web.Auth, projectID, projectFileID int64) (pf *ProjectFile, err error) {
	s := db.NewSession()
	defer s.Close()

	project := &Project{ID: projectID}
	canRead, _, err := project.CanRead(s, a)
	if err != nil {
		_ = s.Rollback()
		return nil, err
	}
	if !canRead {
		_ = s.Rollback()
		return nil, ErrGenericForbidden{}
	}

	pf = &ProjectFile{ID: projectFileID, ProjectID: projectID}
	if err := pf.ReadOne(s, a); err != nil {
		_ = s.Rollback()
		return nil, err
	}

	if err := s.Commit(); err != nil {
		_ = s.Rollback()
		return nil, err
	}

	if err := pf.File.LoadFileByID(); err != nil {
		return nil, err
	}
	return pf, nil
}

// DeleteProjectFile deletes a project file and its underlying stored file.
func DeleteProjectFile(s *xorm.Session, a web.Auth, projectID, projectFileID int64) error {
	project := &Project{ID: projectID}
	can, err := project.CanWrite(s, a)
	if err != nil {
		return err
	}
	if !can {
		return ErrGenericForbidden{}
	}

	pf := &ProjectFile{ID: projectFileID, ProjectID: projectID}
	if err := pf.ReadOne(s, a); err != nil && !files.IsErrFileDoesNotExist(err) {
		return err
	}

	if _, err := s.Where("project_id = ? AND id = ?", projectID, projectFileID).Delete(&ProjectFile{}); err != nil {
		return err
	}

	err = pf.File.Delete(s)
	if err != nil && files.IsErrFileDoesNotExist(err) {
		return nil
	}
	return err
}
