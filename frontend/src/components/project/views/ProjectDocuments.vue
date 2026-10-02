<template>
	<ProjectWrapper
		class="project-documents"
		:is-loading-project="isLoadingProject"
		:project-id="projectId"
		:view-id="viewId"
	>
		<div class="loader-container is-max-width-desktop" :class="{'is-loading': projectFileService.loading}">
			<Card class="project-documents-card">
				<!-- Header Section -->
				<div class="header-section mbe-4">
					<div class="header-titles">
						<div class="title-with-badge">
							<h2>Documentos do Projeto</h2>
							<div class="stats-pills">
								<span class="stat-pill">
									<Icon icon="file" class="stat-icon" />
									{{ filteredFiles.length }} de {{ files.length }} arquivo(s)
								</span>
								<span class="stat-pill">
									<Icon icon="folder" class="stat-icon" />
									{{ allFolders.length }} pasta(s)
								</span>
								<span class="stat-pill">
									<Icon icon="database" class="stat-icon" />
									{{ totalSizeFormatted }}
								</span>
							</div>
						</div>
						<p class="has-text-grey is-size-7">
							Organize seus contratos, briefings, entregáveis e notas fiscais com pastas, preview nativo e busca rápida.
						</p>
					</div>

					<div class="header-actions">
						<BaseButton
							v-if="canWrite"
							class="action-header-btn new-folder-btn"
							@click="openNewFolderModal"
						>
							<Icon icon="folder-plus" />
							<span>Nova Pasta</span>
						</BaseButton>

						<BaseButton
							v-if="canWrite"
							is-primary
							class="action-header-btn upload-btn"
							:disabled="projectFileService.loading || undefined"
							@click="filesRef?.click()"
						>
							<Icon icon="cloud-upload-alt" />
							<span>Enviar Arquivos</span>
						</BaseButton>
					</div>
				</div>

				<!-- Hidden input for file upload -->
				<input
					ref="filesRef"
					multiple
					type="file"
					:disabled="projectFileService.loading || undefined"
					@change="handleFileInputChange"
				>

				<!-- Breadcrumb & Path Navigation -->
				<div class="navigation-bar mbe-4">
					<div class="breadcrumbs">
						<button
							class="crumb-btn"
							:class="{'is-active': currentFolder === ''}"
							@click="currentFolder = ''"
						>
							<Icon icon="home" />
							<span>Início (Todas as Pastas)</span>
						</button>
						<template v-if="currentFolder">
							<span class="crumb-separator">/</span>
							<button class="crumb-btn is-active">
								<Icon icon="folder-open" class="crumb-folder-icon" />
								<span>{{ currentFolder }}</span>
							</button>
						</template>
					</div>

					<!-- Search Input -->
					<div class="search-box">
						<Icon icon="search" class="search-icon" />
						<input
							v-model="searchQuery"
							type="text"
							placeholder="Buscar documento..."
							class="search-input"
						>
						<button
							v-if="searchQuery"
							class="clear-search-btn"
							@click="searchQuery = ''"
						>
							<Icon icon="times" />
						</button>
					</div>
				</div>

				<!-- Quick Folder Switcher Chips -->
				<div class="folder-chips-bar mbe-4">
					<button
						class="folder-chip"
						:class="{'is-selected': currentFolder === ''}"
						@click="currentFolder = ''"
					>
						<Icon icon="layer-group" />
						<span>Todas as Pastas</span>
						<span class="chip-count">{{ files.length }}</span>
					</button>

					<button
						v-for="folder in allFolders"
						:key="folder.name"
						class="folder-chip"
						:class="{'is-selected': currentFolder === folder.name}"
						@click="currentFolder = folder.name"
					>
						<Icon :icon="currentFolder === folder.name ? 'folder-open' : 'folder'" />
						<span>{{ folder.name }}</span>
						<span class="chip-count">{{ folder.count }}</span>
					</button>

					<button
						v-if="canWrite"
						class="folder-chip add-chip"
						@click="openNewFolderModal"
					>
						<Icon icon="plus" />
						<span>Pasta</span>
					</button>
				</div>

				<!-- Folder Cards Grid (shown when at Root and folders exist) -->
				<div v-if="currentFolder === '' && allFolders.length > 0 && !searchQuery" class="folders-grid mbe-4">
					<div
						v-for="folder in allFolders"
						:key="folder.name"
						class="folder-card"
						@click="currentFolder = folder.name"
					>
						<div class="folder-card-top">
							<div class="folder-card-icon-wrapper">
								<Icon icon="folder" class="folder-card-icon" />
							</div>
							<div class="folder-card-actions" @click.stop>
								<button
									v-if="canWrite"
									class="folder-menu-btn"
									title="Excluir pasta vazia"
									@click="deleteFolder(folder.name)"
								>
									<Icon icon="trash-alt" />
								</button>
							</div>
						</div>
						<div class="folder-card-info">
							<span class="folder-card-title">{{ folder.name }}</span>
							<span class="folder-card-subtitle">
								{{ folder.count }} {{ folder.count === 1 ? 'arquivo' : 'arquivos' }} • {{ getHumanSize(folder.size) }}
							</span>
						</div>
					</div>
				</div>

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
						<strong>{{ currentFolder ? `Enviar para pasta "${currentFolder}"` : 'Enviar para pasta Raiz' }}</strong>
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

				<!-- Type Filters -->
				<div class="type-filter-row mbe-4">
					<div class="type-filter-chips">
						<button
							v-for="tf in typeFilters"
							:key="tf.key"
							class="type-filter-btn"
							:class="{'is-active': selectedTypeFilter === tf.key}"
							@click="selectedTypeFilter = tf.key"
						>
							<Icon :icon="tf.icon" />
							<span>{{ tf.label }}</span>
						</button>
					</div>

					<div v-if="currentFolder" class="folder-active-indicator">
						<span>Filtrando por: <strong>{{ currentFolder }}</strong></span>
						<button class="clear-folder-btn" @click="currentFolder = ''">
							<Icon icon="times" />
						</button>
					</div>
				</div>

				<!-- Empty State -->
				<Nothing v-if="filteredFiles.length === 0 && !projectFileService.loading">
					<template v-if="searchQuery">
						Nenhum documento encontrado para "{{ searchQuery }}".
					</template>
					<template v-else-if="currentFolder">
						Nenhum documento na pasta "{{ currentFolder }}". Arraste arquivos para adicioná-los aqui!
					</template>
					<template v-else>
						Nenhum documento enviado ainda neste projeto.
					</template>
				</Nothing>

				<!-- Document List -->
				<div v-else class="project-documents-list">
					<div
						v-for="projectFile in filteredFiles"
						:key="projectFile.id"
						class="project-document-row"
					>
						<div class="document-info" @click="openPreview(projectFile)">
							<div class="doc-icon-container" :class="`ext-${getFileCategory(getCleanFileName(projectFile))}`">
								<Icon :icon="getFileIcon(getCleanFileName(projectFile))" class="doc-icon" />
							</div>
							<div class="doc-meta">
								<div class="doc-title-row">
									<span class="doc-name" :title="getCleanFileName(projectFile)">
										{{ getCleanFileName(projectFile) }}
									</span>
									<span v-if="getFileFolder(projectFile)" class="doc-folder-tag" :title="`Pasta: ${getFileFolder(projectFile)}`">
										<Icon icon="folder" class="tag-folder-icon" />
										{{ getFileFolder(projectFile) }}
									</span>
								</div>
								<div class="doc-details has-text-grey is-size-7">
									<span>{{ getHumanSize(projectFile.file?.size || 0) }}</span>
									<span v-if="projectFile.created" class="detail-sep">•</span>
									<span v-if="projectFile.created">{{ formatDate(projectFile.created) }}</span>
									<span v-if="projectFile.createdBy?.username" class="detail-sep">•</span>
									<span v-if="projectFile.createdBy?.username" class="doc-author">
										<Icon icon="user" class="author-icon" />
										{{ projectFile.createdBy.username }}
									</span>
								</div>
							</div>
						</div>

						<div class="project-document-actions">
							<!-- Visualizar / Preview -->
							<button
								class="action-btn preview-btn"
								title="Visualizar documento"
								@click="openPreview(projectFile)"
							>
								<Icon icon="eye" />
								<span class="btn-label-desktop">Ver</span>
							</button>

							<!-- Download -->
							<button
								class="action-btn download-btn"
								title="Baixar arquivo"
								@click="downloadProjectFile(projectFile)"
							>
								<Icon icon="download" />
							</button>

							<!-- Move to folder -->
							<button
								v-if="canWrite"
								class="action-btn move-btn"
								title="Mover para outra pasta"
								@click="openMoveModal(projectFile)"
							>
								<Icon icon="exchange-alt" />
							</button>

							<!-- Delete -->
							<button
								v-if="canWrite"
								class="action-btn delete-btn"
								title="Excluir arquivo"
								@click="deleteProjectFile(projectFile)"
							>
								<Icon icon="trash-alt" />
							</button>
						</div>
					</div>
				</div>
			</Card>
		</div>

		<!-- MODAL: NOVA PASTA -->
		<div v-if="showNewFolderModal" class="custom-modal-overlay" @click.self="showNewFolderModal = false">
			<div class="custom-modal-content">
				<div class="modal-header">
					<div class="modal-header-title">
						<Icon icon="folder-plus" class="modal-title-icon" />
						<h3>Criar Nova Pasta</h3>
					</div>
					<button class="modal-close-btn" @click="showNewFolderModal = false">
						<Icon icon="times" />
					</button>
				</div>
				<div class="modal-body">
					<label class="modal-label">Nome da Pasta</label>
					<input
						v-model="newFolderName"
						type="text"
						placeholder="Ex: Contratos, Briefings, Entregáveis, Financeiro..."
						class="modal-input"
						autofocus
						@keyup.enter="createFolder"
					>

					<div class="folder-presets mbs-3">
						<span class="preset-label">Sugestões rápidas:</span>
						<div class="preset-chips">
							<button
								v-for="preset in ['Contratos', 'Briefings', 'Entregáveis', 'Financeiro & NF', 'Assets & Design', 'Comprovantes']"
								:key="preset"
								class="preset-chip"
								@click="newFolderName = preset"
							>
								+ {{ preset }}
							</button>
						</div>
					</div>
				</div>
				<div class="modal-footer">
					<BaseButton class="modal-cancel-btn" @click="showNewFolderModal = false">
						Cancelar
					</BaseButton>
					<BaseButton is-primary :disabled="!newFolderName.trim()" @click="createFolder">
						Criar Pasta
					</BaseButton>
				</div>
			</div>
		</div>

		<!-- MODAL: MOVER ARQUIVO -->
		<div v-if="showMoveModal && fileToMove" class="custom-modal-overlay" @click.self="showMoveModal = false">
			<div class="custom-modal-content">
				<div class="modal-header">
					<div class="modal-header-title">
						<Icon icon="exchange-alt" class="modal-title-icon" />
						<h3>Mover Documento</h3>
					</div>
					<button class="modal-close-btn" @click="showMoveModal = false">
						<Icon icon="times" />
					</button>
				</div>
				<div class="modal-body">
					<p class="mbe-3">
						Mover <strong>{{ getCleanFileName(fileToMove) }}</strong> para:
					</p>

					<div class="folder-select-list mbe-3">
						<button
							class="folder-select-item"
							:class="{'is-selected': targetMoveFolder === ''}"
							@click="targetMoveFolder = ''"
						>
							<Icon icon="home" />
							<span>Raiz (Sem Pasta)</span>
						</button>
						<button
							v-for="folder in allFolders"
							:key="folder.name"
							class="folder-select-item"
							:class="{'is-selected': targetMoveFolder === folder.name}"
							@click="targetMoveFolder = folder.name"
						>
							<Icon icon="folder" />
							<span>{{ folder.name }}</span>
						</button>
					</div>

					<div class="or-divider"><span>ou criar uma nova pasta</span></div>
					<input
						v-model="customMoveFolder"
						type="text"
						placeholder="Digite o nome de uma nova pasta..."
						class="modal-input mbs-2"
						@input="targetMoveFolder = customMoveFolder"
					>
				</div>
				<div class="modal-footer">
					<BaseButton class="modal-cancel-btn" @click="showMoveModal = false">
						Cancelar
					</BaseButton>
					<BaseButton is-primary @click="confirmMoveFile">
						Confirmar Mover
					</BaseButton>
				</div>
			</div>
		</div>

		<!-- MODAL: VISUALIZAÇÃO / PREVIEW NATIVO -->
		<div v-if="showPreviewModal && previewFile" class="custom-modal-overlay preview-modal-overlay" @click.self="closePreview">
			<div class="custom-modal-content preview-modal-content">
				<div class="modal-header preview-header">
					<div class="preview-title-info">
						<Icon :icon="getFileIcon(getCleanFileName(previewFile))" class="preview-header-icon" />
						<div>
							<h3 class="preview-file-name">{{ getCleanFileName(previewFile) }}</h3>
							<span class="has-text-grey is-size-7">
								{{ getHumanSize(previewFile.file?.size || 0) }} • {{ formatDate(previewFile.created) }}
							</span>
						</div>
					</div>
					<div class="preview-header-actions">
						<BaseButton class="action-btn download-btn" @click="downloadProjectFile(previewFile)">
							<Icon icon="download" />
							<span>Baixar</span>
						</BaseButton>
						<button class="modal-close-btn" @click="closePreview">
							<Icon icon="times" />
						</button>
					</div>
				</div>

				<div class="preview-body">
					<!-- Loading preview -->
					<div v-if="isLoadingPreview" class="preview-loader">
						<Icon icon="spinner" class="fa-spin preview-spinner" />
						<span>Carregando visualização...</span>
					</div>

					<!-- Image Preview -->
					<div v-else-if="previewType === 'image'" class="image-preview-wrapper">
						<img :src="previewBlobUrl" :alt="getCleanFileName(previewFile)" class="preview-img">
					</div>

					<!-- PDF Preview -->
					<div v-else-if="previewType === 'pdf'" class="pdf-preview-wrapper">
						<iframe :src="previewBlobUrl" class="pdf-iframe" title="PDF Preview"></iframe>
					</div>

					<!-- Text / Code Preview -->
					<div v-else-if="previewType === 'text'" class="text-preview-wrapper">
						<pre class="code-view"><code>{{ previewTextContent }}</code></pre>
					</div>

					<!-- Generic Fallback -->
					<div v-else class="generic-preview-wrapper">
						<Icon :icon="getFileIcon(getCleanFileName(previewFile))" class="generic-preview-icon" />
						<h4>{{ getCleanFileName(previewFile) }}</h4>
						<p class="has-text-grey is-size-7 mbe-4">
							Este tipo de arquivo não possui pré-visualização inline no navegador.
						</p>
						<BaseButton is-primary @click="downloadProjectFile(previewFile)">
							<Icon icon="download" />
							<span>Baixar Documento</span>
						</BaseButton>
					</div>
				</div>
			</div>
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
import Card from '@/components/misc/Card.vue'
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

