<template>
  <main class="main-content">
    <div class="form-section">
      <h2 class="section-title">{{ t('form.serviceConfig') }}</h2>
      <div class="form-grid">
        <div class="form-group">
          <label>{{ t('form.serviceName') }}</label>
          <input v-model="config.serviceName" :placeholder="t('form.serviceNamePlaceholder')" />
        </div>
        <div class="form-group">
          <label>{{ t('form.displayName') }}</label>
          <input v-model="config.displayName" :placeholder="t('form.displayNamePlaceholder')" @focus="autoFillFromServiceName('displayName')" />
        </div>
        <div class="form-group full-width">
          <label>{{ t('form.description') }}</label>
          <textarea v-model="config.description" rows="2" :placeholder="t('form.descriptionPlaceholder')" @focus="autoFillFromServiceName('description')"></textarea>
        </div>
      </div>
    </div>

    <div class="form-section">
      <h2 class="section-title">{{ t('form.application') }}</h2>
      <div class="form-grid">
        <div class="form-group full-width">
          <label>{{ t('form.appPath') }}</label>
          <div class="input-with-btn">
            <textarea v-model="config.appPath" rows="2" :placeholder="t('form.appPathPlaceholder')"></textarea>
            <button class="btn-sm btn-secondary btn-browse" @click="$emit('browse-app')">{{ t('form.browse') }}</button>
          </div>
          <span class="field-hint">{{ t('form.appPathHint') }}</span>
        </div>
        <div class="form-group full-width">
          <label>{{ t('form.arguments') }}</label>
          <textarea v-model="config.arguments" rows="4" :placeholder="t('form.argumentsPlaceholder')"></textarea>
          <span class="field-hint">{{ t('form.argumentsHint') }}</span>
        </div>
        <div class="form-group full-width">
          <label>{{ t('form.workDir') }}</label>
          <div class="input-with-btn">
            <input v-model="config.workDir" :placeholder="t('form.workDirPlaceholder')" />
            <button class="btn-sm btn-secondary btn-browse" @click="$emit('browse-dir')">{{ t('form.browse') }}</button>
          </div>
        </div>
      </div>
    </div>

    <div class="form-section">
      <h2 class="section-title">{{ t('form.startup') }}</h2>
      <div class="form-grid">
        <div class="form-group">
          <label>{{ t('form.account') }}</label>
          <input v-model="config.account" :placeholder="t('form.accountPlaceholder')" />
        </div>
        <div class="form-group">
          <label>{{ t('form.password') }}</label>
          <input v-model="config.password" type="password" :placeholder="t('form.passwordPlaceholder')" />
        </div>
        <div class="form-group">
          <label>{{ t('form.startType') }}</label>
          <select v-model="config.startType">
            <option value="auto">{{ t('form.startAuto') }}</option>
            <option value="demand">{{ t('form.startManual') }}</option>
            <option value="disabled">{{ t('form.startDisabled') }}</option>
          </select>
        </div>
      </div>
    </div>

    <div class="form-section">
      <h2 class="section-title">{{ t('form.logging') }}</h2>
      <div class="form-grid">
        <div class="form-group">
          <label>{{ t('form.logStdout') }}</label>
          <input v-model="config.logStdout" :placeholder="t('form.logStdoutPlaceholder')" />
        </div>
        <div class="form-group">
          <label>{{ t('form.logStderr') }}</label>
          <input v-model="config.logStderr" :placeholder="t('form.logStderrPlaceholder')" />
        </div>
        <div class="form-group checkbox-group">
          <label>
            <input type="checkbox" v-model="config.rotateLog" />
            {{ t('form.rotateLog') }}
          </label>
        </div>
      </div>
    </div>

    <div class="form-section">
      <h2 class="section-title">{{ t('form.recovery') }}</h2>
      <div class="form-grid">
        <div class="form-group">
          <label>{{ t('form.restartDelay') }}</label>
          <input v-model.number="config.restartDelay" type="number" min="0" />
        </div>
        <div class="form-group">
          <label>{{ t('form.restartTimeout') }}</label>
          <input v-model.number="config.restartTimeout" type="number" min="0" />
        </div>
      </div>
    </div>

    <div class="form-section">
      <h2 class="section-title">{{ t('form.environment') }}</h2>
      <div class="kv-editor">
        <div v-for="(val, key, idx) in config.environment" :key="key" class="kv-row">
          <input v-model="envKeys[idx]" class="kv-key" :placeholder="t('form.envKey')" @change="updateEnvKey(idx, envKeys[idx])" />
          <span class="kv-sep">=</span>
          <input v-model="config.environment[key]" class="kv-val" :placeholder="t('form.envValue')" />
          <button class="btn-sm btn-danger" @click="removeEnvKey(key)">&times;</button>
        </div>
        <div class="kv-row kv-add">
          <input v-model="newEnvKey" class="kv-key" :placeholder="t('form.envKey')" @keydown.enter="addEnvKey" />
          <span class="kv-sep">=</span>
          <input v-model="newEnvValue" class="kv-val" :placeholder="t('form.envValue')" @keydown.enter="addEnvKey" />
          <button class="btn-sm btn-primary" @click="addEnvKey">+</button>
        </div>
      </div>
    </div>

    <div class="form-section">
      <h2 class="section-title">{{ t('form.dependencies') }}</h2>
      <div class="dep-editor">
        <div v-for="(dep, idx) in config.dependencies" :key="idx" class="dep-row">
          <input v-model="config.dependencies[idx]" class="dep-input" :placeholder="t('form.depName')" />
          <button class="btn-sm btn-danger" @click="removeDependency(idx)">&times;</button>
        </div>
        <div class="dep-row dep-add">
          <input v-model="newDependency" class="dep-input" :placeholder="t('form.depName')" @keydown.enter="addDependency" />
          <button class="btn-sm btn-primary" @click="addDependency">+</button>
        </div>
      </div>
    </div>
  </main>
