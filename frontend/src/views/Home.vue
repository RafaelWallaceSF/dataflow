<template>
	<div class="dataflow-overview-wrapper">
		<!-- HEADER SECTION -->
		<header class="overview-header">
			<div class="header-left">
				<h1 class="welcome-title">
					Olá<span v-if="displayName">, <span class="highlight-user">{{ displayName }}</span></span>! 👋
				</h1>
				<p class="welcome-subtitle">
					Aqui está o panorama das suas demandas.
				</p>
			</div>

			<div class="header-right">
				<div class="current-date-pill">
					<Icon :icon="['far', 'calendar-alt']" class="date-icon" />
					<span>{{ formattedCurrentDate }}</span>
				</div>
				<button
					type="button"
					class="btn-primary-action"
					@click="openNewDemandModal"
				>
					<Icon icon="plus" class="btn-icon" />
					<span>Nova demanda</span>
				</button>
			</div>
		</header>

		<!-- 4 INDICATOR METRIC CARDS -->
		<section class="metrics-grid">
			<!-- ATRASADAS -->
			<div class="metric-card is-danger">
				<div class="metric-header">
					<span class="metric-tag">ATRASADAS</span>
					<div class="metric-icon-wrap danger-bg">
						<Icon icon="exclamation-circle" />
					</div>
				</div>
				<div class="metric-value">{{ metrics.overdue }}</div>
				<div class="metric-footer">
					<span v-if="metrics.overdue > 0" class="metric-subtext danger-text">
						Demandas com prazo vencido
					</span>
					<span v-else class="metric-subtext success-text">
						Nenhuma demanda em atraso
					</span>
				</div>
			</div>

			<!-- EM ANDAMENTO -->
			<div class="metric-card is-primary">
				<div class="metric-header">
					<span class="metric-tag">EM ANDAMENTO</span>
					<div class="metric-icon-wrap primary-bg">
						<Icon icon="spinner" />
					</div>
				</div>
				<div class="metric-value">{{ metrics.inProgress }}</div>
				<div class="metric-footer">
					<span class="metric-subtext">
						Em execução ativa
					</span>
				</div>
			</div>

			<!-- A FAZER -->
			<div class="metric-card is-warning">
				<div class="metric-header">
					<span class="metric-tag">A FAZER</span>
					<div class="metric-icon-wrap warning-bg">
						<Icon icon="clipboard-list" />
					</div>
				</div>
				<div class="metric-value">{{ metrics.todo }}</div>
				<div class="metric-footer">
					<span class="metric-subtext">
						Aguardando início
					</span>
				</div>
			</div>

			<!-- CONCLUÍDAS -->
			<div class="metric-card is-success">
				<div class="metric-header">
					<span class="metric-tag">CONCLUÍDAS</span>
					<div class="metric-icon-wrap success-bg">
						<Icon icon="check-circle" />
					</div>
				</div>
				<div class="metric-value">{{ metrics.done }}</div>
				<div class="metric-footer">
					<span class="metric-subtext success-text">
						Finalizadas com sucesso
					</span>
				</div>
			</div>
		</section>

		<!-- MAIN SECTION (65% / 35%) -->
		<div class="main-dashboard-grid">
			<!-- LEFT: PRÓXIMAS ENTREGAS (65%) -->
			<section class="dashboard-panel deliveries-panel">
				<div class="panel-header">
					<div class="panel-title-area">
						<h2 class="panel-title">Próximas entregas</h2>
						<span class="panel-badge">{{ upcomingDeliveries.length }}</span>
					</div>
					<RouterLink :to="{ name: 'tasks.range' }" class="panel-link-action">
						<span>Ver todas</span>
						<Icon icon="arrow-right" class="arrow-icon" />
					</RouterLink>
				</div>

				<div v-if="isLoading" class="panel-loading">
					<div class="loading-spinner" />
					<span>Carregando demandas...</span>
				</div>

				<div v-else-if="upcomingDeliveries.length === 0" class="panel-empty-state">
					<div class="empty-icon-wrap">
						<Icon icon="check-double" />
					</div>
					<h3 class="empty-title">Tudo em dia por aqui!</h3>
					<p class="empty-description">
						Nenhuma demanda pendente ou com prazo nos próximos dias.
					</p>
					<button
						type="button"
						class="btn-secondary-action"
						@click="openNewDemandModal"
					>
						<Icon icon="plus" />
						<span>Criar primeira demanda</span>
					</button>
				</div>

				<div v-else class="table-container">
					<table class="dataflow-table">
						<thead>
							<tr>
								<th class="col-checkbox" />
								<th class="col-task">Demanda</th>
								<th class="col-project">Projeto</th>
								<th class="col-due">Prazo</th>
								<th class="col-status">Status</th>
							</tr>
						</thead>
						<tbody>
							<tr
								v-for="task in upcomingDeliveries"
								:key="task.id"
								class="task-row"
								:class="{'is-completed': task.done}"
							>
								<!-- Checkbox -->
								<td class="col-checkbox">
									<button
										type="button"
										class="task-check-button"
										:class="{'is-checked': task.done}"
										:aria-label="task.done ? 'Marcar como não concluída' : 'Marcar como concluída'"
										@click.stop="toggleTaskStatus(task)"
									>
										<Icon v-if="task.done" icon="check" class="check-icon" />
									</button>
								</td>

								<!-- Demanda Title + Labels -->
								<td class="col-task">
									<div class="task-info">
										<RouterLink
											:to="{ name: 'task.detail', params: { id: task.id } }"
											class="task-title-link"
										>
											{{ task.title }}
										</RouterLink>
										<div v-if="task.labels && task.labels.length" class="task-labels-row">
											<span
												v-for="label in task.labels"
												:key="label.id"
												class="task-label-chip"
												:style="getLabelStyle(label)"
											>
												{{ label.title }}
											</span>
										</div>
									</div>
								</td>

								<!-- Projeto -->
								<td class="col-project">
									<RouterLink
										v-if="getProject(task.projectId)"
										:to="{ name: 'project.index', params: { projectId: task.projectId } }"
										class="project-pill"
									>
										<span
											class="project-dot"
											:style="{ background: getProjectColor(task.projectId) }"
										/>
										<span class="project-name">{{ getProject(task.projectId)?.title }}</span>
									</RouterLink>
									<span v-else class="project-pill is-empty">—</span>
								</td>

								<!-- Prazo -->
								<td class="col-due">
									<span
										class="due-date-badge"
										:class="getDueDateClass(task)"
									>
										<Icon :icon="getDueDateIcon(task)" class="due-icon" />
										{{ formatDueBadgeText(task) }}
									</span>
								</td>

								<!-- Status -->
								<td class="col-status">
									<span class="status-pill" :class="getTaskStatusClass(task)">
										{{ getTaskStatusLabel(task) }}
									</span>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</section>

			<!-- RIGHT: DEMANDAS POR PROJETO (35%) -->
			<section class="dashboard-panel chart-panel">
				<div class="panel-header">
					<div class="panel-title-area">
						<h2 class="panel-title">Demandas por projeto</h2>
						<span class="panel-subtitle-text">Distribuição da carga</span>
					</div>
				</div>

				<!-- Donut Chart -->
				<div class="donut-chart-container">
					<div class="donut-svg-wrapper">
						<svg
							viewBox="0 0 160 160"
							class="donut-svg"
							width="160"
							height="160"
						>
							<!-- Background Circle -->
							<circle
								cx="80"
								cy="80"
								r="60"
								class="donut-bg-ring"
								fill="none"
								stroke="#E7EAF0"
								stroke-width="18"
							/>

							<!-- Segment Slices -->
							<circle
								v-for="(seg, idx) in donutSegments"
								:key="idx"
								cx="80"
								cy="80"
								r="60"
								fill="none"
								:stroke="seg.color"
								stroke-width="18"
								:stroke-dasharray="seg.dasharray"
								:stroke-dashoffset="seg.dashoffset"
								class="donut-segment"
								transform="rotate(-90 80 80)"
							/>
						</svg>

						<!-- Central Label -->
						<div class="donut-center-content">
							<span class="donut-center-number">{{ metrics.total }}</span>
							<span class="donut-center-label">Demandas</span>
						</div>
					</div>

					<!-- Legend List -->
					<div class="donut-legend-list">
						<div
							v-for="(projItem, idx) in projectDistribution"
							:key="idx"
							class="legend-item"
						>
							<div class="legend-meta">
								<span
									class="legend-color-dot"
									:style="{ background: projItem.color }"
								/>
								<span class="legend-name" :title="projItem.title">
									{{ projItem.title }}
								</span>
							</div>
							<div class="legend-counts">
								<span class="legend-count">{{ projItem.count }}</span>
								<span class="legend-pct">({{ projItem.percentage }}%)</span>
							</div>
							<div class="legend-bar-track">
								<div
									class="legend-bar-fill"
									:style="{
										width: projItem.percentage + '%',
										background: projItem.color
									}"
								/>
							</div>
						</div>

						<div v-if="projectDistribution.length === 0" class="donut-legend-empty">
							Nenhum projeto com demandas ativas.
						</div>
					</div>
				</div>
			</section>
		</div>

		<!-- SECOND ROW: MEUS PROJETOS -->
		<section class="projects-overview-section">
			<div class="section-heading">
				<div class="heading-text">
					<h2 class="section-title">Meus projetos</h2>
					<p class="section-subtitle">Acompanhe o andamento geral e entregas por projeto</p>
				</div>
				<RouterLink :to="{ name: 'projects.index' }" class="heading-action-link">
					<span>Ver todos os projetos</span>
					<Icon icon="arrow-right" />
				</RouterLink>
			</div>

			<div class="project-cards-grid">
				<div
					v-for="proj in activeProjectsList"
					:key="proj.id"
					class="project-summary-card"
				>
					<div
						class="card-accent-bar"
						:style="{ background: proj.hexColor || '#3157F6' }"
					/>
					<div class="card-inner">
						<div class="project-card-header">
							<RouterLink
								:to="{ name: 'project.index', params: { projectId: proj.id } }"
								class="project-card-title"
							>
								{{ proj.title }}
							</RouterLink>
						</div>

						<div class="project-metrics-row">
							<div class="proj-metric">
								<span class="proj-metric-val">{{ proj.openTasksCount }}</span>
								<span class="proj-metric-label">abertas</span>
							</div>
							<div class="proj-metric">
								<span
									class="proj-metric-val"
									:class="{'text-danger': proj.overdueTasksCount > 0}"
								>
									{{ proj.overdueTasksCount }}
								</span>
								<span class="proj-metric-label">atrasadas</span>
							</div>
							<div class="proj-metric">
								<span class="proj-metric-val text-success">{{ proj.doneTasksCount }}</span>
								<span class="proj-metric-label">concluídas</span>
							</div>
						</div>

						<!-- Progress Bar -->
						<div class="project-progress-wrap">
							<div class="progress-info-row">
								<span class="progress-title">Progresso</span>
								<span class="progress-percent">{{ proj.progressPercent }}%</span>
							</div>
							<div class="progress-track">
								<div
									class="progress-fill"
									:style="{
										width: proj.progressPercent + '%',
										background: proj.progressPercent === 100 ? '#12B76A' : (proj.hexColor || '#3157F6')
									}"
								/>
							</div>
						</div>
					</div>
				</div>

				<div v-if="activeProjectsList.length === 0" class="no-projects-card">
					<Icon icon="layer-group" class="no-proj-icon" />
					<p>Nenhum projeto cadastrado no momento.</p>
					<RouterLink :to="{ name: 'project.create' }" class="btn-create-proj">
						<Icon icon="plus" />
						<span>Criar primeiro projeto</span>
					</RouterLink>
				</div>
			</div>
		</section>

		<!-- QUICK ADD DEMAND MODAL -->
		<Transition name="fade-modal">
			<div
				v-if="showNewDemandModal"
				class="dataflow-modal-backdrop"
				@click.self="closeNewDemandModal"
			>
				<div class="dataflow-modal-box">
					<div class="modal-header">
						<div class="modal-title-wrap">
							<div class="modal-badge-icon">
								<Icon icon="plus" />
							</div>
							<h3 class="modal-title">Nova Demanda</h3>
						</div>
						<button
							type="button"
							class="modal-close-btn"
							aria-label="Fechar"
							@click="closeNewDemandModal"
						>
							<Icon icon="times" />
						</button>
					</div>

					<form class="modal-form" @submit.prevent="submitNewDemand">
						<!-- Title -->
						<div class="form-group">
							<label for="demand-title" class="form-label">
								Título da demanda <span class="required">*</span>
							</label>
							<input
								id="demand-title"
								v-model="newDemandForm.title"
								type="text"
								class="form-input"
								placeholder="Ex: Atualizar relatório financeiro mensal"
								required
								autofocus
							>
						</div>

						<!-- Project -->
						<div class="form-group">
							<label for="demand-project" class="form-label">
								Projeto de destino <span class="required">*</span>
							</label>
							<select
								id="demand-project"
								v-model="newDemandForm.projectId"
								class="form-select"
								required
							>
								<option
									v-for="p in availableProjects"
									:key="p.id"
									:value="p.id"
								>
									{{ p.title }}
								</option>
							</select>
						</div>

						<!-- Due Date Quick Presets & Input -->
						<div class="form-group">
							<label class="form-label">Prazo de entrega</label>
							<div class="preset-buttons-row">
								<button
									type="button"
									class="preset-btn"
									:class="{'is-active': selectedPreset === 'today'}"
									@click="setPresetDate('today')"
								>
									Hoje
								</button>
								<button
									type="button"
									class="preset-btn"
									:class="{'is-active': selectedPreset === 'tomorrow'}"
									@click="setPresetDate('tomorrow')"
								>
									Amanhã
								</button>
								<button
									type="button"
									class="preset-btn"
									:class="{'is-active': selectedPreset === 'next-week'}"
									@click="setPresetDate('next-week')"
								>
									Em 1 semana
								</button>
								<button
									type="button"
									class="preset-btn"
									:class="{'is-active': selectedPreset === 'none'}"
									@click="setPresetDate('none')"
								>
									Sem data
								</button>
							</div>
							<input
								v-model="newDemandForm.dueDate"
								type="date"
								class="form-input m-top-xs"
							>
						</div>

						<div v-if="modalError" class="modal-error-alert">
							{{ modalError }}
						</div>

						<!-- Modal Actions -->
						<div class="modal-actions">
							<button
								type="button"
								class="btn-cancel"
								@click="closeNewDemandModal"
							>
								Cancelar
							</button>
							<button
								type="submit"
								class="btn-submit"
								:disabled="isSubmitting || !newDemandForm.title.trim()"
							>
								<Icon v-if="isSubmitting" icon="spinner" class="spin" />
								<span>{{ isSubmitting ? 'Salvando...' : 'Criar demanda' }}</span>
							</button>
						</div>
					</form>
				</div>
			</div>
		</Transition>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'