// Folder navigation state
const currentFolder = ref('')
const customFolders = ref<string[]>([])
const searchQuery = ref('')
const selectedTypeFilter = ref('all')

// Modals
const showNewFolderModal = ref(false)
const newFolderName = ref('')

const showMoveModal = ref(false)
const fileToMove = ref<IProjectFile | null>(null)
const targetMoveFolder = ref('')
const customMoveFolder = ref('')

// Preview State
const showPreviewModal = ref(false)
const previewFile = ref<IProjectFile | null>(null)
const previewBlobUrl = ref('')
const previewTextContent = ref('')
const previewType = ref<'image' | 'pdf' | 'text' | 'generic'>('generic')
const isLoadingPreview = ref(false)

const project = computed(() => projectStore.projects[props.projectId] ?? baseStore.currentProject)
const canWrite = computed(() => {
	const perm = project.value?.maxPermission ?? baseStore.currentProject?.maxPermission
	if (perm === undefined || perm === null) return true
	return perm > Permissions.READ
})

const typeFilters = [
	{key: 'all', label: 'Todos', icon: 'layer-group'},
	{key: 'pdf', label: 'PDFs', icon: 'file-pdf'},
	{key: 'image', label: 'Imagens', icon: 'file-image'},
	{key: 'doc', label: 'Planilhas & Docs', icon: 'file-alt'},
	{key: 'archive', label: 'Compactados', icon: 'archive'},
]

