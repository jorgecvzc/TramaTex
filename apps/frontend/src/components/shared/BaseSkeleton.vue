<script setup lang="ts">
/**
 * BaseSkeleton.vue
 * Highly parameterizable component to show loading states.
 * Uses CSS classes defined in _skeletons.css.
 */
interface Props {
  type?: 'text' | 'title' | 'circle' | 'button' | 'row'
  width?: string
  height?: string
  borderRadius?: string
  animated?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  animated: true
})

const classes = [
  'skeleton',
  `skeleton-${props.type}`,
  { 'no-animation': !props.animated }
]

const styles = {
  width: props.width,
  height: props.height,
  borderRadius: props.borderRadius
}
</script>

<template>
  <div :class="classes" :style="styles"></div>
</template>

<style scoped>
.no-animation {
  animation: none;
}
</style>