import Icon from '@/components/misc/Icon'
import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/projects'
import { useTaskStore } from '@/stores/tasks'
import TaskService from '@/services/task'
import type { ITask } from '@/modelTypes/ITask'
import type { IProject } from '@/modelTypes/IProject'
import type { ILabel } from '@/modelTypes/ILabel'

const router = useRouter()
const authStore = useAuthStore()
const projectStore = useProjectStore()
const taskStore = useTaskStore()
const taskService = new TaskService()

// State
const isLoading = ref(true)
const isSubmitting = ref(false)
const showNewDemandModal = ref(false)
const modalError = ref('')
const selectedPreset = ref<'today' | 'tomorrow' | 'next-week' | 'none' | 'custom'>('none')

// All active tasks loaded for dashboard computation
const allTasks = ref<ITask[]>([])

// New demand form
const newDemandForm = reactive({
	title: '',
	projectId: 0,
	dueDate: '',
})

// Current user greeting
const displayName = computed(() => {
	const name = authStore.userDisplayName || authStore.user?.name || authStore.user?.username
	return name || ''
})

// Current date formatted in Portuguese
const formattedCurrentDate = computed(() => {
	const now = new Date()
	const formatted = now.toLocaleDateString('pt-BR', {
		weekday: 'long',
		day: 'numeric',
		month: 'long',
		year: 'numeric',
	})
	// Capitalize first letter
	return formatted.charAt(0).toUpperCase() + formatted.slice(1)
})

