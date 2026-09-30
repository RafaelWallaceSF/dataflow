<template>
	<ProjectWrapper
		class="project-documents"
		:is-loading-project="isLoadingProject"
		:project-id="projectId"
		:view-id="viewId"
	>
		<div class="loader-container is-max-width-desktop" :class="{'is-loading': projectFileService.loading}">
			<Card class="project-documents-card">
				<div class="header-section mbe-4">
					<div class="header-titles">
						<h2>{{ t('project.documents.title') }}</h2>
						<p class="has-text-grey">
							{{ t('project.documents.description') }}
						</p>
					</div>
					<div v-if="files.length > 0" class="file-count-badge">
						{{ files.length }} {{ files.length === 1 ? 'documento' : 'documentos' }}
					</div>
				</div>

				<!-- Hidden input -->
				<input
					ref="filesRef"
					multiple
					type="file"
					:disabled="projectFileService.loading || undefined"
					@change="handleFileInputChange"
				>

				<!-- Drop Zone -->
				<div
					v-if="canWrite"
					class="drop-zone mbe-4"
					:class="{'is-dragging': isDragging, 'is-busy': projectFileService.loading}"
					@dragover.prevent="onDragOver"
					@dragleave.prevent="onDragLeave"
					@drop.prevent="onDrop"
					@click="filesRef?.click()"
				>
					<Icon icon="cloud-upload-alt" class="drop-zone-icon" />
					<div class="drop-zone-text">
						<strong>{{ t('project.documents.upload') }}</strong>
						<span class="has-text-grey is-size-7">Arraste e solte arquivos aqui ou clique para selecionar</span>
					</div>
				</div>

				<!-- Progress bar -->
				<ProgressBar
					v-if="projectFileService.uploadProgress > 0"
					:value="projectFileService.uploadProgress"
					is-primary
					class="mbe-4"
				/>

				<!-- Empty State -->
				<Nothing v-if="files.length === 0 && !projectFileService.loading">
					{{ t('project.documents.empty') }}
				</Nothing>

				<!-- Document List -->
				<div v-else class="project-documents-list">
					<div
						v-for="projectFile in files"
						:key="projectFile.id"
						class="project-document-row"
					>
						<div class="document-info">
							<Icon :icon="getFileIcon(projectFile.file?.name)" class="doc-icon" />
							<div class="doc-meta">
								<span class="doc-name" :title="projectFile.file?.name">
									{{ projectFile.file?.name || 'Documento sem nome' }}
								</span>
								<div class="doc-details has-text-grey is-size-7">
									<span>{{ getHumanSize(projectFile.file?.size || 0) }}</span>
									<span v-if="projectFile.created" class="detail-sep">•</span>
									<span v-if="projectFile.created">{{ formatDate(projectFile.created) }}</span>
									<span v-if="projectFile.createdBy?.username" class="detail-sep">•</span>
									<span v-if="projectFile.createdBy?.username">por {{ projectFile.createdBy.username }}</span>
								</div>
							</div>
						</div>
						<div class="project-document-actions">
							<BaseButton
								:aria-label="t('project.documents.download')"
								class="action-btn download-btn"
								@click="downloadProjectFile(projectFile)"
							>
								<Icon icon="download" />
							</BaseButton>
							<BaseButton
								v-if="canWrite"
								:aria-label="t('project.documents.delete')"
								class="action-btn delete-btn"
								@click="deleteProjectFile(projectFile)"
							>
								<Icon icon="trash-alt" />
							</BaseButton>
						</div>
					</div>
				</div>
			</Card>
		</div>
	</ProjectWrapper>
</template>

<script setup lang="ts">
import {computed, onMounted, ref, shallowReactive, watch} from 'vue'
import {useI18n} from 'vue-i18n'

import ProjectWrapper from '@/components/project/ProjectWrapper.vue'
import BaseButton from '@/components/base/BaseButton.vue'
import Nothing from '@/components/misc/Nothing.vue'
import ProgressBar from '@/components/misc/ProgressBar.vue'
import Icon from '@/components/misc/Icon'

import ProjectFileService from '@/services/projectFile'
import ProjectFileModel from '@/models/projectFile'
import type {IProjectFile} from '@/modelTypes/IProjectFile'
import type {IProject} from '@/modelTypes/IProject'
import type {IProjectView} from '@/modelTypes/IProjectView'
import {useProjectStore} from '@/stores/projects'
import {useBaseStore} from '@/stores/base'
import {PERMISSIONS as Permissions} from '@/constants/permissions'
import {getHumanSize} from '@/helpers/getHumanSize'
import {success as notifySuccess, error as notifyError} from '@/message'

const props = defineProps<{
	isLoadingProject: boolean,
	projectId: IProject['id'],
	viewId: IProjectView['id'],
}>()

const {t} = useI18n({useScope: 'global'})
const projectStore = useProjectStore()
const baseStore = useBaseStore()
const projectFileService = shallowReactive(new ProjectFileService())
const filesRef = ref<HTMLInputElement | null>(null)
const files = ref<IProjectFile[]>([])
const isDragging = ref(false)

const project = computed(() => projectStore.projects[props.projectId] ?? baseStore.currentProject)
const canWrite = computed(() => {
	const perm = project.value?.maxPermission ?? baseStore.currentProject?.maxPermission
	if (perm === undefined || perm === null) return true
	return perm > Permissions.READ
})

function onDragOver() {
	if (!projectFileService.loading) {
		isDragging.value = true
	}
}

