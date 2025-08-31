export default deineNuxtRouteMiddleware((to) => {
    const { $router } = useNuxtApp()

    // check if usr is authenticated

    if (process.client) {
        const token = localStorage.getItem("auth-token")


        if (!token) {
            // redirect ot home witha uth moal trigger

            return naviagteTo("/?auth=login")
        }
    }
})