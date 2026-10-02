<script setup lang="ts">
import { computed } from 'vue'
import { useColorScheme } from '@/composables/useColorScheme'
import logoIconUrl from '@/assets/logo-icon.png'

const props = withDefaults(defineProps<{
	variant?: 'default' | 'white' | 'icon'
}>(), {
	variant: 'default',
})

const { isDark } = useColorScheme()
const isWhiteVariant = computed(() => props.variant === 'white' || isDark.value)
</script>

<template>
	<div class="dataflow-brand-logo" :class="{ 'is-white': isWhiteVariant, 'is-icon-only': variant === 'icon' }">
		<img :src="logoIconUrl" alt="Data Flow" class="brand-symbol" />
		<div v-if="variant !== 'icon'" class="brand-type">
			<span class="type-data">Data</span>
			<span class="type-flow">Flow</span>
		</div>
	</div>
</template>

<style lang="scss" scoped>
.dataflow-brand-logo {
	display: inline-flex;
	align-items: center;
	gap: 0.7rem;
	text-decoration: none;
	user-select: none;

	.brand-symbol {
		height: 38px;
		width: auto;
		object-fit: contain;
		display: block;
		filter: drop-shadow(0 4px 10px rgba(0, 0, 0, 0.15));
	}

	.brand-type {
		display: flex;
		align-items: baseline;
		gap: 0.25rem;
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Inter", sans-serif;
		font-weight: 800;
		font-size: 1.65rem;
		line-height: 1;
		letter-spacing: -0.03em;

		.type-data {
			color: #0F172A;
			transition: color 0.15s ease;
		}

		.type-flow {
			background: linear-gradient(135deg, #00C6FF 0%, #7C3AED 100%);
			-webkit-background-clip: text;
			-webkit-text-fill-color: transparent;
		}
	}

	&.is-white {
		.brand-type {
			.type-data {
				color: #FFFFFF;
			}
		}
	}
}
</style>
