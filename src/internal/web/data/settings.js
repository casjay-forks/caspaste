// CasPaste Settings Manager with localStorage persistence
// Handles theme and language preferences without requiring user accounts

(function() {
	'use strict';

	const STORAGE_KEY = 'caspaste_settings';
	
	// Load settings from localStorage
	function loadSettings() {
		try {
			const stored = localStorage.getItem(STORAGE_KEY);
			return stored ? JSON.parse(stored) : {};
		} catch (e) {
			console.warn('Failed to load settings from localStorage:', e);
			return {};
		}
	}

	// Save settings to localStorage
	function saveSettings(settings) {
		try {
			localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
		} catch (e) {
			console.warn('Failed to save settings to localStorage:', e);
		}
	}

	// Initialize on page load - only restore saved values
	document.addEventListener('DOMContentLoaded', function() {
		const settings = loadSettings();

		// Restore theme selection in dropdown
		const themeSelect = document.querySelector('select[name="theme"]');
		if (themeSelect && settings.theme) {
			themeSelect.value = settings.theme;
		}

		// Restore language selection in dropdown
		const langSelect = document.querySelector('select[name="lang"]');
		if (langSelect && settings.lang) {
			langSelect.value = settings.lang;
		}

		// Save to localStorage when form is submitted
		const form = document.querySelector('form[action="/settings"]');
		if (form) {
			form.addEventListener('submit', function() {
				const newSettings = {};
				if (themeSelect) newSettings.theme = themeSelect.value;
				if (langSelect) newSettings.lang = langSelect.value;
				saveSettings(newSettings);
			});
		}
	});
})();
