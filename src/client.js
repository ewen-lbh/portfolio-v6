import * as sapper from '@sapper/app';
import { register, init, getLocaleFromNavigator } from 'svelte-i18n'

register('en', () => import('../messages/en.json'))
register('fr', () => import('../messages/fr.json'))

init({
    fallbackLocale: 'en',
    initialLocale: getLocaleFromNavigator()
})

sapper.start({
	target: document.querySelector('#sapper')
});
