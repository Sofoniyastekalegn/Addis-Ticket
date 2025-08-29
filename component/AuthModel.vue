<script setup>

const message = ref("Nuxt")

const hello = () => sayHello(message.value)


</script>

<template>
    <main>
        <h1>Demo with auto imports</h1>

        <form @submit.prevent="hello">
            <MyInput v-model="message" />
            <button type="submit">Say Hello</button>
        </form>
    </main>
</template>