// Available projects list
const availableProjects = computed(() => {
	return Object.values(projectStore.projects).filter(p => !p.isArchived && p.id > 0)
})

// Helper to find project
function getProject(projectId: number): IProject | undefined {
	return projectStore.projects[projectId]
}

function getProjectColor(projectId: number): string {
	const proj = projectStore.projects[projectId]
	return (proj && proj.hexColor) ? proj.hexColor : '#3157F6'
}

// Compute Metrics
const metrics = computed(() => {
	const now = new Date()
	let overdue = 0
	let inProgress = 0
	let todo = 0
	let done = 0

	for (const t of allTasks.value) {
		if (t.done) {
			done++
		} else {
			// Check if overdue
			if (t.dueDate) {
				const due = new Date(t.dueDate)
				if (due < now) {
					overdue++
				}
			}

			// Check status
			const bucketTitle = t.bucket?.title?.toLowerCase() || ''
			if (
				bucketTitle.includes('andamento') ||
				bucketTitle.includes('doing') ||
				bucketTitle.includes('progresso') ||
				(t.percentDone && t.percentDone > 0)
			) {
				inProgress++
			} else {
				todo++
			}
		}
	}

	return {
		overdue,
		inProgress,
		todo,
		done,
		total: allTasks.value.length,
	}
})

