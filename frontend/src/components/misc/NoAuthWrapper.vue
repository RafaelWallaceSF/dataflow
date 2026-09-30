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
						width="180"
						height="50"
						:dark="true"
					/>
					
					<div class="hero-text">
						<h1 class="image-title">
							Organização<br>
							que impulsiona<br>
							<span class="highlight">resultados.</span>
						</h1>
						<p class="image-subtitle">
							Gestão de projetos e demandas em um<br>
							só lugar, para equipes mais focadas e<br>
							empresas mais produtivas.
						</p>
					</div>
					
					<div class="feature-chips">
						<div class="chip">
							<div class="chip-icon"><i class="fas fa-check-square"></i></div>
							<span>Mais<br>organização</span>
						</div>
						<div class="chip">
							<div class="chip-icon"><i class="fas fa-chart-bar"></i></div>
							<span>Mais<br>produtividade</span>
						</div>
						<div class="chip">
							<div class="chip-icon"><i class="fas fa-users"></i></div>
							<span>Mais<br>resultados</span>
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
					<div class="lang-selector-top">
						<!-- Optional language selector can go here if needed -->
					</div>
					
					<h2 class="welcome-text">Bem-vindo ao</h2>
					<h1 class="brand-text">Data Flow</h1>
					<p class="welcome-subtext">
						Entre na sua conta para continuar gerenciando<br>
						suas demandas e projetos.
					</p>

					<ApiConfig v-if="shouldShowApiConfig" />
					
					<Message
						v-if="motd !== ''"
						class="is-hidden-tablet mbe-4"
					>
						{{ motd }}
					</Message>
					
					<slot />
				</div>
				<Legal class="legal-footer" />
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
import ApiConfig from '@/components/misc/ApiConfig.vue'

import { useTitle } from '@/composables/useTitle'
import { useConfigStore } from '@/stores/config'
import { isDesktopApp } from '@/helpers/desktopAuth'

const props = withDefaults(
	defineProps<{
		showApiConfig?: boolean;
	}>(),
	{
		showApiConfig: false,
	},
)

const isDesktop = isDesktopApp()
const hasStoredApiUrl = isDesktop && localStorage.getItem('API_URL') !== null
const shouldShowApiConfig = computed(() => props.showApiConfig && (!isDesktop || hasStoredApiUrl))

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
	background: #f4f7fe;
	min-block-size: 100vh;
	display: flex;
	place-items: center;
	justify-content: center;
	padding: 2rem;

	@media screen and (max-width: $tablet) {
		padding: 1rem;
	}
}

.noauth-container {
	max-inline-size: 1200px;
	inline-size: 100%;
	min-block-size: 650px;
	display: flex;
	background-color: var(--white);
	border-radius: 24px;
	box-shadow: 0 20px 40px rgba(0, 0, 0, 0.08);
	overflow: hidden;
	
	@media screen and (max-width: $tablet) {
		flex-direction: column;
		border-radius: 16px;
	}
}

.image-panel {
	inline-size: 50%;
	padding: 3rem;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	background: linear-gradient(135deg, #0f172a 0%, #1e3a8a 100%);
	position: relative;
	overflow: hidden;

	@media screen and (max-width: $tablet) {
		display: none; // hide on mobile to save space
	}

	&::before {
		content: "";
		position: absolute;
		top: 0; right: 0; bottom: 0; left: 0;
		background: url("https://images.unsplash.com/photo-1542224566-6e85f2e6772f?q=80&w=2000&auto=format&fit=crop") center/cover no-repeat;
		opacity: 0.25;
		mix-blend-mode: overlay;
		pointer-events: none;
	}

	.image-panel-content {
		position: relative;
		z-index: 2;
		display: flex;
		flex-direction: column;
		height: 100%;
	}

	.logo {
		margin-bottom: auto;
		filter: brightness(0) invert(1);
	}

	.hero-text {
		margin-top: 4rem;
		margin-bottom: 3rem;

		.image-title {
			color: #ffffff;
			font-size: 3.5rem;
			font-weight: 800;
			line-height: 1.1;
			margin-bottom: 1.5rem;
			letter-spacing: -0.02em;

			.highlight {
				color: #00f2fe;
			}
		}

		.image-subtitle {
			color: rgba(255, 255, 255, 0.85);
			font-size: 1.15rem;
			line-height: 1.5;
			max-width: 90%;
		}
	}

	.feature-chips {
		display: flex;
		gap: 1rem;
		margin-top: auto;

		.chip {
			display: flex;
			align-items: center;
			gap: 0.75rem;
			background: rgba(255, 255, 255, 0.1);
			border: 1px solid rgba(255, 255, 255, 0.2);
			border-radius: 12px;
			padding: 0.75rem 1rem;
			color: white;
			font-size: 0.85rem;
			font-weight: 500;
			line-height: 1.2;
			backdrop-filter: blur(10px);

			.chip-icon {
				display: flex;
				align-items: center;
				justify-content: center;
				width: 32px;
				height: 32px;
				border-radius: 8px;
				background: rgba(255, 255, 255, 0.2);
				font-size: 1rem;
			}
		}
	}
}

.content-panel {
	inline-size: 50%;
	display: flex;
	flex-direction: column;
	padding: 4rem;
	background: #ffffff;
	position: relative;

	@media screen and (max-width: $desktop) {
		padding: 3rem 2rem;
	}

	@media screen and (max-width: $tablet) {
		inline-size: 100%;
		padding: 2rem 1.5rem;
	}

	.content-wrapper {
		margin: auto 0;
		max-width: 420px;
		width: 100%;
	}

	.welcome-text {
		font-size: 1.2rem;
		color: var(--grey-600);
		font-weight: 500;
		margin-bottom: 0.25rem;
	}

	.brand-text {
		font-size: 3rem;
		font-weight: 800;
		color: var(--primary);
		margin-bottom: 1rem;
		letter-spacing: -0.02em;
	}

	.welcome-subtext {
		color: var(--grey-500);
		font-size: 1rem;
		line-height: 1.5;
		margin-bottom: 2.5rem;
	}
}

.legal-footer {
	margin-top: 2rem;
	text-align: center;
	opacity: 0.7;
}

// Reset the global logo from the body, as we put it inside the left pane
:deep(.logo) {
	margin: 0 !important;
}
</style>
