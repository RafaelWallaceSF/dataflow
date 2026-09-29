<script setup lang="ts">
import { computed } from 'vue'
import { useColorScheme } from '@/composables/useColorScheme'

const { isDark } = useColorScheme()

const CustomLogo = computed(() => {
	const lightLogo = window.CUSTOM_LOGO_URL
	const darkLogo = window.CUSTOM_LOGO_URL_DARK

	if (!lightLogo && !darkLogo) return ''
	if (!darkLogo) return lightLogo
	if (!lightLogo) return darkLogo

	return isDark.value ? darkLogo : lightLogo
})
</script>

<template>
	<div>
		<span
			v-if="!CustomLogo"
			class="logo dataflow-logo"
			aria-label="DataFlow"
		>
			DataFlow
		</span>
		<img
			v-show="CustomLogo"
			:src="CustomLogo"
			alt="DataFlow"
			class="logo"
		>
	</div>
</template>

<style lang="scss" scoped>
.logo {
	color: var(--logo-text-color);
	max-inline-size: 168px;
	max-block-size: 48px;
}

.dataflow-logo {
	display: inline-flex;
	align-items: center;
	font-family: $vikunja-font;
	font-weight: 800;
	font-size: 1.75rem;
	letter-spacing: -.04em;
	line-height: 1;
}
</style>
