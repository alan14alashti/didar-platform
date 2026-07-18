import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue'),
  },
  {
    path: '/experiences/:slug',
    name: 'experience-detail',
    component: () => import('../views/ExperienceDetailView.vue'),
    props: true,
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('../views/ProfileView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/trips',
    name: 'my-trips',
    component: () => import('../views/MyTripsView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/host',
    component: () => import('../views/host/HostLayout.vue'),
    meta: { requiresHost: true },
    children: [
      { path: '', name: 'host-dashboard', component: () => import('../views/host/HostDashboardView.vue') },
      { path: 'experiences', name: 'host-experiences', component: () => import('../views/host/HostExperiencesView.vue') },
      { path: 'experiences/new', name: 'host-experience-new', component: () => import('../views/host/HostExperienceFormView.vue') },
      { path: 'experiences/:id', name: 'host-experience-edit', component: () => import('../views/host/HostExperienceFormView.vue'), props: true },
      { path: 'bookings', name: 'host-bookings', component: () => import('../views/host/HostBookingsView.vue') },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../views/NotFoundView.vue'),
  },
]

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.requiresHost && !auth.isHost) {
    return auth.isAuthenticated ? { name: 'profile' } : { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { name: 'home' }
  }
  return true
})
