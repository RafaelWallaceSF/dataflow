<template>
	<div class="no-auth-wrapper">
		<div class="noauth-container">
			<section
				class="image-panel"
				:class="{ 'has-message': motd !== '' }"
			>
				<div class="image-panel-content">
					<Logo
						class="logo"
						width="140"
						height="38"
						:dark="true"
					/>
					
					<div class="hero-text">
						<h1 class="image-title">
							Organização<br>
							que impulsiona<br>
							<span class="highlight">resultados.</span>
						</h1>
						<p class="image-subtitle">
							Gestão de projetos e demandas em um só lugar, para equipes mais focadas e empresas mais produtivas.
						</p>
					</div>
					
					<div class="feature-chips">
						<div class="chip">
							<div class="chip-icon"><i class="fas fa-check-square"></i></div>
							<span>Mais organização</span>
						</div>
						<div class="chip">
							<div class="chip-icon"><i class="fas fa-chart-bar"></i></div>
							<span>Mais produtividade</span>
						</div>
						<div class="chip">
							<div class="chip-icon"><i class="fas fa-users"></i></div>
							<span>Mais resultados</span>
						</div>
					</div>
				</div>
				<Message v-if="motd !== ''" class="motd-message">
					{{ motd }}
				</Message>
			</section>
			
			<main
				id="main-content"
				tabindex="-1"
				class="content-panel"
			>
				<div class="content-wrapper">
					<div class="brand-header">
						<h2 class="welcome-text">Bem-vindo ao</h2>
						<h1 class="brand-text">Data Flow</h1>
						<p class="welcome-subtext">
							Entre na sua conta para continuar gerenciando suas demandas e projetos.
						</p>
					</div>

					<Message
						v-if="motd !== ''"
						class="is-hidden-tablet mbe-3"
					>
						{{ motd }}
					</Message>
					
					<slot />

					<Legal class="legal-footer" />
				</div>
			</main>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

import Logo from '@/components/home/Logo.vue'
import Message from '@/components/misc/Message.vue'
import Legal from '@/components/misc/Legal.vue'

import { useTitle } from '@/composables/useTitle'
import { useConfigStore } from '@/stores/config'

defineProps<{
	showApiConfig?: boolean;
}>()

const configStore = useConfigStore()
const motd = computed(() => configStore.motd)

const route = useRoute()
const { t } = useI18n({ useScope: 'global' })
const title = computed(() =>
	route.meta?.title ? t(route.meta.title as string) : '',
)
useTitle(() => title.value)
</script>

<style lang="scss" scoped>
.no-auth-wrapper {
	background: #f0f4f9;
	min-block-size: 100vh;
	display: flex;
	place-items: center;
	justify-content: center;
	padding: 1rem;
	box-sizing: border-box;

	@media screen and (max-width: $tablet) {
		padding: 0.5rem;
	}
}

.noauth-container {
	max-inline-size: 880px;
	inline-size: 100%;
	display: flex;
	background-color: #ffffff;
	border-radius: 16px;
	box-shadow: 0 12px 32px rgba(15, 23, 42, 0.08);
	overflow: hidden;
	
	@media screen and (max-width: $tablet) {
		flex-direction: column;
		border-radius: 12px;
	}
}

.image-panel {
	inline-size: 50%;
	padding: 1.75rem 2rem;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	background: linear-gradient(135deg, #0b1528 0%, #1e3a8a 100%);
	position: relative;
	overflow: hidden;

	@media screen and (max-width: $tablet) {
		display: none;
	}

	&::before {
		content: "";
		position: absolute;
		top: 0; right: 0; bottom: 0; left: 0;
		background: url("https://images.unsplash.com/photo-1542224566-6e85f2e6772f?q=80&w=2000&auto=format&fit=crop") center/cover no-repeat;
		opacity: 0.22;
		mix-blend-mode: overlay;
		pointer-events: none;
	}

	.image-panel-content {
		position: relative;
		z-index: 2;
		display: flex;
		flex-direction: column;
		height: 100%;
		justify-content: space-between;
	}

	.logo {
		margin-bottom: 1rem;
		filter: brightness(0) invert(1);
	}

	.hero-text {
		margin: auto 0;
		padding: 0.75rem 0;

		.image-title {
			color: #ffffff;
			font-size: 2.1rem;
			font-weight: 800;
			line-height: 1.15;
			margin-bottom: 0.75rem;
			letter-spacing: -0.02em;

			.highlight {
				color: #00f2fe;
			}
		}

		.image-subtitle {
			color: rgba(255, 255, 255, 0.85);
			font-size: 0.875rem;
			line-height: 1.45;
			max-width: 95%;
		}
	}

	.feature-chips {
		display: flex;
		flex-direction: column;
		gap: 0.45rem;
		margin-top: 1rem;

		.chip {
			display: flex;
			align-items: center;
			gap: 0.5rem;
			background: rgba(255, 255, 255, 0.08);
			border: 1px solid rgba(255, 255, 255, 0.15);
			border-radius: 8px;
			padding: 0.35rem 0.65rem;
			color: white;
			font-size: 0.78rem;
			font-weight: 500;
			backdrop-filter: blur(8px);

			.chip-icon {
				display: flex;
				align-items: center;
				justify-content: center;
				width: 22px;
				height: 22px;
				border-radius: 5px;
				background: rgba(255, 255, 255, 0.2);
				font-size: 0.75rem;
				flex-shrink: 0;
			}
		}
	}
}

.content-panel {
	inline-size: 50%;
	display: flex;
	flex-direction: column;
	padding: 1.75rem 2rem;
	background: #ffffff;
	position: relative;
	overflow-y: auto;

	@media screen and (max-width: $desktop) {
		padding: 1.5rem 1.5rem;
	}

	@media screen and (max-width: $tablet) {
		inline-size: 100%;
		padding: 1.5rem 1.25rem;
	}

	.content-wrapper {
		margin: auto;
		max-width: 340px;
		width: 100%;
		display: flex;
		flex-direction: column;
	}

	.brand-header {
		margin-bottom: 0.85rem;
	}

	.welcome-text {
		font-size: 0.875rem;
		color: var(--grey-600);
		font-weight: 500;
		margin-bottom: 0.1rem;
	}

	.brand-text {
		font-size: 1.85rem;
		font-weight: 800;
		color: var(--primary);
		margin-bottom: 0.25rem;
		letter-spacing: -0.02em;
		line-height: 1.1;
	}

	.welcome-subtext {
		color: var(--grey-500);
		font-size: 0.825rem;
		line-height: 1.4;
		margin-bottom: 0;
	}
}

.legal-footer {
	margin-top: 0.85rem;
	text-align: center;
	opacity: 0.65;
	font-size: 0.75rem;
}

:deep(.logo) {
	margin: 0 !important;
}
</style>