// Extract clean file name without folder prefix
function getCleanFileName(pf: IProjectFile | null): string {
	if (!pf?.file?.name) return 'Documento sem nome'
	const name = pf.file.name
	if (name.includes('/')) {
		return name.split('/').pop() || name
	}
	return name
}

// Extract folder from filename
function getFileFolder(pf: IProjectFile | null): string {
	if (!pf?.file?.name) return ''
	const name = pf.file.name
	if (name.includes('/')) {
		const parts = name.split('/')
		parts.pop()
		return parts.join('/')
	}
	return ''
}

// LocalStorage helpers for custom empty folders
function getFolderStorageKey() {
	return `dataflow_folders_${props.projectId}`
}

function loadCustomFolders() {
	try {
		const stored = localStorage.getItem(getFolderStorageKey())
		if (stored) {
			customFolders.value = JSON.parse(stored)
		} else {
			// default preset folders for projects
			customFolders.value = ['Contratos', 'Briefings', 'Entregáveis', 'Financeiro & NF']
		}
	} catch {
		customFolders.value = ['Contratos', 'Briefings', 'Entregáveis', 'Financeiro & NF']
	}
}

function saveCustomFolders() {
	try {
		localStorage.setItem(getFolderStorageKey(), JSON.stringify(customFolders.value))
	} catch (e) {
		console.error('Failed to save folders to storage', e)
	}
}

