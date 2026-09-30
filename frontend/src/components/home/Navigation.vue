<template>
	<aside
		:class="{'is-active': baseStore.menuActive, 'is-resizing': isResizing}"
		class="menu-container dataflow-sidebar"
		:style="{'--sidebar-width': sidebarWidth}"
	>
		<!-- On mobile only: show brand logo inside drawer -->
		<div class="sidebar-brand-mobile">
			<RouterLink
				:to="{name: 'home'}"
				class="brand-link"
				aria-label="Data Flow Home"
			>
				<Logo
					variant="white"
					class="brand-logo"
				/>
			</RouterLink>
		</div>

		<nav
			class="menu top-menu"
			:aria-label="$t('navigation.main')"
		>
			<menu class="menu-list other-menu-items">
				<li>
					<RouterLink
						v-shortcut="SHORTCUTS.navigation.overview"
						:to="{ name: 'home'}"
						class="nav-item-link"
						:class="{'is-current': isRouteActive('home')}"
					>
						<span class="menu-item-icon icon">
							<Icon icon="chart-pie" />
						</span>
						<span class="nav-label">Visão geral</span>
					</RouterLink>
				</li>
				<li>
					<RouterLink
						v-shortcut="SHORTCUTS.navigation.upcoming"
						:to="{ name: 'tasks.range'}"
						class="nav-item-link"
						:class="{'is-current': isRouteActive('tasks.range')}"
					>
						<span class="menu-item-icon icon">
							<Icon :icon="['far', 'calendar-alt']" />
						</span>
						<span class="nav-label">Em breve</span>
					</RouterLink>
				</li>
				<li>
					<RouterLink
						v-shortcut="SHORTCUTS.navigation.projects"
						:to="{ name: 'projects.index'}"
						class="nav-item-link"
						:class="{'is-current': isRouteActive('projects.index')}"
					>
						<span class="menu-item-icon icon">
							<Icon icon="layer-group" />
						</span>
						<span class="nav-label">Projetos</span>
					</RouterLink>
				</li>
				<li>
					<RouterLink
						:to="{ name: 'tasks.index'}"
						class="nav-item-link"
						:class="{'is-current': isRouteActive('tasks.index')}"
					>
						<span class="menu-item-icon icon">
							<Icon icon="tasks" />
						</span>
						<span class="nav-label">Demandas</span>
					</RouterLink>
				</li>
				<li>
					<RouterLink
						v-shortcut="SHORTCUTS.navigation.labels"
						:to="{ name: 'labels.index'}"
						class="nav-item-link"
						:class="{'is-current': isRouteActive('labels.index')}"
					>
						<span class="menu-item-icon icon">
							<Icon icon="tags" />
						</span>
						<span class="nav-label">Marcadores</span>
					</RouterLink>
				</li>
				<li>
					<RouterLink
						v-shortcut="SHORTCUTS.navigation.teams"
						:to="{ name: 'teams.index'}"
						class="nav-item-link"
						:class="{'is-current': isRouteActive('teams.index')}"
					>
						<span class="menu-item-icon icon">
							<Icon icon="users" />
						</span>
						<span class="nav-label">Equipes</span>
					</RouterLink>
				</li>
			</menu>
		</nav>

		<div class="sidebar-divider" />

		<!-- PROJECTS SECTION -->
		<div class="sidebar-projects-section">
			<div class="projects-section-title">
				<span>MEUS PROJETOS</span>
			</div>
			
			<Loading
				v-if="projectStore.isLoading"
				variant="small"
			/>
			<template v-else>
				<nav
					v-if="favoriteProjects.length"
					class="menu projects-nav-sub"
					:aria-label="$t('project.pseudo.favorites.title')"
				>
					<ProjectsNavigation
						:model-value="favoriteProjects"
						:can-edit-order="false"
						:can-collapse="false"
					/>
				</nav>
				
				<nav
					v-if="savedFilterProjects.length"
					class="menu projects-nav-sub"
					:aria-label="$t('navigation.savedFilters')"
				>
					<ProjectsNavigation
						:model-value="savedFilterProjects"
						:can-edit-order="false"
						:can-collapse="false"
					/>
				</nav>

				<nav
					class="menu projects-nav-sub"
					:aria-label="$t('project.projects')"
				>
					<ProjectsNavigation
						:model-value="projects"
						:can-edit-order="true"
						:can-collapse="true"
					/>
				</nav>
			</template>
		</div>

		<!-- RESIZE HANDLE -->
		<div
			v-if="!isMobile"
			class="resize-handle"
			@mousedown="startResize"
			@touchstart="startResize"
		/>
	</aside>
</template>

<script setup lang="ts">
import {computed} from 'vue'
import {useRoute} from 'vue-router'

import {SHORTCUTS} from '@/constants/shortcuts'
import Logo from '@/components/home/Logo.vue'
import Loading from '@/components/misc/Loading.vue'
import Icon from '@/components/misc/Icon'

