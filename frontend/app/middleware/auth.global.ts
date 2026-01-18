export default defineNuxtRouteMiddleware(async (to, from) => {
    const auth = useAuthStore()

    if (!auth.initialized) {
        await auth.fetchUser()
    }

    if (to.path === '/login' || to.path === '/register') {
        if (auth.isLoggedIn) return navigateTo("/")
        return
    }

    if (!auth.isLoggedIn) {
        return navigateTo('/login')
    }
})
