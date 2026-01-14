// This file is part of CasPaste.
// CasPaste is free software released under the MIT License.
// See LICENSE file for details.

document.addEventListener("DOMContentLoaded", () => {
	var editor = document.getElementById("editor");

	editor.addEventListener("keydown", (e) => {
		// If TAB pressed
		if (e.keyCode === 9) {
			e.preventDefault();

			let startOrig = editor.selectionStart;
			let endOrig = editor.selectionEnd;

			editor.value = editor.value.substring(0, startOrig) + "\t" + editor.value.substring(endOrig);

			editor.selectionStart = editor.selectionEnd = startOrig + 1;
		}
	});

	// Add HTML and CSS code for line numbers support
	var editorContainer = document.getElementById("editor-container");
	if (editorContainer) {
		editorContainer.insertAdjacentHTML("afterbegin", "<textarea id='editorLines' wrap='off' tabindex=-1 readonly>1</textarea>");
	} else {
		editor.insertAdjacentHTML("beforebegin", "<textarea id='editorLines' wrap='off' tabindex=-1 readonly>1</textarea>");
	}
	var editorLines = document.getElementById("editorLines");
	editorLines.rows = editor.rows;
	
	var styleSheet = document.createElement("style");
	styleSheet.innerText = `
	.form-group {
		position: relative;
	}

	#editor {
		margin-left: 60px;
		resize: none;
		width: calc(100% - 60px);
		min-width: calc(100% - 60px);
		max-width: calc(100% - 60px);
		line-height: 1.6;
		padding: 1.25rem;
		font-size: 15px;
		font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace;
	}

	#editorLines {
		display: block;
		user-select: none;
		text-align: right;
		position: absolute;
		left: 0;
		top: 0;
		resize: none;
		overflow: hidden;
		width: 60px;
		max-width: 60px;
		min-width: 60px;
		padding: 1.25rem 0.5rem 1.25rem 0.25rem;
		line-height: 1.6;
		font-size: 15px;
		font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace;
		border: none;
		border-right: 1px solid;
		border-radius: 8px 0 0 8px;
		background: inherit;
		color: inherit;
		opacity: 0.6;
		pointer-events: none;
		box-sizing: border-box;
	}

	#editor:focus-visible, #editorLines:focus-visible {
		outline: 0;
	}

	/* Container for editor needs relative positioning */
	#editor-container {
		position: relative;
		width: 100%;
	}
	
	.char-counter-container {
		margin-top: 0.5rem;
		text-align: right;
		font-size: 0.875rem;
		opacity: 0.7;
	}
	
	.file-input {
		display: none;
	}
	
	.file-label {
		display: inline-block;
		padding: 0.75rem 1.5rem;
		background: var(--color-element, #44475A);
		color: var(--color-font, #F8F8F2);
		border-radius: 4px;
		cursor: pointer;
		transition: background 0.2s;
	}
	
	.file-label:hover {
		background: var(--color-input-hover, #6272A4);
	}
	
	.form-help {
		margin-top: 0.5rem;
		font-size: 0.875rem;
		opacity: 0.7;
	}
`;
	document.head.appendChild(styleSheet);

	editorLines.addEventListener("focus", () => {
		editor.focus();
	});

	// Sync height of line numbers with editor
	function syncEditorHeight() {
		editorLines.style.height = editor.offsetHeight + 'px';
	}
	syncEditorHeight();

	// Use ResizeObserver if available for dynamic height sync
	if (window.ResizeObserver) {
		new ResizeObserver(syncEditorHeight).observe(editor);
	}

	// Add JS code for line numbers
	editor.addEventListener("scroll", () => {
		editorLines.scrollTop = editor.scrollTop;
		editorLines.scrollLeft = editor.scrollLeft;
	});

	var lineCountCache = 0;
	editor.addEventListener("input", () => {
		let lineCount = editor.value.split("\n").length;

		if (lineCountCache != lineCount) {
			editorLines.value = "";
			
			for (var i = 0; i < lineCount; i++) {
				editorLines.value = editorLines.value + (i + 1) + "\n";
			}
			
			lineCountCache = lineCount;
		}
	});

	// Add symbol counter
	document.getElementById("symbolCounterContainer").innerHTML = "<span id='symbolCounter' class='text-grey'></span>";
	var symbolCounter = document.getElementById("symbolCounter");

	function updateSymbolCounter() {
		var length = editor.value.length;

		if (editor.maxLength !== -1) {
			symbolCounter.textContent = length + "/" + editor.maxLength;
		} else {
			symbolCounter.textContent = length + "/∞";
		}
	}

	editor.addEventListener("input", updateSymbolCounter);
	updateSymbolCounter();
});

// Handle file upload and textarea mutual exclusivity
document.addEventListener('DOMContentLoaded', function() {
var fileInput = document.getElementById('paste-file');
var textarea = document.getElementById('editor');

if (fileInput && textarea) {
// When file is selected, disable textarea
fileInput.addEventListener('change', function() {
if (this.files && this.files.length > 0) {
textarea.disabled = true;
textarea.required = false;
textarea.classList.add('disabled');
} else {
textarea.disabled = false;
textarea.required = false;
textarea.classList.remove('disabled');
}
});

// When text is entered, disable file input
textarea.addEventListener('input', function() {
if (this.value.trim().length > 0) {
fileInput.disabled = true;
} else {
fileInput.disabled = false;
}
});
}
});
