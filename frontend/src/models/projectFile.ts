import AbstractModel from './abstractModel'
import UserModel from './user'
import FileModel from './file'
import type {IUser} from '@/modelTypes/IUser'
import type {IFile} from '@/modelTypes/IFile'
import type {IProjectFile} from '@/modelTypes/IProjectFile'

export default class ProjectFileModel extends AbstractModel<IProjectFile> implements IProjectFile {
	id = 0
	projectId = 0
	createdBy: IUser = UserModel
	file: IFile = FileModel
	created: Date = null

	constructor(data: Partial<IProjectFile> = {}) {
		super()
		this.assignData(data)

		this.createdBy = new UserModel(this.createdBy)
		this.file = new FileModel(this.file)
		this.created = new Date(this.created)
	}
}
