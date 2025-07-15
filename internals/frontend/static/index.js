import ajax from '@imacrayon/alpine-ajax';
import Alpine from 'alpinejs';

window.Alpine = Alpine;

Alpine.plugin(ajax);

Alpine.start();
