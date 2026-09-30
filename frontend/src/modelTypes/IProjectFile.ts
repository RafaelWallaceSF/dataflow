import type {IAbstract} from './IAbstract'
import type {IFile} from './IFile'
import type {IUser} from './IUser'

export interface IProjectFile extends IAbstract {
	id: number
	projectId: number
	createdBy: IUser
	file: IFile
	created: Date
}
