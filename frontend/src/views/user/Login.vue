<template>
	<div class="login-page">
		<Message
			v-if="confirmedEmailSuccess"
			variant="success"
			text-align="center"
			class="mbe-2"
		>
			{{ $t('user.auth.confirmEmailSuccess') }}
		</Message>
		<Message
			v-if="errorMessage"
			variant="danger"
			class="mbe-2"
		>
			{{ errorMessage }}
		</Message>

		<DesktopLogin v-if="isDesktop" />

		<form
			v-if="!isDesktop && (localAuthEnabled || ldapAuthEnabled)"
			id="loginform"
			@submit.prevent="submit"
			class="login-form-content"
		>
			<div class="field-wrapper mbe-3">
				<label class="label">{{ $t('user.auth.usernameEmail') }}</label>
				<div class="input-container">
					<i class="far fa-user input-icon"></i>
					<FormField
						id="username"
						ref="usernameRef"
						v-focus
						name="username"
						:placeholder="$t('user.auth.usernamePlaceholder')"
						required
						type="text"
						autocomplete="username"
						:error="usernameValid ? null : $t('user.auth.usernameRequired')"
						@keyup.enter="submit"
						@focusout="validateUsernameField()"
					/>
				</div>
			</div>
			
			<div class="field-wrapper mbe-3">
				<div class="label-with-link">
					<label
						class="label"
						for="password"
					>{{ $t('user.auth.password') }}</label>
					<RouterLink
						v-if="localAuthEnabled"
						:to="{ name: 'user.password-reset.request' }"
						class="reset-password-link"
					>
						{{ $t('user.auth.forgotPassword') }}
					</RouterLink>
				</div>
				<div class="input-container">
					<i class="fas fa-lock input-icon"></i>
					<Password
						v-model="password"
						:validate-initially="validatePasswordInitially"
						:validate-min-length="false"
						@submit="submit"
					/>
				</div>
			</div>
			
			<FormField
				v-if="needsTotpPasscode"
				id="totpPasscode"
				ref="totpPasscode"
				v-focus
				:label="$t('user.auth.totpTitle')"
				autocomplete="one-time-code"
				:placeholder="$t('user.auth.totpPlaceholder')"
				required
				type="text"
				inputmode="numeric"
				@keyup.enter="submit"
			/>
			
			<div class="checkbox-wrapper mbe-3">
				<FormCheckbox
					v-model="rememberMe"
					:label="$t('user.auth.remember')"
				/>
			</div>

			<button
				type="submit"
				:disabled="isLoading"
				class="submit-button"
			>
				<span>{{ isLoading ? 'Entrando...' : 'Entrar' }}</span>
				<i class="fas fa-arrow-right arrow-icon"></i>
			</button>
			
			<p
				v-if="registrationEnabled"
				class="create-account-wrapper mbs-3"
			>
				{{ $t('user.auth.noAccountYet') }}
				<RouterLink
					:to="{ name: 'user.register' }"
					class="create-account-link"
				>
					{{ $t('user.auth.createAccount') }}
				</RouterLink>
			</p>
		</form>

		<div
			v-if="!isDesktop && hasOpenIdProviders"
			class="sso-providers mbs-3"
		>
			<div class="divider"><span>Ou entre com</span></div>
			<XButton
				v-for="(p, k) in openidConnect.providers"
				:key="k"
				variant="secondary"
				class="is-fullwidth sso-btn mbs-2"
				@click="redirectToProvider(p)"
			>
				{{ $t('user.auth.loginWith', {provider: p.name}) }}
			</XButton>
		</div>
	</div>
</template>