// Upcoming Deliveries (sorted by due_date ASC, pending first)
const upcomingDeliveries = computed(() => {
	const pending = allTasks.value.filter(t => !t.done)

	return pending
		.slice()
		.sort((a, b) => {
			if (!a.dueDate && !b.dueDate) return b.id - a.id
			if (!a.dueDate) return 1
			if (!b.dueDate) return -1
			return new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime()
		})
		.slice(0, 15) // Top 15 deliveries
})

// Due Date formatting & helpers
function getDueDateClass(task: ITask): string {
	if (!task.dueDate) return 'is-neutral'
	const due = new Date(task.dueDate)
	const now = new Date()
	const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate())
	const endOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59)

	if (due < startOfToday) {
		return 'is-overdue'
	}
	if (due >= startOfToday && due <= endOfToday) {
		return 'is-today'
	}
	return 'is-future'
}

function getDueDateIcon(task: ITask): string {
	if (!task.dueDate) return 'calendar-minus'
	const due = new Date(task.dueDate)
	const now = new Date()
	if (due < now) return 'exclamation-triangle'
	return 'calendar-alt'
}

function formatDueBadgeText(task: ITask): string {
	if (!task.dueDate) return 'Sem prazo'
	const due = new Date(task.dueDate)
	const now = new Date()
	const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate())
	const endOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59)

	if (due < startOfToday) {
		return `Atrasada: ${due.toLocaleDateString('pt-BR')}`
	}
	if (due >= startOfToday && due <= endOfToday) {
		return 'Hoje'
	}
	return due.toLocaleDateString('pt-BR')
}

// Task Status label and class
function getTaskStatusClass(task: ITask): string {
	if (task.done) return 'status-done'
	const bucketTitle = task.bucket?.title?.toLowerCase() || ''
	if (
		bucketTitle.includes('andamento') ||
		bucketTitle.includes('doing') ||
		(task.percentDone && task.percentDone > 0)
	) {
		return 'status-doing'
	}
	return 'status-todo'
}

function getTaskStatusLabel(task: ITask): string {
	if (task.done) return 'Concluído'
	const bucketTitle = task.bucket?.title?.toLowerCase() || ''
	if (
		bucketTitle.includes('andamento') ||
		bucketTitle.includes('doing') ||
		(task.percentDone && task.percentDone > 0)
	) {
		return 'Em andamento'
	}
	return 'A fazer'
}

// Labels
function getLabelStyle(label: ILabel) {
	if (!label.hexColor) {
		return { background: '#F2F4F7', color: '#344054' }
	}
	return {
		background: `${label.hexColor}20`,
		color: label.hexColor,
		borderColor: `${label.hexColor}40`,
	}
}

// Donut Chart & Project Distribution
const PALETTE = ['#3157F6', '#7F56D9', '#12B76A', '#F79009', '#06AED4', '#667085']

const projectDistribution = computed(() => {
	const countsByProject: Record<number, number> = {}
	for (const t of allTasks.value) {
		countsByProject[t.projectId] = (countsByProject[t.projectId] || 0) + 1
	}

	const total = allTasks.value.length
	const sorted = Object.entries(countsByProject)
		.map(([pId, count]) => {
			const proj = projectStore.projects[Number(pId)]
			return {
				id: Number(pId),
				title: proj ? proj.title : `Projeto #${pId}`,
				count,
				percentage: total > 0 ? Math.round((count / total) * 100) : 0,
			}
		})
		.sort((a, b) => b.count - a.count)

	if (sorted.length <= 5) {
		return sorted.map((item, idx) => ({
			...item,
			color: PALETTE[idx % PALETTE.length],
		}))
	}

	// Top 4 + Outros
	const top4 = sorted.slice(0, 4)
	const othersCount = sorted.slice(4).reduce((sum, it) => sum + it.count, 0)
	const othersPct = total > 0 ? Math.round((othersCount / total) * 100) : 0

	const res = top4.map((item, idx) => ({
		...item,
		color: PALETTE[idx % PALETTE.length],
	}))
	res.push({
		id: -999,
		title: 'Outros projetos',
		count: othersCount,
		percentage: othersPct,
		color: '#98A2B3',
	})
	return res
})

// Compute SVG Donut segments
const donutSegments = computed(() => {
	const circumference = 2 * Math.PI * 60 // ~376.99
	const total = metrics.value.total
	if (total === 0) return []

	let currentOffset = 0
	return projectDistribution.value.map(item => {
		const segmentLength = (item.count / total) * circumference
		const dasharray = `${segmentLength} ${circumference - segmentLength}`
		const dashoffset = -currentOffset
		currentOffset += segmentLength

		return {
			color: item.color,
			dasharray,
			dashoffset,
		}
	})
})

// Active Projects List with metrics
const activeProjectsList = computed(() => {
	const now = new Date()
	return availableProjects.value.map(proj => {
		const projTasks = allTasks.value.filter(t => t.projectId === proj.id)
		const total = projTasks.length
		const doneTasks = projTasks.filter(t => t.done)
		const doneCount = doneTasks.length
		const openCount = total - doneCount
		const overdueCount = projTasks.filter(t => !t.done && t.dueDate && new Date(t.dueDate) < now).length
		const progressPercent = total > 0 ? Math.round((doneCount / total) * 100) : 0

		return {
			...proj,
			tasksCount: total,
			openTasksCount: openCount,
			doneTasksCount: doneCount,
			overdueTasksCount: overdueCount,
			progressPercent,
		}
	})
})

// Actions
async function loadDashboardData() {
	isLoading.value = true
	try {
		// Load all projects first if needed
		if (!projectStore.hasProjects) {
			await projectStore.loadProjects()
		}

		// Fetch all user tasks (both pending and done)
		const tasks = await taskService.getAll({}, {
			sort_by: ['due_date', 'id'],
			order_by: ['asc', 'desc'],
			per_page: 250,
			expand: 'subtasks',
		})
		allTasks.value = tasks || []
	} catch (err) {
		console.error('Error loading dashboard data:', err)
	} finally {
		isLoading.value = false
	}
}

async function toggleTaskStatus(task: ITask) {
	try {
		const updated = { ...task, done: !task.done }
		await taskStore.update(updated)
		task.done = !task.done
	} catch (err) {
		console.error('Error toggling task:', err)
	}
}

