<template>
	<div class="login-page">
		<Message
			v-if="confirmedEmailSuccess"
			variant="success"
			text-align="center"
			class="mbe-4"
		>
			{{ $t('user.auth.confirmEmailSuccess') }}
		</Message>
		<Message
			v-if="errorMessage"
			variant="danger"
			class="mbe-4"
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
			<div class="field-wrapper mbe-4">
				<label class="label">{{ $t('user.auth.usernameEmail') }}</label>
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
			
			<div class="field-wrapper mbe-5">
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
				<Password
					v-model="password"
					:validate-initially="validatePasswordInitially"
					:validate-min-length="false"
					@submit="submit"
				/>
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
			
			<div class="checkbox-wrapper mbe-5">
				<FormCheckbox
					v-model="rememberMe"
					:label="$t('user.auth.remember')"
				/>
			</div>

			<XButton
				:loading="isLoading"
				class="is-fullwidth login-action-btn"
				@click="submit"
			>
				{{ $t('user.auth.login') }}
				<i class="fas fa-arrow-right icon-right"></i>
			</XButton>
			
			<p
				v-if="registrationEnabled"
				class="create-account-wrapper mbs-5"
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
			class="sso-providers mbs-5"
		>
			<div class="divider"><span>Ou entre com</span></div>
			<XButton
				v-for="(p, k) in openidConnect.providers"
				:key="k"
				variant="secondary"
				class="is-fullwidth sso-btn mbs-3"
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

import {getErrorText} from '@/message'
import {getAutoRedirectProvider, redirectToProvider} from '@/helpers/redirectToProvider'
import {useRedirectToLastVisited} from '@/composables/useRedirectToLastVisited'
import {isDesktopApp} from '@/helpers/desktopAuth'
import {REDIRECT_HASH_PREFIX} from '@/constants/redirectHash'

import {useAuthStore, JUST_LOGGED_OUT_KEY} from '@/stores/auth'
import {useConfigStore} from '@/stores/config'

import {useTitle} from '@/composables/useTitle'

const {t} = useI18n({useScope: 'global'})
useTitle(() => t('user.auth.login'))

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const configStore = useConfigStore()
const {redirectIfSaved} = useRedirectToLastVisited()

const registrationEnabled = computed(() => configStore.auth.local.registrationEnabled)
const localAuthEnabled = computed(() => configStore.auth.local.enabled)
const ldapAuthEnabled = computed(() => configStore.auth.ldap.enabled)

const openidConnect = computed(() => configStore.auth.openidConnect)
const hasOpenIdProviders = computed(() => openidConnect.value.enabled && openidConnect.value.providers?.length > 0)

const isLoading = computed(() => authStore.isLoading)
const isDesktop = isDesktopApp()

const confirmedEmailSuccess = ref(false)
const errorMessage = ref('')
const password = ref('')
const validatePasswordInitially = ref(false)
const rememberMe = ref(false)

const authenticated = computed(() => authStore.authenticated)

onBeforeMount(() => {
	authStore.verifyEmail().then((confirmed) => {
		confirmedEmailSuccess.value = confirmed
	}).catch((e: Error) => {
		errorMessage.value = e.message
	})

	if (authenticated.value) {
		router.push({name: 'home'})
		return
	}

	const justLoggedOut = sessionStorage.getItem(JUST_LOGGED_OUT_KEY) !== null
	if (justLoggedOut) {
		sessionStorage.removeItem(JUST_LOGGED_OUT_KEY)
	}

	const autoRedirectProvider = getAutoRedirectProvider({
		localAuthEnabled: localAuthEnabled.value,
		ldapAuthEnabled: ldapAuthEnabled.value,
		openIdEnabled: openidConnect.value.enabled,
		providers: openidConnect.value.providers ?? [],
		isDesktopApp: isDesktop,
		justLoggedOut,
		hasCopyableRedirect: route.hash.startsWith(REDIRECT_HASH_PREFIX),
	})
	if (autoRedirectProvider) {
		redirectToProvider(autoRedirectProvider)
	}
})

const usernameValid = ref(true)
const usernameRef = ref<HTMLInputElement | null>(null)
const validateUsernameField = useDebounceFn(() => {
	usernameValid.value = usernameRef.value?.value !== ''
}, 100)


const needsTotpPasscode = computed(() => authStore.needsTotpPasscode)
const totpPasscode = ref<HTMLInputElement | null>(null)

async function submit() {
	errorMessage.value = ''
	const credentials: any = {
		username: usernameRef.value?.value,
		password: password.value,
		longToken: rememberMe.value,
	}

	if (credentials.username === '' || credentials.password === '') {
		validateUsernameField()
		validatePasswordInitially.value = true
		return
	}

	if (needsTotpPasscode.value) {
		credentials.totpPasscode = totpPasscode.value?.value
	}

	try {
		await authStore.login(credentials)
		authStore.setNeedsTotpPasscode(false)

		redirectIfSaved()
	} catch (e: any) {
		if (e.response?.data.code === 1017 && !credentials.totpPasscode) {
			return
		}

		errorMessage.value = getErrorText(e)
	}
}
</script>

<style lang="scss" scoped>
.login-page {
	width: 100%;
}

.login-form-content {
	width: 100%;
}

.field-wrapper {
	:deep(.label) {
		font-weight: 600;
		color: var(--grey-800);
		margin-bottom: 0.5rem;
		font-size: 0.95rem;
	}
	
	:deep(input) {
		padding: 0.75rem 1rem;
		border-radius: 8px;
		border: 1px solid var(--grey-300);
		transition: all 0.2s ease;
		
		&:focus {
			border-color: var(--primary);
			box-shadow: 0 0 0 3px rgba(25, 115, 255, 0.1);
		}
	}
}

.label-with-link {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-block-end: .5rem;

	.label {
		margin-block-end: 0;
	}
	
	.reset-password-link {
		color: var(--primary);
		font-size: 0.9rem;
		font-weight: 500;
		text-decoration: none;
		
		&:hover {
			text-decoration: underline;
		}
	}
}

.checkbox-wrapper {
	:deep(.checkbox) {
		font-size: 0.95rem;
		color: var(--grey-700);
	}
}

.login-action-btn {
	padding: 1.25rem;
	font-size: 1.05rem;
	font-weight: 600;
	border-radius: 8px;
	box-shadow: 0 4px 12px rgba(25, 115, 255, 0.2);
	display: flex;
	justify-content: center;
	align-items: center;
	gap: 0.5rem;
	
	.icon-right {
		font-size: 0.9rem;
	}
}

.create-account-wrapper {
	text-align: center;
	color: var(--grey-600);
	font-size: 0.95rem;
	
	.create-account-link {
		color: var(--primary);
		font-weight: 600;
		margin-left: 0.25rem;
		text-decoration: none;
		
		&:hover {
			text-decoration: underline;
		}
	}
}

.divider {
	display: flex;
	align-items: center;
	text-align: center;
	color: var(--grey-400);
	font-size: 0.9rem;
	margin: 2rem 0;
	
	&::before, &::after {
		content: '';
		flex: 1;
		border-bottom: 1px solid var(--grey-200);
	}
	
	span {
		padding: 0 1rem;
	}
}

.sso-btn {
	border-radius: 8px;
	padding: 1rem;
	font-weight: 500;
}
</style>
