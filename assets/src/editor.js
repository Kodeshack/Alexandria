let cm = CodeMirror.fromTextArea(document.getElementById("editor"), {
    mode:  "markdown",
    theme: "xq-light",
    indentUnit: 4,
    keyMap: "vim",
    lineWrapping: true,
    spellcheck: true
});

cm.getWrapperElement().classList.add("textarea")