// Modal controls
function openNewDemandModal() {
	modalError.value = ''
	newDemandForm.title = ''
	newDemandForm.dueDate = ''
	selectedPreset.value = 'none'

	if (availableProjects.value.length > 0) {
		newDemandForm.projectId = availableProjects.value[0].id
	}
	showNewDemandModal.value = true
}

function closeNewDemandModal() {
	showNewDemandModal.value = false
}

function setPresetDate(preset: 'today' | 'tomorrow' | 'next-week' | 'none') {
	selectedPreset.value = preset
	const now = new Date()

	if (preset === 'today') {
		newDemandForm.dueDate = now.toISOString().split('T')[0]
	} else if (preset === 'tomorrow') {
		const tom = new Date(now)
		tom.setDate(tom.getDate() + 1)
		newDemandForm.dueDate = tom.toISOString().split('T')[0]
	} else if (preset === 'next-week') {
		const nw = new Date(now)
		nw.setDate(nw.getDate() + 7)
		newDemandForm.dueDate = nw.toISOString().split('T')[0]
	} else {
		newDemandForm.dueDate = ''
	}
}

async function submitNewDemand() {
	if (!newDemandForm.title.trim()) {
		modalError.value = 'Por favor, informe o título da demanda.'
		return
	}
	if (!newDemandForm.projectId) {
		modalError.value = 'Selecione um projeto para a demanda.'
		return
	}

	isSubmitting.value = true
	modalError.value = ''

	try {
		const payload: any = {
			title: newDemandForm.title.trim(),
			projectId: newDemandForm.projectId,
		}

		if (newDemandForm.dueDate) {
			const d = new Date(newDemandForm.dueDate + 'T23:59:59')
			payload.dueDate = d.toISOString()
		}

		const createdTask = await taskService.create(payload)
		if (createdTask) {
			allTasks.value.unshift(createdTask)
			closeNewDemandModal()
		}
	} catch (err: any) {
		console.error('Error creating demand:', err)
		modalError.value = err?.message || 'Erro ao criar demanda. Tente novamente.'
	} finally {
		isSubmitting.value = false
	}
}

onMounted(() => {
	loadDashboardData()
})
</script>

<style lang="scss" scoped>
/* DATA FLOW SAAS CORPORATIVO PALETTE */
$bg-app: #F6F8FC;
$card-bg: #FFFFFF;
$text-main: #162033;
$text-secondary: #667085;
$border-color: #E7EAF0;
$primary: #3157F6;
$success: #12B76A;
$warning: #F79009;
$danger: #F04438;
$purple: #7F56D9;

.dataflow-overview-wrapper {
	width: 100%;
	max-width: 1500px;
	margin: 0 auto;
	padding: 0 0 2rem 0;
	box-sizing: border-box;
	color: $text-main;
	font-family: inherit;

	@media (max-width: 768px) {
		padding: 0 0 1.5rem 0;
	}
}

/* HEADER */
.overview-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 2rem;
	flex-wrap: wrap;
	gap: 1.25rem;

	.header-left {
		.welcome-title {
			font-size: 1.75rem;
			font-weight: 700;
			letter-spacing: -0.02em;
			margin: 0 0 0.35rem 0;
			color: $text-main;

			.highlight-user {
				color: $primary;
			}
		}

		.welcome-subtitle {
			font-size: 0.95rem;
			color: $text-secondary;
			margin: 0;
		}
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 1rem;

		.current-date-pill {
			display: flex;
			align-items: center;
			gap: 0.5rem;
			background: #FFFFFF;
			border: 1px solid $border-color;
			border-radius: 9999px;
			padding: 0.55rem 1rem;
			font-size: 0.875rem;
			font-weight: 500;
			color: $text-secondary;
			box-shadow: 0 1px 2px rgba(16, 24, 40, 0.05);

			.date-icon {
				color: $primary;
			}
		}

		.btn-primary-action {
			display: inline-flex;
			align-items: center;
			gap: 0.5rem;
			background: $primary;
			color: #FFFFFF;
			border: none;
			border-radius: 8px;
			padding: 0.65rem 1.25rem;
			font-size: 0.925rem;
			font-weight: 600;
			cursor: pointer;
			transition: all 0.2s ease;
			box-shadow: 0 4px 12px rgba(49, 87, 246, 0.25);

			&:hover {
				background: darken($primary, 6%);
				box-shadow: 0 6px 16px rgba(49, 87, 246, 0.35);
				transform: translateY(-1px);
			}

			&:active {
				transform: translateY(0);
			}

			.btn-icon {
				font-size: 0.85rem;
			}
		}
	}
}

/* 4 METRICS CARDS */
.metrics-grid {
	display: grid;
	grid-template-columns: repeat(4, minmax(0, 1fr));
	gap: 1rem;
	margin-bottom: 1.5rem;

	@media (max-width: 1100px) {
		grid-template-columns: repeat(2, minmax(0, 1fr));
	}

	@media (max-width: 600px) {
		grid-template-columns: 1fr;
	}

	.metric-card {
		background: $card-bg;
		border: 1px solid $border-color;
		border-radius: 12px;
		padding: 1.35rem 1.5rem;
		box-shadow: 0 1px 3px rgba(16, 24, 40, 0.05);
		transition: transform 0.2s ease, box-shadow 0.2s ease;

		&:hover {
			transform: translateY(-2px);
			box-shadow: 0 8px 24px rgba(16, 24, 40, 0.08);
		}

		.metric-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			margin-bottom: 0.75rem;

			.metric-tag {
				font-size: 0.75rem;
				font-weight: 700;
				letter-spacing: 0.06em;
				color: $text-secondary;
			}

			.metric-icon-wrap {
				width: 36px;
				height: 36px;
				border-radius: 8px;
				display: flex;
				align-items: center;
				justify-content: center;
				font-size: 1rem;

				&.danger-bg {
					background: #FEF3F2;
					color: $danger;
				}
				&.primary-bg {
					background: #EFF4FF;
					color: $primary;
				}
				&.warning-bg {
					background: #FFFAEB;
					color: $warning;
				}
				&.success-bg {
					background: #ECFDF3;
					color: $success;
				}
			}
		}

		.metric-value {
			font-size: 2rem;
			font-weight: 700;
			line-height: 1.1;
			margin-bottom: 0.5rem;
			color: $text-main;
		}

		.metric-footer {
			.metric-subtext {
				font-size: 0.8rem;
				color: $text-secondary;

				&.danger-text { color: $danger; font-weight: 500; }
				&.success-text { color: $success; font-weight: 500; }
			}
		}

		&.is-danger {
			border-left: 4px solid $danger;
		}
		&.is-primary {
			border-left: 4px solid $primary;
		}
		&.is-warning {
			border-left: 4px solid $warning;
		}
		&.is-success {
			border-left: 4px solid $success;
		}
	}
}

