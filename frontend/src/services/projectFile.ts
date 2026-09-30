import AbstractService from '@/services/abstractService'
import ProjectFileModel from '@/models/projectFile'
import type {IProjectFile} from '@/modelTypes/IProjectFile'
import {downloadBlob} from '@/helpers/downloadBlob'
import {apiV2Url} from '@/helpers/fetcher'

export default class ProjectFileService extends AbstractService<IProjectFile> {
	constructor() {
		super({
			create: apiV2Url('projects/{projectId}/files'),
			getAll: apiV2Url('projects/{projectId}/files'),
			delete: apiV2Url('projects/{projectId}/files/{id}'),
		})
	}

	getRouteParameterPattern(): RegExp {
		return /(?:{|%7B)([^}%]+)(?:}|%7D)/
	}

	processModel(model: IProjectFile) {
		return {
			...model,
			created: new Date(model.created).toISOString(),
		}
	}

	useCreateInterceptor() {
		return false
	}

	modelFactory(data: Partial<IProjectFile>) {
		return new ProjectFileModel(data)
	}

	modelCreateFactory(data) {
		data.success = (data.success === null ? [] : data.success).map(a => this.modelFactory(a))
		return data
	}

	modelGetAllFactory(data: Partial<IProjectFile>) {
		return this.modelFactory(data)
	}

	async getAll(model: IProjectFile = new ProjectFileModel({}), params = {}, page = 1): Promise<IProjectFile[]> {
		if (!model.projectId || Number(model.projectId) <= 0) {
			return []
		}
		const cancel = this.setLoading()
		params.page = page
		const finalUrl = this.getReplacedRoute(this.paths.getAll, this.beforeGet(model))

		try {
			const {data} = await this.http.get(finalUrl, {params})
			return (data.items ?? []).map(entry => this.modelGetAllFactory(entry))
		} finally {
			cancel()
		}
	}

	getProjectFileBlobUrl(model: IProjectFile) {
		return AbstractService.prototype.getBlobUrl.call(this, apiV2Url(`projects/${model.projectId}/files/${model.id}`))
	}

	async download(model: IProjectFile) {
		const url = await this.getProjectFileBlobUrl(model)
		return downloadBlob(url, model.file.name)
	}

	async upload(model: IProjectFile, files: File[] | FileList) {
		const data = new FormData()
		for (let i = 0; i < files.length; i++) {
			data.append('files', files[i], files[i].name)
		}

		const cancel = this.setLoading()
		const finalUrl = this.getReplacedRoute(this.paths.create, model)

		try {
			const response = await this.http.post(finalUrl, data, {
				onUploadProgress: ({progress}) => {
					this.uploadProgress = progress ? Math.round(progress * 100) : 0
				},
			})
			return this.modelCreateFactory(response.data)
		} finally {
			this.uploadProgress = 0
			cancel()
		}
	}
}
