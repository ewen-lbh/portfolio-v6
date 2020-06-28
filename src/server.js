import sirv from 'sirv';
import polka from 'polka';
import compression from 'compression';
import * as sapper from '@sapper/server';

const { PORT, NODE_ENV } = process.env;
const dev = NODE_ENV === 'development';

import { register, init, getLocaleFromNavigator } from 'svelte-i18n'

register('en', () => import('../messages/en.json'))
register('fr', () => import('../messages/fr.json'))

init({
    fallbackLocale: 'en',
    initialLocale: getLocaleFromNavigator()
})

polka() // You can also use Express
	.use(
		compression({ threshold: 0 }),
		sirv('static', { dev }),
		sapper.middleware()
	)
	.listen(PORT, err => {
		if (err) console.log('error', err);
	});