// Calculate all distinct folders with metrics
const allFolders = computed(() => {
	const folderMap = new Map<string, {name: string, count: number, size: number}>()

	// add custom folders first
	for (const cf of customFolders.value) {
		folderMap.set(cf, {name: cf, count: 0, size: 0})
	}

	// add folders from actual files
	for (const f of files.value) {
		const folderName = getFileFolder(f)
		if (folderName) {
			if (!folderMap.has(folderName)) {
				folderMap.set(folderName, {name: folderName, count: 0, size: 0})
			}
			const existing = folderMap.get(folderName)!
			existing.count += 1
			existing.size += (f.file?.size || 0)
		}
	}

	return Array.from(folderMap.values()).sort((a, b) => a.name.localeCompare(b.name))
})

// Filtered files according to folder, search, and type
const filteredFiles = computed(() => {
	return files.value.filter(f => {
		// Folder filter
		const folder = getFileFolder(f)
		if (currentFolder.value !== '' && folder !== currentFolder.value) {
			return false
		}

		// Search filter
		if (searchQuery.value.trim()) {
			const cleanName = getCleanFileName(f).toLowerCase()
			if (!cleanName.includes(searchQuery.value.trim().toLowerCase())) {
				return false
			}
		}

		// Type filter
		if (selectedTypeFilter.value !== 'all') {
			const cat = getFileCategory(getCleanFileName(f))
			if (selectedTypeFilter.value !== cat) {
				return false
			}
		}

		return true
	})
})

const totalFilesSize = computed(() => {
	return files.value.reduce((acc, f) => acc + (f.file?.size || 0), 0)
})

const totalSizeFormatted = computed(() => {
	return getHumanSize(totalFilesSize.value)
})

function getFileCategory(filename: string = ''): string {
	const ext = filename.split('.').pop()?.toLowerCase() || ''
	if (['png', 'jpg', 'jpeg', 'gif', 'svg', 'webp', 'bmp'].includes(ext)) {
		return 'image'
	}
	if (ext === 'pdf') {
		return 'pdf'
	}
	if (['zip', 'tar', 'gz', 'rar', '7z'].includes(ext)) {
		return 'archive'
	}
	if (['xls', 'xlsx', 'csv', 'doc', 'docx', 'ppt', 'pptx', 'txt'].includes(ext)) {
		return 'doc'
	}
	if (['js', 'ts', 'py', 'go', 'json', 'html', 'css', 'sql', 'sh'].includes(ext)) {
		return 'code'
	}
	return 'generic'
}

