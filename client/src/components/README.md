## name: Đặt tên cho component.
```
name: 'MyComponent'
```
## props: Định nghĩa các props mà component có thể nhận.
```
props: {
  title: String,
  count: {
    type: Number,
    required: true
  }
}
```
## data: Trả về một đối tượng chứa dữ liệu riêng của component.
```
data() {
  return {
    message: 'Hello, world!'
  };
}
```
## methods: Định nghĩa các phương thức (methods) của component.
```
methods: {
  increment() {
    this.count++;
  }
}
```
## computed: Định nghĩa các thuộc tính tính toán (computed properties).
```
computed: {
  doubleCount() {
    return this.count * 2;
  }
}
```
## watch: Theo dõi các thay đổi của thuộc tính hoặc dữ liệu.
```
watch: {
  count(newVal, oldVal) {
    console.log(`Count changed from ${oldVal} to ${newVal}`);
  }
}
```
## emits: Khai báo các sự kiện mà component có thể phát ra.
```
emits: ['update', 'delete']
```
## components: Đăng ký các component con.
```
components: {
  ChildComponent
}
```
## setup: Định nghĩa thành phần logic sử dụng Composition API.

```
import { ref, computed } from 'vue';

export default defineComponent({
  name: 'MyComponent',
  props: {
    title: String,
    count: {
      type: Number,
      required: true
    }
  },
  setup(props) {
    const message = ref('Hello, world!');
    const doubleCount = computed(() => props.count * 2);

    function increment() {
      props.count++;
    }

    return {
      message,
      doubleCount,
      increment
    };
  }
});
```

## template: Định nghĩa template trực tiếp trong JavaScript (ít phổ biến hơn).
```
template: `
  <div>
    <h1>{{ title }}</h1>
    <p>{{ message }}</p>
  </div>
`
```