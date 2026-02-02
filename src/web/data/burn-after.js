/**
 * This file is part of CasPaste.
 * CasPaste is free software released under the MIT License.
 * See LICENSE.md file for details.
 */

// Handle burn-after custom input visibility
document.addEventListener('DOMContentLoaded', function() {
	var burnSelect = document.getElementById('burn-after');
	var customInput = document.getElementById('burn-custom');

	if (burnSelect && customInput) {
		burnSelect.addEventListener('change', function() {
			if (this.value === 'custom') {
				customInput.classList.add('active');
				customInput.required = true;
			} else {
				customInput.classList.remove('active');
				customInput.required = false;
				customInput.value = '';
			}
		});
	}
});
