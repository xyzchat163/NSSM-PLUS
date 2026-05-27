<template>
  <div class="app-container">
    <!-- Header -->
    <header class="app-header">
      <div class="header-left">
        <h1 class="app-title">{{ t('app.title') }}</h1>
        <span class="app-subtitle">{{ t('app.subtitle') }}</span>
      </div>
      <div class="header-actions">
        <!-- Language Switcher -->
        <div class="lang-switcher">
          <button
            class="btn-sm"
            :class="locale === 'en' ? 'btn-primary' : 'btn-secondary'"
            @click="switchLang('en')"
          >{{ t('lang.en') }}</button>
          <button
            class="btn-sm"
            :class="locale === 'zh' ? 'btn-primary' : 'btn-secondary'"
            @click="switchLang('zh')"
          >{{ t('lang.zh') }}</button>
        </div>
        <span class="header-sep"></span>
        <template v-if="configFilePath">
          <span class="config-file-label" :title="configFilePath">
            &#x1F4C4; {{ configFilePath.split(/[\\/]/).pop() }}
          </span>
          <span v-if="dirty" class="dirty-indicator" :title="t('app.unsaved')">&#x25CF;</span>
          <button class="btn-sm btn-secondary" @click="openInExplorer" title="Open in File Explorer">
            &#x1F4C2;
          </button>
        </template>
        <span class="header-sep" v-if="configFilePath"></span>
        <a class="header-link" href="https://zhengkai.blog.csdn.net/" target="_blank" title="CSDN Blog">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M3.5 8.4a1.4 1.4 0 0 0-1.4 1.4v4.2a1.4 1.4 0 0 0 1.4 1.4h1.4a1.4 1.4 0 0 0 1.4-1.4V9.8a1.4 1.4 0 0 0-1.4-1.4H3.5zm7.7-4.2a1.4 1.4 0 0 0-1.4 1.4v8.4a1.4 1.4 0 0 0 1.4 1.4h1.4a1.4 1.4 0 0 0 1.4-1.4V5.6a1.4 1.4 0 0 0-1.4-1.4h-1.4zm7.7 2.8a1.4 1.4 0 0 0-1.4 1.4v5.6a1.4 1.4 0 0 0 1.4 1.4h1.4a1.4 1.4 0 0 0 1.4-1.4V8.4a1.4 1.4 0 0 0-1.4-1.4h-1.4z"/></svg>
          CSDN
        </a>
        <a class="header-link" href="https://github.com/moshowgame/NSSM-PLUS" target="_blank" title="GitHub">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z"/></svg>
          GitHub
        </a>
      </div>
    </header>

    <div class="app-body">
      <ServiceList
        :display-services="displayServices"
        :selected-service="selectedService"
        @new="newConfig"
        @refresh="refreshServices"
        @select="selectService"
        @copy="copyService"
      />

      <ConfigForm
        :config="config"
        @browse-app="browseAppPath"
        @browse-dir="browseWorkDir"
      />
    </div>

    <ActionBar
      :config="config"
      :is-editing="isEditing"
      :config-file-path="configFilePath"
      :source="selectedServiceSource"
      @new="newConfig"
      @load="loadConfig"
      @save="saveConfig"
      @save-service="saveService"
      @install="installNewService"
      @reconfigure="reconfigureService"
      @start="startService"
      @stop="stopService"
      @restart="restartService"
      @check="checkService"
      @uninstall="removeService"
      @delete="deleteConfig"
    />

    <!-- Toast Notification -->
    <div v-if="toast.show" class="toast" :class="'toast-' + toast.type">
      {{ toast.message }}
    </div>

    <!-- Modal Overlay -->
    <div v-if="modal.show" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h3>{{ modal.title }}</h3>
        <p>{{ modal.message }}</p>
        <div class="modal-actions">
          <button class="btn-secondary" @click="closeModal">{{ t('modal.cancel') }}</button>
          <button :class="modal.confirmClass || 'btn-danger'" @click="modal.onConfirm">{{ t('modal.confirm') }}</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import ServiceList from './components/ServiceList.vue'
import ConfigForm from './components/ConfigForm.vue'
import ActionBar from './components/ActionBar.vue'

