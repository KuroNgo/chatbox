# HOOKS
## ref và reactive: Để tạo ra các biến phản ứng.
```
import { ref, reactive } from 'vue';

const count = ref(0);
const state = reactive({
  message: 'Hello, world!'
});
```
## computed: Để tạo ra các thuộc tính tính toán.
```
import { computed } from 'vue';

const doubleCount = computed(() => count.value * 2);
```
## watch và watchEffect: Để theo dõi các thay đổi trong các biến hoặc trạng thái.
```
import { watch, watchEffect } from 'vue';

watch(count, (newVal, oldVal) => {
  console.log(`Count changed from ${oldVal} to ${newVal}`);
});

watchEffect(() => {
  console.log(`Count is now ${count.value}`);
});
```
## Lifecycle Hooks: Để xử lý các sự kiện trong vòng đời của component.
```
import { onMounted, onUnmounted, onUpdated } from 'vue';

onMounted(() => {
  console.log('Component mounted');
});

onUnmounted(() => {
  console.log('Component unmounted');
});

onUpdated(() => {
  console.log('Component updated');
});
```
## provide và inject: Để chia sẻ dữ liệu giữa các component tổ tiên và con cháu.
```
import { provide, inject } from 'vue';

// Trong component tổ tiên
provide('message', 'Hello from ancestor');

// Trong component con
const message = inject('message');
console.log(message); // 'Hello from ancestor'
```

# LIFECYCLE HOOK
## onBeforeMount: Được gọi trước khi component được gắn vào DOM.
```
import { onBeforeMount } from 'vue';

onBeforeMount(() => {
  console.log('Component is about to be mounted');
});
```
## onMounted: Được gọi sau khi component đã được gắn vào DOM.
```
import { onMounted } from 'vue';

onMounted(() => {
  console.log('Component has been mounted');
});
```
## onBeforeUpdate: Được gọi trước khi component cập nhật do các thay đổi trong reactive data hoặc props.
```
import { onBeforeUpdate } from 'vue';

onBeforeUpdate(() => {
  console.log('Component is about to update');
});
```
## onUpdated: Được gọi sau khi component đã cập nhật do các thay đổi trong reactive data hoặc props.
```
import { onUpdated } from 'vue';

onUpdated(() => {
  console.log('Component has been updated');
});
```
## onBeforeUnmount: Được gọi ngay trước khi component bị hủy bỏ và gỡ khỏi DOM.
```
import { onBeforeUnmount } from 'vue';

onBeforeUnmount(() => {
  console.log('Component is about to be unmounted');
});
```
## onUnmounted: Được gọi sau khi component bị hủy bỏ và gỡ khỏi DOM.
```
import { onUnmounted } from 'vue';

onUnmounted(() => {
  console.log('Component has been unmounted');
});
```
## onErrorCaptured: Được gọi khi một lỗi không được bắt xảy ra trong cây component con.
```
import { onErrorCaptured } from 'vue';

onErrorCaptured((err, instance, info) => {
  console.log('Error captured:', err);
  return false; // Trả về false để ngăn chặn lỗi lan truyền thêm
});
```
## onRenderTracked: Được gọi khi sự theo dõi reactive được theo dõi trong quá trình render.
```
import { onRenderTracked } from 'vue';

onRenderTracked((e) => {
  console.log('Render tracked:', e);
});
```
## onRenderTriggered: Được gọi khi sự theo dõi reactive kích hoạt lại render.
```
import { onRenderTriggered } from 'vue';

onRenderTriggered((e) => {
  console.log('Render triggered:', e);
});
```