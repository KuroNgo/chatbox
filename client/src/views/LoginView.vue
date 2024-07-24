<script lang="ts"></script>

<template>
    <div>
        <h2>Login</h2>
        <form @submit.prevent="">
            <div class="form-group">
                <label for="email">Email</label>
                <input type="text" v-model="email" name="username" class="form-control" :class="{'is-invalid': submitted && !email}"/>
                <div v-show="submitted && !email" class="invalid-feedback">Email is required</div>
            </div>
        </form>
    </div>
</template>


<style scoped>

</style>