function toggleEdit(messageId) {
    const textDiv = document.getElementById('content-' + messageId);
    const formDiv = document.getElementById('edit-form-' + messageId);
    const textarea = document.getElementById('textarea-' + messageId);

    if (formDiv.style.display === "none" || formDiv.style.display === "") {
        formDiv.style.display = "block";
        textDiv.style.display = "none";
        autoGrow(textarea);
    } else {
        formDiv.style.display = "none";
        textDiv.style.display = "block";
    }
}

function autoGrow(element) {
    element.style.height = "5px";
    element.style.height = element.scrollHeight + "px";
}
