// Event handler, emits "send-message" if enter is pressed with no modifier key
const sendMessageOnEnter = (event) => {
    if (event.key === "Enter" && !(event.shiftKey || event.altKey || event.ctrlKey || event.metaKey)) {
        event.preventDefault();
        htmx.trigger(event.target, "send-message");
    }
};
