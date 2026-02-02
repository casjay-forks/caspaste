/**
 * This file is part of CasPaste.
 * CasPaste is free software released under the MIT License.
 * See LICENSE.md file for details.
 */

function copyToClipboard(text) {
	var tmp = document.createElement("textarea");
	var focus = document.activeElement;

	tmp.value = text;

	document.body.appendChild(tmp);
	tmp.select();
	document.execCommand("copy");
	document.body.removeChild(tmp);
	focus.focus();
}

function copyButton(element) {
	var result = "";

	var strings = element.parentNode.getElementsByTagName("code")[0].textContent.split("\n");
	var stringsLen = strings.length;
	var cutLen = stringsLen.toString().length;

	for (var i = 0; stringsLen > i; i++) {
		if (i !== 0) {
			result = result + "\n";
		}

		result = result + strings[i].slice(cutLen);
	}

	result = result.trim() + "\n";
	copyToClipboard(result);
}

document.addEventListener("DOMContentLoaded", function() {
	// Add CSS for copy button
	var newStyleSheet = "\
		pre {\
			position: relative;\
			overflow: auto;\
		}\
		pre button {\
			visibility: hidden;\
		}\
		pre:hover > button {\
			visibility: visible;\
		}\
		.copy-btn {\
			position: absolute;\
			top: 16px;\
			right: 16px;\
			margin: 0;\
			animation: fadeout .2s both;\
		}\
	";

	var styleSheet = document.createElement("style");
	styleSheet.innerText = newStyleSheet;
	document.head.appendChild(styleSheet);

	// Add copy button to all pre tags
	var preElements = document.getElementsByTagName("pre");

	for (var i = 0; preElements.length > i; i++) {
		var btn = document.createElement("button");
		btn.className = "button-green copy-btn";
		btn.textContent = "{{call .Translate `codeJS.Paste`}}";
		btn.addEventListener("click", function() {
			copyButton(this);
		});
		preElements[i].appendChild(btn);
	}
});