function onDragLeave() {
	isDragging.value = false
}

async function onDrop(e: DragEvent) {
	isDragging.value = false
	if (projectFileService.loading) return
	if (e.dataTransfer?.files?.length) {
		await uploadFiles(e.dataTransfer.files)
	}
}

async function handleFileInputChange() {
	if (!filesRef.value?.files?.length) return
	await uploadFiles(filesRef.value.files)
	filesRef.value.value = ''
}

function getFileIcon(filename: string = ''): string {
	const ext = filename.split('.').pop()?.toLowerCase() || ''
	if (['png', 'jpg', 'jpeg', 'gif', 'svg', 'webp', 'bmp'].includes(ext)) {
		return 'file-image'
	}
	if (ext === 'pdf') {
		return 'file-pdf'
	}
	if (['zip', 'tar', 'gz', 'rar', '7z'].includes(ext)) {
		return 'archive'
	}
	if (['js', 'ts', 'py', 'go', 'json', 'html', 'css', 'sql', 'sh'].includes(ext)) {
		return 'code'
	}
	return 'file'
}

function formatDate(date: string | Date | number): string {
	try {
		return new Date(date).toLocaleDateString('pt-BR', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		})
	} catch {
		return ''
	}
}

async function loadProjectFiles() {
	if (!props.projectId || Number(props.projectId) <= 0) return
	try {
		files.value = await projectFileService.getAll(new ProjectFileModel({projectId: Number(props.projectId)}))
	} catch (e) {
		console.error('Failed to load project files:', e)
	}
}

async function uploadFiles(fileList: FileList | File[]) {
	if (!props.projectId || Number(props.projectId) <= 0) return
	try {
		const uploaded = await projectFileService.upload(new ProjectFileModel({projectId: Number(props.projectId)}), fileList)
		const newFiles = (uploaded.success ?? []) as IProjectFile[]
		if (newFiles.length > 0) {
			files.value = [...files.value, ...newFiles]
			notifySuccess('Documento(s) enviado(s) com sucesso!')
		}
		if (uploaded.errors?.length) {
			notifyError(new Error('Alguns arquivos não puderam ser enviados.'))
		}
	} catch (e) {
		notifyError(e)
	}
}

async function downloadProjectFile(projectFile: IProjectFile) {
	try {
		await projectFileService.download(projectFile)
	} catch (e) {
		notifyError(e)
	}
}

async function deleteProjectFile(projectFile: IProjectFile) {
	if (!confirm('Deseja realmente excluir este documento?')) {
		return
	}
	try {
		await projectFileService.delete(projectFile)
		files.value = files.value.filter(({id}) => id !== projectFile.id)
		notifySuccess('Documento excluído com sucesso.')
	} catch (e) {
		notifyError(e)
	}
}

onMounted(loadProjectFiles)
watch(() => props.projectId, loadProjectFiles)
</script>

<style lang="scss" scoped>
.project-documents-card {
	padding: 1.5rem;
	border-radius: $radius;

	input[type='file'] {
		display: none;
	}
}

.header-section {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	gap: 1rem;
}

.file-count-badge {
	padding: .25rem .75rem;
	background: var(--grey-100);
	color: var(--grey-700);
	border-radius: 9999px;
	font-size: .8rem;
	font-weight: 600;
	white-space: nowrap;
}

.drop-zone {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: .75rem;
	padding: 2rem 1.5rem;
	border: 2px dashed var(--grey-300);
	border-radius: $radius;
	background: var(--grey-50);
	cursor: pointer;
	transition: all .2s ease;
	text-align: center;

	&:hover,
	&.is-dragging {
		border-color: var(--primary);
		background: rgba(25, 115, 255, 0.05);
	}

	&.is-busy {
		pointer-events: none;
		opacity: .6;
	}

	.drop-zone-icon {
		font-size: 2rem;
		color: var(--primary);
	}

	.drop-zone-text {
		display: flex;
		flex-direction: column;
		gap: .25rem;
	}
}

.project-documents-list {
	display: flex;
	flex-direction: column;
	gap: .75rem;
}

.project-document-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 1rem;
	padding: .85rem 1rem;
	border: 1px solid var(--grey-200);
	border-radius: $radius;
	background: var(--card-background);
	transition: transform .15s ease, box-shadow .15s ease;

	&:hover {
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
	}

	.document-info {
		display: flex;
		align-items: center;
		gap: 1rem;
		min-width: 0;

		.doc-icon {
			font-size: 1.5rem;
			color: var(--primary);
			flex-shrink: 0;
		}

		.doc-meta {
			display: flex;
			flex-direction: column;
			min-width: 0;

			.doc-name {
				font-weight: 600;
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}

			.doc-details {
				display: flex;
				align-items: center;
				gap: .4rem;

				.detail-sep {
					opacity: .5;
				}
			}
		}
	}

	.project-document-actions {
		display: flex;
		align-items: center;
		gap: .5rem;
		flex-shrink: 0;

		.action-btn {
			padding: .4rem .6rem;
			border-radius: $radius;
			border: 1px solid var(--grey-200);
			background: transparent;
			cursor: pointer;
			transition: background .2s ease, color .2s ease;

			&.download-btn:hover {
				background: var(--primary);
				color: white;
				border-color: var(--primary);
			}

			&.delete-btn:hover {
				background: #e53e3e;
				color: white;
				border-color: #e53e3e;
			}
		}
	}
}
</style>