<script setup lang="ts">
import {computed, onBeforeMount, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {useDebounceFn} from '@vueuse/core'

import Message from '@/components/misc/Message.vue'
import Password from '@/components/input/Password.vue'
import FormField from '@/components/input/FormField.vue'
import FormCheckbox from '@/components/input/FormCheckbox.vue'
import DesktopLogin from '@/views/user/DesktopLogin.vue'
import XButton from '@/components/input/Button.vue'

import {getErrorText} from '@/message'
import {useAuthStore} from '@/stores/auth'
import {useConfigStore} from '@/stores/config'
import type {IError} from '@/types/IError'
import type {IOpenIDConnectProvider} from '@/modelTypes/IOpenIDConnectProvider'

const authStore = useAuthStore()
const configStore = useConfigStore()

const route = useRoute()
const router = useRouter()
const {t} = useI18n()

const username = ref('')
const password = ref('')
const totpPasscode = ref('')
const rememberMe = ref(true)

const isLoading = ref(false)
const errorMessage = ref('')
const usernameValid = ref(true)
const validatePasswordInitially = ref(false)
const needsTotpPasscode = ref(false)

const isDesktop = computed(() => false)
const localAuthEnabled = computed(() => configStore.localAuthEnabled)
const ldapAuthEnabled = computed(() => configStore.ldapAuthEnabled)
const openidConnect = computed(() => configStore.openidConnect)
const hasOpenIdProviders = computed(() => openidConnect.value.enabled && openidConnect.value.providers?.length > 0)
const registrationEnabled = computed(() => configStore.registrationEnabled)
const confirmedEmailSuccess = computed(() => route.query.emailConfirmed === 'true')

onBeforeMount(() => {
	if (authStore.isAuthenticated) {
		router.push({name: 'home'})
	}
})

const validateUsernameField = useDebounceFn(() => {
	usernameValid.value = username.value.trim().length > 0
}, 200)

async function submit() {
	if (isLoading.value) return
	validatePasswordInitially.value = true
	validateUsernameField()

	if (!usernameValid.value || !password.value) return

	isLoading.value = true
	errorMessage.value = ''

	try {
		await authStore.login({
			username: username.value,
			password: password.value,
			totpPasscode: totpPasscode.value,
			longToken: rememberMe.value,
		})

		const redirect = (route.query.redirect as string) || {name: 'home'}
		router.push(redirect)
	} catch (e: any) {
		const err = e as IError
		if (err.error_code === 1018) {
			needsTotpPasscode.value = true
		} else {
			errorMessage.value = getErrorText(err, t('user.auth.loginFailed'))
		}
	} finally {
		isLoading.value = false
	}
}

function redirectToProvider(p: IOpenIDConnectProvider) {
	window.location.href = p.auth_url
}
</script>

<style lang="scss" scoped>
.login-page {
	width: 100%;
}

.login-form-content {
	display: flex;
	flex-direction: column;
}

.field-wrapper {
	display: flex;
	flex-direction: column;
	gap: 0.35rem;

	.label {
		font-size: 0.875rem;
		font-weight: 600;
		color: #334155;
		margin: 0;
	}

	.label-with-link {
		display: flex;
		align-items: center;
		justify-content: space-between;

		.reset-password-link {
			font-size: 0.8rem;
			font-weight: 500;
			color: #2563EB;
			text-decoration: none;

			&:hover {
				text-decoration: underline;
			}
		}
	}
}

.input-container {
	position: relative;
	display: flex;
	align-items: center;

	.input-icon {
		position: absolute;
		left: 0.85rem;
		top: 50%;
		transform: translateY(-50%);
		color: #94A3B8;
		font-size: 0.9rem;
		pointer-events: none;
		z-index: 2;
	}

	:deep(.input-container),
	:deep(.form-field) {
		width: 100%;
		margin: 0;
	}

	:deep(input) {
		width: 100%;
		padding-left: 2.35rem !important;
		height: 44px;
		border-radius: 8px;
		border: 1px solid #CBD5E1;
		font-size: 0.9rem;
		color: #1E293B;
		transition: all 0.2s ease;

		&:focus {
			border-color: #2563EB;
			box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
			outline: none;
		}

		&::placeholder {
			color: #94A3B8;
		}
	}

	:deep(.password-input-wrapper) {
		width: 100%;
		position: relative;

		input {
			padding-left: 2.35rem !important;
			padding-right: 2.5rem !important;
		}

		button {
			position: absolute;
			right: 0.75rem;
			top: 50%;
			transform: translateY(-50%);
			background: transparent;
			border: none;
			color: #94A3B8;
			cursor: pointer;
		}
	}
}

.checkbox-wrapper {
	display: flex;
	align-items: center;
	font-size: 0.875rem;
	color: #475569;
	margin-top: 0.25rem;
}

.submit-button {
	width: 100%;
	display: inline-flex;
	align-items: center;
	justify-content: center;
	gap: 0.5rem;
	height: 44px;
	padding: 0 1.5rem;
	background: #2563EB;
	color: #ffffff;
	border: none;
	border-radius: 8px;
	font-size: 0.95rem;
	font-weight: 600;
	cursor: pointer;
	transition: all 0.2s ease;
	box-shadow: 0 4px 14px rgba(37, 99, 235, 0.35);

	&:hover:not(:disabled) {
		background: #1D4ED8;
		box-shadow: 0 6px 18px rgba(37, 99, 235, 0.45);
		transform: translateY(-1px);
	}

	&:disabled {
		opacity: 0.65;
		cursor: not-allowed;
	}

	.arrow-icon {
		font-size: 0.85rem;
		transition: transform 0.2s ease;
	}

	&:hover .arrow-icon {
		transform: translateX(3px);
	}
}

.create-account-wrapper {
	text-align: center;
	font-size: 0.85rem;
	color: #64748B;
	margin: 0;

	.create-account-link {
		color: #2563EB;
		font-weight: 600;
		text-decoration: none;
		margin-left: 0.25rem;

		&:hover {
			text-decoration: underline;
		}
	}
}

.divider {
	display: flex;
	align-items: center;
	text-align: center;
	margin: 1rem 0;
	color: #94A3B8;
	font-size: 0.75rem;

	&::before, &::after {
		content: '';
		flex: 1;
		border-bottom: 1px solid #E2E8F0;
	}

	span {
		padding: 0 0.5rem;
	}
}
</style>