</template>

<script>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'ConfigForm',
  props: {
    config: { type: Object, required: true },
  },
  emits: ['browse-app', 'browse-dir'],
  setup(props) {
    const { t } = useI18n()

    const envKeys = ref([])
    const newEnvKey = ref('')
    const newEnvValue = ref('')
    const newDependency = ref('')

    function syncEnvKeys() {
      envKeys.value = Object.keys(props.config.environment || {})
    }

    watch(() => props.config.environment, () => syncEnvKeys(), { deep: true })
    syncEnvKeys()

    function addEnvKey() {
      if (!newEnvKey.value) return
      if (!props.config.environment) props.config.environment = {}
      props.config.environment[newEnvKey.value] = newEnvValue.value
      newEnvKey.value = ''
      newEnvValue.value = ''
      syncEnvKeys()
    }

    function removeEnvKey(key) {
      delete props.config.environment[key]
      syncEnvKeys()
    }

    function updateEnvKey(idx, newKey) {
      const oldKey = envKeys.value[idx]
      if (oldKey === newKey) return
      const val = props.config.environment[oldKey]
      delete props.config.environment[oldKey]
      props.config.environment[newKey] = val
      syncEnvKeys()
    }

    function addDependency() {
      if (!newDependency.value) return
      if (!props.config.dependencies) props.config.dependencies = []
      props.config.dependencies.push(newDependency.value)
      newDependency.value = ''
    }

    function removeDependency(idx) {
      props.config.dependencies.splice(idx, 1)
    }

    function autoFillFromServiceName(field) {
      if (props.config.serviceName && !props.config[field]) {
        props.config[field] = props.config.serviceName
      }
    }

    return {
      t, envKeys, newEnvKey, newEnvValue, newDependency,
      addEnvKey, removeEnvKey, updateEnvKey,
      addDependency, removeDependency,
      autoFillFromServiceName,
    }
  },
}
</script>
