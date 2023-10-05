const handleKeypress = (event) => {
    if (event.key === "Enter" && !(event.shiftKey || event.altKey || event.ctrlKey || event.metaKey)) {
        event.preventDefault();
        htmx.trigger(event.target, "send-message");
    }
};
