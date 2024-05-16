import { createApp } from 'vue'
import App from './App.vue'
import router from './config/Router'
import i18n from './config/i18n';
import './style.scss';

createApp(App)
    .use(router)
    .use(i18n)
    .mount('#app')
