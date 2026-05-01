<template>
  <div 
    class="toast-item" 
    :class="type"
    role="alert"
    @mouseenter="pauseTimer"
    @mouseleave="resumeTimer"
  >
    <div class="toast-icon">
      <component :is="iconComponent" :size="20" />
    </div>
    <div class="toast-content">
      <p>{{ message }}</p>
    </div>
    <button class="toast-close" @click="$emit('close')" aria-label="Cerrar">
      <X :size="18" />
    </button>
    <div class="toast-progress" :style="{ width: progress + '%' }"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { 
  CheckCircle2, 
  AlertCircle, 
  AlertTriangle, 
  Info, 
  X 
} from 'lucide-vue-next';

const props = defineProps({
  message: {
    type: String,
    required: true,
  },
  type: {
    type: String,
    default: 'info',
    validator: (value: string) => ['success', 'error', 'warning', 'info'].includes(value),
  },
  duration: {
    type: Number,
    default: 5000,
  }
});

const emit = defineEmits(['close']);

const iconComponent = computed(() => {
  switch (props.type) {
    case 'success': return CheckCircle2;
    case 'error': return AlertCircle;
    case 'warning': return AlertTriangle;
    default: return Info;
  }
});

const progress = ref(100);
const remaining = ref(props.duration);
let lastUpdate = Date.now();
let timerId: any = null;

const updateProgress = () => {
  const now = Date.now();
  const delta = now - lastUpdate;
  lastUpdate = now;
  
  remaining.value -= delta;
  progress.value = (remaining.value / props.duration) * 100;
  
  if (remaining.value <= 0) {
    progress.value = 0;
    stopTimer();
    emit('close');
  } else {
    timerId = requestAnimationFrame(updateProgress);
  }
};

const startTimer = () => {
  lastUpdate = Date.now();
  timerId = requestAnimationFrame(updateProgress);
};

const stopTimer = () => {
  if (timerId) {
    cancelAnimationFrame(timerId);
    timerId = null;
  }
};

const pauseTimer = () => stopTimer();
const resumeTimer = () => {
  if (props.duration > 0) startTimer();
};

onMounted(() => {
  if (props.duration > 0) startTimer();
});

onUnmounted(() => stopTimer());
</script>

<style scoped>
.toast-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  margin-bottom: 0.75rem;
  overflow: hidden;
  min-width: 320px;
  max-width: 450px;
  border-left: 4px solid #cbd5e1;
  transition: all 0.3s ease;
}

.toast-item.success { border-left-color: #22c55e; }
.toast-item.error { border-left-color: #ef4444; }
.toast-item.warning { border-left-color: #f59e0b; }
.toast-item.info { border-left-color: #3b82f6; }

.toast-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.success .toast-icon { color: #22c55e; }
.error .toast-icon { color: #ef4444; }
.warning .toast-icon { color: #f59e0b; }
.info .toast-icon { color: #3b82f6; }

.toast-content {
  flex-grow: 1;
}

.toast-content p {
  margin: 0;
  font-size: 0.9375rem;
  color: #1e293b;
  font-weight: 500;
  line-height: 1.4;
}

.toast-close {
  background: transparent;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.toast-close:hover {
  background: #f1f5f9;
  color: #475569;
}

.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 3px;
  background: rgba(0, 0, 0, 0.05);
  width: 100%;
}

.success .toast-progress { background: rgba(34, 197, 94, 0.2); }
.error .toast-progress { background: rgba(239, 68, 68, 0.2); }
.warning .toast-progress { background: rgba(245, 158, 11, 0.2); }
.info .toast-progress { background: rgba(59, 130, 246, 0.2); }
</style>
