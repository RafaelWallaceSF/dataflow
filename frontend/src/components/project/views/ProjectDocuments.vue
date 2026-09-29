<template>
	<ProjectWrapper
		class="project-documents"
		:is-loading-project="isLoadingProject"
		:project-id="projectId"
		:view-id="viewId"
	>
		<div class="loader-container is-max-width-desktop" :class="{'is-loading': projectFileService.loading}">
			<Card class="project-documents-card">
				<h2>{{ t('project.documents.title') }}</h2>
				<p class="has-text-grey mbe-4">
					{{ t('project.documents.description') }}
				</p>

				<input
					v-if="canWrite"
					ref="filesRef"
					multiple
					type="file"
					:disabled="projectFileService.loading || undefined"
					@change="uploadProjectFiles"
				>

				<ProgressBar
					v-if="projectFileService.uploadProgress > 0"
					:value="projectFileService.uploadProgress"
					is-primary
				/>

				<XButton
					v-if="canWrite"
					:disabled="projectFileService.loading"
					class="mbe-4"
					icon="cloud-upload-alt"
					variant="secondary"
					:shadow="false"
					@click="filesRef?.click()"
				>
					{{ t('project.documents.upload') }}
				</XButton>

				<Nothing v-if="files.length === 0 && !projectFileService.loading">
					{{ t('project.documents.empty') }}
				</Nothing>

				<div v-else class="project-documents-list">
					<div
						v-for="projectFile in files"
						:key="projectFile.id"
						class="project-document-row"
					>
						<div>
							<strong>{{ projectFile.file.name }}</strong>
							<p class="has-text-grey is-size-7">
								{{ getHumanSize(projectFile.file.size) }}
							</p>
						</div>
						<div class="project-document-actions">
							<BaseButton
								:aria-label="t('project.documents.download')"
								@click="downloadProjectFile(projectFile)"
							>
								<Icon icon="download" />
							</BaseButton>
							<BaseButton
								v-if="canWrite"
								:aria-label="t('project.documents.delete')"
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
import {PERMISSIONS as Permissions} from '@/constants/permissions'
import {getHumanSize} from '@/helpers/getHumanSize'

const props = defineProps<{
	isLoadingProject: boolean,
	projectId: IProject['id'],
	viewId: IProjectView['id'],
}>()

const {t} = useI18n({useScope: 'global'})
const projectStore = useProjectStore()
const projectFileService = shallowReactive(new ProjectFileService())
const filesRef = ref<HTMLInputElement | null>(null)
const files = ref<IProjectFile[]>([])

const project = computed(() => projectStore.projects[props.projectId])
const canWrite = computed(() => project.value?.maxPermission > Permissions.READ && project.value?.id > 0)

async function loadProjectFiles() {
	files.value = await projectFileService.getAll(new ProjectFileModel({projectId: props.projectId}))
}

async function uploadProjectFiles() {
	if (!filesRef.value?.files?.length) {
		return
	}

	const uploaded = await projectFileService.upload(new ProjectFileModel({projectId: props.projectId}), filesRef.value.files)
	files.value = [
		...files.value,
		...(uploaded.success ?? []),
	]
	filesRef.value.value = ''
}

async function downloadProjectFile(projectFile: IProjectFile) {
	await projectFileService.download(projectFile)
}

async function deleteProjectFile(projectFile: IProjectFile) {
	await projectFileService.delete(projectFile)
	files.value = files.value.filter(({id}) => id !== projectFile.id)
}

onMounted(loadProjectFiles)
watch(() => props.projectId, loadProjectFiles)
</script>

<style lang="scss" scoped>
.project-documents-card {
	input[type='file'] {
		display: none;
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
	padding: .75rem;
	border: 1px solid var(--grey-200);
	border-radius: $radius;
}

.project-document-actions {
	display: flex;
	gap: .5rem;
}
</style>