function getFileIcon(filename: string = ''): string {
	const cat = getFileCategory(filename)
	switch (cat) {
		case 'image': return 'file-image'
		case 'pdf': return 'file-pdf'
		case 'archive': return 'archive'
		case 'doc': return 'file-alt'
		case 'code': return 'file-code'
		default: return 'file'
	}
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

// Drag & Drop
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

// Operations
async function loadProjectFiles() {
	if (!props.projectId || Number(props.projectId) <= 0) return
	try {
		loadCustomFolders()
		files.value = await projectFileService.getAll(new ProjectFileModel({projectId: Number(props.projectId)}))
	} catch (e) {
		console.error('Failed to load project files:', e)
	}
}

async function uploadFiles(fileList: FileList | File[]) {
	if (!props.projectId || Number(props.projectId) <= 0) return
	try {
		const uploaded = await projectFileService.upload(
			new ProjectFileModel({projectId: Number(props.projectId)}),
			fileList,
			currentFolder.value
		)
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
	const cleanName = getCleanFileName(projectFile)
	if (!confirm(`Deseja realmente excluir o documento "${cleanName}"?`)) {
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

// Folders Management
function openNewFolderModal() {
	newFolderName.value = ''
	showNewFolderModal.value = true
}

function createFolder() {
	const name = newFolderName.value.trim().replace(/\/+/g, '')
	if (!name) return
	if (!customFolders.value.includes(name)) {
		customFolders.value.push(name)
		saveCustomFolders()
	}
	currentFolder.value = name
	showNewFolderModal.value = false
	notifySuccess(`Pasta "${name}" criada com sucesso!`)
}

function deleteFolder(folderName: string) {
	const folderFiles = files.value.filter(f => getFileFolder(f) === folderName)
	if (folderFiles.length > 0) {
		alert(`A pasta "${folderName}" contém ${folderFiles.length} documento(s). Mova ou exclua os documentos antes de remover a pasta.`)
		return
	}
	if (!confirm(`Deseja remover a pasta vazia "${folderName}"?`)) {
		return
	}
	customFolders.value = customFolders.value.filter(f => f !== folderName)
	saveCustomFolders()
	if (currentFolder.value === folderName) {
		currentFolder.value = ''
	}
	notifySuccess(`Pasta "${folderName}" removida.`)
}

// Move File to Folder
function openMoveModal(pf: IProjectFile) {
	fileToMove.value = pf
	targetMoveFolder.value = getFileFolder(pf)
	customMoveFolder.value = ''
	showMoveModal.value = true
}

async function confirmMoveFile() {
	if (!fileToMove.value) return
	const target = targetMoveFolder.value.trim().replace(/\/+/g, '')
	try {
		await projectFileService.move(fileToMove.value, target)
		// Update local state
		const baseName = getCleanFileName(fileToMove.value)
		const newFullName = target ? `${target}/${baseName}` : baseName
		if (fileToMove.value.file) {
			fileToMove.value.file.name = newFullName
		}
		if (target && !customFolders.value.includes(target)) {
			customFolders.value.push(target)
			saveCustomFolders()
		}
		showMoveModal.value = false
		notifySuccess(`Documento movido para "${target || 'Raiz'}" com sucesso!`)
	} catch (e) {
		notifyError(e)
	}
}

// Inline Preview
async function openPreview(projectFile: IProjectFile) {
	previewFile.value = projectFile
	const cleanName = getCleanFileName(projectFile)
	const ext = cleanName.split('.').pop()?.toLowerCase() || ''

	isLoadingPreview.value = true
	showPreviewModal.value = true

	try {
		if (['png', 'jpg', 'jpeg', 'gif', 'webp', 'svg'].includes(ext)) {
			previewType.value = 'image'
			const blob = await projectFileService.getBlob(projectFile)
			previewBlobUrl.value = URL.createObjectURL(blob)
		} else if (ext === 'pdf') {
			previewType.value = 'pdf'
			const blob = await projectFileService.getBlob(projectFile)
			const pdfBlob = new Blob([blob], {type: 'application/pdf'})
			previewBlobUrl.value = URL.createObjectURL(pdfBlob)
		} else if (['txt', 'json', 'csv', 'sql', 'js', 'ts', 'py', 'sh', 'md'].includes(ext)) {
			previewType.value = 'text'
			const blob = await projectFileService.getBlob(projectFile)
			previewTextContent.value = await blob.text()
		} else {
			previewType.value = 'generic'
		}
	} catch (e) {
		console.error('Failed to load preview blob', e)
		previewType.value = 'generic'
	} finally {
		isLoadingPreview.value = false
	}
}

function closePreview() {
	showPreviewModal.value = false
	if (previewBlobUrl.value) {
		URL.revokeObjectURL(previewBlobUrl.value)
		previewBlobUrl.value = ''
	}
	previewFile.value = null
	previewTextContent.value = ''
}

onMounted(loadProjectFiles)
watch(() => props.projectId, loadProjectFiles)
</script>

<style lang="scss" scoped>
.project-documents-card {
	padding: 1.5rem;
	border-radius: $radius;
	border: 1px solid var(--grey-200);
	background: var(--card-background);

	input[type='file'] {
		display: none;
	}
}

.header-section {
	display: flex;
	justify-content: space-between;
	align-items: flex-start;
	gap: 1.5rem;
	flex-wrap: wrap;
}

.title-with-badge {
	display: flex;
	align-items: center;
	gap: 1rem;
	flex-wrap: wrap;

	h2 {
		margin: 0;
		font-size: 1.35rem;
		font-weight: 700;
	}
}

.stats-pills {
	display: flex;
	gap: .5rem;
	flex-wrap: wrap;

	.stat-pill {
		display: inline-flex;
		align-items: center;
		gap: .35rem;
		padding: .2rem .65rem;
		background: var(--grey-100);
		color: var(--grey-700);
		border-radius: 9999px;
		font-size: .75rem;
		font-weight: 600;

		.stat-icon {
			font-size: .7rem;
			color: var(--primary);
		}
	}
}

.header-actions {
	display: flex;
	align-items: center;
	gap: .75rem;

	.action-header-btn {
		display: inline-flex;
		align-items: center;
		gap: .4rem;
		font-weight: 600;
		font-size: .85rem;
		padding: .5rem 1rem;
		border-radius: $radius;
		cursor: pointer;
		transition: all .2s ease;
	}

	.new-folder-btn {
		border: 1px solid var(--grey-300);
		background: var(--card-background);
		color: var(--grey-800);

		&:hover {
			border-color: var(--primary);
			color: var(--primary);
			background: rgba(25, 115, 255, 0.04);
		}
	}
}

/* Navigation bar: Breadcrumbs & Search */
.navigation-bar {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 1rem;
	flex-wrap: wrap;
	padding: .75rem 1rem;
	background: var(--grey-50);
	border: 1px solid var(--grey-200);
	border-radius: $radius;
}

.breadcrumbs {
	display: flex;
	align-items: center;
	gap: .4rem;
	font-size: .85rem;

	.crumb-btn {
		display: inline-flex;
		align-items: center;
		gap: .35rem;
		background: none;
		border: none;
		color: var(--grey-600);
		font-weight: 500;
		cursor: pointer;
		padding: .2rem .4rem;
		border-radius: 4px;
		transition: all .15s ease;

		&:hover {
			color: var(--primary);
			background: var(--grey-200);
		}

		&.is-active {
			color: var(--primary);
			font-weight: 700;
		}

		.crumb-folder-icon {
			color: #f59e0b;
		}
	}

	.crumb-separator {
		color: var(--grey-400);
		font-size: .8rem;
	}
}

.search-box {
	position: relative;
	display: flex;
	align-items: center;

	.search-icon {
		position: absolute;
		left: .75rem;
		color: var(--grey-400);
		font-size: .8rem;
		pointer-events: none;
	}

	.search-input {
		padding: .4rem .75rem .4rem 2rem;
		border: 1px solid var(--grey-300);
		border-radius: 9999px;
		font-size: .8rem;
		background: var(--card-background);
		width: 220px;
		transition: all .2s ease;

		&:focus {
			outline: none;
			border-color: var(--primary);
			box-shadow: 0 0 0 2px rgba(25, 115, 255, 0.15);
			width: 260px;
		}
	}

	.clear-search-btn {
		position: absolute;
		right: .5rem;
		background: none;
		border: none;
		color: var(--grey-400);
		cursor: pointer;
		font-size: .75rem;

		&:hover {
			color: var(--grey-700);
		}
	}
}

/* Folder Chips bar */
.folder-chips-bar {
	display: flex;
	gap: .5rem;
	overflow-x: auto;
	padding-bottom: .25rem;

	.folder-chip {
		display: inline-flex;
		align-items: center;
		gap: .4rem;
		padding: .35rem .75rem;
		border-radius: 9999px;
		border: 1px solid var(--grey-200);
		background: var(--card-background);
		font-size: .8rem;
		font-weight: 500;
		color: var(--grey-700);
		cursor: pointer;
		white-space: nowrap;
		transition: all .2s ease;

		&:hover {
			border-color: var(--primary);
			color: var(--primary);
		}

		&.is-selected {
			background: var(--primary);
			color: white;
			border-color: var(--primary);

			.chip-count {
				background: rgba(255, 255, 255, 0.25);
				color: white;
			}
		}

		.chip-count {
			padding: .1rem .45rem;
			background: var(--grey-100);
			color: var(--grey-600);
			border-radius: 9999px;
			font-size: .7rem;
			font-weight: 700;
		}

		&.add-chip {
			border-style: dashed;
			color: var(--primary);

			&:hover {
				background: rgba(25, 115, 255, 0.05);
			}
		}
	}
}

/* Folders Grid Cards */
.folders-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
	gap: 1rem;
}

.folder-card {
	padding: 1rem;
	background: var(--grey-50);
	border: 1px solid var(--grey-200);
	border-radius: $radius;
	cursor: pointer;
	transition: all .2s ease;
	display: flex;
	flex-direction: column;
	gap: .75rem;

	&:hover {
		border-color: var(--primary);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
		background: var(--card-background);
	}

	.folder-card-top {
		display: flex;
		justify-content: space-between;
		align-items: center;

		.folder-card-icon-wrapper {
			width: 36px;
			height: 36px;
			border-radius: 8px;
			background: rgba(245, 158, 11, 0.12);
			display: flex;
			align-items: center;
			justify-content: center;

			.folder-card-icon {
				color: #f59e0b;
				font-size: 1.25rem;
			}
		}

		.folder-menu-btn {
			background: none;
			border: none;
			color: var(--grey-400);
			padding: .25rem;
			cursor: pointer;
			border-radius: 4px;
			font-size: .8rem;

			&:hover {
				color: #e53e3e;
				background: rgba(229, 62, 62, 0.1);
			}
		}
	}

	.folder-card-info {
		display: flex;
		flex-direction: column;
		gap: .2rem;

		.folder-card-title {
			font-weight: 600;
			font-size: .9rem;
			color: var(--grey-900);
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}

		.folder-card-subtitle {
			font-size: .75rem;
			color: var(--grey-500);
		}
	}
}

/* Drop Zone */
.drop-zone {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: .75rem;
	padding: 1.75rem 1.5rem;
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
		font-size: 1.85rem;
		color: var(--primary);
	}

	.drop-zone-text {
		display: flex;
		flex-direction: column;
		gap: .25rem;
	}
}

/* Type Filter Row */
.type-filter-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 1rem;
	flex-wrap: wrap;

	.type-filter-chips {
		display: flex;
		gap: .4rem;
		flex-wrap: wrap;

		.type-filter-btn {
			display: inline-flex;
			align-items: center;
			gap: .35rem;
			padding: .25rem .65rem;
			border-radius: 6px;
			border: 1px solid var(--grey-200);
			background: var(--card-background);
			font-size: .75rem;
			font-weight: 500;
			color: var(--grey-600);
			cursor: pointer;
			transition: all .15s ease;

			&:hover {
				border-color: var(--grey-400);
			}

			&.is-active {
				background: var(--grey-800);
				color: white;
				border-color: var(--grey-800);
			}
		}
	}

	.folder-active-indicator {
		display: inline-flex;
		align-items: center;
		gap: .4rem;
		font-size: .75rem;
		color: var(--grey-600);

		.clear-folder-btn {
			background: none;
			border: none;
			color: var(--grey-400);
			cursor: pointer;

			&:hover {
				color: var(--grey-800);
			}
		}
	}
}

/* Document List */
.project-documents-list {
	display: flex;
	flex-direction: column;
	gap: .65rem;
}

.project-document-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 1rem;
	padding: .85rem 1.15rem;
	border: 1px solid var(--grey-200);
	border-radius: $radius;
	background: var(--card-background);
	transition: all .15s ease;

	&:hover {
		box-shadow: 0 3px 10px rgba(0, 0, 0, 0.04);
		border-color: var(--grey-300);
	}

	.document-info {
		display: flex;
		align-items: center;
		gap: 1rem;
		min-width: 0;
		cursor: pointer;
		flex: 1;

		.doc-icon-container {
			width: 40px;
			height: 40px;
			border-radius: 8px;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 1.25rem;
			flex-shrink: 0;

			&.ext-image { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
			&.ext-pdf { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
			&.ext-doc { background: rgba(16, 185, 129, 0.1); color: #10b981; }
			&.ext-archive { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
			&.ext-code { background: rgba(139, 92, 246, 0.1); color: #8b5cf6; }
			&.ext-generic { background: rgba(107, 114, 128, 0.1); color: #6b7280; }
		}

		.doc-meta {
			display: flex;
			flex-direction: column;
			gap: .25rem;
			min-width: 0;

			.doc-title-row {
				display: flex;
				align-items: center;
				gap: .6rem;
				flex-wrap: wrap;

				.doc-name {
					font-weight: 600;
					font-size: .9rem;
					white-space: nowrap;
					overflow: hidden;
					text-overflow: ellipsis;
					color: var(--grey-900);

					&:hover {
						color: var(--primary);
					}
				}

				.doc-folder-tag {
					display: inline-flex;
					align-items: center;
					gap: .25rem;
					padding: .1rem .45rem;
					background: rgba(245, 158, 11, 0.1);
					color: #d97706;
					border-radius: 4px;
					font-size: .7rem;
					font-weight: 600;

					.tag-folder-icon {
						font-size: .65rem;
					}
				}
			}

			.doc-details {
				display: flex;
				align-items: center;
				gap: .4rem;

				.detail-sep {
					opacity: .5;
				}

				.doc-author {
					display: inline-flex;
					align-items: center;
					gap: .2rem;

					.author-icon {
						font-size: .65rem;
					}
				}
			}
		}
	}

	.project-document-actions {
		display: flex;
		align-items: center;
		gap: .4rem;
		flex-shrink: 0;

		.action-btn {
			padding: .4rem .65rem;
			border-radius: 6px;
			border: 1px solid var(--grey-200);
			background: transparent;
			cursor: pointer;
			font-size: .85rem;
			display: inline-flex;
			align-items: center;
			gap: .35rem;
			transition: all .2s ease;

			&.preview-btn {
				background: rgba(25, 115, 255, 0.05);
				border-color: rgba(25, 115, 255, 0.2);
				color: var(--primary);
				font-weight: 600;

				&:hover {
					background: var(--primary);
					color: white;
					border-color: var(--primary);
				}
			}

			&.download-btn:hover {
				background: var(--grey-100);
				color: var(--primary);
				border-color: var(--primary);
			}

			&.move-btn:hover {
				background: var(--grey-100);
				color: #f59e0b;
				border-color: #f59e0b;
			}

			&.delete-btn:hover {
				background: #e53e3e;
				color: white;
				border-color: #e53e3e;
			}
		}
	}
}

/* Modals */
.custom-modal-overlay {
	position: fixed;
	inset: 0;
	background: rgba(0, 0, 0, 0.55);
	backdrop-filter: blur(4px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 9999;
	padding: 1rem;
}

.custom-modal-content {
	background: var(--card-background);
	border: 1px solid var(--grey-300);
	border-radius: 12px;
	width: 100%;
	max-width: 480px;
	box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
	overflow: hidden;
	display: flex;
	flex-direction: column;
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 1.15rem 1.5rem;
	border-bottom: 1px solid var(--grey-200);

	.modal-header-title {
		display: flex;
		align-items: center;
		gap: .5rem;

		.modal-title-icon {
			color: var(--primary);
			font-size: 1.15rem;
		}

		h3 {
			margin: 0;
			font-size: 1.1rem;
			font-weight: 700;
		}
	}

	.modal-close-btn {
		background: none;
		border: none;
		color: var(--grey-500);
		font-size: 1.1rem;
		cursor: pointer;
		padding: .25rem;

		&:hover {
			color: var(--grey-900);
		}
	}
}

.modal-body {
	padding: 1.5rem;

	.modal-label {
		display: block;
		font-size: .85rem;
		font-weight: 600;
		margin-bottom: .4rem;
	}

	.modal-input {
		width: 100%;
		padding: .65rem .85rem;
		border: 1px solid var(--grey-300);
		border-radius: 6px;
		font-size: .9rem;

		&:focus {
			outline: none;
			border-color: var(--primary);
			box-shadow: 0 0 0 2px rgba(25, 115, 255, 0.15);
		}
	}

	.folder-presets {
		.preset-label {
			display: block;
			font-size: .75rem;
			color: var(--grey-500);
			margin-bottom: .35rem;
		}

		.preset-chips {
			display: flex;
			gap: .35rem;
			flex-wrap: wrap;

			.preset-chip {
				padding: .25rem .5rem;
				background: var(--grey-100);
				border: 1px solid var(--grey-200);
				border-radius: 4px;
				font-size: .75rem;
				color: var(--grey-700);
				cursor: pointer;

				&:hover {
					border-color: var(--primary);
					color: var(--primary);
					background: rgba(25, 115, 255, 0.05);
				}
			}
		}
	}

	.folder-select-list {
		display: flex;
		flex-direction: column;
		gap: .35rem;
		max-height: 200px;
		overflow-y: auto;

		.folder-select-item {
			display: flex;
			align-items: center;
			gap: .5rem;
			padding: .6rem .85rem;
			border: 1px solid var(--grey-200);
			border-radius: 6px;
			background: var(--grey-50);
			cursor: pointer;
			text-align: left;
			font-size: .85rem;
			font-weight: 500;
			transition: all .15s ease;

			&:hover {
				border-color: var(--primary);
				background: rgba(25, 115, 255, 0.05);
			}

			&.is-selected {
				border-color: var(--primary);
				background: rgba(25, 115, 255, 0.1);
				color: var(--primary);
				font-weight: 700;
			}
		}
	}

	.or-divider {
		display: flex;
		align-items: center;
		justify-content: center;
		margin: .75rem 0 .5rem 0;
		color: var(--grey-400);
		font-size: .75rem;
		text-transform: uppercase;
		letter-spacing: .5px;
	}
}

.modal-footer {
	display: flex;
	justify-content: flex-end;
	gap: .75rem;
	padding: 1rem 1.5rem;
	border-top: 1px solid var(--grey-200);
	background: var(--grey-50);

	.modal-cancel-btn {
		background: transparent;
		border: 1px solid var(--grey-300);
		color: var(--grey-700);
	}
}

/* Preview Modal */
.preview-modal-content {
	max-width: 900px;
	max-height: 90vh;
	height: 85vh;
}

.preview-header {
	.preview-title-info {
		display: flex;
		align-items: center;
		gap: .75rem;

		.preview-header-icon {
			font-size: 1.5rem;
			color: var(--primary);
		}

		.preview-file-name {
			margin: 0;
			font-size: 1rem;
			font-weight: 600;
		}
	}

	.preview-header-actions {
		display: flex;
		align-items: center;
		gap: .5rem;
	}
}

.preview-body {
	flex: 1;
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: auto;
	background: #0f172a;
	padding: 1rem;

	.preview-loader {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: .75rem;
		color: white;

		.preview-spinner {
			font-size: 2rem;
		}
	}

	.image-preview-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;

		.preview-img {
			max-width: 100%;
			max-height: 100%;
			object-fit: contain;
			border-radius: 4px;
			box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
		}
	}

	.pdf-preview-wrapper {
		width: 100%;
		height: 100%;

		.pdf-iframe {
			width: 100%;
			height: 100%;
			border: none;
			border-radius: 4px;
		}
	}

	.text-preview-wrapper {
		width: 100%;
		height: 100%;
		overflow: auto;

		.code-view {
			background: #1e293b;
			color: #e2e8f0;
			padding: 1rem;
			border-radius: 6px;
			font-family: monospace;
			font-size: .85rem;
			white-space: pre-wrap;
			margin: 0;
		}
	}

	.generic-preview-wrapper {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
		color: white;
		gap: .5rem;

		.generic-preview-icon {
			font-size: 3rem;
			color: #94a3b8;
			margin-bottom: .5rem;
		}
	}
}
</style>