/* MAIN DASHBOARD GRID (65% / 35%) */
.main-dashboard-grid {
	display: grid;
	grid-template-columns: minmax(0, 1.85fr) minmax(0, 1fr);
	gap: 1.25rem;
	margin-bottom: 2rem;

	@media (max-width: 1024px) {
		grid-template-columns: 1fr;
	}

	.dashboard-panel {
		background: $card-bg;
		border: 1px solid $border-color;
		border-radius: 12px;
		padding: 1.5rem;
		box-shadow: 0 1px 3px rgba(16, 24, 40, 0.05);

		.panel-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			margin-bottom: 1.25rem;
			padding-bottom: 0.85rem;
			border-bottom: 1px solid $border-color;

			.panel-title-area {
				display: flex;
				align-items: center;
				gap: 0.65rem;

				.panel-title {
					font-size: 1.125rem;
					font-weight: 700;
					color: $text-main;
					margin: 0;
				}

				.panel-badge {
					background: #F2F4F7;
					color: $text-secondary;
					font-size: 0.75rem;
					font-weight: 600;
					padding: 0.2rem 0.55rem;
					border-radius: 9999px;
				}

				.panel-subtitle-text {
					font-size: 0.8rem;
					color: $text-secondary;
				}
			}

			.panel-link-action {
				display: inline-flex;
				align-items: center;
				gap: 0.35rem;
				font-size: 0.85rem;
				font-weight: 600;
				color: $primary;
				text-decoration: none;
				transition: color 0.15s ease;

				&:hover {
					color: darken($primary, 10%);
					.arrow-icon {
						transform: translateX(2px);
					}
				}

				.arrow-icon {
					font-size: 0.75rem;
					transition: transform 0.15s ease;
				}
			}
		}
	}
}

/* DELIVERIES TABLE */
.table-container {
	overflow-x: auto;

	.dataflow-table {
		width: 100%;
		border-collapse: collapse;

		th {
			text-align: left;
			font-size: 0.75rem;
			font-weight: 700;
			letter-spacing: 0.05em;
			color: $text-secondary;
			padding: 0.6rem 0.75rem;
			border-bottom: 1px solid $border-color;
		}

		td {
			padding: 0.85rem 0.75rem;
			border-bottom: 1px solid #F2F4F7;
			vertical-align: middle;
			font-size: 0.875rem;
		}

		.task-row {
			transition: background 0.15s ease;

			&:hover {
				background: #F8FAFC;
			}

			&.is-completed {
				opacity: 0.6;
				.task-title-link {
					text-decoration: line-through;
					color: $text-secondary;
				}
			}
		}

		.col-checkbox {
			width: 36px;
			padding-left: 0.25rem;

			.task-check-button {
				width: 20px;
				height: 20px;
				border-radius: 6px;
				border: 1.5px solid #CBD5E1;
				background: transparent;
				cursor: pointer;
				display: flex;
				align-items: center;
				justify-content: center;
				padding: 0;
				transition: all 0.15s ease;

				&:hover {
					border-color: $primary;
					background: rgba(49, 87, 246, 0.08);
				}

				&.is-checked {
					background: $success;
					border-color: $success;
					color: #FFFFFF;
				}

				.check-icon {
					font-size: 0.65rem;
				}
			}
		}

		.col-task {
			.task-info {
				display: flex;
				flex-direction: column;
				gap: 0.25rem;

				.task-title-link {
					font-weight: 600;
					color: $text-main;
					text-decoration: none;
					transition: color 0.15s ease;

					&:hover {
						color: $primary;
					}
				}

				.task-labels-row {
					display: flex;
					flex-wrap: wrap;
					gap: 0.35rem;

					.task-label-chip {
						font-size: 0.7rem;
						font-weight: 500;
						padding: 0.1rem 0.45rem;
						border-radius: 4px;
						border: 1px solid transparent;
					}
				}
			}
		}

		.col-project {
			.project-pill {
				display: inline-flex;
				align-items: center;
				gap: 0.45rem;
				background: #F2F4F7;
				padding: 0.25rem 0.65rem;
				border-radius: 6px;
				font-size: 0.8rem;
				font-weight: 500;
				color: $text-main;
				text-decoration: none;
				transition: background 0.15s ease;

				&:hover {
					background: #E4E7EC;
				}

				.project-dot {
					width: 8px;
					height: 8px;
					border-radius: 50%;
				}

				&.is-empty {
					background: transparent;
					color: $text-secondary;
				}
			}
		}

		.col-due {
			white-space: nowrap;

			.due-date-badge {
				display: inline-flex;
				align-items: center;
				gap: 0.35rem;
				font-size: 0.775rem;
				font-weight: 600;
				padding: 0.25rem 0.65rem;
				border-radius: 6px;

				.due-icon {
					font-size: 0.75rem;
				}

				&.is-overdue {
					background: #FEF3F2;
					color: $danger;
					border: 1px solid #FECDCA;
				}

				&.is-today {
					background: #FFFAEB;
					color: $warning;
					border: 1px solid #FEDF89;
				}

				&.is-future {
					background: #F2F4F7;
					color: $text-secondary;
				}

				&.is-neutral {
					color: #98A2B3;
				}
			}
		}

		.col-status {
			white-space: nowrap;

			.status-pill {
				display: inline-block;
				font-size: 0.75rem;
				font-weight: 600;
				padding: 0.2rem 0.55rem;
				border-radius: 9999px;

				&.status-done {
					background: #ECFDF3;
					color: $success;
				}
				&.status-doing {
					background: #EFF4FF;
					color: $primary;
				}
				&.status-todo {
					background: #FFFAEB;
					color: $warning;
				}
			}
		}
	}
}