import {useBaseStore} from '@/stores/base'
import {useProjectStore} from '@/stores/projects'
import ProjectsNavigation from '@/components/home/ProjectsNavigation.vue'
import type {IProject} from '@/modelTypes/IProject'
import {useSidebarResize} from '@/composables/useSidebarResize'

const route = useRoute()
const baseStore = useBaseStore()
const projectStore = useProjectStore()

const {
	sidebarWidth,
	isResizing,
	startResize,
	isMobile,
} = useSidebarResize()

const projects = computed(() => projectStore.notArchivedRootProjects as IProject[])
const favoriteProjects = computed(() => projectStore.favoriteProjects as IProject[])
const savedFilterProjects = computed(() => projectStore.savedFilterProjects as IProject[])

function isRouteActive(name: string): boolean {
	return route.name === name
}
</script>

<style lang="scss" scoped>
.menu-container.dataflow-sidebar {
	--sidebar-width: 250px;

	display: flex;
	flex-direction: column;
	background: #0B1633;
	color: #94A3B8;
	padding: 0.85rem 0.65rem;
	transition: transform $transition-duration ease-in;
	position: fixed;
	inset-block-start: $navbar-height;
	inset-block-end: 0;
	inset-inline-start: 0;
	transform: translateX(-100%);
	inline-size: var(--sidebar-width);
	overflow-y: auto;
	overflow-x: hidden;
	z-index: 20;
	border-inline-end: 1px solid rgba(255, 255, 255, 0.08);

	[dir="rtl"] & {
		transform: translateX(100%);
	}

	@media screen and (max-width: $tablet) {
		inset-block-start: 0;
		inline-size: 78vw;
		z-index: 40;
	}

	&.is-active {
		transform: translateX(0) !important;
		transition: transform $transition-duration ease-out;

		[dir="rtl"] & {
			transform: translateX(0) !important;
		}
	}

	&.is-resizing {
		transition: none;
	}
}

.sidebar-brand-mobile {
	padding: 0.5rem 0.5rem 1rem 0.5rem;
	display: none;
	align-items: center;

	@media screen and (max-width: $tablet) {
		display: flex;
	}

	.brand-link {
		display: flex;
		align-items: center;
		text-decoration: none;
	}

	.brand-logo {
		max-inline-size: 140px;
	}
}

.menu.top-menu {
	margin-bottom: 0.25rem;

	.menu-list {
		margin: 0;
		padding: 0;
		list-style: none;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;

		.nav-item-link {
			display: flex;
			align-items: center;
			gap: 0.65rem;
			padding: 0.55rem 0.75rem;
			color: #94A3B8;
			font-size: 0.875rem;
			font-weight: 500;
			border-radius: 8px;
			text-decoration: none;
			transition: all 0.15s ease;

			.menu-item-icon {
				display: flex;
				align-items: center;
				justify-content: center;
				width: 1.25rem;
				font-size: 0.95rem;
				color: #64748B;
				transition: color 0.15s ease;
			}

			.nav-label {
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}

			&:hover {
				background: rgba(255, 255, 255, 0.06);
				color: #F8FAFC;

				.menu-item-icon {
					color: #38BDF8;
				}
			}

			&.router-link-active,
			&.is-current {
				background: #3157F6;
				color: #FFFFFF;
				font-weight: 600;
				box-shadow: 0 4px 12px rgba(49, 87, 246, 0.35);

				.menu-item-icon {
					color: #FFFFFF;
				}
			}
		}
	}
}

.sidebar-divider {
	height: 1px;
	background: rgba(255, 255, 255, 0.08);
	margin: 0.65rem 0.5rem 0.85rem 0.5rem;
}

.sidebar-projects-section {
	flex: 1 1 auto;
	display: flex;
	flex-direction: column;

	.projects-section-title {
		padding: 0.25rem 0.75rem 0.4rem 0.75rem;
		font-size: 0.68rem;
		letter-spacing: 0.08em;
		font-weight: 700;
		color: #64748B;
	}

	.projects-nav-sub {
		:deep(.menu-list) {
			a {
				color: #94A3B8;
				border-radius: 6px;
				font-size: 0.85rem;
				padding: 0.4rem 0.65rem;

				&:hover {
					background: rgba(255, 255, 255, 0.05);
					color: #F8FAFC;
				}

				&.router-link-active {
					background: rgba(49, 87, 246, 0.15);
					color: #60A5FA;
					font-weight: 600;
				}
			}
		}
	}
}

.resize-handle {
	position: absolute;
	inset-block-start: 0;
	inset-block-end: 0;
	inset-inline-end: 0;
	inline-size: 4px;
	cursor: ew-resize;
	background: transparent;
	transition: background-color 0.2s ease;
	touch-action: none;

	&:hover,
	&:active {
		background-color: #3157F6;
	}
}
</style>
