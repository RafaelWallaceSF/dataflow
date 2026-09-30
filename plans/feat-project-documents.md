# Feature: project-level documents/files

## Goal

Add an MVP "Documents" capability to Vikunja projects so users can attach files directly to a project, instead of only attaching files inside individual tasks.

## Current state

- Production LUAH MATRIZ is running `vikunja/vikunja:latest` at Vikunja `v2.6.0`.
- Source checkout is based on upstream tag `v2.6.0` / revision `a4218eeb8fb3ada95e14a2331731524c9120d5bd`.
- Existing task attachments already provide the closest storage, upload, list, download, delete, permission, preview, and test patterns.

## MVP scope

Backend first:

- New database table `project_files`.
- New model `ProjectFile` linked to `projects.id` and existing `files.id`.
- API v2 routes:
  - `GET /api/v2/projects/{project}/files`
  - `POST /api/v2/projects/{project}/files`
  - `GET /api/v2/projects/{project}/files/{file}`
  - `DELETE /api/v2/projects/{project}/files/{file}`
- Permissions inherit from the project:
  - read/list/download requires project read access;
  - upload/delete requires project write access.
- Uploads use existing file storage and configured max-file-size behavior.

Frontend follow-up:

- Add a project route/view for "Documents".
- Add a Documents button/tab in the project wrapper.
- Reuse attachment list/upload UI patterns where practical.

## Non-goals for MVP

- Folder hierarchy.
- Rename/move.
- Full document preview UI.
- Wiki pages.
- Search inside files.

## Test plan

- Add v2 web tests proving project file upload/list/download/delete.
- Add permission test proving inaccessible projects cannot list/upload.
- Run targeted `go test ./pkg/webtests -run TestProjectFilesV2`.
- Run `mage test:filter TestProjectFilesV2` after backend implementation.

## Deployment plan

- Build a custom image from this source.
- Change `/opt/vikunja-luah/docker-compose.yml` from `vikunja/vikunja:latest` to the custom image.
- Preserve existing Postgres database and `./files` volume.
- Run migrations automatically on container startup.
- Verify public URL and project documents workflow.