export default {
  name: 'App',
  components: { ServiceList, ConfigForm, ActionBar },
  setup() {
    const { t, locale } = useI18n()

    // State
    const services = ref([])
    const loadedServices = ref([])
    const configFilePath = ref('')
    const selectedService = ref('')
    const selectedServiceSource = ref('')
    const isEditing = ref(false)
    const dirty = ref(false)
    const STORAGE_CONFIG_KEY = 'nssm-plus-last-config'
    let autoRefreshTimer = null

    // Merged service list for sidebar display.
    // When a config file is loaded, only show services from that file (with real status if installed).
    // When no config file is loaded, show all installed services.
    const displayServices = computed(() => {
      const installedMap = {}
      for (const svc of services.value) {
        installedMap[svc.name] = svc
      }

      if (configFilePath.value) {
        // Config file mode: show only services from the file, enriched with real status
        return loadedServices.value.map(ls => {
          const installed = installedMap[ls.serviceName]
          if (installed) {
            return {
              name: installed.name,
              displayName: installed.displayName,
              status: installed.status,
              startType: installed.startType,
              appPath: installed.appPath || ls.appPath,
              source: 'installed',
            }
          }
          return {
            name: ls.serviceName,
            displayName: ls.displayName || ls.serviceName,
            status: 'Not Installed',
            startType: ls.startType || '-',
            appPath: ls.appPath,
            source: 'file',
          }
        })
      }

      // No config file: show all installed NSSM-Plus services
      return services.value.map(svc => ({
        name: svc.name,
        displayName: svc.displayName,
        status: svc.status,
        startType: svc.startType,
        appPath: svc.appPath,
        source: 'installed',
      }))
    })

    const defaultConfig = () => ({
      serviceName: '',
      displayName: '',
      description: '',
      appPath: '',
      arguments: '',
      workDir: '',
      startType: 'auto',
      account: '',
      password: '',
      environment: {},
      logStdout: '',
      logStderr: '',
      rotateLog: false,
      restartDelay: 0,
      restartTimeout: 30,
      dependencies: [],
    })

    const config = reactive(defaultConfig())

    // Track unsaved changes: set dirty flag when form fields change while a service is selected
    watch(config, () => {
      if (isEditing.value) {
        dirty.value = true
      }
    }, { deep: true })

    const toast = reactive({ show: false, message: '', type: 'info' })
    const modal = reactive({
      show: false,
      title: '',
      message: '',
      confirmClass: '',
      onConfirm: () => {},
      onCancel: () => {},
    })

    function call(method, ...args) {
      if (window.go) {
        return window.go.main.App[method](...args)
      }
      return Promise.reject(new Error('Wails runtime not available'))
    }

    function errorMsg(e) {
      if (typeof e === 'string') return e
      return e?.message || String(e)
    }

    function showToast(message, type = 'info') {
      toast.message = message
      toast.type = type
      toast.show = true
      setTimeout(() => { toast.show = false }, 3000)
    }

    function showModal(title, message, confirmClass, onConfirm, onCancel) {
      modal.title = title
      modal.message = message
      modal.confirmClass = confirmClass
      modal.onConfirm = () => {
        modal.show = false
        if (onConfirm) onConfirm()
      }
      modal.onCancel = () => {
        if (onCancel) onCancel()
      }
      modal.show = true
    }

    function closeModal() {
      modal.show = false
      modal.onCancel()
    }

    // Guard an action: if there are unsaved changes, confirm before discarding.
    function guardAction(action) {
      if (!dirty.value) { action(); return }
      showModal(
        t('modal.unsavedTitle'),
        t('modal.unsavedMessage'),
        'btn-warning',
        () => {
          dirty.value = false
          action()
        }
      )
    }

    // #17 Language switcher
    function switchLang(lang) {
      locale.value = lang
      localStorage.setItem('nssm-plus-lang', lang)
    }

    // #10 File pickers
    async function browseAppPath() {
      try {
        const filePath = await call('ShowOpenAppDialog', t('form.appPath'))
        if (filePath) {
          config.appPath = filePath
        }
      } catch (e) {
        // Dialog cancelled
      }
    }

    async function browseWorkDir() {
      try {
        const dirPath = await call('ShowOpenDirectoryDialog', t('form.workDir'))
        if (dirPath) {
          config.workDir = dirPath
        }
      } catch (e) {
        // Dialog cancelled
      }
    }

    // Service operations
    async function refreshServices() {
      try {
        // Always fetch installed services for status display
        try {
          const result = await call('GetInstalledServices')
          services.value = result || []
        } catch (e) {
          services.value = []
        }
        // If config file is loaded, also reload it for up-to-date config data
        if (configFilePath.value) {
          const configs = await call('LoadConfigFromFile', configFilePath.value)
          loadedServices.value = configs || []
        }
      } catch (e) {
        showToast(t('toast.refreshFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    async function selectService(svc) {
      guardAction(async () => {
        // Disable dirty watch during config load to avoid false positives
        isEditing.value = false
        if (svc.source === 'file') {
          const cached = loadedServices.value.find(s => s.serviceName === svc.name)
          if (cached) {
            Object.assign(config, cached)
          }
        } else {
          try {
            const cfg = await call('GetServiceConfig', svc.name)
            if (cfg) {
              Object.assign(config, cfg)
            }
          } catch (e) {
            showToast(t('toast.loadConfigFailed') + ': ' + errorMsg(e), 'error')
          }
        }
        selectedService.value = svc.name
        selectedServiceSource.value = svc.source
        isEditing.value = true
        dirty.value = false
        await refreshServices()
      })
    }

    function statusClass(status) {
      switch (status) {
        case 'Running': return 'status-running'
        case 'Stopped': return 'status-stopped'
        case 'Not Installed': return 'status-file'
        default: return 'status-other'
      }
    }

    async function installNewService() {
      if (!config.serviceName || !config.appPath) {
        showToast(t('toast.nameAndPathRequired'), 'warning')
        return
      }
      try {
        await call('InstallService', JSON.parse(JSON.stringify(config)))
        showToast(t('toast.installed'), 'success')
        await refreshServices()
        selectedService.value = config.serviceName
        selectedServiceSource.value = 'installed'
        isEditing.value = true
      } catch (e) {
        showToast(t('toast.installFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    // #16 Enhanced error handling - no longer silently swallowing errors
    async function reconfigureService() {
      if (!config.serviceName || !config.appPath) {
        showToast(t('toast.nameAndPathRequired'), 'warning')
        return
      }
      try {
        const name = config.serviceName || selectedService.value
        try {
          await call('StopService', name)
        } catch (e) {
          console.log('StopService note:', errorMsg(e))
        }
        await call('ModifyService', selectedService.value, JSON.parse(JSON.stringify(config)))
        try {
          await call('StartService', name)
          showToast(t('toast.reconfiguredStarted'), 'success')
        } catch (e) {
          showToast(t('toast.reconfiguredNoStart') + ': ' + errorMsg(e), 'warning')
        }
        await refreshServices()
      } catch (e) {
        showToast(t('toast.reconfigureFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    async function startService() {
      try {
        const name = config.serviceName || selectedService.value
        await call('StartService', name)
        showToast(t('toast.started'), 'success')
        await refreshServices()
      } catch (e) {
        showToast(t('toast.startFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    async function stopService() {
      try {
        const name = config.serviceName || selectedService.value
        await call('StopService', name)
        showToast(t('toast.stopped'), 'success')
        await refreshServices()
      } catch (e) {
        showToast(t('toast.stopFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    async function restartService() {
      try {
        const name = config.serviceName || selectedService.value
        await call('RestartService', name)
        showToast(t('toast.restarted'), 'success')
        await refreshServices()
      } catch (e) {
        showToast(t('toast.restartFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    async function removeService() {
      const name = config.serviceName || selectedService.value
      showModal(
        t('action.uninstall'),
        `Are you sure you want to uninstall service "${name}"?`,
        'btn-danger',
        async () => {
          try {
            await call('RemoveService', name)
            const snapshot = JSON.parse(JSON.stringify(config))
            loadedServices.value = loadedServices.value.filter(s => s.serviceName !== name)
            loadedServices.value.push(snapshot)
            selectedServiceSource.value = 'file'
            showToast(t('toast.uninstalled'), 'success')
            await refreshServices()
          } catch (e) {
            const msg = errorMsg(e)
            showToast(t('toast.uninstallFailed') + ': ' + msg, 'error')
            if (msg.includes('marked for deletion')) {
              setTimeout(() => {
                showToast(t('toast.markedForDeletion'), 'warning')
              }, 500)
            }
          }
        }
      )
    }

    function deleteConfig() {
      const name = config.serviceName || selectedService.value
      if (!name) {
        showToast(t('toast.nameRequired'), 'warning')
        return
      }
      showModal(
        t('action.delete'),
        `Are you sure you want to delete "${name}"? This will remove the config entirely.`,
        'btn-danger',
        async () => {
          loadedServices.value = loadedServices.value.filter(s => s.serviceName !== name)
          newConfig()
          await refreshServices()
          showToast(t('toast.deleted'), 'success')
        }
      )
    }

    async function copyService(svc) {
      if (svc.source === 'file') {
        const cached = loadedServices.value.find(s => s.serviceName === svc.name)
        if (cached) {
          Object.assign(config, cached)
        }
      } else {
        try {
          const cfg = await call('GetServiceConfig', svc.name)
          if (cfg) {
            Object.assign(config, cfg)
          }
        } catch (e) {
          showToast(t('toast.loadConfigFailed') + ': ' + errorMsg(e), 'error')
          return
        }
      }
      config.serviceName = ''
      config.displayName = ''
      selectedService.value = ''
      selectedServiceSource.value = ''
      isEditing.value = true
      showToast(t('toast.copied'), 'info')
    }

    function newConfig() {
      guardAction(() => {
        Object.assign(config, defaultConfig())
        selectedService.value = ''
        selectedServiceSource.value = ''
        isEditing.value = false
        dirty.value = false
        loadedServices.value = []
        configFilePath.value = ''
        services.value = []
        localStorage.removeItem(STORAGE_CONFIG_KEY)
      })
    }

    async function checkService() {
      const name = config.serviceName || selectedService.value
      if (!name) return
      try {
        const status = await call('GetServiceStatus', name)
        showToast(`Service "${name}" exists, status: ${status}`, 'success')
        selectedService.value = name
        selectedServiceSource.value = 'installed'
        isEditing.value = true
      } catch (e) {
        showToast(`Service "${name}" does not exist: ${errorMsg(e)}`, 'warning')
      }
    }

    // --- Config file operations ---
    async function saveConfig() {
      try {
        const defaultName = configFilePath.value ? configFilePath.value.split(/[\\/]/).pop() : 'services.json'
        const filePath = await call('ShowSaveDialog', 'Save Config', defaultName)
        if (!filePath) return
        const current = JSON.parse(JSON.stringify(config))
        const allConfigs = loadedServices.value
          .filter(s => s.serviceName !== current.serviceName)
        if (current.serviceName) {
          allConfigs.unshift(current)
        }
        await call('SaveConfigToFile', filePath, allConfigs)
        configFilePath.value = filePath
        dirty.value = false
        localStorage.setItem(STORAGE_CONFIG_KEY, filePath)
        showToast(t('toast.saved', { count: allConfigs.length, file: filePath.split(/[\\/]/).pop() }), 'success')
      } catch (e) {
        showToast(t('toast.saveFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    async function saveService() {
      if (!configFilePath.value) {
        showToast(t('toast.noConfigFile'), 'warning')
        return
      }
      if (!config.serviceName) {
        showToast(t('toast.nameRequired'), 'warning')
        return
      }
      try {
        const current = JSON.parse(JSON.stringify(config))
        const allConfigs = loadedServices.value
          .filter(s => s.serviceName !== current.serviceName)
        allConfigs.unshift(current)
        await call('SaveConfigToFile', configFilePath.value, allConfigs)
        loadedServices.value = allConfigs
        dirty.value = false
        showToast(t('toast.saved', { count: 1, file: configFilePath.value.split(/[\\/]/).pop() }), 'success')
      } catch (e) {
        showToast(t('toast.saveFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    function loadConfig() {
      guardAction(async () => {
        try {
          const filePath = await call('ShowOpenDialog', 'Open Config File')
          if (!filePath) return
          const configs = await call('LoadConfigFromFile', filePath)
          if (!configs || configs.length === 0) {
            showToast(t('toast.noConfigs'), 'warning')
            return
          }
          isEditing.value = false
          Object.assign(config, configs[0])
          loadedServices.value = configs
          configFilePath.value = filePath
          selectedService.value = configs[0].serviceName
          selectedServiceSource.value = 'file'
          isEditing.value = true
          dirty.value = false
          localStorage.setItem(STORAGE_CONFIG_KEY, filePath)
          await refreshServices()
          showToast(t('toast.loaded', { count: configs.length, file: filePath.split(/[\\/]/).pop() }), 'success')
        } catch (e) {
          showToast(t('toast.loadConfigFailed') + ': ' + errorMsg(e), 'error')
        }
      })
    }

    async function openInExplorer() {
      if (!configFilePath.value) return
      try {
        await call('OpenInExplorer', configFilePath.value)
      } catch (e) {
        showToast(t('toast.fileOpenFailed') + ': ' + errorMsg(e), 'error')
      }
    }

    async function debugInfo() {
      console.group('[NSSM-Plus Debug]')
      console.log('Config:', JSON.parse(JSON.stringify(config)))
      console.log('Selected:', selectedService.value, '| Source:', selectedServiceSource.value)
      console.log('Installed Services:', JSON.parse(JSON.stringify(services.value)))
      console.log('Loaded Services:', JSON.parse(JSON.stringify(loadedServices.value)))
      console.log('Display Services:', JSON.parse(JSON.stringify(displayServices.value)))
      console.log('Config File:', configFilePath.value)
      console.groupEnd()
      showToast(t('toast.debugInfo'), 'info')
    }

    function autoFillFromServiceName(field) {
      if (config.serviceName && !config[field]) {
        config[field] = config.serviceName
      }
    }

    onMounted(async () => {
      const lastPath = localStorage.getItem(STORAGE_CONFIG_KEY)
      if (lastPath) {
        try {
          const configs = await call('LoadConfigFromFile', lastPath)
          if (configs && configs.length > 0) {
            loadedServices.value = configs
            configFilePath.value = lastPath
            Object.assign(config, configs[0])
            selectedService.value = configs[0].serviceName
            selectedServiceSource.value = 'file'
            isEditing.value = true
            dirty.value = false
          }
        } catch (e) {
          localStorage.removeItem(STORAGE_CONFIG_KEY)
        }
      }
      await refreshServices()
      autoRefreshTimer = setInterval(refreshServices, 10000)
    })

    onUnmounted(() => {
      if (autoRefreshTimer) {
        clearInterval(autoRefreshTimer)
        autoRefreshTimer = null
      }
    })

    return {
      services, loadedServices, displayServices,
      configFilePath, selectedService, selectedServiceSource,
      config, isEditing, dirty, toast, modal,
      locale, t, switchLang,
      refreshServices, selectService, copyService,
      installNewService, reconfigureService, startService, stopService, restartService, removeService,
      newConfig, deleteConfig, checkService, saveConfig, saveService, loadConfig, openInExplorer, debugInfo,
      browseAppPath, browseWorkDir,
      closeModal, guardAction,
    }
  }
}
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

/* Header */
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.app-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--accent);
}

.app-subtitle {
  font-size: 13px;
  color: var(--text-muted);
}

.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.config-file-label {
  font-size: 12px;
  color: var(--text-secondary);
  background: var(--bg-hover);
  padding: 3px 10px;
  border-radius: var(--radius);
  max-width: 240px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dirty-indicator {
  color: var(--warning);
  font-size: 14px;
  margin-left: 2px;
  animation: pulse 1.2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

/* Language Switcher */
.lang-switcher {
  display: flex;
  gap: 2px;
}
.lang-switcher .btn-sm {
  padding: 3px 8px;
  font-size: 11px;
  border-radius: 3px;
}

/* Body */
.app-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.header-sep {
  width: 1px;
  height: 18px;
  background: var(--border);
  margin: 0 4px;
}

.header-link {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-muted);
  text-decoration: none;
  padding: 2px 6px;
  border-radius: var(--radius);
  transition: color 0.15s, background 0.15s;
}
.header-link:hover {
  color: var(--accent);
  background: var(--bg-hover);
}

/* Sidebar */
.sidebar {
  width: 280px;
  min-width: 280px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  font-weight: 600;
  font-size: 14px;
  border-bottom: 1px solid var(--border);
}

.sidebar-actions {
  display: flex;
  gap: 6px;
}

.service-list {
  flex: 1;
  overflow-y: auto;
  padding: 6px;
}

.service-item {
  padding: 10px 12px;
  border-radius: var(--radius);
  cursor: pointer;
  transition: background 0.15s;
  margin-bottom: 2px;
}
.service-item:hover { background: var(--bg-hover); }
.service-item:hover .btn-copy { opacity: 1; }
.service-item.active { background: var(--bg-active); }

.service-item-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
}

.service-item-info {
  flex: 1;
  min-width: 0;
}

.btn-copy {
  opacity: 0;
  transition: opacity 0.15s;
  font-size: 14px;
  padding: 2px 6px;
  flex-shrink: 0;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-muted);
  cursor: pointer;
  line-height: 1;
}
.btn-copy:hover {
  background: var(--bg-hover);
  color: var(--accent);
  border-color: var(--accent);
}

.service-item-name {
  font-weight: 500;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 6px;
}

.source-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 5px;
  border-radius: 3px;
  background: rgba(33, 150, 243, 0.2);
  color: var(--info, #2196F3);
}

.service-item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.status-badge {
  padding: 1px 6px;
  border-radius: 3px;
  font-size: 11px;
  font-weight: 600;
}
.status-running { background: rgba(76, 175, 80, 0.2); color: var(--success); }
.status-stopped { background: rgba(244, 67, 54, 0.2); color: var(--danger); }
.status-file { background: rgba(158, 158, 158, 0.2); color: var(--text-muted); }
.status-other { background: rgba(255, 152, 0, 0.2); color: var(--warning); }
.start-type { color: var(--text-muted); }

.empty-state {
  padding: 40px 20px;
  text-align: center;
  color: var(--text-muted);
}
.empty-state .hint {
  font-size: 12px;
  margin-top: 6px;
}

/* Main Content */
.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.form-section {
  margin-bottom: 20px;
  background: var(--bg-secondary);
  border-radius: var(--radius);
  padding: 16px 20px;
  border: 1px solid var(--border);
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--accent);
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.form-group.full-width { grid-column: 1 / -1; }
.form-group.two-thirds { grid-column: 1 / 3; }

.form-group label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}

.field-hint {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
}

/* Input with browse button */
.input-with-btn {
  display: flex;
  gap: 8px;
  align-items: stretch;
}
.input-with-btn textarea,
.input-with-btn input {
  flex: 1;
  min-width: 0;
}
.btn-browse {
  align-self: flex-end;
  flex-shrink: 0;
  white-space: nowrap;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding-top: 6px;
}
.checkbox-group input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: var(--accent);
}

/* KV Editor (Environment Variables) */
.kv-editor {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.kv-row {
  display: flex;
  align-items: center;
  gap: 6px;
}
.kv-key {
  flex: 2;
  min-width: 0;
}
.kv-sep {
  color: var(--text-muted);
  font-weight: 600;
  flex-shrink: 0;
}
.kv-val {
  flex: 3;
  min-width: 0;
}

/* Dependency Editor */
.dep-editor {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.dep-row {
  display: flex;
  align-items: center;
  gap: 6px;
}
.dep-input {
  flex: 1;
  min-width: 0;
}

/* Action Bar */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border);
  flex-shrink: 0;
}

.action-left, .action-right {
  display: flex;
  gap: 8px;
}

button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Toast */
.toast {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: var(--radius);
  font-size: 13px;
  font-weight: 500;
  z-index: 1000;
  animation: slideIn 0.3s ease;
  box-shadow: var(--shadow);
  max-width: 500px;
}
.toast-success { background: var(--success); color: white; }
.toast-error { background: var(--danger); color: white; }
.toast-warning { background: var(--warning); color: white; }
.toast-info { background: var(--accent); color: white; }

@keyframes slideIn {
  from { transform: translateX(100%); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}

.modal {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 24px;
  min-width: 400px;
  max-width: 500px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}
.modal h3 { font-size: 16px; margin-bottom: 12px; }
.modal p { color: var(--text-secondary); margin-bottom: 20px; line-height: 1.5; }

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.icon {
  font-size: 15px;
}
</style>