/* DONUT CHART CONTAINER */
.donut-chart-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 1.5rem;

	.donut-svg-wrapper {
		position: relative;
		width: 160px;
		height: 160px;

		.donut-svg {
			transform: rotate(0deg);
		}

		.donut-center-content {
			position: absolute;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
			pointer-events: none;

			.donut-center-number {
				font-size: 1.65rem;
				font-weight: 800;
				line-height: 1;
				color: $text-main;
			}

			.donut-center-label {
				font-size: 0.75rem;
				font-weight: 600;
				color: $text-secondary;
				margin-top: 0.15rem;
			}
		}
	}

	.donut-legend-list {
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 0.85rem;

		.legend-item {
			display: flex;
			flex-direction: column;
			gap: 0.25rem;

			.legend-meta {
				display: flex;
				align-items: center;
				gap: 0.5rem;
				font-size: 0.85rem;
				font-weight: 600;
				color: $text-main;

				.legend-color-dot {
					width: 10px;
					height: 10px;
					border-radius: 50%;
					flex-shrink: 0;
				}

				.legend-name {
					white-space: nowrap;
					overflow: hidden;
					text-overflow: ellipsis;
					flex: 1;
				}
			}

			.legend-counts {
				display: flex;
				align-items: center;
				justify-content: flex-end;
				gap: 0.35rem;
				font-size: 0.775rem;
				color: $text-secondary;

				.legend-count {
					font-weight: 600;
					color: $text-main;
				}
			}

			.legend-bar-track {
				width: 100%;
				height: 5px;
				background: #F2F4F7;
				border-radius: 9999px;
				overflow: hidden;

				.legend-bar-fill {
					height: 100%;
					border-radius: 9999px;
					transition: width 0.3s ease;
				}
			}
		}

		.donut-legend-empty {
			font-size: 0.85rem;
			color: $text-secondary;
			text-align: center;
			padding: 1rem;
		}
	}
}

/* PROJECTS OVERVIEW ROW */
.projects-overview-section {
	.section-heading {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		margin-bottom: 1.25rem;
		flex-wrap: wrap;
		gap: 0.5rem;

		.heading-text {
			.section-title {
				font-size: 1.25rem;
				font-weight: 700;
				color: $text-main;
				margin: 0 0 0.25rem 0;
			}

			.section-subtitle {
				font-size: 0.875rem;
				color: $text-secondary;
				margin: 0;
			}
		}

		.heading-action-link {
			display: inline-flex;
			align-items: center;
			gap: 0.35rem;
			font-size: 0.875rem;
			font-weight: 600;
			color: $primary;
			text-decoration: none;

			&:hover {
				color: darken($primary, 10%);
			}
		}
	}

	.project-cards-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		gap: 1.25rem;

		.project-summary-card {
			background: $card-bg;
			border: 1px solid $border-color;
			border-radius: 12px;
			overflow: hidden;
			box-shadow: 0 1px 3px rgba(16, 24, 40, 0.05);
			transition: transform 0.2s ease, box-shadow 0.2s ease;

			&:hover {
				transform: translateY(-2px);
				box-shadow: 0 8px 24px rgba(16, 24, 40, 0.08);
			}

			.card-accent-bar {
				height: 4px;
				width: 100%;
			}

			.card-inner {
				padding: 1.25rem;
				display: flex;
				flex-direction: column;
				gap: 1rem;

				.project-card-header {
					.project-card-title {
						font-size: 1rem;
						font-weight: 700;
						color: $text-main;
						text-decoration: none;
						transition: color 0.15s ease;

						&:hover {
							color: $primary;
						}
					}
				}

				.project-metrics-row {
					display: flex;
					align-items: center;
					justify-content: space-between;
					background: #F8FAFC;
					border-radius: 8px;
					padding: 0.65rem 0.85rem;

					.proj-metric {
						display: flex;
						flex-direction: column;
						align-items: center;

						.proj-metric-val {
							font-size: 1rem;
							font-weight: 700;
							color: $text-main;

							&.text-danger { color: $danger; }
							&.text-success { color: $success; }
						}

						.proj-metric-label {
							font-size: 0.7rem;
							color: $text-secondary;
						}
					}
				}

				.project-progress-wrap {
					display: flex;
					flex-direction: column;
					gap: 0.35rem;

					.progress-info-row {
						display: flex;
						align-items: center;
						justify-content: space-between;
						font-size: 0.775rem;

						.progress-title {
							color: $text-secondary;
							font-weight: 500;
						}

						.progress-percent {
							color: $text-main;
							font-weight: 700;
						}
					}

					.progress-track {
						width: 100%;
						height: 6px;
						background: #F2F4F7;
						border-radius: 9999px;
						overflow: hidden;

						.progress-fill {
							height: 100%;
							border-radius: 9999px;
							transition: width 0.3s ease;
						}
					}
				}
			}
		}

		.no-projects-card {
			grid-column: 1 / -1;
			background: $card-bg;
			border: 1px dashed $border-color;
			border-radius: 12px;
			padding: 2.5rem;
			text-align: center;
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 0.75rem;

			.no-proj-icon {
				font-size: 2.5rem;
				color: #98A2B3;
			}

			p {
				color: $text-secondary;
				font-size: 0.95rem;
				margin: 0;
			}

			.btn-create-proj {
				display: inline-flex;
				align-items: center;
				gap: 0.45rem;
				background: $primary;
				color: #FFFFFF;
				text-decoration: none;
				border-radius: 8px;
				padding: 0.55rem 1rem;
				font-size: 0.875rem;
				font-weight: 600;
				margin-top: 0.5rem;
			}
		}
	}
}

