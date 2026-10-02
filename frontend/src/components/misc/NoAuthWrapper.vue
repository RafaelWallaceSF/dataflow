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
						variant="white"
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
				<div class="top-nav-lang">
					<div class="lang-selector">
						<span class="flag">🇧🇷</span>
						<span class="lang-label">Português</span>
						<i class="fas fa-chevron-down chevron"></i>
					</div>
				</div>

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
	min-height: 100vh;
	background: #EFF3F8;
	background-image: 
		radial-gradient(circle at 10% 20%, rgba(99, 102, 241, 0.07) 0%, transparent 45%),
		radial-gradient(circle at 90% 80%, rgba(168, 85, 247, 0.07) 0%, transparent 45%);
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 1.5rem 1rem;
	box-sizing: border-box;

	@media screen and (max-width: $tablet) {
		padding: 0.5rem;
	}
}

.noauth-container {
	max-width: 980px;
	width: 100%;
	display: flex;
	background-color: #ffffff;
	border-radius: 20px;
	box-shadow: 0 20px 50px rgba(15, 23, 42, 0.12), 0 1px 3px rgba(15, 23, 42, 0.05);
	overflow: hidden;
	
	@media screen and (max-width: $tablet) {
		flex-direction: column;
		border-radius: 14px;
	}
}

.image-panel {
	width: 50%;
	min-height: 580px;
	padding: 2.5rem 2.25rem;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	position: relative;
	background: linear-gradient(180deg, rgba(8, 16, 40, 0.88) 0%, rgba(11, 22, 55, 0.72) 100%), url('@/assets/login_bg_mountains.jpg') center/cover no-repeat;
	color: #ffffff;
	
	@media screen and (max-width: $tablet) {
		width: 100%;
		min-height: auto;
		padding: 2rem 1.5rem;
	}

	.image-panel-content {
		height: 100%;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
	}

	.logo {
		margin-bottom: 2.5rem;
	}

	.hero-text {
		margin: auto 0;
	}

	.image-title {
		font-size: 2.15rem;
		font-weight: 800;
		line-height: 1.2;
		color: #ffffff;
		margin-bottom: 1rem;
		letter-spacing: -0.02em;

		.highlight {
			background: linear-gradient(135deg, #00D2FF 0%, #A855F7 100%);
			-webkit-background-clip: text;
			-webkit-text-fill-color: transparent;
			display: inline-block;
		}
	}

	.image-subtitle {
		font-size: 0.95rem;
		line-height: 1.55;
		color: #CBD5E1;
		max-width: 90%;
		margin: 0;
	}

	.feature-chips {
		display: flex;
		gap: 0.5rem;
		margin-top: 2rem;

		@media screen and (max-width: $desktop) {
			flex-wrap: wrap;
		}

		.chip {
			flex: 1;
			display: flex;
			align-items: center;
			gap: 0.45rem;
			padding: 0.55rem 0.6rem;
			background: rgba(255, 255, 255, 0.12);
			backdrop-filter: blur(10px);
			border: 1px solid rgba(255, 255, 255, 0.18);
			border-radius: 8px;
			color: #ffffff;
			font-size: 0.75rem;
			font-weight: 500;
			white-space: nowrap;

			.chip-icon {
				display: flex;
				align-items: center;
				justify-content: center;
				width: 22px;
				height: 22px;
				background: rgba(255, 255, 255, 0.2);
				border-radius: 5px;
				font-size: 0.72rem;
				flex-shrink: 0;
			}
		}
	}
}

.content-panel {
	width: 50%;
	padding: 2rem 2.75rem 2.5rem 2.75rem;
	display: flex;
	flex-direction: column;
	background-color: #ffffff;
	position: relative;
	
	@media screen and (max-width: $tablet) {
		width: 100%;
		padding: 2rem 1.5rem;
	}

	.top-nav-lang {
		display: flex;
		justify-content: flex-end;
		margin-bottom: 1.25rem;

		.lang-selector {
			display: inline-flex;
			align-items: center;
			gap: 0.4rem;
			padding: 0.35rem 0.65rem;
			background: #F8FAFC;
			border: 1px solid #E2E8F0;
			border-radius: 8px;
			font-size: 0.8rem;
			font-weight: 500;
			color: #475569;
			cursor: default;

			.flag {
				font-size: 0.95rem;
			}

			.chevron {
				font-size: 0.65rem;
				color: #94A3B8;
			}
		}
	}

	.content-wrapper {
		width: 100%;
		margin: auto 0;
	}

	.brand-header {
		margin-bottom: 1.5rem;

		.welcome-text {
			font-size: 1.15rem;
			font-weight: 700;
			color: #1E293B;
			margin: 0 0 0.25rem 0;
		}

		.brand-text {
			font-size: 2.35rem;
			font-weight: 900;
			line-height: 1.1;
			letter-spacing: -0.02em;
			margin: 0 0 0.5rem 0;
			background: linear-gradient(135deg, #2563EB 0%, #7C3AED 100%);
			-webkit-background-clip: text;
			-webkit-text-fill-color: transparent;
		}

		.welcome-subtext {
			font-size: 0.875rem;
			color: #64748B;
			line-height: 1.45;
			margin: 0;
		}
	}

	.legal-footer {
		margin-top: 1.5rem;
		text-align: center;
		font-size: 0.75rem;
		color: #94A3B8;
	}
}
</style>
