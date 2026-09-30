<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useColorScheme } from '@/composables/useColorScheme'

import LogoFull from '@/assets/logo-full.svg?component'
import LogoFullWhite from '@/assets/logo-full-white.svg?component'
import LogoIcon from '@/assets/logo.svg?component'

const props = withDefaults(defineProps<{
	variant?: 'default' | 'white' | 'icon'
}>(), {
	variant: 'default',
})

const authStore = useAuthStore()
const { isDark } = useColorScheme()

const CustomLogo = computed(() => {
	const lightLogo = (window as any).CUSTOM_LOGO_URL
	const darkLogo = (window as any).CUSTOM_LOGO_URL_DARK

	if (!lightLogo && !darkLogo) return ''
	if (!darkLogo) return lightLogo
	if (!lightLogo) return darkLogo

	return (props.variant === 'white' || isDark.value) ? darkLogo : lightLogo
})

const ActiveLogo = computed(() => {
	if (props.variant === 'icon') {
		return LogoIcon
	}
	if (props.variant === 'white' || isDark.value) {
		return LogoFullWhite
	}
	return LogoFull
})
</script>

<template>
	<div class="dataflow-logo-wrapper">
		<component
			:is="ActiveLogo"
			v-if="!CustomLogo"
			alt="Data Flow"
			class="dataflow-brand-logo"
		/>
		<img
			v-show="CustomLogo"
			:src="CustomLogo"
			alt="Data Flow"
			class="dataflow-brand-logo"
		>
	</div>
</template>

<style lang="scss" scoped>
.dataflow-logo-wrapper {
	display: inline-flex;
	align-items: center;
	line-height: 1;
}

.dataflow-brand-logo {
	display: block;
	max-inline-size: 180px;
	max-block-size: 46px;
	width: auto;
	height: 38px;
	object-fit: contain;
}
</style>