/* EMPTY & LOADING STATES */
.panel-loading {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 3rem 1rem;
	gap: 0.75rem;
	color: $text-secondary;
	font-size: 0.9rem;

	.loading-spinner {
		width: 32px;
		height: 32px;
		border: 3px solid #E4E7EC;
		border-top-color: $primary;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}
}

.panel-empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 3rem 1rem;
	text-align: center;

	.empty-icon-wrap {
		width: 52px;
		height: 52px;
		border-radius: 12px;
		background: #ECFDF3;
		color: $success;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
		margin-bottom: 1rem;
	}

	.empty-title {
		font-size: 1.1rem;
		font-weight: 700;
		color: $text-main;
		margin: 0 0 0.35rem 0;
	}

	.empty-description {
		font-size: 0.875rem;
		color: $text-secondary;
		max-width: 320px;
		margin: 0 0 1.25rem 0;
	}

	.btn-secondary-action {
		display: inline-flex;
		align-items: center;
		gap: 0.45rem;
		background: #FFFFFF;
		border: 1px solid $border-color;
		border-radius: 8px;
		padding: 0.55rem 1rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: $text-main;
		cursor: pointer;
		transition: all 0.15s ease;

		&:hover {
			background: #F8FAFC;
			border-color: $primary;
			color: $primary;
		}
	}
}

/* QUICK ADD MODAL */
.dataflow-modal-backdrop {
	position: fixed;
	inset: 0;
	background: rgba(11, 22, 51, 0.6);
	backdrop-filter: blur(4px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: 1rem;

	.dataflow-modal-box {
		background: #FFFFFF;
		border-radius: 16px;
		width: 100%;
		max-width: 520px;
		box-shadow: 0 20px 40px rgba(11, 22, 51, 0.2);
		overflow: hidden;
		animation: modalZoomIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);

		.modal-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			padding: 1.25rem 1.5rem;
			border-bottom: 1px solid $border-color;

			.modal-title-wrap {
				display: flex;
				align-items: center;
				gap: 0.75rem;

				.modal-badge-icon {
					width: 32px;
					height: 32px;
					border-radius: 8px;
					background: #EFF4FF;
					color: $primary;
					display: flex;
					align-items: center;
					justify-content: center;
					font-size: 0.9rem;
				}

				.modal-title {
					font-size: 1.15rem;
					font-weight: 700;
					margin: 0;
					color: $text-main;
				}
			}

			.modal-close-btn {
				background: transparent;
				border: none;
				color: $text-secondary;
				font-size: 1.1rem;
				cursor: pointer;
				padding: 0.25rem;
				border-radius: 6px;
				transition: all 0.15s ease;

				&:hover {
					background: #F2F4F7;
					color: $text-main;
				}
			}
		}

		.modal-form {
			padding: 1.5rem;
			display: flex;
			flex-direction: column;
			gap: 1.25rem;

			.form-group {
				display: flex;
				flex-direction: column;
				gap: 0.4rem;

				.form-label {
					font-size: 0.85rem;
					font-weight: 600;
					color: $text-main;

					.required {
						color: $danger;
					}
				}

				.form-input,
				.form-select {
					width: 100%;
					padding: 0.65rem 0.85rem;
					border: 1px solid $border-color;
					border-radius: 8px;
					font-size: 0.925rem;
					color: $text-main;
					background: #FFFFFF;
					transition: border-color 0.15s ease, box-shadow 0.15s ease;

					&:focus {
						outline: none;
						border-color: $primary;
						box-shadow: 0 0 0 3px rgba(49, 87, 246, 0.15);
					}
				}

				.m-top-xs {
					margin-top: 0.5rem;
				}

				.preset-buttons-row {
					display: flex;
					flex-wrap: wrap;
					gap: 0.5rem;

					.preset-btn {
						background: #F2F4F7;
						border: 1px solid transparent;
						border-radius: 6px;
						padding: 0.35rem 0.75rem;
						font-size: 0.8rem;
						font-weight: 500;
						color: $text-secondary;
						cursor: pointer;
						transition: all 0.15s ease;

						&:hover {
							background: #E4E7EC;
							color: $text-main;
						}

						&.is-active {
							background: #EFF4FF;
							border-color: #D1E0FF;
							color: $primary;
							font-weight: 600;
						}
					}
				}
			}

			.modal-error-alert {
				background: #FEF3F2;
				border: 1px solid #FECDCA;
				color: $danger;
				padding: 0.65rem 0.85rem;
				border-radius: 8px;
				font-size: 0.85rem;
			}

			.modal-actions {
				display: flex;
				align-items: center;
				justify-content: flex-end;
				gap: 0.75rem;
				margin-top: 0.5rem;

				.btn-cancel {
					background: #FFFFFF;
					border: 1px solid $border-color;
					border-radius: 8px;
					padding: 0.6rem 1.15rem;
					font-size: 0.875rem;
					font-weight: 600;
					color: $text-secondary;
					cursor: pointer;
					transition: all 0.15s ease;

					&:hover {
						background: #F8FAFC;
						color: $text-main;
					}
				}

				.btn-submit {
					display: inline-flex;
					align-items: center;
					gap: 0.45rem;
					background: $primary;
					border: none;
					border-radius: 8px;
					padding: 0.6rem 1.35rem;
					font-size: 0.875rem;
					font-weight: 600;
					color: #FFFFFF;
					cursor: pointer;
					transition: all 0.15s ease;

					&:hover:not(:disabled) {
						background: darken($primary, 6%);
					}

					&:disabled {
						opacity: 0.6;
						cursor: not-allowed;
					}

					.spin {
						animation: spin 0.8s linear infinite;
					}
				}
			}
		}
	}
}

@keyframes spin {
	from { transform: rotate(0deg); }
	to { transform: rotate(360deg); }
}

@keyframes modalZoomIn {
	from {
		opacity: 0;
		transform: scale(0.95);
	}
	to {
		opacity: 1;
		transform: scale(1);
	}
}
</style>
